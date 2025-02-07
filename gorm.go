package f8n

import "gorm.io/gorm"

// Gorm Добавление в ORM всех разобранных значений.
func (f8n *impl) Gorm(orm *gorm.DB) (ret *gorm.DB, err error) {
	var filter *gorm.DB

	if ret, filter, err = f8n.GormSepFilter(orm); err != nil {
		return
	}
	// Объединение с фильтрами.
	if len(f8n.Filter) > 0 {
		ret = orm.Where(filter)
	}

	return
}

// GormSepFilter Добавление в ORM всех разобранных значений и возвращение фильтрации отдельным объектом.
func (f8n *impl) GormSepFilter(orm *gorm.DB) (ret *gorm.DB, filter *gorm.DB, err error) {
	var n int

	if orm == nil {
		err = f8n.Errors().OrmIsNil()
		return
	}
	defer func() { ret = orm }()
	// Создание фильтров простой фильтрации.
	filter = f8n.gormFilter(orm)
	// Создание фильтров сложной фильтрации.

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
func (f8n *impl) gormFilter(orm *gorm.DB) (filter *gorm.DB) {
	const keyInstance = "original-orm-object"
	var (
		n    int
		t, q string
		v    []any
		or   bool
	)

	// Создание чистой сессии для отделения всех фильтров от основного запроса.
	filter = orm.Session(&gorm.Session{NewDB: true})
	// Устаревший режим фильтрации применяется только если нет MAP.
	or = f8n.Map == nil && f8n.Tie == tieOr
	// Добавление фильтров.
	for n = range f8n.Filter {
		t = f8n.fieldName(f8n.Filter[n].Field)
		q, v = f8n.Filter[n].queryGorm()
		switch {
		case or && n > 0:
			filter = filter.Or(t+q, v...)
		default:
			filter = filter.Where(t+q, v...)
		}
	}

	return
}
