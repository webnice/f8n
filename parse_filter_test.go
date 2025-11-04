package f8n

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// Тестирование простой фильтрации - фильтрация отсутствует.
func TestParseFilterEmpty(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	// Пустой срез фильтрации.
	if rq, err = http.NewRequest("GET", "http://localhost", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseFilter(rq); len(ers) > 0 {
		t.Errorf("ParseFilter() = %v, ошибка не ожидалась.", ers)
		return
	}
	if len(obj.Filter) > 0 {
		t.Errorf("Filter[] = %v, ожидался пустой срез.", obj.Filter)
		return
	}
}

// Тестирование простой фильтрации - ошибка формата.
//
//goland:noinspection ALL
func TestParseFilterFormat(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	// Не верный формат значения сортировки.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", "http://localhost?filter=id-asc:s&filter=a-b-c", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseFilter(rq); len(ers) == 0 {
		t.Errorf("ParseFilter() = nil, ожидалась ошибка.")
		return
	}
	for n = range ers {
		if !Errors().FilterFormat.Is(ers[n].Ei) {
			t.Errorf("ParseFilter() = %q, ожидалось: %q", ers[n].Ei, Errors().FilterFormat.Bind())
			return
		}
	}
}

// Тестирование простой фильтрации - неизвестный метод.
//
//goland:noinspection ALL
func TestParseFilterUnknownMethod(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	// Не верный формат значения сортировки.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", "http://localhost?filter=id:asc:value", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseFilter(rq); len(ers) == 0 {
		t.Errorf("ParseFilter() = nil, ожидалась ошибка.")
		return
	}
	for n = range ers {
		if !Errors().FilterUnknownMethod.Is(ers[n].Ei) {
			t.Errorf("ParseFilter() = %q, ожидалось: %q", ers[n].Ei, Errors().FilterUnknownMethod.Bind())
			return
		}
	}
}

// Тестирование простой фильтрации - строгие и не строгие поля.
func TestParseFilterField(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	// Не строгие поля.
	obj = New().(*impl)
	if rq, err = http.NewRequest(
		"GET",
		"http://localhost?filter=aaa:eq:asc&filter=bbb:in:desc",
		nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseFilter(rq); len(ers) > 0 {
		t.Errorf("ParseFilter() = %v, ошибка не ожидалась.", ers)
		return
	}
	// Строгие поля.
	obj = New().FieldSet("id", "name").(*impl)
	if rq, err = http.NewRequest(
		"GET",
		"http://localhost?filter=aaa:eq:asc&filter=bbb:in:desc&filter=id:in:1&filter=name:eq:desc",
		nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseFilter(rq); len(ers) != 2 {
		t.Errorf("ParseFilter() = nil, ожидалось две ошибки.")
		return
	}
	for n = range ers {
		if !Errors().FilterUnknownField.Is(ers[n].Ei) {
			t.Errorf("ParseFilter() = %q, ожидалось: %q", ers[n].Ei, Errors().FilterUnknownField.Bind())
			return
		}
	}
}

// Тестирование простой фильтрации - проверка значений и типов значений.
func TestParseFilterValueAndType(t *testing.T) {
	const expectedType = 5
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
		nOk uint64
	)

	// Не строгие поля.
	obj = New().
		FieldDatatype("uint64", TypeUint64).
		FieldDatatype("int64", TypeInt64).
		FieldDatatype("float64", TypeFloat64).
		FieldDatatype("bool", TypeBool).
		FieldDatatype("time", TypeTime).(*impl)
	if rq, err = http.NewRequest(
		"GET",
		"http://localhost?filter=uint64:le:11&filter=int64:ge:-32&filter=float64:eq:42.3&filter=bool:eq:true"+
			"&filter=time:ne:"+url.QueryEscape("2018-08-03T03:00:00+03:00"),
		nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseFilter(rq); len(ers) > 0 {
		t.Errorf("ParseFilter() = %v, ошибка не ожидалась.", ers)
		return
	}
	for n = range obj.Filter {
		switch {
		case obj.Filter[n].Field == "uint64" &&
			obj.Filter[n].Value.Source == "11" &&
			obj.Filter[n].Value.Type == TypeUint64:
			nOk++
		case obj.Filter[n].Field == "int64" &&
			obj.Filter[n].Value.Source == "-32" &&
			obj.Filter[n].Value.Type == TypeInt64:
			nOk++
		case obj.Filter[n].Field == "float64" &&
			obj.Filter[n].Value.Source == "42.3" &&
			obj.Filter[n].Value.Type == TypeFloat64:
			nOk++
		case obj.Filter[n].Field == "bool" &&
			obj.Filter[n].Value.Source == "true" &&
			obj.Filter[n].Value.Type == TypeBool:
			nOk++
		case obj.Filter[n].Field == "time" &&
			obj.Filter[n].Value.Source == "2018-08-03T03:00:00+03:00" &&
			obj.Filter[n].Value.Type == TypeTime:
			nOk++
		}
	}
	if nOk != expectedType {
		t.Errorf("Field = %d, ожидалось %d типов.", nOk, expectedType)
		return
	}
	if rq, err = http.NewRequest(
		"GET",
		"http://localhost?filter=uint64:le:-1&filter=int64:ge:-o2&filter=float64:eq:42,3&filter=bool:eq:no"+
			"&filter=time:ne:"+url.QueryEscape("2018-08-03T03:00:00"),
		nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseFilter(rq); len(ers) != expectedType {
		t.Errorf("ParseFilter() = %d ошибок, ожидалось %d ошибок.", len(ers), expectedType)
		return
	}
}

// Тестирование простой фильтрации - проверка известных методов фильтрации.
func TestParseFilterMethod(t *testing.T) {
	var (
		err     error
		obj     *impl
		rq      *http.Request
		methods []FilterMethod
		n       int
		ers     []*ParseError
	)

	methods = []FilterMethod{filterEquivalent, filterLessThan, filterLessThanOrEquivalent, filterGreaterThan,
		filterGreaterThanOrEquivalent, filterNotEquivalent, filterLikeThan, filterNotLikeThan, filterIn,
		filterNotIn}
	for n = range methods {
		obj = New().FieldDatatype("a", TypeUint64).(*impl)
		if rq, err = http.NewRequest(
			"GET",
			"http://localhost?filter=a:"+methods[n].String()+":1234",
			nil); err != nil {
			t.Errorf("Ошибка создания запроса: %s", err)
		}
		if ers = obj.ParseFilter(rq); len(ers) > 0 {
			t.Errorf("ParseFilter() = %v, ошибка не ожидалась.", ers)
			return
		}
		if len(obj.Filter) != 1 {
			t.Errorf("[]Filter = %d, ожидалось %d.", len(obj.Filter), 1)
			return
		}
		if obj.Filter[0].Method != methods[n] {
			t.Errorf("[]Filter.Method = %q, ожидалось %q.", obj.Filter[0].Method.String(), methods[n].String())
			continue
		}
	}
}

// Тестирование простой фильтрации.
func TestExportFilter(t *testing.T) {
	var (
		obj *impl
		key string
		val []string
	)

	obj = New().(*impl)
	// Пустая фильтрация.
	if key, val = obj.exportFilter(); key != "" || len(val) > 0 {
		t.Errorf(
			"exportFilter() вернул начения: %q = %v, ожидалось: %q, %v",
			key, val,
			"", []string{},
		)
		return
	}
	// Настроенная фильтрация.
	obj.Filter = []Filter{
		{Field: "id", Method: "ge", Value: FilterValue{Source: "7", Type: "uint64"}},
		{Field: "di", Method: "ne", Value: FilterValue{Source: "1", Type: "int64"}},
	}
	if key, val = obj.exportFilter(); key != keyFilter || !reflect.DeepEqual(val, []string{"id:ge:7", "di:ne:1"}) {
		t.Errorf(
			"exportFilter() вернул начения: %q = %v, ожидалось: %q, %v",
			key, val,
			keyFilter, []string{"id:ge:7", "di:ne:1"},
		)
		return
	}
}
