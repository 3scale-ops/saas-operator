package util

import "fmt"

// WrapError can be used to provide information about the
// context where the error occured
func WrapError(context string, err error) error {
	return fmt.Errorf("[%s] %s", context, err.Error())
}
