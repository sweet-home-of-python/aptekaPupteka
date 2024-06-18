package storage

import "errors"

// ! Полагаю встроенных нет?
var ( // ошибки
	ErrDrugNotFound = errors.New("drug not found")
	ErrDrugExist    = errors.New("drug exist")
)
