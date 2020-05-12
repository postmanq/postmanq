package errors

import "fmt"

func MXServersIsNotFound(hostname string) error {
	return fmt.Errorf("mx servers is not found by hostname=%s", hostname)
}
