package errors

import (
	"fmt"
	"github.com/postmanq/postmanq/module/pipe/entity"
)

func ComponentDescriptorAlreadyDefined(name string) error {
	return fmt.Errorf("component_descriptor=%s is already defined", name)
}

func ComponentDescriptorNotDefined(name string) error {
	return fmt.Errorf("component_descriptor=%s is not defined", name)
}

func UnknownStageType(s *entity.Stage) error {
	return fmt.Errorf("stage_type=%s is unknown", s.Type)
}

func ComponentNotDefinedForStage(s *entity.Stage) error {
	return fmt.Errorf("component is not defined for stage=%s", s.Type)
}

func ComponentsNotDefinedForStage(s *entity.Stage) error {
	return fmt.Errorf("components is not defined for stage=%s", s.Type)
}

func ConstructNotDefinedForStage(s *entity.Stage) error {
	return fmt.Errorf("construct is not defined for stage=%s", s.Type)
}
