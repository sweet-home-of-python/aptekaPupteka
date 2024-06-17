package storage

import "errors"

var ( // ошибки
	ErrDrugNotFound = errors.New("drug not found")
	ErrDrugExist   = errors.New("drug exist")
)
