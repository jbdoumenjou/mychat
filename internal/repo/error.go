package repo

import "errors"

// ErrPhoneNumberAlreadyRegistered is returned when the phone number is already registered.
var ErrPhoneNumberAlreadyRegistered = errors.New("phone number already registered")
