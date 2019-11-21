package controller

import "fmt"

type BaseResp struct {
	Data interface{} `json:"data"`
	Desc string      `json:"desc"`
}

// HTTPError is an error with status code
type HTTPError struct {
	Err  error
	Code int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf(e.Err.Error())
}
