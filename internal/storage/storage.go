package storage

import "errors"

var (
	ErrDrugNotFound = errors.New("drug not found")
	ErrDrugExist   = errors.New("drug exist")
)
