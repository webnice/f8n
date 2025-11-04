package f8n

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// Тестирование загрузки пустой карты условий сложной фильтрации.
func TestParseMapZero(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	obj = New().(*impl)
	// Отсутствие карты.
	if rq, err = http.NewRequest("GET", "http://localhost?f1=id:eq:1", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf("ParseMap() != nil, ошибка не ожидалась.")
		return
	}
	if obj.Map != nil {
		t.Errorf("Map != nil, ожидался пустой узел карты.")
		return
	}
	obj.ResetMap()
	// Пустая карта.
	if rq, err = http.NewRequest("GET", "http://localhost?map=&f1=id:eq:1", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf("ParseMap() != nil, ошибка не ожидалась.")
		return
	}
	if obj.Map != nil {
		t.Errorf("Map != nil, ожидался пустой узел карты.")
		return
	}
}

// Тестирование загрузки и разбора простейшей карты условий сложной фильтрации.
func TestParseMapSimplified(t *testing.T) {
	const fg1 = "filter_group_one"
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", "http://localhost?&map="+fg1+"&"+fg1+"=id:eq:1", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf("ParseMap() != nil, ошибка не ожидалась.")
		return
	}
	if obj.Map == nil {
		t.Errorf("Map == nil, ожидался не пустой узел карты.")
		return
	}
	if obj.Map.Origin != OriginFiltration {
		t.Errorf("Map.Origin != OriginFiltration, ожидался узел фильтрации.")
		return
	}
	if obj.Map.Content != fg1 {
		t.Errorf("Map.Content != %q, ожидалось значение узла фильтрации %q.", obj.Map.Content, fg1)
		return
	}
}

// Тестирование загрузки и разбора простой одноуровневой карты условий сложной фильтрации.
func TestParseMapSimple(t *testing.T) {
	const (
		fg1, fg2 = "group1", "group2"
		vMap     = "map=" + fg1 + ":or:" + fg2 + ":or:" + fg1 + ":and:" + fg2
		filters  = "&" + fg1 + "=id:eq:1&" + fg2 + "=key:eq:1"
		size     = 37
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", "http://localhost?"+vMap+filters, nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf("ParseMap() != nil, ошибка не ожидалась.")
		return
	}
	if obj.Map == nil {
		t.Errorf("Map == nil, ожидался не пустой узел карты.")
		return
	}
	if !obj.Map.IsEmpty() && obj.Map.Origin != OriginUnknown {
		t.Errorf("Map.Origin != OriginUnknown, ожидался пустой корневой узел содержащий узлы.")
		return
	}
	if obj.Map.SumSize() != size {
		t.Errorf("Map.SumSize == %d, ожидался суммарный размер карты узлов равный %d", obj.Map.SumSize(), size)
		return
	}
}

// Тестирование загрузки и разбора карты условий сложной фильтрации содержащей ошибку в операторных скобках.
func TestParseMapWithError(t *testing.T) {
	const data = `http://localhost
?map=(group1:and:group2:or:group3):or:(group4:or:group5):or:(group4:or:group5)
&map=group4:and:group5
&map=(bbb:or:ccc):and:bbb)
&group1=field1:eq:value1
&group2=field2:ke:value2
&group3=field3:ke:value3
&group4=field1:ke:value4
&group5=field2:ke:value5
`
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", strings.ReplaceAll(data, "\n", ""), nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) <= 0 {
		t.Errorf("ParseMap() == nil, ожидались ошибки.")
		return
	}
	for n = range ers {
		if !Errors().PairedTagNotMatch.Is(ers[n].Ei) {
			t.Errorf("ParseMap() = %q, ожидалось: %q", ers[n].Ei, Errors().PairedTagNotMatch.Bind())
			return
		}
	}
	// Проверка ошибки паники.
	obj = New().(*impl)
	if err = obj.ParseRequest(rq, func(a []byte, b error) {
		panic(string(a) + b.Error())
	}); err == nil {
		t.Errorf("ParseRequest() == nil, ожидалась ошибка.")
		return
	}
	if !Errors().Panic.Is(err) {
		t.Errorf("ParseRequest() = %q, ожидалось: %q", err, Errors().Panic.Bind())
		return
	}
}

// Тестирование экспорта карты условий сложной фильтрации.
func TestExportMap(t *testing.T) {
	var (
		obj *impl
		key string
		val string
		flt url.Values
	)

	obj = New().(*impl)
	// Пустая фильтрация.
	if key, val, flt = obj.exportMap(); key != "" || val != "" || len(flt) > 0 {
		t.Errorf(
			"exportMap() вернул начения: %q = %v, %v, ожидалось: %q, %v, %v",
			key, val, flt,
			"", "", flt,
		)
		return
	}
	// Настроенная фильтрация.

	// TODO: Сделать тест проверки экспорта сложной фильтрации.

}
