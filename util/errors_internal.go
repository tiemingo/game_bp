package util

import "errors"

var (
	IErrUnsetPhases            = errors.New("phases not set")
	IErrInvalidInitialDuration = errors.New("initial duration must be greater than 0")
	IErrMissingInitialPhase    = errors.New("initial phase not found in phases map")
	IErrEmptyPhase             = errors.New("phase cannot be empty string")
)
