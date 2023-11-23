package models

import "net/url"

type Request struct {
	Method      string
	Path        string
	QueryValues url.Values
	ReqBody     any
}

type Response struct {
	Errors *[]string `json:"error,omitempty"`

	Updated  *bool `json:"updated,omitempty"`
	Inserted *bool `json:"inserted,omitempty"`

	Item *any   `json:",omitempty"`
	Ids  *[]int `json:"ids,omitempty"`

	Transactions *[]Transaction `json:"transactions,omitempty"`
	Categories   *[]Category    `json:"categories,omitempty"`
	Assets       *[]Asset       `json:"assets,omitempty"`
}
