package schema

type BaseResp struct {
	Data interface{} `json:"data"`
	Desc string      `json:"desc"`
}
