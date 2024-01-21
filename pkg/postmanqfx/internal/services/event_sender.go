package services

import (
	"context"
	"errors"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
	"time"
)

type EventSenderFactoryParams struct {
	fx.In
	Ctx                     context.Context
	Logger                  log.Logger
	WorkflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
	ActivityExecutorFactory temporal.ActivityExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
}

func NewFxEventSenderFactory(params EventSenderFactoryParams) postmanq.EventSenderFactory {
	return &eventSenderFactory{
		ctx:                     params.Ctx,
		logger:                  params.Logger,
		workflowExecutorFactory: params.WorkflowExecutorFactory,
		activityExecutorFactory: params.ActivityExecutorFactory,
	}
}

type eventSenderFactory struct {
	ctx                     context.Context
	logger                  log.Logger
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
	activityExecutorFactory temporal.ActivityExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
}

func (f *eventSenderFactory) Create(pipeline *postmanq.Pipeline) postmanq.EventSender {
	return &eventSender{
		ctx:                     f.ctx,
		logger:                  f.logger,
		workflowExecutorFactory: f.workflowExecutorFactory,
		activityExecutorFactory: f.activityExecutorFactory,
		pipeline:                pipeline,
	}
}

type eventSender struct {
	ctx                     context.Context
	logger                  log.Logger
	workflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
	activityExecutorFactory temporal.ActivityExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
	pipeline                *postmanq.Pipeline
}

func (s *eventSender) SendEvent(ctx workflow.Context, event *postmanqv1.Event) error {
	err := s.executeOrResend(ctx, s.pipeline.Middlewares, event)
	if err != nil {
		return err
	}

	err = s.executeOrResend(ctx, s.pipeline.Senders, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventSender) executeOrResend(
	ctx workflow.Context,
	activities collection.Slice[postmanq.WorkflowPlugin],
	event *postmanqv1.Event,
) error {
	activityErr := s.execute(ctx, activities, event)
	if activityErr != nil {
		resendErr := s.resend(event)
		if resendErr != nil {
			return errors.Join(activityErr, resendErr)
		}

		return activityErr
	}

	return nil
}

func (s *eventSender) execute(
	ctx workflow.Context,
	activities collection.Slice[postmanq.WorkflowPlugin],
	event *postmanqv1.Event,
) error {
	for _, plugin := range activities.Entries() {
		executor := s.activityExecutorFactory.Create(plugin.GetActivityDescriptor().GetActivityType())
		_, err := executor.Execute(ctx, event)
		if err != nil {
			s.logger.Error(err)
			return err
		}
	}

	return nil
}

func (s *eventSender) resend(event *postmanqv1.Event) error {
	event.AttemptsCount++
	executor := s.workflowExecutorFactory.Create(
		temporal.WithWorkflowType(temporal.WorkflowTypeSendEvent),
		temporal.WithWorkflowID(temporal.WorkflowTypeSendEvent, event.Uuid),
		temporal.WithWorkflowExecutionTimeout(time.Minute),
		temporal.WithStartDelay(time.Duration(event.AttemptsCount)*time.Hour),
	)
	_, err := executor.Execute(s.ctx, event)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}
