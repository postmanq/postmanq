package errors

import "fmt"

func CantCastTypeToComponent(component interface{}) error {
	return fmt.Errorf("can`t cast type=%T to component", component)
}
