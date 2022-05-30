package weapp

import (
	"fmt"
	"time"

	"github.com/gookit/color"
)

type Migrations struct {
	migrations []IMigration
}

func (m *Migrations) Init() {
	m.migrations = make([]IMigration, 0)
}

func (m *Migrations) Migration(migrations ...IMigration) {
	m.migrations = append(m.migrations, migrations...)
}

func (m *Migrations) Commit(names ...string) {
	for _, migration := range m.migrations {
		if len(names) > 0 {
			exists := false
			for _, name := range names {
				if name == migration.Name() {
					exists = true
					break
				}
			}
			if !exists {
				continue
			}
		}
		db := migration.Database()
		if db.driver == "mysql" {
			db.Set("gorm:table_options", "ENGINE=InnoDB")
		}
		migrator := db.Migrator()
		db.Begin()
		if !migrator.HasTable(new(MigrationRecord)) {
			if err := migrator.CreateTable(new(MigrationRecord)); err != nil {
				db.Rollback()
				panic(err)
			}
		}
		var count int64 = 0
		if err := db.Model(&MigrationRecord{}).Where(&MigrationRecord{
			Migration: migration.Name(),
			Version:   migration.Version(),
		}).Count(&count).Error; err != nil {
			db.Rollback()
			panic(err)
		}
		if count > 0 {
			continue
		}
		startTime := time.Now()
		fmt.Printf("%s: %s\r\n", color.Warn.Render("migrating"), migration.Name())
		if err := migration.Commit(); err != nil {
			db.Rollback()
			panic(err)
		}
		db.Create(&MigrationRecord{
			Version:   migration.Version(),
			Migration: migration.Name(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		db.Commit()
		fmt.Printf("%s: %s(%2fs)\r\n", color.Success.Render("migrated"), migration.Name(), time.Now().Sub(startTime).Seconds())
	}
}

func (m *Migrations) Rollback(names ...string) {
	for _, migration := range m.migrations {
		if len(names) > 0 {
			exists := false
			for _, name := range names {
				if name == migration.Name() {
					exists = true
					break
				}
			}
			if !exists {
				continue
			}
		}
		db := migration.Database()
		if db.driver == "mysql" {
			db.Set("gorm:table_options", "ENGINE=InnoDB")
		}
		migrator := db.Migrator()
		db.Begin()
		if !migrator.HasTable(new(MigrationRecord)) {
			migrator.CreateTable(new(MigrationRecord))
		}
		var count int64 = 0
		if err := db.Model(&MigrationRecord{}).Where(&MigrationRecord{
			Migration: migration.Name(),
			Version:   migration.Version(),
		}).Count(&count).Error; err != nil {
			db.Rollback()
			panic(err)
		}
		if count == 0 {
			continue
		}
		startTime := time.Now()
		fmt.Printf("%s: %s\r\n", color.Warn.Render("migrating"), migration.Name())
		if err := migration.Rollback(); err != nil {
			db.Rollback()
			panic(err)
		}
		if err := db.Where("version = ?", migration.Version()).Where("migration = ?", migration.Name()).Delete(&MigrationRecord{}).Error; err != nil {
			db.Rollback()
			panic(err)
		}
		db.Commit()
		fmt.Printf("%s: %s(%2fs)\r\n", color.Success.Render("migrated"), migration.Name(), time.Now().Sub(startTime).Seconds())
	}
}

type MigrationRecord struct {
	ID        uint64 `gorm:"primaryKey"`
	Version   string
	Migration string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (MigrationRecord) TableName() string {
	return "migrations"
}

type IMigration interface {
	Commit() error
	Rollback() error
	Name() string
	Version() string
	Database() *Database
}

type Migration struct {
	*Component `inject:"component"`
}

func (m *Migration) Database() *Database {
	return m.DB("default")
}
