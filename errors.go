package stat

import "errors"

var (
	ErrLowerGreaterUpper = errors.New("lower value is greater than upper")
	ErrSpansEmpty        = errors.New("an empty list of spans was specified")
)
