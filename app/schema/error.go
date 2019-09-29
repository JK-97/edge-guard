package schema

import "fmt"

// HTTPError is an error with status code
type HTTPError struct {
	Err  error
	Code int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf(e.Err.Error())
}
