package sqlite

import (
	"aptekaPupteka/internal/storage"
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	db, err := sql.Open("sqlite3", storagePath) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем таблицу, если ее еще нет
	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS med(
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        count INTEGER);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveDrug(drugToSave string) (int64, error) {
	const op = "storage.sqlite.SaveDrug"

	// Подготавливаем запрос
	stmt, err := s.db.Prepare("INSERT INTO med(name) values(?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	// Выполняем запрос
	res, err := stmt.Exec(drugToSave)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrDrugExist)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	// Возвращаем ID
	return id, nil
}
func (s *Storage) AddDrugCount(drug string, count int)( int64 , error){
	const op = "storage.sqlite.addDrugCount"
	var resultCount int
	_ = resultCount
	stmt, err := s.db.Prepare("UPDATE med SET count = count +? WHERE name = ?")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	res, err := stmt.Exec(count, drug)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return  0, fmt.Errorf("%s: %w", op, storage.ErrDrugExist)
		}

		return  0, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return  0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	// Возвращаем ID
	return id, nil
}
// func (s *Storage) GetDrugCount() (int, error) {
// 	const op = "storage.sqlite.GetDrugCount"


// 	return 0, nil
// }