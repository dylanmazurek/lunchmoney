package lunchmoney

import (
	"fmt"
	"strconv"

	"github.com/Rhymond/go-money"
)

// ParseCurrency turns two strings into a money struct.
func ParseCurrency(amount, currency string) (*money.Money, error) {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return nil, fmt.Errorf("%q is not valid float: %w", amount, err)
	}

	v := int64(100 * f)
	return money.New(v, currency), nil
}
