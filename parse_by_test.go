package f8n

import (
	"net/http"
	"reflect"
	"testing"
)

// Тестирование сортировки.
//
//goland:noinspection DuplicatedCode
func TestParseBy(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	// Пустой массив сортировки.
	if rq, err = http.NewRequest("GET", "http://localhost", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseBy(rq); len(ers) > 0 {
		t.Errorf("ParseBy() = %v, ошибка не ожидалась.", ers)
		return
	}
	if len(obj.By) > 0 {
		t.Errorf("By[] = %v, ожидался пустой срез.", obj.By)
		return
	}
	// Два значения сортировки.
	if rq, err = http.NewRequest("GET", "http://localhost?by=id:asc&by=name:desc", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseBy(rq); len(ers) > 0 {
		t.Errorf("ParseBy() = %v, ошибка не ожидалась.", ers)
		return
	}
	if len(obj.By) != 2 {
		t.Errorf("By[] = %v, ожидался срез из двух значений.", obj.By)
		return
	}
	// Не верный формат значения сортировки.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", "http://localhost?by=id-asc", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseBy(rq); len(ers) == 0 {
		t.Errorf("ParseBy() = nil, ожидалась ошибка.")
		return
	}
	for n = range ers {
		if !Errors().ByFormat.Is(ers[n].Ei) {
			t.Errorf("ParseBy() = %q, ожидалось: %q", ers[n].Ei, Errors().ByFormat.Bind())
			return
		}
	}
	// Не верная константа сортировки.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", "http://localhost?by=id:one&by=name:two", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseBy(rq); len(ers) != 2 {
		t.Errorf("ParseBy() = nil, ожидалось две ошибки.")
		return
	}
	for n = range ers {
		if !Errors().ByDirectionUnknown.Is(ers[n].Ei) {
			t.Errorf("ParseBy() = %q, ожидалось: %q", ers[n].Ei, Errors().ByDirectionUnknown.Bind())
			return
		}
	}
}

// Тестирование сортировки. Строгие и не строгие поля.
//
//goland:noinspection ALL
func TestParseByField(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	// Не строгие поля.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", "http://localhost?by=aaa:asc&by=bbb:desc", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseBy(rq); len(ers) > 0 {
		t.Errorf("ParseBy() = %v, ошибка не ожидалась.", ers)
		return
	}
	// Строгие поля.
	obj = New().FieldSet("id", "name").(*impl)
	if rq, err = http.
		NewRequest("GET", "http://localhost?by=aaa:asc&by=bbb:desc&by=id:asc&by=name:desc", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseBy(rq); len(ers) != 2 {
		t.Errorf("ParseBy() = nil, ожидалось две ошибки.")
		return
	}
	for n = range ers {
		if !Errors().ByDirectionField.Is(ers[n].Ei) {
			t.Errorf("ParseBy() = %q, ожидалось: %q", ers[n].Ei, Errors().ByDirectionField.Bind())
			return
		}
	}
}

// Тестирование экспорта сортировки.
func TestExportBy(t *testing.T) {
	var (
		obj *impl
		key string
		val []string
	)

	obj = New().(*impl)
	// Пустая сортировка.
	if key, val = obj.exportBy(); key != "" || len(val) > 0 {
		t.Errorf(
			"exportBy() вернул начения: %q = %v, ожидалось: %q, %v",
			key, val,
			"", []string{},
		)
		return
	}
	// Настроенная сортировка.
	obj.By = []Direction{
		{Field: "id", By: "asc"},
		{Field: "di", By: "desc"},
	}
	if key, val = obj.exportBy(); key != keyBy || !reflect.DeepEqual(val, []string{"id:asc", "di:desc"}) {
		t.Errorf(
			"exportBy() вернул начения: %q = %v, ожидалось: %q, %v",
			key, val,
			keyBy, []string{"id:asc", "di:desc"},
		)
		return
	}
}
