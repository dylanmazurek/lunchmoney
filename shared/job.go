package shared

import "time"

type Job struct {
	Action   string     `json:"action"`
	AssetID  *int       `json:"assetId,omitempty"`
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
}
