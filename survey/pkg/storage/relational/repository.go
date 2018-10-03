package relational

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"xii/survey/pkg/report"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(uri string, tablePrefix string) *Storage {
	uri += "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open("mysql", uri)
	if err != nil {
		log.Printf("storage connection failed to database. Err: %v\nretry loop just started...", err)

		for i := 1; i <= 16; i++ {
			db, err = gorm.Open("mysql", uri)
			if err == nil {
				break
			}

			time.Sleep(time.Second * time.Duration(i))
		}

		if err != nil {
			log.Fatal("storage connection failed to database after many attempts")
		}
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + "_" + defaultTableName
	}

	db.AutoMigrate(&Report{})

	return &Storage{db}
}

func (s *Storage) FindReports(r report.Report, limit int) []report.Report {
	var reports []report.Report
	s.db.
		Where(&Report{
			Device:   r.Device,
			Location: r.Location,
		}).
		Limit(limit).
		Find(&reports)

	return reports
}

func (s *Storage) CreateReport(r report.Report) error {
	s.db.
		Create(&Report{
			Text:     r.Text,
			Device:   r.Device,
			Location: r.Location,
		})

	if s.db.Error != nil {
		return s.db.Error
	}

	return nil
}

func (s Storage) Close() {
	s.db.Close()
}
