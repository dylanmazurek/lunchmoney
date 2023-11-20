package models

type RequestOptions struct {
	Method      string
	Path        string
	QueryValues any
	ReqBody     any
}

type Response struct {
	Errors      *[]string `json:"errors,omitempty"`
	Updated     bool      `json:"updated,omitempty"`
	InsertedIds []int64   `json:"ids,omitempty"`

	Split *[]int64 `json:"split,omitempty"`

	Transaction  *Transaction   `json:"_,omitempty"`
	Transactions *[]Transaction `json:"transactions,omitempty"`

	Category   *Category   `json:",omitempty"`
	Categories *[]Category `json:"categories,omitempty"`

	Asset  *Asset   `json:"asset,omitempty"`
	Assets *[]Asset `json:"assets,omitempty"`
}
