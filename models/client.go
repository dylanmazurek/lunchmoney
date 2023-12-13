package models

import "net/http"

// type Response struct {
// 	Errors *[]string `json:"error,omitempty"`

// 	Updated  *bool `json:"updated,omitempty"`
// 	Inserted *bool `json:"inserted,omitempty"`

// 	Item *any   `json:",omitempty"`
// 	Ids  *[]int `json:"ids,omitempty"`

// 	Transactions *[]Transaction `json:"transactions,omitempty"`
// 	Categories   *[]Category    `json:"categories,omitempty"`
// 	Assets       *[]Asset       `json:"assets,omitempty"`
// }

type Request struct {
	HTTPRequest *http.Request
}

// type Response struct {
// 	Errors []Error `json:"error,omitempty"`

// 	Updated  *bool `json:"updated,omitempty"`
// 	Inserted *bool `json:"inserted,omitempty"`

// 	Item *any   `json:",omitempty"`
// 	Ids  *[]int `json:"ids,omitempty"`
// }

type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
