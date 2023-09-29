package services_test

import (
	"bufio"
	"bytes"
	"context"
	"github.com/postmanq/postmanq/pkg/log"
	"github.com/postmanq/postmanq/pkg/log/internal/services"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerSuite))
}

type LoggerSuite struct {
	suite.Suite
}

func (s *LoggerSuite) TestAll() {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	logger := services.NewTestLogger(writer)

	logger.With("key", "value").Debug("test")
	writer.Flush()
	s.Contains(buffer.String(), `"key": "value"`)

	logger.Named("test_logger").Debugf("test2")
	writer.Flush()
	s.Contains(buffer.String(), `test_logger`)

	ctx := context.WithValue(context.Background(), log.CorrelationID, "1234567890")
	logger.WithCtx(ctx).Debug("test")
	writer.Flush()
	s.Contains(buffer.String(), `"`+log.CorrelationID+`": "1234567890"`)

	logger.WithCtx(ctx, "key1", "value1").Debug("test")
	writer.Flush()
	s.Contains(buffer.String(), `"key1": "value1"`)
	s.Contains(buffer.String(), `"`+log.CorrelationID+`": "1234567890"`)
}
