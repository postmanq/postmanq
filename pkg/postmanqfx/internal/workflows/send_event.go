package workflows

import (
	"github.com/postmanq/postmanq/pkg/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
	"go.temporal.io/sdk/workflow"
)

type SendEventWorkflow func(ctx workflow.Context, event *postmanqv1.Event) error

func (w SendEventWorkflow) GetWorkflowType() temporal.WorkflowType {
	return temporal.WorkflowTypeSendEvent
}

func NewFxSendEventWorkflow(
	eventSender postmanq.EventSender,
) SendEventWorkflow {
	return func(ctx workflow.Context, event *postmanqv1.Event) error {
		return eventSender.SendEvent(ctx, event)
	}
}
