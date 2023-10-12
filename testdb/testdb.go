// Package testdb
package testdb

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	sqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	var (
		err   error
		sqlDB *sql.DB
	)

	if DB, err = OpenTestConnection(&gorm.Config{}); err != nil {
		log.Printf("failed to connect database, got error %v", err)
		os.Exit(1)
	} else {
		if sqlDB, err = DB.DB(); err != nil {
			log.Printf("failed to connect database, got error %v", err)
			os.Exit(1)
		}
		err = sqlDB.Ping()
		if err != nil {
			log.Printf("failed to ping sqlDB, got error %v", err)
			os.Exit(1)
		}
		Migrations()
		if DB.Dialector.Name() == "sqlite" {
			DB.Exec("PRAGMA foreign_keys = ON")
		}
	}
}

// OpenTestConnection Открытие тестового соединения с тестовой базой данных.
func OpenTestConnection(cfg *gorm.Config) (db *gorm.DB, err error) {
	log.Println("testing sqlite3")
	if db, err = gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), cfg); err != nil {
		return
	}

	return
}

func Migrations() {
	var err error

	allModels := []interface{}{&Parent{}, &Child{}}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allModels), func(i, j int) { allModels[i], allModels[j] = allModels[j], allModels[i] })
	if err = DB.Migrator().DropTable("user_friends", "user_speaks"); err != nil {
		log.Printf("Failed to drop table, got error %v\n", err)
		os.Exit(1)
	}
	if err = DB.Migrator().DropTable(allModels...); err != nil {
		log.Printf("Failed to drop table, got error %v\n", err)
		os.Exit(1)
	}
	if err = DB.AutoMigrate(allModels...); err != nil {
		log.Printf("Failed to auto migrate, but got error %v\n", err)
		os.Exit(1)
	}
	for _, m := range allModels {
		if !DB.Migrator().HasTable(m) {
			log.Printf("Failed to create table for %#v\n", m)
			os.Exit(1)
		}
	}
}
