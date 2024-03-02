package services

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"github.com/postmanq/postmanq/pkg/plugins/serverfx/server"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"time"
)

func NewFxEventServiceServerFactory(
	logger log.Logger,
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanq.Event, *postmanq.Event],
) server.EventServiceServerFactory {
	return &eventServiceServerFactory{
		logger:                  logger,
		workflowExecutorFactory: workflowExecutorFactory,
	}
}

type eventServiceServerFactory struct {
	logger                  log.Logger
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanq.Event, *postmanq.Event]
}

func (f *eventServiceServerFactory) Create(ctx context.Context, pipeline postmanq.Pipeline) server.EventServiceServer {
	return &eventServiceServer{
		pipeline:                pipeline,
		logger:                  f.logger,
		workflowExecutorFactory: f.workflowExecutorFactory,
	}
}

type eventServiceServer struct {
	pipeline                postmanq.Pipeline
	logger                  log.Logger
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanq.Event, *postmanq.Event]
}

func (s *eventServiceServer) ReceiveEvent(ctx context.Context, event *postmanq.Event) (*postmanq.Event, error) {
	executor := s.workflowExecutorFactory.Create(
		temporal.WithWorkflowType(fmt.Sprintf("WorkflowType%s", s.pipeline.Queue)),
		temporal.WithWorkflowID(fmt.Sprintf("WorkflowType%s_%s_%d", s.pipeline.Queue, event.Uuid, event.AttemptsCount)),
		temporal.WithWorkflowExecutionTimeout(time.Minute),
	)

	return executor.Execute(ctx, event)
}
