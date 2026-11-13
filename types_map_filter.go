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

	// Признак значения NULL для типов полей допускающих NULL.
	valueNil := isFieldTypeNil(filter.Value.Type) && filter.Value.IsValueNil()

	switch filter.Method {
	case filterEquivalent:
		switch {
		case valueNil:
			_, _ = q.WriteString(sqlIsNull)
		default:
			_, _ = q.WriteString(sqlEqual)
			values = append(values, filter.Value.Value())
		}
	case filterLessThan:
		// Фильтрация '< ?', для значения NULL не имеет смысла.
		if !valueNil {
			_, _ = q.WriteString(sqlLessThan)
			values = append(values, filter.Value.Value())
		}
	case filterLessThanOrEquivalent:
		// Фильтрация '<= ?', для значения NULL не имеет смысла.
		if !valueNil {
			_, _ = q.WriteString(sqlLessOrEqual)
			values = append(values, filter.Value.Value())
		}
	case filterGreaterThan:
		// Фильтрация '> ?', для значения NULL не имеет смысла.
		if !valueNil {
			_, _ = q.WriteString(sqlGreaterThan)
			values = append(values, filter.Value.Value())
		}
	case filterGreaterThanOrEquivalent:
		// Фильтрация '>= ?', для значения NULL не имеет смысла.
		if !valueNil {
			_, _ = q.WriteString(sqlGreaterOrEq)
			values = append(values, filter.Value.Value())
		}
	case filterNotEquivalent:
		switch {
		case valueNil:
			_, _ = q.WriteString(sqlIsNotNull)
		default:
			_, _ = q.WriteString(sqlNotEqual)
			values = append(values, filter.Value.Value())
		}
	case filterLikeThan:
		// Фильтрация 'LIKE ?', для значения NULL не имеет смысла.
		if !valueNil {
			_, _ = q.WriteString(sqlLike)
			values = append(values, filter.queryGormValueLike())
		}
	case filterNotLikeThan:
		// Фильтрация 'NOT LIKE ?', для значения NULL не имеет смысла.
		if !valueNil {
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
