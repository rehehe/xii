package relational

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"xii/dashboard/pkg/aggregator"
	"xii/dashboard/pkg/reporter"
)

type Storage struct {
	db *gorm.DB
}

// NewStorage returns a new database storage
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

// Set saves the given reporter to the memory
func (s *Storage) FindReports(r reporter.Report, limit int) []reporter.Report {
	var reports []reporter.Report
	s.db.
		Where(&Report{
			Device:   r.Device,
			Location: r.Location,
		}).
		Limit(limit).
		Find(&reports)

	return reports
}

// Set saves the given reporter to the memory
func (s *Storage) CreateReports(rs []aggregator.Report) error {
	transaction := s.db.Begin()
	for _, r := range rs {
		transaction.
			//Set("gorm:insert_option", "ON CONFLICT IGNORE").
			Create(&Report{
				Text:              r.Text,
				Device:            r.Device,
				Location:          r.Location,
				Provider:          r.Provider,
				ProviderID:        r.ID,
				ProviderCreatedAt: r.CreatedAt,
			})
	}
	transaction.Commit()

	if s.db.Error != nil {
		return s.db.Error
	}

	return nil
}

func (s *Storage) LastResultIDbyProvider(provider string) uint {
	r := Report{}
	s.db.Where(&Report{Provider: provider}).Last(&r)
	return r.ProviderID
}

func (s Storage) Close() {
	s.db.Close()
}
