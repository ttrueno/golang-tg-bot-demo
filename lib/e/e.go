package e

import "fmt"

func Wrap(msg string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf(msg+": %v", e)
}
