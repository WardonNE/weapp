package weapp

import (
	"time"
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

func (m *Migrations) Commit() {
	for _, migration := range m.migrations {
		db := migration.Database()
		migrator := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator()
		db.Begin()
		if !migrator.HasTable(new(MigrationRecord)) {
			if err := migrator.CreateTable(new(MigrationRecord)); err != nil {
				db.Rollback()
				panic(err)
			}
		}
		var count int64 = 0
		if err := db.Where(&MigrationRecord{
			Migration: migration.Name(),
			Version:   migration.Version(),
		}).Count(&count).Error; err != nil {
			db.Rollback()
			panic(err)
		}
		if count > 0 {
			continue
		}
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
	}
}

func (m *Migrations) Rollback() {
	for _, migration := range m.migrations {
		db := migration.Database()
		migrator := db.Migrator()
		db.Begin()
		if !migrator.HasTable(new(MigrationRecord)) {
			migrator.CreateTable(new(MigrationRecord))
		}
		var count int64 = 0
		if err := db.Where(&MigrationRecord{
			Migration: migration.Name(),
			Version:   migration.Version(),
		}).Count(&count).Error; err != nil {
			db.Rollback()
			panic(err)
		}
		if count == 0 {
			continue
		}
		if err := migration.Rollback(); err != nil {
			db.Rollback()
			panic(err)
		}
		if err := db.Where("version = ?", migration.Version()).Where("migration = ?", migration.Name()).Delete(&MigrationRecord{}).Error; err != nil {
			db.Rollback()
			panic(err)
		}
		db.Commit()
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
	SetApplication(app *Application)
	Commit() error
	Rollback() error
	Name() string
	Version() string
	Database() *Database
}

type Migration struct {
	app *Application
}

func (m *Migration) SetApplication(app *Application) {
	m.app = app
}
