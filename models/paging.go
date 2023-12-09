package models

const (
	DefaultFetchLimit    = 40
	DefaultFetchPageSize = 10
)

type Paging struct {
	Page     int
	PageSize int
}

// func (t *TransactionFilter) Set(assetID int, pageSize *int, limit *int) {
// 	t.assetID = assetID

// 	t.pageSize = DefaultFetchPageSize
// 	if pageSize != nil {
// 		t.pageSize = *pageSize
// 	}

// 	t.limit = DefaultFetchLimit
// 	if limit != nil {
// 		t.limit = *limit
// 	}

// 	t.page = 1
// 	t.offset = 0
// }

// func (t *TransactionFilter) NextPage(currCount int) bool {
// 	t.page++
// 	t.offset = t.page * t.pageSize

// 	reachedLimit := t.ReachedLimit(currCount)
// 	morePages := currCount >= t.pageSize

// 	return !reachedLimit && morePages
// }

// func (t *TransactionFilter) ReachedLimit(idx int) bool {
// 	return idx >= t.limit
// }

// func (t *TransactionFilter) CurrentPage() int {
// 	return t.page
// }

// func (t TransactionFilter) ToVals() (url.Values, error) {
// 	vals, err := query.Values(t)
// 	if err != nil {
// 		return nil, err
// 	}

// 	vals.Set("asset_id", fmt.Sprintf("%d", t.assetID))

// 	if t.limit > 0 {
// 		vals.Set("limit", fmt.Sprintf("%d", t.limit))
// 	}

// 	if t.offset > 0 {
// 		vals.Set("offset", fmt.Sprintf("%d", t.offset))
// 	}

// 	return vals, nil
// }
