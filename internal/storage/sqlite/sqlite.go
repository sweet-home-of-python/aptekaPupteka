package sqlite

import (
	"aptekaPupteka/internal/storage"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Drugs struct {
    gorm.Model
    Name   string `gorm:"size:255; unique"`
    Count    uint16    `gorm:"default:0; type:int"`
}

type Storage struct {
	db *gorm.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" 

	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{TranslateError: true})
	if err != nil{
		fmt.Errorf("failed to connect database: %v", op, err)
	}
	err = db.AutoMigrate(&Drugs{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) NewDrug(drugToSave string) (uint, error) {
	const op = "storage.sqlite.NewDrug"
	drug := Drugs{Name: drugToSave, Count: 0}
	err := s.db.Create(&drug).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey  {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrDrugExist)
		}else{
			return 0, fmt.Errorf("%s: error save drug %w", op, err)
		}
		
	}
	id := drug.ID	
	return id, nil	
}

func (s *Storage) AddDrug(drugName string, count uint16) (uint, error) {
	const op = "storage.sqlite.AddDrug"
	var drug Drugs
	err := s.db.Where("name = ?", drugName).First(&drug).Error
	if err != nil {
		
		return 0, fmt.Errorf("%s: error save drug %w", op, err)
	}
	
	newCount := drug.Count + count
	err = s.db.Model(&drug).Where("name = ?", drugName).Update("count", newCount).Error
	if err != nil {
		
			return 0, fmt.Errorf("%s: error save drug %w", op, err)
		}
	return drug.ID, nil	
}

// func (s *Storage) AddDrugCount(drug string, count int)( int64 , error){
// 	const op = "storage.sqlite.addDrugCount"
// 	var resultCount int
// 	_ = resultCount
// 	stmt, err := s.db.Prepare("UPDATE med SET count = count +? WHERE name = ?")
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
// 	}
// 	res, err := stmt.Exec(count, drug)
// 	if err != nil {
// 		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
// 			return  0, fmt.Errorf("%s: %w", op, storage.ErrDrugExist)
// 		}

// 		return  0, fmt.Errorf("%s: execute statement: %w", op, err)
// 	}
// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		return  0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
// 	}

// 	// Возвращаем ID
// 	return id, nil
// }

// func (s *Storage) TakeDrugCount(drug string, count int)( int64 , error){
// 	const op = "storage.sqlite.TakeDrugCount"
// 	stmt, err := s.db.Prepare("UPDATE med SET count = count - ? WHERE name = ?")
// 	if err != nil{
// 		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
// 	}

// 	res, err:= stmt.Exec(count, drug)
// 	if err != nil{
// 		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
// 	}
// 	id, err := res.LastInsertId()
// 	if err != nil{
// 		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
// 	}
// 	return id, nil

// }

// func (s *Storage) DeleteDrug(drug string)(int64, error){
// 	const op = "storage.sqlite.DeleteDrug"
	
// 	stmt, err := s.db.Prepare("DELETE FROM med WHERE name = ?")
// 	if err != nil{
// 		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
// 	}

// 	res, err := stmt.Exec(drug)
// 	if err != nil{
// 		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
// 	}

// 	id, err := res.LastInsertId()
// 	if err != nil{
// 		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
// 	}
// 	return id, nil
// }

// func (s *Storage) GetDrugCount() (int, error) {
// 	const op = "storage.sqlite.GetDrugCount"


// 	return 0, nil
// }