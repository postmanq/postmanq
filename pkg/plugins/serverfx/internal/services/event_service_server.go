package services

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"github.com/postmanq/postmanq/pkg/plugins/serverfx/server"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"time"
)

func NewFxEventServiceServerFactory(
	logger log.Logger,
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event],
) server.EventServiceServerFactory {
	return &eventServiceServerFactory{
		logger:                  logger,
		workflowExecutorFactory: workflowExecutorFactory,
	}
}

type eventServiceServerFactory struct {
	logger                  log.Logger
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
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
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
}

func (s *eventServiceServer) ReceiveEvent(ctx context.Context, event *postmanqv1.Event) (*postmanqv1.Event, error) {
	executor := s.workflowExecutorFactory.Create(
		temporal.WithWorkflowType(fmt.Sprintf("WorkflowType%s", s.pipeline.Queue)),
		temporal.WithWorkflowID(fmt.Sprintf("WorkflowType%s_%s_%d", s.pipeline.Queue, event.Uuid, event.AttemptsCount)),
		temporal.WithWorkflowExecutionTimeout(time.Minute),
	)

	return executor.Execute(ctx, event)
}
