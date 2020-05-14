package errors

import "fmt"

func IPsIsNotFoundByMX(hostname string) error {
	return fmt.Errorf("ips is not found by mx=%s", hostname)
}
