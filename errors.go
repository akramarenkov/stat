package stat

import "errors"

var (
	ErrItemsQuantityNegative = errors.New("items quantity is negative")
	ErrItemsQuantityZero     = errors.New("items quantity is zero")
	ErrLowerGreaterUpper     = errors.New("lower value is greater than upper")
	ErrSpansListEmpty        = errors.New("an empty list of spans was specified")
)
