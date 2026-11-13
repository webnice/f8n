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

const (
	tplConnectError = "Ошибка подключения к тестовой базе данных: %v."
	tplPingError    = "Ошибка проверки соединения с тестовой базой данных: %v."
	tplDropError    = "Ошибка удаления таблицы, ошибка: %v."
	tplMigrateError = "Ошибка автоматической миграции, ошибка: %v."
	tplTableError   = "Ошибка создания таблицы для %#v."
	tplTestDatabase = "тестирование sqlite3"
	pragForeignKeys = "PRAGMA foreign_keys = ON"
)

func init() {
	var (
		err   error
		sqlDB *sql.DB
	)

	if DB, err = OpenTestConnection(&gorm.Config{}); err != nil {
		log.Printf(tplConnectError, err)
		os.Exit(1)
	}
	if sqlDB, err = DB.DB(); err != nil {
		log.Printf(tplConnectError, err)
		os.Exit(1)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Printf(tplPingError, err)
		os.Exit(1)
	}
	Migrations()
	if DB.Dialector.Name() == "sqlite" {
		if err = DB.Exec(pragForeignKeys).Error; err != nil {
			log.Printf(tplConnectError, err)
			os.Exit(1)
		}
	}
}

// OpenTestConnection Открытие тестового соединения с тестовой базой данных.
func OpenTestConnection(cfg *gorm.Config) (db *gorm.DB, err error) {
	log.Println(tplTestDatabase)
	if db, err = gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), cfg); err != nil {
		return
	}

	return
}

// Migrations Пересоздание таблиц тестовой базы данных.
func Migrations() {
	var err error

	// Случайный порядок миграций моделей, для выявления скрытых зависимостей между таблицами.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	allModels := []any{&Parent{}, &Child{}}
	rnd.Shuffle(len(allModels), func(i int, j int) { allModels[i], allModels[j] = allModels[j], allModels[i] })
	if err = DB.Migrator().DropTable("user_friends", "user_speaks"); err != nil {
		log.Printf(tplDropError, err)
		os.Exit(1)
	}
	if err = DB.Migrator().DropTable(allModels...); err != nil {
		log.Printf(tplDropError, err)
		os.Exit(1)
	}
	if err = DB.AutoMigrate(allModels...); err != nil {
		log.Printf(tplMigrateError, err)
		os.Exit(1)
	}
	for _, m := range allModels {
		if !DB.Migrator().HasTable(m) {
			log.Printf(tplTableError, m)
			os.Exit(1)
		}
	}
}
