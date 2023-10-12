// Package f8n
package f8n

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/webnice/f8n/testdb"

	ddd "github.com/webnice/debug"

	"gorm.io/gorm"
)

func testGorm(t *testing.T) (ret *gorm.DB) {
	var (
		err    error
		sqlDB  *sql.DB
		rows   *sql.Rows
		tmp    string
		tables []string
	)

	ret = testdb.DB
	if sqlDB, err = ret.DB(); err != nil {
		t.Errorf("DB(). Ошибка не ожидалась. Ошибка: %s", err)
	}
	if rows, err = sqlDB.Query("SELECT `name` FROM `sqlite_master` WHERE type=? ORDER BY `name`", "table"); err != nil {
		t.Errorf("Query(). Ошибка не ожидалась. Ошибка: %s", err)
	}
	for rows.Next() {
		if err = rows.Scan(&tmp); err != nil {
			t.Errorf("Scan(). Ошибка не ожидалась. Ошибка: %s", err)
		}
		tables = append(tables, tmp)
	}
	// Проверка создались ли в базе данных таблицы для тестирования.
	if !reflect.DeepEqual(tables, []string{"children", "parents"}) {
		t.Errorf(
			"Не верные таблицы для тестирования. В баде данных таблицы: %v, ожидались: %v.",
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
		ero Err
		ok  bool
	)

	f8n = New()
	if orm, err = f8n.Gorm(orm); err == nil {
		t.Errorf("Gorm(). Ожидалась ошибка.")
	}
	if orm != nil {
		t.Errorf("Gorm(). Ожидался объект равный nil.")
	}
	if ero, ok = err.(Err); !ok {
		t.Errorf("Gorm() = %q, не верный тип ошибки.", err)
		return
	}
	if ero.Anchor() != f8n.Errors().OrmIsNil().Anchor() {
		t.Errorf("Gorm(). Ошибка: %q, ожидалась: %q", ero.Error(), f8n.Errors().OrmIsNil().Error())
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
		}
		if !strings.EqualFold(rexQuery[3], tests[n].Query) {
			t.Log("Ошибка формирования запроса.")
			t.Errorf("Получен запрос:\n %q, ожидался:\n %q", rexQuery[3], tests[n].Query)
		}
	}

	ddd.Nop()
}
