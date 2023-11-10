package models

type RequestOptions struct {
	Method      string
	Path        string
	QueryValues any
	ReqBody     any
}
