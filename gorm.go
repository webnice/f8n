// Package f8n
package f8n

import "gorm.io/gorm"

// Gorm Добавление в ORM всех разобранных значений.
func (f8n *impl) Gorm(orm *gorm.DB) (ret *gorm.DB, err error) {
	var n int

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

	// Добавление сортировки.
	for n = range f8n.By {
		switch f8n.By[n].By {
		case byAsc:
			orm = orm.Order(f8n.fieldName(f8n.By[n].Field) + " ASC")
		case byDesc:
			orm = orm.Order(f8n.fieldName(f8n.By[n].Field) + " DESC")
		}
	}
	// Добавление лимитов.
	if f8n.Offset > 0 {
		orm = orm.Offset(int(f8n.Offset))
	}
	if f8n.Limit > 0 {
		orm = orm.Limit(int(f8n.Limit))
	}

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
