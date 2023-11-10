package models

// AssetUpdateRequest
type AssetUpdateRequest struct {
	Balance string `json:"balance"`
}

type AssetsResponse struct {
	Errors *[]string `json:"error,omitempty"`

	// update
	Updated bool     `json:"updated,omitempty"`
	Split   *[]int64 `json:"split,omitempty"`
}
