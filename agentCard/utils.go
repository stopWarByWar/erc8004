package agentcard

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"trpc.group/trpc-go/trpc-a2a-go/server"
)

const defaultIPFSGateway = "https://ipfs.io/ipfs/"

func GetAgentProfile(tokenURL string) (*TokenURLResponse, error) {
	if tokenURL == "" {
		return nil, nil
	}
	if strings.HasPrefix(tokenURL, "https://") || strings.HasPrefix(tokenURL, "ipfs://") {
		return getAgentProfileFromTokenURL(tokenURL)
	}

	if strings.HasPrefix(tokenURL, "data:") {
		return getAgentProfileFromEncodedData(tokenURL)
	}

	if strings.Contains(tokenURL, "name") {
		return decodeAgentProfileData([]byte(tokenURL))
	}

	return nil, fmt.Errorf("invalid token URL: %s", tokenURL)
}

func getAgentProfileFromTokenURL(tokenURL string) (*TokenURLResponse, error) {
	body, err := fetchTokenURLBody(tokenURL)
	if err != nil {
		return nil, err
	}
	return decodeAgentProfileData(body)
}

func decodeAgentProfileData(data []byte) (*TokenURLResponse, error) {
	var tokenURLResponse TokenURLResponse
	err := json.Unmarshal(data, &tokenURLResponse)
	if err != nil {
		return nil, err
	}

	for i, service := range tokenURLResponse.Services {
		tokenURLResponse.Services[i].Name = strings.ToLower(service.Name)
	}
	return &tokenURLResponse, nil
}

// DataURLResult 表示 RFC 2397 data URL 解码结果。
type DataURLResult struct {
	MediaType string            // 如 "application/json"
	Params    map[string]string // 如 charset, enc 等
	Data      []byte            // 解码后的原始字节（已处理 base64 / percent-encode，若 enc=gzip 已解压）
}

// ParseDataURL 按 RFC 2397 解析任意 data: URL。
// 语法: data:[<mediatype>][;base64],<data>
// - mediatype 可含参数，如 application/json;charset=utf-8;enc=gzip
// - 若含 ;base64 则 payload 为 Base64；否则为 %xx 编码
// - 若参数 enc=gzip 则对解码后的字节做 gzip 解压
func ParseDataURL(dataURL string) (*DataURLResult, error) {
	const prefix = "data:"
	if !strings.HasPrefix(dataURL, prefix) {
		return nil, fmt.Errorf("invalid data URL: missing data: prefix")
	}
	s := dataURL[len(prefix):]
	idx := strings.Index(s, ",")
	if idx < 0 {
		idx = strings.Index(s, " ")
		if idx < 0 {
			return nil, fmt.Errorf("invalid data URL: no comma or space")
		}

	}
	header, payload := strings.TrimSpace(s[:idx]), s[idx+1:]
	if payload == "" {
		return nil, fmt.Errorf("invalid data URL: empty payload")
	}

	// 解析 header：mediatype 与若干 ;parameter，末尾可能为 ;base64
	parts := splitDataURLHeader(header)
	var base64Enc bool
	var mediatypeParts []string
	for i, p := range parts {
		if i == len(parts)-1 && strings.EqualFold(p, "base64") {
			base64Enc = true
			mediatypeParts = parts[:i]
			break
		}
		mediatypeParts = parts
	}
	mediaType := ""
	params := make(map[string]string)
	for i, p := range mediatypeParts {
		eq := strings.Index(p, "=")
		if eq < 0 {
			if i == 0 {
				mediaType = strings.TrimSpace(p)
			}
			continue
		}
		k := strings.TrimSpace(strings.ToLower(p[:eq]))
		v := strings.TrimSpace(p[eq+1:])
		if k == "" {
			continue
		}
		if i == 0 && mediaType == "" && !strings.Contains(p, "/") {
			mediaType = strings.TrimSpace(p)
			continue
		}
		if v != "" {
			params[k] = v
		}
	}
	if mediaType == "" && len(mediatypeParts) > 0 {
		mediaType = strings.TrimSpace(mediatypeParts[0])
	}
	if mediaType == "" {
		mediaType = "text/plain"
	}

	// 解码 payload
	var raw []byte
	var err error
	if base64Enc {
		raw, err = base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return nil, fmt.Errorf("invalid data URL base64: %w", err)
		}
	} else {
		// RFC 2397：非 base64 时使用 %xx 编码，不用 query 的 + 规则
		decoded, err := url.PathUnescape(payload)
		if err != nil {
			return nil, fmt.Errorf("invalid data URL percent-encoding: %w", err)
		}
		raw = []byte(decoded)
	}

	// 可选 gzip 解压（常见扩展 enc=gzip）
	if strings.EqualFold(params["enc"], "gzip") {
		zr, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("invalid data URL gzip: %w", err)
		}
		defer zr.Close()
		raw, err = io.ReadAll(zr)
		if err != nil {
			return nil, fmt.Errorf("data URL gzip read: %w", err)
		}
	}

	return &DataURLResult{MediaType: mediaType, Params: params, Data: raw}, nil
}

// splitDataURLHeader 按分号分割 header，但保留 type/subtype 为一段（第一个不含 = 的段）。
func splitDataURLHeader(header string) []string {
	var parts []string
	for {
		header = strings.TrimLeft(header, " \t")
		if header == "" {
			break
		}
		i := 0
		for i < len(header) && header[i] != ';' {
			i++
		}
		parts = append(parts, strings.TrimSpace(header[:i]))
		if i >= len(header) {
			break
		}
		header = header[i+1:]
	}
	return parts
}

// decodeDataURLPayload 解析 data URL 并返回解码后的 JSON 字节；仅当 mediatype 为 application/json（或兼容）时用于 Agent Profile。
func decodeDataURLPayload(dataURL string) (rawJSON []byte, err error) {
	res, err := ParseDataURL(dataURL)
	if err != nil {
		return nil, err
	}
	mt := strings.ToLower(res.MediaType)
	if mt != "application/json" && !strings.HasPrefix(mt, "application/json;") {
		return nil, fmt.Errorf("invalid data URL for agent profile: expected application/json, got %s", res.MediaType)
	}
	return res.Data, nil
}

func getAgentProfileFromEncodedData(endpoint string) (*TokenURLResponse, error) {
	body, err := decodeDataURLPayload(endpoint)
	if err != nil {
		return nil, err
	}

	var tokenURLResponse TokenURLResponse
	if err := json.Unmarshal(body, &tokenURLResponse); err != nil {
		return nil, fmt.Errorf("invalid encoded data: %w", err)
	}

	for i, service := range tokenURLResponse.Services {
		tokenURLResponse.Services[i].Name = strings.ToLower(service.Name)
	}
	return &tokenURLResponse, nil
}

func getAgentCardFromA2AEndpoint(endpoint string) (*server.AgentCard, error) {
	if !validateA2AEndpoint(endpoint) {
		return nil, fmt.Errorf("invalid A2A endpoint: %s", endpoint)
	}

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to get A2A endpoint response: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed, status: %s, status code: %d", response.Status, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to read A2A endpoint response body: %v", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body")
	}
	agentCard, err := unmarshalAgentCard(body)
	if err != nil {
		return nil, fmt.Errorf("data error: failed to unmarshal A2A endpoint response body: %v", err)
	}
	return agentCard, nil
}

func unmarshalAgentCard(body []byte) (*server.AgentCard, error) {
	var agentCard *server.AgentCard
	err := json.Unmarshal(body, &agentCard)
	if err != nil {
		return nil, fmt.Errorf("data error: failed to unmarshal A2A endpoint response body: %v", err)
	}
	return agentCard, nil
}

func validateA2AEndpoint(endpoint string) bool {
	//检查https://agent.example/.well-known/agent-card.json这个格式
	if !strings.HasSuffix(endpoint, "/.well-known/agent-card.json") || !strings.HasPrefix(endpoint, "https://") {
		return false
	}
	return true
}

func fetchTokenURLBody(tokenURL string) ([]byte, error) {
	resolvedURL, err := resolveTokenURL(tokenURL)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(resolvedURL)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to get token URL response: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed, status: %s, status code: %d", response.Status, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("network error: failed to read token URL response body: %v", err)
	}
	if len(body) == 0 {
		return nil, errors.New("empty token URL response body")
	}
	return body, nil
}

func resolveTokenURL(tokenURL string) (string, error) {
	if strings.HasPrefix(tokenURL, "ipfs://") {
		cidPath := strings.TrimPrefix(tokenURL, "ipfs://")
		if len(cidPath) == 0 {
			return "", fmt.Errorf("invalid ipfs token URL: %s", tokenURL)
		}
		gateway := os.Getenv("IPFS_GATEWAY_URL")
		if len(gateway) == 0 {
			gateway = defaultIPFSGateway
		} else if !strings.HasSuffix(gateway, "/") {
			gateway += "/"
		}
		return fmt.Sprintf("%s%s", gateway, cidPath), nil
	}
	if strings.HasPrefix(tokenURL, "https://") {
		return tokenURL, nil
	}
	return "", fmt.Errorf("unsupported token URL scheme: %s", tokenURL)
}

func encodeAgentProfileData(data []byte) string {
	return "data:application/json;base64," + base64.StdEncoding.EncodeToString(data)
}
