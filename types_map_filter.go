package f8n

import "strings"

// Экспорт значения простой фильтрации.
func (filter Filter) exportAsString() (ret string) {
	var b strings.Builder

	if filter.Field == "" {
		return
	}
	_, _ = b.WriteString(filter.Field)
	_, _ = b.WriteString(keyDelimiter)
	_, _ = b.WriteString(filter.Method.String())
	_, _ = b.WriteString(keyDelimiter)
	_, _ = b.WriteString(filter.Value.String())
	ret = b.String()

	return
}

// Выполнение замены символов в значении фильтрации для запросов SQL LIKE.
func (filter Filter) queryGormValueLike() interface{} {
	var ret = strings.Replace(filter.Value.String(), "?", "_", -1)
	ret = strings.Replace(ret, "*", "%", -1)
	return ret
}

// Формирование SQL запроса и параметров для ORM GORM.
func (filter Filter) queryGorm() (query string, values []interface{}) {
	var q strings.Builder

	switch filter.Method {
	case filterEquivalent:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				_, _ = q.WriteString(" IS NULL")
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" = ?")
			values = append(values, filter.Value.Value())
		}
	case filterLessThan:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				// Фильтрация '< ?', для значения NULL не имеет смысла.
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" < ?")
			values = append(values, filter.Value.Value())
		}
	case filterLessThanOrEquivalent:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				// Фильтрация '<= ?', для значения NULL не имеет смысла.
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" <= ?")
			values = append(values, filter.Value.Value())
		}
	case filterGreaterThan:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				// Фильтрация '> ?', для значения NULL не имеет смысла.
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" > ?")
			values = append(values, filter.Value.Value())
		}
	case filterGreaterThanOrEquivalent:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				// Фильтрация '>= ?', для значения NULL не имеет смысла.
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" >= ?")
			values = append(values, filter.Value.Value())
		}
	case filterNotEquivalent:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				_, _ = q.WriteString(" IS NOT NULL")
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" != ?")
			values = append(values, filter.Value.Value())
		}
	case filterLikeThan:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				// Фильтрация 'LIKE ?', для значения NULL не имеет смысла.
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" LIKE ?")
			values = append(values, filter.queryGormValueLike())
		}
	case filterNotLikeThan:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				// Фильтрация 'NOT LIKE ?', для значения NULL не имеет смысла.
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(" NOT LIKE ?")
			values = append(values, filter.queryGormValueLike())
		}
	//case filterIn:
	//	_, _ = q.WriteString()
	//case filterNotIn:
	//	_, _ = q.WriteString()
	default:
		_, _ = q.WriteString("Метод sql запроса `")
		_, _ = q.WriteString(filter.Method.String())
		_, _ = q.WriteString("` не реализован.")
	}
	query = q.String()

	return
}
