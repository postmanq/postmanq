package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
	"go.uber.org/fx"
	"time"
)

type EventServiceServerParams struct {
	fx.In
	Logger                  log.Logger
	WorkflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
}

func NewFxEventServiceServer(params *EventServiceServerParams) postmanqv1.EventServiceServer {
	return &eventServiceServer{
		logger:                  params.Logger,
		workflowExecutorFactory: params.WorkflowExecutorFactory,
	}
}

type eventServiceServer struct {
	logger                  log.Logger
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
}

func (s *eventServiceServer) ReceiveEvent(ctx context.Context, event *postmanqv1.Event) (*postmanqv1.Event, error) {
	executor := s.workflowExecutorFactory.Create(
		temporal.WithWorkflowType(temporal.WorkflowTypeSendEvent),
		temporal.WithWorkflowID(temporal.WorkflowTypeSendEvent, event.Uuid),
		temporal.WithWorkflowExecutionTimeout(time.Minute),
	)

	return executor.Execute(ctx, event)
}
