package models

import "time"

type Hash string

func (h Hash) String() string {
	return string(h)
}

func (h Hash) FromString(s string) Hash {
	return Hash(s)
}

// TransactionCache is a cache of transactions.
type TransactionCache struct {
	IsInitialised bool

	StartDate *time.Time
	EndDate   *time.Time

	Cache map[Hash]Transaction
}
