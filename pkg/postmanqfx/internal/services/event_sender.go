package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
	"time"
)

type EventSenderFactoryParams struct {
	fx.In
	Ctx                          context.Context
	Logger                       log.Logger
	ChildWorkflowExecutorFactory temporal.ChildWorkflowExecutorFactory[*postmanq.Event, *postmanq.Event]
	ActivityExecutorFactory      temporal.ActivityExecutorFactory[*postmanq.Event, *postmanq.Event]
}

func NewFxEventSenderFactory(params EventSenderFactoryParams) postmanq.EventSenderFactory {
	return &eventSenderFactory{
		ctx:                          params.Ctx,
		logger:                       params.Logger,
		childWorkflowExecutorFactory: params.ChildWorkflowExecutorFactory,
		activityExecutorFactory:      params.ActivityExecutorFactory,
	}
}

type eventSenderFactory struct {
	ctx                          context.Context
	logger                       log.Logger
	childWorkflowExecutorFactory temporal.ChildWorkflowExecutorFactory[*postmanq.Event, *postmanq.Event]
	activityExecutorFactory      temporal.ActivityExecutorFactory[*postmanq.Event, *postmanq.Event]
}

func (f *eventSenderFactory) Create(pipeline *postmanq.Pipeline) postmanq.EventSender {
	return &eventSender{
		ctx:                          f.ctx,
		logger:                       f.logger,
		childWorkflowExecutorFactory: f.childWorkflowExecutorFactory,
		activityExecutorFactory:      f.activityExecutorFactory,
		pipeline:                     pipeline,
	}
}

type eventSender struct {
	ctx                          context.Context
	logger                       log.Logger
	childWorkflowExecutorFactory temporal.ChildWorkflowExecutorFactory[*postmanq.Event, *postmanq.Event]
	activityExecutorFactory      temporal.ActivityExecutorFactory[*postmanq.Event, *postmanq.Event]
	pipeline                     *postmanq.Pipeline
}

func (s *eventSender) SendEvent(ctx workflow.Context, event *postmanq.Event) (*postmanq.Event, error) {
	ctx = temporal.WithActivityOptions(ctx, temporal.NewActivityOptions(
		temporal.WithStartToCloseTimeout(time.Minute*30),
		temporal.WithRetryPolicy(time.Second*5, 3),
	))

	err := s.executeOrResend(ctx, s.pipeline.Middlewares, event)
	if err != nil {
		return nil, err
	}

	err = s.executeOrResend(ctx, s.pipeline.Senders, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventSender) executeOrResend(
	ctx workflow.Context,
	activities collection.Slice[postmanq.WorkflowPlugin],
	event *postmanq.Event,
) error {
	activityErr := s.execute(ctx, activities, event)
	if activityErr != nil {
		resendErr := s.resend(ctx, event)
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
	event *postmanq.Event,
) error {
	for _, plugin := range activities.Entries() {
		executor := s.activityExecutorFactory.Create(plugin.GetType())
		_, err := executor.Execute(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *eventSender) resend(ctx workflow.Context, event *postmanq.Event) error {
	event.AttemptsCount++
	executor := s.childWorkflowExecutorFactory.Create(
		temporal.WithWorkflowType(fmt.Sprintf("WorkflowType%s", s.pipeline.Queue)),
		temporal.WithWorkflowId(fmt.Sprintf("WorkflowType%s_%s_%d", s.pipeline.Queue, event.Uuid, event.AttemptsCount)),
		temporal.WithWorkflowExecutionTimeout(time.Minute),
		temporal.WithStartDelay(time.Duration(event.AttemptsCount)*time.Hour),
	)
	_, err := executor.Execute(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
