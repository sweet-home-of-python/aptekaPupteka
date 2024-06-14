package sqlite

import (
	"aptekaPupteka/internal/storage"
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// объявляем структуру наркотика представляющую собой схему для ОРМ

type Drugs struct {
	gorm.Model        // наследуем параметры схемы от стандартных gorm
	Name       string `gorm:"size:255; unique"`      // имя наркотика
	Count      int32  `gorm:"default:0; type:int32"` // количество наркотика
}

// создаем стуктуру харилища, будет хранить db
type Storage struct {
	db *gorm.DB
}

// инициализируем базу, функция возвращает указатель на хранилище
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // метка для отладки

	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{TranslateError: true}) // открываем базу если нет то создаем и открываем, конфигом задаем трансляцию ошибок для sqlite

	if err != nil {
		return nil, fmt.Errorf("failed to create database: %s: %w", op, err)
	}

	err = db.AutoMigrate(&Drugs{}) // мигрируем таблицу

	if err != nil {
		return nil, fmt.Errorf("failed to migrate table: %s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// метод хранилища для создания нового наркотика, добавляет в базу наркотик с новым именем drugToSave в количестве count

func (s *Storage) NewDrug(drugName string, count int32) (uint, error) {
	const op = "storage.sqlite.NewDrug"
	var drug Drugs
	
	s.db.Unscoped().Where("name = ?", drugName).Delete(&drug)

	if count > 0 && count < 65500{
		drug.Count = count
	}else {
		drug.Count = 0
	}
	
	drug.Name = drugName
	err := s.db.Create(&drug).Error

	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrDrugExist)
		} else {
			return 0, fmt.Errorf("%s: error save drug %w", op, err)
		}

	}
	id := drug.ID

	return id, nil // возвращаем id записи
}

// метод хранилища для того чтобы увеличить количество конкретного наркотика

func (s *Storage) AddDrug(drugName string, count int32) (uint, error) {
	const op = "storage.sqlite.AddDrug"

	var drug Drugs

	err := s.db.Where("name = ?", drugName).First(&drug).Error

	if err != nil {

		return 0, fmt.Errorf("%s: error add drug %w", op, err)
	}

	var newCount int32 = drug.Count
	if count > 0{
		newCount = drug.Count + count
	}else
	{
		err := errors.New("can not take, count is not correct!")
		fmt.Errorf("%s: error count : %w", op, err)
		return 0, err
	}
	 // новое, пересчитанное количество наркотиков

	if newCount > 65500 && newCount < 0 { // ограничиваем верхний порог
		err := errors.New("can not add, drugs will full!")
		fmt.Errorf("%s: drug add error:  %w", op, err)
		return 0, err
	}

	err = s.db.Model(&drug).Where("name = ?", drugName).Update("count", newCount).Error

	if err != nil {

		return 0, fmt.Errorf("%s: error updating after add: %w", op, err)
	}

	return drug.ID, nil // также возвращаем id как и всегда
}

// метод хранилища для того чтобы забрать некоторое количество конкретного наркотика

func (s *Storage) SubDrug(drugName string, count int32) (uint, error) {
	const op = "storage.sqlite.SubDrug"

	var drug Drugs

	err := s.db.Where("name = ?", drugName).First(&drug).Error

	if err != nil {
		return 0, fmt.Errorf("%s: error find drug for substraction: %w", op, err)
	}

	var newCount int32 = drug.Count
	if count > 0{
		newCount = drug.Count - count
	}else {
		err := errors.New("can not take, count is not correct!")
		fmt.Errorf("%s: error count : %w", op, err)
		return 0, err
	}
	

	if newCount < 0  && newCount < 65500{ // нельзя забрать больше 0
		err := errors.New("can not take, drugs will empty!")
		fmt.Errorf("%s: drug take error: %w", op, err)
		return 0, err
	}

	err = s.db.Model(&drug).Where("name = ?", drugName).Update("count", newCount).Error

	if err != nil {
		fmt.Errorf("%s: error updating after sub: %w", op, err)
		return 0, err
	}

	return drug.ID, nil
}

// метод хранилища для получения всех наркотиков

func (s *Storage) GetAllDrugs() ([]Drugs, error) {
	const op = "storage.sqlite.SubDrug"

	var drugs []Drugs // список всех наркотиков для вывода

	err := s.db.Find(&drugs)
	if err != nil {

		return drugs, fmt.Errorf("%s: error find drugs in db: %w", op, err)
	}

	return drugs, nil
}

// метод хранилища для удаления конкретного наркотика

func (s *Storage) DeleteDrug(drugToDelete string) (uint, error) {
	const op = "storage.sqlite.NewDrug"

	var drug Drugs

	err := s.db.Where("name = ?", drugToDelete).First(&drug).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Errorf("%s: drug not exist: %w", op, err)
			return 0, err
		} else {
			return 0, err
		}
	}

	err = s.db.Where("name = ?", drugToDelete).Delete(&drug).Error
	if err != nil {
		return 0, fmt.Errorf("%s: error delete drug: %w", op, err)
	}

	id := drug.ID
	return id, nil
}

// метод хранилища для получения страницы
func (s *Storage) GetPage(page int, limit int) ([]Drugs, error) {
	const op = "storage.sqlite.GetPage"

	var drugs []Drugs

	err := s.db.Offset(page * limit).Limit(limit).Find(&drugs).Error

	if err != nil {

		return drugs, fmt.Errorf("%s: error get page: %w", op, err)
	}
	return drugs, nil
}

func (s *Storage) SearchDrug(drugName string) ([]Drugs, error) {
	const op = "storage.sqlite.GetPage"

	var drugs []Drugs

	err := s.db.Where("name LIKE ?", drugName+"%").Find(&drugs).Error
	if err != nil{
	}

	return drugs, nil
}