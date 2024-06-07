package sqlite

import (
	"aptekaPupteka/internal/storage"
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Drugs struct {
	gorm.Model
	Name  string `gorm:"size:255; unique"`
	Count int32  `gorm:"default:0; type:int32"`
}

type Storage struct {
	db *gorm.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %s: %w", op, err)
	}
	err = db.AutoMigrate(&Drugs{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate table: %s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) NewDrug(drugToSave string, count int32) (uint, error) {
	const op = "storage.sqlite.NewDrug"
	drug := Drugs{Name: drugToSave, Count: count}
	err := s.db.Create(&drug).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrDrugExist)
		} else {
			return 0, fmt.Errorf("%s: error save drug %w", op, err)
		}

	}
	id := drug.ID
	return id, nil
}

func (s *Storage) AddDrug(drugName string, count int32) (uint, error) {
	const op = "storage.sqlite.AddDrug"
	var drug Drugs
	err := s.db.Where("name = ?", drugName).First(&drug).Error
	if err != nil {

		return 0, fmt.Errorf("%s: error add drug %w", op, err)
	}

	newCount := drug.Count + count
	if newCount > 65500 {
		err := errors.New("can not add, drugs will empty")
		fmt.Errorf("%s: can not add, drugs will full! %w", op, err)
		return 0, err
	}
	err = s.db.Model(&drug).Where("name = ?", drugName).Update("count", newCount).Error
	if err != nil {

		return 0, fmt.Errorf("%s: error add drug %w", op, err)
	}
	return drug.ID, nil
}
func (s *Storage) SubDrug(drugName string, count int32) (uint, error) {
	const op = "storage.sqlite.SubDrug"
	var drug Drugs
	err := s.db.Where("name = ?", drugName).First(&drug).Error
	if err != nil {

		return 0, fmt.Errorf("%s: error sub drug %w", op, err)
	}

	newCount := drug.Count - count
	if newCount < 0 {
		err := errors.New("can not take, drugs will empty")
		fmt.Errorf("%s: can not take, drugs will empty! %w", op, err)
		return 0, err
	}
	err = s.db.Model(&drug).Where("name = ?", drugName).Update("count", newCount).Error
	if err != nil {
		fmt.Errorf("%s: error sub drug %w", op, err)
		return 0, err
	}
	return drug.ID, nil
}

func (s *Storage) GetAllDrugs() ([]Drugs, error) {
	const op = "storage.sqlite.SubDrug"
	var drugs []Drugs
	err := s.db.Find(&drugs)
	drugList := map[string]int{}
	for _, i := range drugs {
		drugList[i.Name] = int(i.Count)
	}
	if err != nil {

		return drugs, fmt.Errorf("%s: error sub drug %w", op, err)
	}

	// newCount := drug.Count - count
	// if(newCount < 0){
	// 	err := errors.New("can not take, drugs will empty")
	// 	fmt.Errorf("%s: can not take, drugs will empty! %w", op, err)
	// 	return 0, err
	// }

	if err != nil {

		return drugs, fmt.Errorf("%s: error sub drug %w", op, err)
	}
	return drugs, nil
}

func (s *Storage) DeleteDrug(drugToDelete string) (uint, error) {
	const op = "storage.sqlite.NewDrug"
	var drug Drugs
	result := s.db.Where("name = ?", drugToDelete).First(&drug)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Errorf("%s: drug not exist %w", op, result.Error)
			return 0, result.Error
		} else {
			return 0, result.Error
		}
	} else {
		err := s.db.Where("name = ?", drugToDelete).Delete(&drug).Error
		if err == gorm.ErrRecordNotFound {
			fmt.Errorf("%s: drug not exist %w", op, err)
			return 0, err
		}
		if err != nil {
			return 0, fmt.Errorf("%s: error delete drug %w", op, err)
		}
		id := drug.ID
		return id, nil
	}

	return 0, nil
}
func (s *Storage) GetPage(page int, limit int) ([]Drugs, error) {
	const op = "storage.sqlite.GetPage"
	var drugs []Drugs
	err := s.db.Offset(page * limit).Limit(limit).Find(&drugs).Error
	if err != nil {

		return drugs, fmt.Errorf("%s: error get page %w", op, err)
	}

	if err != nil {

		return drugs, fmt.Errorf("%s: error get page %w", op, err)
	}
	return drugs, nil
}
