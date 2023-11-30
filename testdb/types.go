package testdb

import (
	"regexp"
	"time"

	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	RexWhere = regexp.MustCompile(`(?mi)^(.*)\s(where)\s(.*)$`)
)

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Parent struct {
	BaseModel
	FavChildID uint
	FavChild   *Child
	Children   []*Child
}

type Child struct {
	BaseModel
	Name     string
	ParentID *uint
	Parent   *Parent
}
