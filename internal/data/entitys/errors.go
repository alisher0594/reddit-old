package entitys

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	// ErrEditConflict ...
	ErrEditConflict = errors.New("edit conflict")
)
