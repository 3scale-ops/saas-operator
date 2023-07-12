package util

import (
	"encoding/json"
	"fmt"
)

// WrapError can be used to provide information about the
// context where the error occured
func WrapError(context string, err error) error {
	return fmt.Errorf("[%s] %s", context, err.Error())
}

// Type MultiError can be used for iterative operations that
// must keep going even if errors are detected. It will return
// the list of all encountered errors. MultiError implements the
// Error interface and can be passed around as a normal error.
type MultiError []error

// Ensure the Error interface is implemented
var _ error = MultiError{}

func (me MultiError) Error() string {

	list := make([]string, 0, len(me))
	for _, err := range me {
		list = append(list, err.Error())
	}

	b, _ := json.Marshal(list)
	return string(b)
}

func (me MultiError) ErrorOrNil() error {
	if len(me) > 0 {
		return me
	}
	return nil
}
