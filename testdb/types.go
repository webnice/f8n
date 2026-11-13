package testdb

import (
	"regexp"
	"time"

	"gorm.io/gorm"
)

var (
	// DB Соединение с тестовой базой данных.
	DB *gorm.DB

	// RexWhere Регулярное выражение выделения WHERE условия из SQL запроса.
	RexWhere = regexp.MustCompile(`(?mi)^(.*)\s(where)\s(.*)$`)
)

// BaseModel Базовая модель тестовой таблицы.
type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Parent Родительская тестовая модель.
type Parent struct {
	BaseModel
	FavChildID uint
	FavChild   *Child
	Children   []*Child
}

// Child Дочерняя тестовая модель.
type Child struct {
	BaseModel
	Name     string
	ParentID *uint
	Parent   *Parent
}
