package stat

import "errors"

var (
	ErrInvalidIndex      = errors.New("invalid index of item is specified")
	ErrLowerGreaterUpper = errors.New("lower value is greater than upper")
	ErrNegInfNotPresent  = errors.New("negative infinity is not present")
	ErrPosInfNotPresent  = errors.New("positive infinity is not present")
)
