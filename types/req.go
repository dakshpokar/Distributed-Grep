package types

type Request struct {
	Req_type string `json:"req_type"`
	Data     string	`json:"data"`
	File	 string `json:"file"`
}
