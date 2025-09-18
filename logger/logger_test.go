package logger

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	config := &Config{
		logrus.TraceLevel,
		true,
		"logrus",
		time.Hour,
		time.Second * time.Duration(5)}
	logger, _ := New(config)
	i := 0
	for {
		i++
		logger.WithFields(logrus.Fields{
			"index": i,
			"age":   18,
		}).Info("info msg")
		time.Sleep(time.Second)
	}
}
