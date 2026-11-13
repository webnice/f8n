package f8n

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"
)

// Тестирование создания SQL условия WHERE для всех поддерживаемых диалектов.
//
//goland:noinspection ALL
func TestSql(t *testing.T) {
	const (
		datUrn = `http://localhost?filter=id:eq:1&filter=key:ge:-3&filter=ratio:lt:2.5` +
			`&filter=flag:eq:true&filter=name:eq:o'clock` +
			`&filter=created:eq:2018-08-03T03:00:00%2B03:00&filter=tags:in:a,b&filter=ids:in:7,8` +
			`&filter=dur:ge:1h30m&filter=blob:eq:aGVsbG8=&filter=ip:eq:192.168.0.1&filter=js:eq:%7B%22a%22%3A1%7D` +
			`&map=m1:and:m2&m1=u:eq:5&m2=deletedAt:eq:nil`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseRequest() != nil, ошибка не ожидалась. Ошибка: %s"
		tplErrorDialect = "%s: Sql() = %q, ожидалось: %q."
		wantMySQL       = "WHERE `id` = 1 AND `key` >= -3 AND `ratio` < 2.5 AND `flag` = 1" +
			" AND `name` = 'o''clock' AND `created` = '2018-08-03T03:00:00+03:00'" +
			" AND `tags` IN ('a', 'b') AND `ids` IN (7, 8)" +
			" AND `dur` >= 5400000000000 AND `blob` = 0x68656c6c6f" +
			" AND `ip` = '192.168.0.1' AND `js` = '{\"a\":1}'" +
			" AND (`u` = 5 AND `deletedAt` IS NULL)"
		wantPostgreSQL = "WHERE \"id\" = 1 AND \"key\" >= -3 AND \"ratio\" < 2.5 AND \"flag\" = TRUE" +
			" AND \"name\" = 'o''clock' AND \"created\" = '2018-08-03T03:00:00+03:00'" +
			" AND \"tags\" IN ('a', 'b') AND \"ids\" IN (7, 8)" +
			" AND \"dur\" >= 5400000000000 AND \"blob\" = '\\x68656c6c6f'::bytea" +
			" AND \"ip\" = '192.168.0.1' AND \"js\" = '{\"a\":1}'" +
			" AND (\"u\" = 5 AND \"deletedAt\" IS NULL)"
		wantSqlite = "WHERE `id` = 1 AND `key` >= -3 AND `ratio` < 2.5 AND `flag` = 1" +
			" AND `name` = 'o''clock' AND `created` = '2018-08-03T03:00:00+03:00'" +
			" AND `tags` IN ('a', 'b') AND `ids` IN (7, 8)" +
			" AND `dur` >= 5400000000000 AND `blob` = X'68656c6c6f'" +
			" AND `ip` = '192.168.0.1' AND `js` = '{\"a\":1}'" +
			" AND (`u` = 5 AND `deletedAt` IS NULL)"
		wantRemapMySQL   = "WHERE `cursor`.`id` = 1"
		wantRemapPostgre = `WHERE "cursor"."id" = 1`
	)
	var (
		err      error
		rq       *http.Request
		f8n      Interface
		ret      string
		n        int
		dialects = []struct {
			Name string
			Key  SqlDialect
			Want string
		}{
			{Name: "mysql", Key: SqlDialectMySQL, Want: wantMySQL},
			{Name: "mariadb", Key: SqlDialectMariaDB, Want: wantMySQL},
			{Name: "postgres", Key: SqlDialectPostgreSQL, Want: wantPostgreSQL},
			{Name: "sqlite3", Key: SqlDialectSqlite3, Want: wantSqlite},
		}
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	f8n = New().
		FieldDatatype("id", TypeUint64).
		FieldDatatype("key", TypeInt64).
		FieldDatatype("ratio", TypeFloat64).
		FieldDatatype("flag", TypeBool).
		FieldDatatype("created", TypeTime).
		FieldDatatype("tags", TypeSliceString).
		FieldDatatype("ids", TypeSliceInt64).
		FieldDatatype("dur", TypeDuration).
		FieldDatatype("blob", TypeBytes).
		FieldDatatype("ip", TypeIP).
		FieldDatatype("js", TypeRawMessage).
		FieldDatatype("u", TypeUint64).
		FieldDatatype("deletedAt", TypeTimeNil)
	if err = f8n.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	for n = range dialects {
		if ret, err = f8n.Sql(dialects[n].Key); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		if ret != dialects[n].Want {
			t.Errorf(tplErrorDialect, dialects[n].Name, ret, dialects[n].Want)
		}
	}
	// Переопределение поля, кавычки диалекта применяются и к наименованию полученному из переопределения.
	f8n = New().FieldDatatype("id", TypeUint64).Redefinition("id", "`cursor`.`id`")
	if rq, err = http.NewRequest("GET", "http://localhost?filter=id:eq:1", nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if err = f8n.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	if ret, err = f8n.Sql(SqlDialectMySQL); err != nil || ret != wantRemapMySQL {
		t.Errorf(tplErrorDialect, "mysql-remap", ret, wantRemapMySQL)
	}
	if ret, err = f8n.Sql(SqlDialectPostgreSQL); err != nil || ret != wantRemapPostgre {
		t.Errorf(tplErrorDialect, "postgres-remap", ret, wantRemapPostgre)
	}
}

// Тестирование создания SQL условия WHERE в устаревшем режиме фильтрации tie=or.
func TestSqlTieOr(t *testing.T) {
	const (
		datUrn          = `http://localhost?filter=id:eq:1&filter=key:eq:2&tie=or`
		wantMySQL       = "WHERE `id` = 1 OR `key` = 2"
		wantPostgreSQL  = `WHERE "id" = 1 OR "key" = 2`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseRequest() != nil, ошибка не ожидалась. Ошибка: %s"
		tplErrorDialect = "Sql() = %q, ожидалось: %q."
	)
	var (
		err error
		rq  *http.Request
		f8n Interface
		ret string
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	f8n = New().FieldDatatype("id", TypeUint64).FieldDatatype("key", TypeInt64)
	if err = f8n.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	if ret, err = f8n.Sql(SqlDialectMySQL); err != nil || ret != wantMySQL {
		t.Errorf(tplErrorDialect, ret, wantMySQL)
	}
	if ret, err = f8n.Sql(SqlDialectPostgreSQL); err != nil || ret != wantPostgreSQL {
		t.Errorf(tplErrorDialect, ret, wantPostgreSQL)
	}
}

// Тестирование создания SQL условия WHERE без условий фильтрации и с не поддерживаемым диалектом.
func TestSqlErrors(t *testing.T) {
	const (
		dialectUnknown = "oracle"
		tplErrorNil    = "Sql() == nil, ожидалась ошибка SqlNotImplemented."
		tplErrorIs     = "Sql() = %q, ожидалось: %q"
		tplErrorEmpty  = "Sql() = %q, ожидалась пустая строка."
		tplErrorNotNil = "Sql() != nil, ошибка не ожидалась."
	)
	var (
		err error
		f8n Interface
		ret string
	)

	f8n = New()
	if ret, err = f8n.Sql(SqlDialectMySQL); err != nil || ret != "" {
		t.Errorf(tplErrorEmpty, ret)
		return
	}
	if ret, err = f8n.Sql(dialectUnknown); err == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if !Errors().SqlNotImplemented.Is(err) {
		t.Errorf(tplErrorIs, err, Errors().SqlNotImplemented.Bind())
		return
	}
	if ret != "" {
		t.Errorf(tplErrorEmpty, ret)
		return
	}
	// Сброс условий возвращает результат к пустому условию.
	f8n.FieldDatatype("id", TypeUint64)
	f8n.(*impl).Filter = []Filter{{Field: "id", Method: filterEquivalent, Value: FilterValue{Source: "1", Type: TypeUint64}}}
	if ret, err = f8n.Sql(SqlDialectMySQL); err != nil || ret == "" {
		t.Errorf(tplErrorNotNil)
	}
	f8n.Reset()
	if ret, err = f8n.Sql(SqlDialectMySQL); err != nil || ret != "" {
		t.Errorf(tplErrorEmpty, ret)
	}
}

// Тестирование форматирования значений литералами диалекта.
func TestSqlLiteral(t *testing.T) {
	const (
		tplLiteral = "sqlLiteral(%v) = %q, ожидалось: %q."
		tplFloat   = "sqlLiteral(%v) = %q, ожидалось: %q."
	)
	type test struct {
		Value any
		Want  string
	}
	var (
		tmo   time.Time
		err   error
		ret   string
		n     int
		cust  struct{ Name string }
		tests []test
	)

	tmo, err = time.Parse(time.RFC3339, "2018-08-03T03:00:00+03:00")
	if err != nil {
		t.Errorf("Ошибка подготовки тестовых данных: %s.", err)
		return
	}
	cust.Name = "custom"
	tests = []test{
		// Строки и экранирование.
		{Value: "просто строка", Want: "'просто строка'"},
		{Value: "o'clock", Want: "'o''clock'"},
		// NULL.
		{Value: nil, Want: "NULL"},
		// Булевы значения.
		{Value: true, Want: "1"},
		{Value: false, Want: "0"},
		// Числа.
		{Value: int64(-9223372036854775808), Want: "-9223372036854775808"},
		{Value: uint64(18446744073709551615), Want: "18446744073709551615"},
		{Value: float64(-3.1415), Want: "-3.1415"},
		// Время.
		{Value: tmo, Want: "'2018-08-03T03:00:00+03:00'"},
		// Срезы.
		{Value: []string{"a", "o'b"}, Want: "'a', 'o''b'"},
		{Value: []int64{3, -2}, Want: "3, -2"},
		{Value: []uint64{7, 8}, Want: "7, 8"},
		// Двоичные данные, JSON, IP адрес, длительность.
		{Value: []byte("hello"), Want: "0x68656c6c6f"},
		{Value: json.RawMessage(`{"a":1}`), Want: `'{"a":1}'`},
		{Value: net.ParseIP("192.168.0.1"), Want: `'192.168.0.1'`},
		{Value: 90 * time.Minute, Want: "5400000000000"},
		// Не известный тип значения.
		{Value: cust, Want: "'{custom}'"},
	}
	for n = range tests {
		if ret = sqlLiteral(SqlDialectMySQL, tests[n].Value); ret != tests[n].Want {
			t.Errorf(tplLiteral, tests[n].Value, ret, tests[n].Want)
			return
		}
	}
	// Булевы значения PostgreSQL.
	if ret = sqlLiteral(SqlDialectPostgreSQL, true); ret != "TRUE" {
		t.Errorf(tplLiteral, true, ret, "TRUE")
		return
	}
	if ret = sqlLiteral(SqlDialectPostgreSQL, false); ret != "FALSE" {
		t.Errorf(tplLiteral, false, ret, "FALSE")
		return
	}
	// Двоичные данные для диалектов PostgreSQL, MariaDB и sqlite3.
	if ret = sqlLiteral(SqlDialectPostgreSQL, []byte("hi")); ret != `'\x6869'::bytea` {
		t.Errorf(tplLiteral, []byte("hi"), ret, `'\x6869'::bytea`)
		return
	}
	if ret = sqlLiteral(SqlDialectMariaDB, []byte("hi")); ret != "0x6869" {
		t.Errorf(tplLiteral, []byte("hi"), ret, "0x6869")
		return
	}
	if ret = sqlLiteral(SqlDialectSqlite3, []byte("hi")); ret != "X'6869'" {
		t.Errorf(tplLiteral, []byte("hi"), ret, "X'6869'")
		return
	}
	// Шестнадцатеричное кодирование.
	if ret = hexEncode([]byte{0x01, 0xab, 0xff}); ret != "01abff" {
		t.Errorf(tplLiteral, []byte{0x01, 0xab, 0xff}, ret, "01abff")
		return
	}
	// Не известный тип значения PostgreSQL.
	if ret = sqlLiteral(SqlDialectPostgreSQL, 5); ret != "'5'" {
		t.Errorf(tplLiteral, 5, ret, "'5'")
		return
	}
	// Подстановка значений в запрос.
	ret = sqlSubstitute("`a` = ? AND `b` IN (?) AND `c` IS NULL", []any{"o'clock", []uint64{1, 2}}, SqlDialectMySQL)
	if ret != "`a` = 'o''clock' AND `b` IN (1, 2) AND `c` IS NULL" {
		t.Errorf(tplFloat, "substitute", ret, "`a` = 'o''clock' AND `b` IN (1, 2) AND `c` IS NULL")
		return
	}
}
