package models

import "github.com/jinzhu/gorm"

type Services struct {
	db *gorm.DB
	Models ModuleService
}

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err !=  nil {
		return nil, err
	}
	return &Services{
		db: db,
		Models: NewModuleService(db),
	}, nil
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&Module{}).Error
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&Module{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
