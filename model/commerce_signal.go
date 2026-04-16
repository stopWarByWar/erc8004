package model

// Commerce action types (from ERC-8183 events)
const (
	ActionJobCreated   = "job_created"
	ActionJobFunded    = "job_funded"
	ActionJobSubmitted = "job_submitted"
	ActionJobCompleted = "job_completed"
	ActionJobRejected  = "job_rejected"
	ActionJobExpired   = "job_expired"
	ActionProviderSet  = "provider_set"
	ActionBudgetSet    = "budget_set"
)

// Job statuses (previous_status context for terminal events)
const (
	StatusOpen      = "open"
	StatusFunded    = "funded"
	StatusSubmitted = "submitted"
	StatusCompleted = "completed"
	StatusRejected  = "rejected"
	StatusExpired   = "expired"
)

// Roles
const (
	RoleClient    = "client"
	RoleProvider  = "provider"
	RoleEvaluator = "evaluator"
	RolePlatform  = "platform"
)

// Signal polarity
const (
	PolarityPositive = "positive"
	PolarityNegative = "negative"
	PolarityNeutral  = "neutral"
)

// Signal certainty
const (
	CertaintyDefinitive = "definitive"
	CertaintyIndicative = "indicative"
)

type Signal struct {
	Polarity    string
	Weight      float64
	Certainty   string
	ShouldWrite bool
}

var noSignal = Signal{}

// DetermineSignal maps (action, previousStatus, role) to a reputation signal
// per the terminal signal matrix (architecture design §4.2) and process signal
// table (§4.3). Returns ShouldWrite=false when the role should not receive a
// record for this event.
func DetermineSignal(action, previousStatus, role string) Signal {
	switch action {
	case ActionJobCompleted:
		return signalCompleted(role)
	case ActionJobRejected:
		return signalRejected(previousStatus, role)
	case ActionJobExpired:
		return signalExpired(previousStatus, role)
	case ActionJobCreated:
		return signalCreated(role)
	case ActionJobFunded:
		return signalFunded(role)
	case ActionJobSubmitted:
		return signalSubmitted(role)
	case ActionProviderSet:
		// Audit-like process event; keep it neutral and low impact.
		if role == RoleProvider {
			return Signal{PolarityNeutral, 0, CertaintyIndicative, true}
		}
		return noSignal
	case ActionBudgetSet:
		// Budget changes are primarily client-side process context.
		if role == RoleClient {
			return Signal{PolarityNeutral, 0, CertaintyIndicative, true}
		}
		return noSignal
	default:
		return noSignal
	}
}

func signalCompleted(role string) Signal {
	switch role {
	case RoleProvider:
		return Signal{PolarityPositive, 1.0, CertaintyDefinitive, true}
	case RoleClient:
		return Signal{PolarityPositive, 1.0, CertaintyDefinitive, true}
	case RoleEvaluator:
		return Signal{PolarityNeutral, 0.5, CertaintyDefinitive, true}
	default:
		return noSignal
	}
}

func signalRejected(prevStatus, role string) Signal {
	switch prevStatus {
	case StatusOpen:
		if role == RoleClient {
			return Signal{PolarityNeutral, 0.2, CertaintyDefinitive, true}
		}
		return noSignal

	case StatusFunded:
		switch role {
		case RoleProvider:
			return Signal{PolarityNegative, 0.5, CertaintyDefinitive, true}
		case RoleClient:
			return Signal{PolarityNeutral, 0.3, CertaintyDefinitive, true}
		case RoleEvaluator:
			return Signal{PolarityNeutral, 0.5, CertaintyDefinitive, true}
		}

	case StatusSubmitted:
		switch role {
		case RoleProvider:
			return Signal{PolarityNegative, 1.0, CertaintyDefinitive, true}
		case RoleClient:
			return Signal{PolarityNeutral, 0.3, CertaintyDefinitive, true}
		case RoleEvaluator:
			return Signal{PolarityNeutral, 0.5, CertaintyDefinitive, true}
		}
	}
	return noSignal
}

func signalExpired(prevStatus, role string) Signal {
	switch prevStatus {
	case StatusFunded:
		switch role {
		case RoleProvider:
			return Signal{PolarityNegative, 1.0, CertaintyDefinitive, true}
		case RoleClient:
			return Signal{PolarityNeutral, 0.3, CertaintyDefinitive, true}
		default:
			return noSignal
		}

	case StatusSubmitted:
		switch role {
		case RoleProvider:
			return Signal{PolarityNeutral, 0.3, CertaintyDefinitive, true}
		case RoleClient:
			return Signal{PolarityNegative, 0.5, CertaintyDefinitive, true}
		case RoleEvaluator:
			return Signal{PolarityNegative, 0.7, CertaintyDefinitive, true}
		}
	}
	return noSignal
}

func signalCreated(role string) Signal {
	if role == RoleClient {
		return Signal{PolarityNeutral, 0.1, CertaintyIndicative, true}
	}
	return noSignal
}

func signalFunded(role string) Signal {
	if role == RoleClient {
		return Signal{PolarityPositive, 0.3, CertaintyIndicative, true}
	}
	return noSignal
}

func signalSubmitted(role string) Signal {
	if role == RoleProvider {
		return Signal{PolarityPositive, 0.3, CertaintyIndicative, true}
	}
	return noSignal
}
