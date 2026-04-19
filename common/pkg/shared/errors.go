package shared

import "errors"

// ErrNotFound indicates that a requested entity was not found in the repository.
var ErrNotFound = errors.New("entity not found")

// ErrAlreadyExists indicates that an entity with the same unique identifier already exists.
var ErrAlreadyExists = errors.New("entity already exists")

// ErrInvalidInput indicates that provided input is invalid.
var ErrInvalidInput = errors.New("invalid input")

// ErrTransactionConflict indicates that a transaction could not be completed due to a conflict or concurrency issue.
var ErrTransactionConflict = errors.New("transaction conflict")

// ErrOperationFailed is a generic error used for unexpected operation failures.
var ErrOperationFailed = errors.New("operation failed")

// ErrUnauthorized indicates that the caller is not authorized to perform the operation.
var ErrUnauthorized = errors.New("unauthorized")

// ErrCanceled indicates that the operation was canceled or rolled back.
var ErrCanceled = errors.New("operation canceled")

var ErrLimitReached = errors.New("limit reached")

var ErrNoContentFound = errors.New("no content")
