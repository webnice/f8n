package f8n

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/webnice/f8n/testdb"

	"gorm.io/gorm"
)

// Подготовка тестовой таблицы и контроль её целостности.
func injectionTestSetup(t *testing.T) (count int64) {
	const (
		tplDelete = "DELETE FROM `children`"
		tplCount  = "SELECT COUNT(*) FROM `children`"
		tplError  = "Ошибка подготовки тестовой таблицы: %s."
	)
	var err error

	if err = testdb.DB.Exec(tplDelete).Error; err != nil {
		t.Errorf(tplError, err)
		return
	}
	if err = testdb.DB.Create(&testdb.Child{Name: "Пётр"}).Error; err != nil {
		t.Errorf(tplError, err)
		return
	}
	if err = testdb.DB.Raw(tplCount).Scan(&count).Error; err != nil {
		t.Errorf(tplError, err)
		return
	}

	return
}

// Проверка целостности тестовой таблицы после попытки инъекции.
func injectionTestIntegrity(t *testing.T, count int64) {
	const (
		tplCount  = "SELECT COUNT(*) FROM `children`"
		tplError  = "Ошибка проверки целостности таблицы: %s."
		tplTable  = "Таблица children не существует после попытки SQL инъекции."
		tplChange = "Количество записей изменилось после попытки SQL инъекции: до=%d, после=%d."
	)
	var (
		err    error
		after  int64
		exists bool
	)

	if exists = testdb.DB.Migrator().HasTable(&testdb.Child{}); !exists {
		t.Errorf(tplTable)
		return
	}
	if err = testdb.DB.Raw(tplCount).Scan(&after).Error; err != nil {
		t.Errorf(tplError, err)
		return
	}
	if after != count {
		t.Errorf(tplChange, count, after)
		return
	}
}

// Тестирование защиты от SQL инъекций через значение простой фильтрации.
// Значения фильтрации всегда передаются в SQL запрос параметрами, инъекция через значение невозможна.
//
//goland:noinspection ALL
func TestSQLInjectionFilterValue(t *testing.T) {
	const (
		qTpl            = `http://localhost?filter=name:eq:%s`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseRequest() != nil, ошибка не ожидалась. Ошибка: %s"
		tplErrorRows    = "Количество найденных записей = %d, ожидалось %d."
	)
	var (
		err      error
		rq       *http.Request
		db       *gorm.DB
		items    []testdb.Child
		n        int
		payloads = []string{
			"' OR '1'='1",
			"' OR '1'='1'--",
			"' OR '1'='1'#",
			`" OR "1"="1`,
			"1; DROP TABLE children;--",
			"1; DELETE FROM children;--",
			"' UNION SELECT id, name FROM children--",
			"' AND SLEEP(10)--",
		}
	)

	count := injectionTestSetup(t)
	for n = range payloads {
		var f8n Interface

		f8n = New().FieldDatatype("name", TypeString)
		if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, url.QueryEscape(payloads[n])), nil); err != nil {
			t.Errorf(tplErrorRequest, err)
		}
		if err = f8n.ParseRequest(rq); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		// Выполнение запроса на реальной базе данных, значение инъекции должно обрабатываться как данные.
		if db, err = f8n.Gorm(testdb.DB); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		items = items[:0]
		if err = db.Find(&items).Error; err != nil {
			// Значение обработано как данные, синтаксической ошибки SQL быть не должно, но допускается ошибка типов.
			t.Logf("Запрос прерван ошибкой: %s.", err)
		}
		if len(items) != 0 {
			t.Errorf(tplErrorRows, len(items), 0)
			continue
		}
		injectionTestIntegrity(t, count)
	}
	// Контрольный запрос, доказывающий что фильтрация реально работает.
	var f8n Interface

	f8n = New().FieldDatatype("name", TypeString)
	if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, url.QueryEscape("Пётр")), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if err = f8n.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	if db, err = f8n.Gorm(testdb.DB); err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	items = items[:0]
	if err = db.Find(&items).Error; err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	if len(items) != 1 {
		t.Errorf(tplErrorRows, len(items), 1)
	}
}

// Тестирование защиты от SQL инъекций через наименование поля простой фильтрации.
// Наименование поля очищается от символов позволяющих выйти за пределы идентификатора поля.
//
//goland:noinspection ALL
func TestSQLInjectionFilterField(t *testing.T) {
	const (
		qTpl            = `http://localhost?filter=%s:eq:1`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseRequest() != nil, ошибка не ожидалась. Ошибка: %s"
	)
	var (
		err      error
		rq       *http.Request
		db       *gorm.DB
		items    []testdb.Child
		n        int
		payloads = []string{
			"`",
			"` IS NULL--",
			"`=1 OR `1",
			"`=`1` OR `a",
			"a`, `b",
			"a' OR '1'='1",
			`a" OR "1"="1`,
			"name` = `x` OR `1",
		}
	)

	count := injectionTestSetup(t)
	for n = range payloads {
		var f8n Interface

		f8n = New().FieldDatatype("name", TypeString)
		if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, url.QueryEscape(payloads[n])), nil); err != nil {
			t.Errorf(tplErrorRequest, err)
		}
		if err = f8n.ParseRequest(rq); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		// Выполнение запроса, очищенное наименование поля не является существующей колонкой.
		if db, err = f8n.Gorm(testdb.DB); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		items = items[:0]
		if err = db.Find(&items).Error; err == nil && len(items) != 0 {
			t.Errorf("Найдено записей: %d, ожидалось %d.", len(items), 0)
			continue
		}
		injectionTestIntegrity(t, count)
	}
}

// Тестирование защиты от SQL инъекций через наименование поля сортировки.
//
//goland:noinspection ALL
func TestSQLInjectionByField(t *testing.T) {
	const (
		qTpl            = `http://localhost?by=%s:asc`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseRequest() != nil, ошибка не ожидалась. Ошибка: %s"
	)
	var (
		err      error
		rq       *http.Request
		db       *gorm.DB
		items    []testdb.Child
		n        int
		payloads = []string{
			"`",
			"` ASC, (SELECT 1)--",
			"a`, (SELECT 1)--",
			"a' OR '1'='1",
			`a" OR "1"="1`,
		}
	)

	count := injectionTestSetup(t)
	for n = range payloads {
		var f8n Interface

		f8n = New()
		if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, url.QueryEscape(payloads[n])), nil); err != nil {
			t.Errorf(tplErrorRequest, err)
		}
		if err = f8n.ParseRequest(rq); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		if db, err = f8n.Gorm(testdb.DB); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		items = items[:0]
		if err = db.Find(&items).Error; err == nil && len(items) != 1 {
			t.Errorf("Найдено записей: %d, ожидалось %d.", len(items), 1)
			continue
		}
		injectionTestIntegrity(t, count)
	}
}

// Тестирование защиты от SQL инъекций через наименование поля фильтрации сложной фильтрации (map).
//
//goland:noinspection ALL
func TestSQLInjectionMapField(t *testing.T) {
	const (
		qTpl            = `http://localhost?map=f1:and:%s&f1=name:eq:Пётр`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseRequest() != nil, ошибка не ожидалась. Ошибка: %s"
	)
	var (
		err      error
		rq       *http.Request
		db       *gorm.DB
		items    []testdb.Child
		n        int
		payloads = []string{
			"f2",
			"a` IS NULL--",
			"a', '1'='1",
			`a" OR "1"="1`,
		}
	)

	count := injectionTestSetup(t)
	for n = range payloads {
		var f8n Interface

		f8n = New().FieldDatatype("name", TypeString)
		if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, url.QueryEscape(payloads[n])), nil); err != nil {
			t.Errorf(tplErrorRequest, err)
		}
		if err = f8n.ParseRequest(rq); err != nil {
			// Имена фильтров, которым не передано значение, отклоняются ошибкой.
			if !Errors().FilterCalledByNameWasNotFound.Is(err) {
				t.Errorf(tplErrorParse, err)
			}
			continue
		}
		if db, err = f8n.Gorm(testdb.DB); err != nil {
			t.Errorf(tplErrorParse, err)
			continue
		}
		items = items[:0]
		if err = db.Find(&items).Error; err == nil && len(items) != 0 {
			t.Errorf("Найдено записей: %d, ожидалось %d.", len(items), 0)
			continue
		}
		injectionTestIntegrity(t, count)
	}
}

// Тестирование защиты от SQL инъекций через способ сравнения простой фильтрации.
// Способ сравнения проверяется по справочнику допустимых значений.
func TestSQLInjectionFilterMethod(t *testing.T) {
	const (
		datUrn          = `http://localhost?filter=id:eq' OR '1'='1:1`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorNil     = "ParseRequest() == nil, ожидалась ошибка FilterUnknownMethod."
		tplErrorIs      = "ParseRequest() = %q, ожидалось: %q"
	)
	var (
		err error
		rq  *http.Request
		f8n Interface
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	f8n = New().FieldDatatype("id", TypeUint64)
	if err = f8n.ParseRequest(rq); err == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if !Errors().FilterUnknownMethod.Is(err) {
		t.Errorf(tplErrorIs, err, Errors().FilterUnknownMethod.Bind())
		return
	}
}

// Тестирование защиты от SQL инъекций через значение лимита.
// Значение лимита конвертируется в число, не числовые значения отклоняются ошибкой.
func TestSQLInjectionLimit(t *testing.T) {
	const (
		qTpl            = `http://localhost?limit=%s`
		datPayload      = `1; DROP TABLE children;--`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorNil     = "ParseRequest() == nil, ожидалась ошибка LimitInvalidValue."
		tplErrorIs      = "ParseRequest() = %q, ожидалось: %q"
	)
	var (
		err error
		rq  *http.Request
		f8n Interface
	)

	if rq, err = http.NewRequest("GET", fmt.Sprintf(qTpl, url.QueryEscape(datPayload)), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	f8n = New()
	if err = f8n.ParseRequest(rq); err == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if !Errors().LimitInvalidValue.Is(err) {
		t.Errorf(tplErrorIs, err, Errors().LimitInvalidValue.Bind())
		return
	}
	injectionTestIntegrity(t, injectionTestSetup(t))
}
