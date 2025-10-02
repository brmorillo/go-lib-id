package idgen

import "errors"

// Common errors returned by the idgen package
var (
	// ErrNotImplemented is returned when a feature is not yet implemented
	ErrNotImplemented = errors.New("feature not implemented")

	// ErrCUIDNotImplemented is returned when CUID generation is called
	ErrCUIDNotImplemented = errors.New("CUID generation not implemented yet")

	// ErrULIDNotImplemented is returned when ULID generation is called
	ErrULIDNotImplemented = errors.New("ULID generation not implemented yet")

	// ErrNanoIDNotImplemented is returned when NanoID generation is called
	ErrNanoIDNotImplemented = errors.New("NanoID generation not implemented yet")

	// ErrShortIDNotImplemented is returned when ShortID generation is called
	ErrShortIDNotImplemented = errors.New("ShortID generation not implemented yet")

	// ErrKSUIDNotImplemented is returned when KSUID generation is called
	ErrKSUIDNotImplemented = errors.New("KSUID generation not implemented yet")

	// ErrXIDNotImplemented is returned when xid generation is called
	ErrXIDNotImplemented = errors.New("xid generation not implemented yet")

	// ErrSonyflakeNotImplemented is returned when Sonyflake generation is called
	ErrSonyflakeNotImplemented = errors.New("sonyflake generation not implemented yet")
)
