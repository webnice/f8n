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
func (filter Filter) queryGormValueLike() (ret string) {
	const (
		likeOne  = "?"
		likeAll  = "*"
		likeOneR = "_"
		likeAllR = "%"
	)

	ret = strings.Replace(filter.Value.String(), likeOne, likeOneR, -1)
	ret = strings.Replace(ret, likeAll, likeAllR, -1)

	return
}

// Формирование SQL запроса и параметров для ORM GORM.
func (filter Filter) queryGorm() (query string, values []any) {
	const (
		sqlIsNull         = " IS NULL"
		sqlIsNotNull      = " IS NOT NULL"
		sqlEqual          = " = ?"
		sqlNotEqual       = " != ?"
		sqlLessThan       = " < ?"
		sqlLessOrEqual    = " <= ?"
		sqlGreaterThan    = " > ?"
		sqlGreaterOrEq    = " >= ?"
		sqlLike           = " LIKE ?"
		sqlNotLike        = " NOT LIKE ?"
		sqlIn             = " IN (?)"
		sqlNotIn          = " NOT IN (?)"
		tplMethodUnknown1 = "Метод sql запроса `"
		tplMethodUnknown2 = "` не реализован."
	)
	var q strings.Builder

	switch filter.Method {
	case filterEquivalent:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				_, _ = q.WriteString(sqlIsNull)
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(sqlEqual)
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
			_, _ = q.WriteString(sqlLessThan)
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
			_, _ = q.WriteString(sqlLessOrEqual)
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
			_, _ = q.WriteString(sqlGreaterThan)
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
			_, _ = q.WriteString(sqlGreaterOrEq)
			values = append(values, filter.Value.Value())
		}
	case filterNotEquivalent:
		switch filter.Value.Type {
		case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
			if filter.Value.IsValueNil() {
				_, _ = q.WriteString(sqlIsNotNull)
				break
			}
			fallthrough
		default:
			_, _ = q.WriteString(sqlNotEqual)
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
			_, _ = q.WriteString(sqlLike)
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
			_, _ = q.WriteString(sqlNotLike)
			values = append(values, filter.queryGormValueLike())
		}
	case filterIn:
		switch filter.Value.Type {
		case TypeSliceString, TypeSliceInt64, TypeSliceUint64:
			_, _ = q.WriteString(sqlIn)
			values = append(values, filter.Value.Value())
		default:
			break
		}
	case filterNotIn:
		switch filter.Value.Type {
		case TypeSliceString, TypeSliceInt64, TypeSliceUint64:
			_, _ = q.WriteString(sqlNotIn)
			values = append(values, filter.Value.Value())
		default:
			break
		}
	default:
		_, _ = q.WriteString(tplMethodUnknown1)
		_, _ = q.WriteString(filter.Method.String())
		_, _ = q.WriteString(tplMethodUnknown2)
	}
	query = q.String()

	return
}
