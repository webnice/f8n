// Package f8n
package f8n

import (
	"gorm.io/gorm"
)

// Gorm Добавление в ORM всех разобранных значений.
func (f8n *impl) Gorm(orm *gorm.DB) (ret *gorm.DB, err error) {

	if orm == nil {
		err = f8n.Errors().OrmIsNil()
		return
	}
	defer func() { ret = orm }()
	// Добавление фильтров простой фильтрации.
	orm = f8n.gormFilter(orm)
	// Добавление фильтров сложной фильтрации.

	// TODO: Добавить фильтры сложной фильтрации.

	//

	return
}

// Добавление фильтров простой фильтрации.
func (f8n *impl) gormFilter(orm *gorm.DB) (ret *gorm.DB) {
	var (
		n    int
		t, q string
		v    []interface{}
	)

	defer func() { ret = orm }()
	for n = range f8n.Filter {
		t = f8n.fieldName(f8n.Filter[n].Field)
		q, v = f8n.Filter[n].queryGorm()
		orm = orm.Where(t+q, v...)
	}

	return
}
