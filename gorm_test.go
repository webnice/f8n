package f8n

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/webnice/f8n/testdb"

	"gorm.io/gorm"
)

// Проверка тестового окружения базы данных.
func testGorm(t *testing.T) (ret *gorm.DB) {
	const tplQueryTables = "SELECT `name` FROM `sqlite_master` WHERE type=? ORDER BY `name`"
	var (
		err        error
		sqlDB      *sql.DB
		rows       *sql.Rows
		tmp        string
		tables     []string
		isChildren bool
		isParents  bool
	)

	ret = testdb.DB
	if sqlDB, err = ret.DB(); err != nil {
		t.Errorf("DB(). Ошибка не ожидалась. Ошибка: %s", err)
	}
	if rows, err = sqlDB.Query(tplQueryTables, "table"); err != nil {
		t.Errorf("Query(). Ошибка не ожидалась. Ошибка: %s", err)
	}
	for rows.Next() {
		if err = rows.Scan(&tmp); err != nil {
			t.Errorf("Scan(). Ошибка не ожидалась. Ошибка: %s", err)
		}
		if tmp == "children" {
			isChildren = true
		}
		if tmp == "parents" {
			isParents = true
		}
		tables = append(tables, tmp)
	}
	// Проверка создались ли в базе данных таблицы для тестирования.
	if !isChildren || !isParents {
		t.Errorf(
			"Не верные таблицы для тестирования. В базе данных таблицы: %v, ожидались: %v.",
			tables, []string{"children", "parents"},
		)
	}

	return
}

// Тестирование вызова функции с не инициализированным объектом orm.
func TestGormNil(t *testing.T) {
	var (
		err error
		f8n Interface
		orm *gorm.DB
	)

	f8n = New()
	if orm, err = f8n.Gorm(orm); err == nil {
		t.Errorf("Gorm(). Ожидалась ошибка.")
	}
	if orm != nil {
		t.Errorf("Gorm(). Ожидался объект равный nil.")
	}
	if !f8n.Errors().OrmIsNil.Is(err) {
		t.Errorf("Gorm(). Ошибка: %q, ожидалась: %q", err, f8n.Errors().OrmIsNil.Bind())
		return
	}
}

// Тестирование создания ORM условий запроса.
func TestGorm(t *testing.T) {
	const qTpl = `http://localhost?%s`
	type (
		ft struct {
			Field string
			Type  FieldType
		}
		gs struct {
			Urn    string
			Fields []ft
			Query  string
		}
	)

	var (
		err      error
		rq       *http.Request
		f8n      Interface
		orm      *gorm.DB
		dryOrm   *gorm.DB
		stmt     *gorm.Statement
		tmp      testdb.Child
		query    string
		rexQuery []string
		n, j     int
		tests    = []gs{
			{
				Urn:    "filter=id:eq:1",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "`id` = 1",
			},
			{
				Urn:    "filter=id:le:11&filter=key:ge:-32",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "`id` <= 11 AND `key` >= -32",
			},
			{
				Urn:    "filter=dur:ge:1h30m&filter=ip:eq:192.168.0.1&filter=blob:eq:aGVsbG8=",
				Fields: []ft{{Field: "dur", Type: TypeDuration}, {Field: "ip", Type: TypeIP}, {Field: "blob", Type: TypeBytes}},
				Query:  "`dur` >= 5400000000000 AND `ip` = \"192.168.0.1\" AND `blob` = \"hello\"",
			},
		}
	)

	f8n, orm = New(), testGorm(t)
	for n = range tests {
		// Создание запроса.
		if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, tests[n].Urn), nil); err != nil {
			t.Errorf("Ошибка создания запроса: %s", err)
		}
		// Подготовка объекта фильтрации.
		f8n.Reset()
		for j = range tests[n].Fields {
			f8n.FieldDatatype(tests[n].Fields[j].Field, tests[n].Fields[j].Type)
		}
		// Разбор запроса.
		if err = f8n.ParseRequest(rq); err != nil {
			t.Errorf("ParseRequest(). Ошибка не ожидалась. Ошибка: %s", err)
		}
		// Добавление условий из запроса в ORM.
		dryOrm = orm.Session(&gorm.Session{DryRun: true})
		if dryOrm, err = f8n.Gorm(dryOrm); err != nil {
			t.Errorf("Gorm(). Ошибка не ожидалась. Ошибка: %s", err)
		}
		stmt = dryOrm.Find(&tmp).Statement
		query = orm.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
		if rexQuery = testdb.RexWhere.FindStringSubmatch(query); len(rexQuery) < 4 {
			t.Errorf("SQL запрос: %s", query)
			continue
		}
		if !strings.EqualFold(rexQuery[3], tests[n].Query) {
			t.Log("Ошибка формирования запроса.")
			t.Errorf("Получен запрос:\n %q, ожидался:\n %q", rexQuery[3], tests[n].Query)
		}
	}
}

// Тестирование создания ORM условий запроса со всеми разобранными значениями.
func TestGormFull(t *testing.T) {
	const (
		qTpl          = `http://localhost?%s`
		datUrn        = "filter=id:eq:1&filter=key:eq:2&tie=or&by=id:asc&by=name:desc&limit=5:10"
		wantQuery     = "`id` = 1 OR `key` = 2 ORDER BY `id` ASC,`name` DESC LIMIT 10 OFFSET 5"
		wantOrder     = "ORDER BY `id` ASC,`name` DESC"
		wantLimit     = "LIMIT 10 OFFSET 5"
		tplErrorUrn   = "Ошибка создания запроса: %s."
		tplErrorParse = "ParseRequest(). Ошибка не ожидалась. Ошибка: %s."
		tplErrorGorm  = "Gorm(). Ошибка не ожидалась. Ошибка: %s."
	)
	var (
		err      error
		rq       *http.Request
		f8n      Interface
		orm      *gorm.DB
		dryOrm   *gorm.DB
		stmt     *gorm.Statement
		tmp      testdb.Child
		query    string
		rexQuery []string
	)

	f8n, orm = New(), testGorm(t)
	f8n.
		FieldDatatype("id", TypeUint64).
		FieldDatatype("key", TypeInt64)
	// Создание запроса.
	if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, datUrn), nil); err != nil {
		t.Errorf(tplErrorUrn, err)
	}
	// Разбор запроса.
	if err = f8n.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
	}
	// Добавление условий из запроса в ORM.
	dryOrm = orm.Session(&gorm.Session{DryRun: true})
	if dryOrm, err = f8n.Gorm(dryOrm); err != nil {
		t.Errorf(tplErrorGorm, err)
	}
	stmt = dryOrm.Find(&tmp).Statement
	query = orm.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	// Проверка фильтрации, режим tie=or.
	if rexQuery = testdb.RexWhere.FindStringSubmatch(query); len(rexQuery) < 4 {
		t.Errorf("SQL запрос: %s", query)
		return
	}
	if !strings.EqualFold(rexQuery[3], wantQuery) {
		t.Errorf("Условие WHERE запроса:\n %q, ожидалось:\n %q", rexQuery[3], wantQuery)
	}
	// Проверка сортировки.
	if !strings.Contains(query, wantOrder) {
		t.Errorf("SQL запрос: %s, ожидалась сортировка: %s", query, wantOrder)
	}
	// Проверка лимитов.
	if !strings.Contains(query, wantLimit) {
		t.Errorf("SQL запрос: %s, ожидались лимиты: %s", query, wantLimit)
	}
}

// Тестирование создания ORM условий запроса с переопределением названий полей.
func TestGormRedefinition(t *testing.T) {
	const (
		qTpl          = `http://localhost?%s`
		datUrn        = "filter=id:eq:1"
		wantQuery     = "`cursor`.`id` = 1"
		tplErrorUrn   = "Ошибка создания запроса: %s."
		tplErrorParse = "ParseRequest(). Ошибка не ожидалась. Ошибка: %s."
		tplErrorGorm  = "Gorm(). Ошибка не ожидалась. Ошибка: %s."
	)
	var (
		err      error
		rq       *http.Request
		f8n      Interface
		orm      *gorm.DB
		dryOrm   *gorm.DB
		stmt     *gorm.Statement
		tmp      testdb.Child
		query    string
		rexQuery []string
	)

	f8n, orm = New(), testGorm(t)
	f8n.
		FieldDatatype("id", TypeUint64).
		Redefinition("id", "`cursor`.`id`")
	// Создание запроса.
	if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, datUrn), nil); err != nil {
		t.Errorf(tplErrorUrn, err)
	}
	// Разбор запроса.
	if err = f8n.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
	}
	// Добавление условий из запроса в ORM.
	dryOrm = orm.Session(&gorm.Session{DryRun: true})
	if dryOrm, err = f8n.Gorm(dryOrm); err != nil {
		t.Errorf(tplErrorGorm, err)
	}
	stmt = dryOrm.Find(&tmp).Statement
	query = orm.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	// Проверка фильтрации с переопределённым названием поля.
	if rexQuery = testdb.RexWhere.FindStringSubmatch(query); len(rexQuery) < 4 {
		t.Errorf("SQL запрос: %s", query)
		return
	}
	if !strings.EqualFold(rexQuery[3], wantQuery) {
		t.Errorf("Условие WHERE запроса:\n %q, ожидалось:\n %q", rexQuery[3], wantQuery)
	}
}

// Тестирование создания ORM условий запроса сложной фильтрации (map).
//
//goland:noinspection ALL
func TestGormMap(t *testing.T) {
	const (
		qTpl          = `http://localhost?%s`
		tplErrorUrn   = "Ошибка создания запроса: %s."
		tplErrorParse = "ParseRequest(). Ошибка не ожидалась. Ошибка: %s."
		tplErrorGorm  = "Gorm(). Ошибка не ожидалась. Ошибка: %s."
		tplErrorSQL   = "SQL запрос: %s"
		tplErrorWhere = "Условие WHERE запроса:\n %q, ожидалось:\n %q"
	)
	type (
		ft struct {
			Field string
			Type  FieldType
		}
		gs struct {
			Name   string
			Urn    string
			Fields []ft
			Query  string
		}
	)
	var (
		err      error
		rq       *http.Request
		orm      *gorm.DB
		dryOrm   *gorm.DB
		stmt     *gorm.Statement
		tmp      testdb.Child
		query    string
		rexQuery []string
		n, j     int
		f8n      Interface
		tests    = []gs{
			{
				Name: "bracket-and-or-with-order-limit",
				Urn: "map=(f1:and:f2):or:f3&f1=id:eq:1&f2=key:eq:2&f3=name:ke:go*" +
					"&by=id:asc&limit=0:10",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "((`id` = 1 AND `key` = 2) OR `name` LIKE \"go%\") ORDER BY `id` ASC LIMIT 10",
			},
			{
				Name:   "or-chain",
				Urn:    "map=f1:or:f2&f1=id:eq:1&f2=key:eq:2",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "(`id` = 1 OR `key` = 2)",
			},
			{
				Name: "nested-brackets",
				Urn:  "map=((f1:or:f2):and:(f3:or:f4))&f1=id:eq:1&f2=key:eq:2&f3=name:eq:вася&f4=age:ge:18",
				Fields: []ft{
					{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64},
					{Field: "age", Type: TypeUint64},
				},
				Query: "(((`id` = 1 OR `key` = 2) AND (`name` = \"вася\" OR `age` >= 18)))",
			},
			{
				Name: "filter-and-map",
				Urn:  "filter=flag:eq:true&map=f1:or:f2&f1=id:eq:1&f2=key:eq:2",
				Fields: []ft{
					{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64},
					{Field: "flag", Type: TypeBool},
				},
				Query: "((`id` = 1 OR `key` = 2)) AND `flag` = true",
			},
			{
				Name:   "in-and-nil",
				Urn:    "map=f1:and:f2&f1=id:in:1,3,5&f2=deletedAt:eq:nil",
				Fields: []ft{{Field: "id", Type: TypeSliceUint64}, {Field: "deletedAt", Type: TypeTimeNil}},
				Query:  "(`id` IN (1,3,5) AND `deletedAt` IS NULL)",
			},
			{
				Name:   "multiple-map-or",
				Urn:    "map=f1:and:f2&map=f3:and:f4&f1=id:eq:1&f2=key:eq:2&f3=u:eq:3&f4=v:eq:4",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "((`id` = 1 AND `key` = 2) OR (`u` = \"3\" AND `v` = \"4\"))",
			},
			{
				// Устаревший режим tie=or игнорируется при использовании map, фильтры объединяются по «И».
				Name: "tie-ignored-with-map",
				Urn:  "filter=id:eq:1&filter=key:eq:2&tie=or&map=f1&f1=u:eq:3",
				Fields: []ft{
					{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64},
				},
				Query: "(`u` = \"3\") AND (`id` = 1 AND `key` = 2)",
			},
			{
				// Полная формула: (filter1 И filter2) И (map1 ИЛИ map2).
				Name: "full-formula",
				Urn: "filter=id:eq:1&filter=key:eq:2&map=f1&map=f2:and:f3" +
					"&f1=u:eq:3&f2=v:eq:4&f3=w:eq:5",
				Fields: []ft{
					{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64},
				},
				Query: "((`u` = \"3\" OR (`v` = \"4\" AND `w` = \"5\"))) AND (`id` = 1 AND `key` = 2)",
			},
			{
				// Операнды объединения «ИЛИ» уже заключены в операторные скобки.
				Name:   "multiple-map-pre-bracketed",
				Urn:    "map=(f1:and:f2)&map=(f3:and:f4)&f1=id:eq:1&f2=key:eq:2&f3=u:eq:3&f4=v:eq:4",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "((`id` = 1 AND `key` = 2) OR (`u` = \"3\" AND `v` = \"4\"))",
			},
		}
	)

	orm = testGorm(t)
	for n = range tests {
		if !t.Run(tests[n].Name, func(t *testing.T) {
			// Подготовка объекта фильтрации.
			f8n = New()
			for j = range tests[n].Fields {
				f8n.FieldDatatype(tests[n].Fields[j].Field, tests[n].Fields[j].Type)
			}
			// Создание запроса.
			if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, tests[n].Urn), nil); err != nil {
				t.Errorf(tplErrorUrn, err)
			}
			// Разбор запроса.
			if err = f8n.ParseRequest(rq); err != nil {
				t.Errorf(tplErrorParse, err)
			}
			// Добавление условий из запроса в ORM.
			dryOrm = orm.Session(&gorm.Session{DryRun: true})
			if dryOrm, err = f8n.Gorm(dryOrm); err != nil {
				t.Errorf(tplErrorGorm, err)
			}
			stmt = dryOrm.Find(&tmp).Statement
			query = orm.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
			// Проверка условия WHERE запроса.
			if rexQuery = testdb.RexWhere.FindStringSubmatch(query); len(rexQuery) < 4 {
				t.Errorf(tplErrorSQL, query)
				return
			}
			if !strings.EqualFold(rexQuery[3], tests[n].Query) {
				t.Errorf(tplErrorWhere, rexQuery[3], tests[n].Query)
			}
		}) {
			return
		}
	}
}

// Тестирование создания ORM условий запроса с пустой картой сложной фильтрации.
func TestGormMapEmpty(t *testing.T) {
	const (
		tplErrorGorm  = "Gorm(). Ошибка не ожидалась. Ошибка: %s."
		tplErrorWhere = "SQL запрос: %s, условие WHERE не ожидалось."
	)
	var (
		err    error
		obj    *impl
		orm    *gorm.DB
		dryOrm *gorm.DB
		stmt   *gorm.Statement
		tmp    testdb.Child
		query  string
	)

	obj = New().(*impl)
	// Пустое DOM дерево карты сложной фильтрации.
	obj.Map = &Map{parent: obj, Origin: OriginUnknown}
	orm = testGorm(t)
	dryOrm = orm.Session(&gorm.Session{DryRun: true})
	if dryOrm, err = obj.Gorm(dryOrm); err != nil {
		t.Errorf(tplErrorGorm, err)
	}
	stmt = dryOrm.Find(&tmp).Statement
	query = orm.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	if rexQuery := testdb.RexWhere.FindStringSubmatch(query); len(rexQuery) >= 4 {
		t.Errorf(tplErrorWhere, query)
		return
	}
}
