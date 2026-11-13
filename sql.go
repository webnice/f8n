package f8n

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Sql Создание SQL запроса с фильтрацией данных.
// В функцию передаётся диалект базы данных, для правильного формирования SQL запроса с ориентировкой на
// тип используемой базы данных.
// Возвращается условие WHERE со всеми разобранными значениями, готовое к подстановке в SQL запрос.
// Наименования полей заключаются в кавычки соответствующие диалекту базы данных, значения подставляются
// экранированными литералами. Если условия фильтрации отсутствуют, возвращается пустая строка.
// Не поддерживаемый диалект - ошибка SqlNotImplemented.
func (f8n *impl) Sql(dialect SqlDialect) (ret string, err error) {
	const (
		keyWhere   = "WHERE "
		sqlAnd     = " AND "
		sqlOr      = " OR "
		quoteMySQL = "`"
		quotePgSQL = `"`
	)
	var (
		quote     string
		q         string
		v         []any
		values    []any
		or        bool
		n         int
		b         strings.Builder
		condition string
	)

	// Кавычки идентификаторов соответствующие диалекту базы данных.
	switch dialect {
	case SqlDialectMySQL, SqlDialectMariaDB, SqlDialectSqlite3:
		quote = quoteMySQL
	case SqlDialectPostgreSQL:
		quote = quotePgSQL
	default:
		err = f8n.Errors().SqlNotImplemented.Bind(dialect)
		return
	}
	or = f8n.Map == nil && f8n.Tie == tieOr
	// Условия простой фильтрации.
	for n = range f8n.Filter {
		if n > 0 {
			if or {
				_, _ = b.WriteString(sqlOr)
			} else {
				_, _ = b.WriteString(sqlAnd)
			}
		}
		q, v = f8n.Filter[n].queryGorm()
		_, _ = b.WriteString(f8n.fieldName(f8n.Filter[n].Field))
		_, _ = b.WriteString(q)
		values = append(values, v...)
	}
	// Условия сложной фильтрации, группируются в скобки, для корректного объединения с условиями простой фильтрации.
	if f8n.Map != nil {
		if q, v = f8n.Map.queryGorm(); q != "" {
			if b.Len() > 0 {
				_, _ = b.WriteString(sqlAnd)
			}
			_, _ = b.WriteString(pairBeg)
			_, _ = b.WriteString(q)
			_, _ = b.WriteString(pairEnd)
			values = append(values, v...)
		}
	}
	if b.Len() == 0 {
		return
	}
	// Перевод кавычек идентификаторов на кавычки диалекта и подстановка литералов значений.
	condition = strings.ReplaceAll(b.String(), quoteMySQL, quote)
	ret = keyWhere + sqlSubstitute(condition, values, dialect)

	return
}

// Подстановка значений вместо плейсхолдеров SQL запроса, значения форматируются литералами диалекта.
func sqlSubstitute(query string, values []any, dialect SqlDialect) (ret string) {
	var (
		b strings.Builder
		n int
		j int
	)

	for n = 0; n < len(query); n++ {
		if query[n] == '?' && j < len(values) {
			_, _ = b.WriteString(sqlLiteral(dialect, values[j]))
			j++
			continue
		}
		_ = b.WriteByte(query[n])
	}
	ret = b.String()

	return
}

// Форматирование значения литералом диалекта базы данных.
func sqlLiteral(dialect SqlDialect, value any) (ret string) {
	const (
		sqlNull     = "NULL"
		sqlTrue     = "TRUE"
		sqlFalse    = "FALSE"
		sqlOne      = "1"
		sqlZero     = "0"
		tplString   = "'%s'"
		tplByteaPg  = "'\\x%s'::bytea"
		tplHexMySQL = "0x%s"
		tplHexLite  = "X'%s'"
		sepItem     = ", "
	)
	var (
		items []string
		n     int
	)

	switch v := value.(type) {
	case nil:
		ret = sqlNull
	case string:
		ret = fmt.Sprintf(tplString, sqlEscape(v))
	case []byte:
		// Двоичные данные форматируются в соответствии с диалектом базы данных.
		switch dialect {
		case SqlDialectPostgreSQL:
			ret = fmt.Sprintf(tplByteaPg, hexEncode(v))
		case SqlDialectMySQL, SqlDialectMariaDB:
			ret = fmt.Sprintf(tplHexMySQL, hexEncode(v))
		default:
			ret = fmt.Sprintf(tplHexLite, hexEncode(v))
		}
	case json.RawMessage:
		ret = fmt.Sprintf(tplString, sqlEscape(string(v)))
	case net.IP:
		ret = fmt.Sprintf(tplString, sqlEscape(v.String()))
	case time.Duration:
		// Длительность хранится в базе данных числом наносекунд.
		ret = strconv.FormatInt(int64(v), 10)
	case []string:
		items = make([]string, 0, len(v))
		for n = range v {
			items = append(items, fmt.Sprintf(tplString, sqlEscape(v[n])))
		}
		ret = strings.Join(items, sepItem)
	case bool:
		// PostgreSQL поддерживает литералы TRUE и FALSE, остальные диалекты используют числа 1 и 0.
		switch {
		case v && dialect == SqlDialectPostgreSQL:
			ret = sqlTrue
		case v:
			ret = sqlOne
		case dialect == SqlDialectPostgreSQL:
			ret = sqlFalse
		default:
			ret = sqlZero
		}
	case int64:
		ret = strconv.FormatInt(v, 10)
	case uint64:
		ret = strconv.FormatUint(v, 10)
	case float64:
		ret = strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		ret = fmt.Sprintf(tplString, sqlEscape(v.Format(time.RFC3339Nano)))
	case []int64:
		items = make([]string, 0, len(v))
		for n = range v {
			items = append(items, strconv.FormatInt(v[n], 10))
		}
		ret = strings.Join(items, sepItem)
	case []uint64:
		items = make([]string, 0, len(v))
		for n = range v {
			items = append(items, strconv.FormatUint(v[n], 10))
		}
		ret = strings.Join(items, sepItem)
	default:
		ret = fmt.Sprintf(tplString, sqlEscape(fmt.Sprint(v)))
	}

	return
}

// Кодирование среза байт в шестнадцатеричную строку.
func hexEncode(b []byte) (ret string) {
	const hexDigits = "0123456789abcdef"
	var n int

	var sb strings.Builder
	sb.Grow(len(b) * 2)
	for n = range b {
		_ = sb.WriteByte(hexDigits[b[n]>>4])
		_ = sb.WriteByte(hexDigits[b[n]&0x0f])
	}
	ret = sb.String()

	return
}

// Экранирование строкового литерала SQL запроса.
func sqlEscape(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
