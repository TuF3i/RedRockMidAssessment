package models

type Commander struct {
	Role string      `json:"role"`
	Msg  interface{} `json:"msg"`
}
