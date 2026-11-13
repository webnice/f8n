package f8n

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// Тестирование загрузки пустой карты условий сложной фильтрации.
func TestParseMapZero(t *testing.T) {
	const (
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorMap     = "ParseMap() != nil, ошибка не ожидалась."
		tplErrorNotNil  = "Map != nil, ожидался пустой узел карты."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	obj = New().(*impl)
	// Отсутствие карты.
	if rq, err = http.NewRequest("GET", "http://localhost?f1=id:eq:1", nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorMap)
		return
	}
	if obj.Map != nil {
		t.Errorf(tplErrorNotNil)
		return
	}
	obj.ResetMap()
	// Пустая карта.
	if rq, err = http.NewRequest("GET", "http://localhost?map=&f1=id:eq:1", nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorMap)
		return
	}
	if obj.Map != nil {
		t.Errorf(tplErrorNotNil)
		return
	}
}

// Тестирование загрузки и разбора простейшей карты условий сложной фильтрации.
func TestParseMapSimplified(t *testing.T) {
	const (
		fg1             = "filter_group_one"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorMap     = "ParseMap() != nil, ошибка не ожидалась."
		tplErrorNil     = "Map == nil, ожидался не пустой узел карты."
		tplErrorOrigin  = "Map.Origin != OriginFiltration, ожидался узел фильтрации."
		tplErrorContent = "Map.Content != %q, ожидалось значение узла фильтрации %q."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", "http://localhost?&map="+fg1+"&"+fg1+"=id:eq:1", nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorMap)
		return
	}
	if obj.Map == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if obj.Map.Origin != OriginFiltration {
		t.Errorf(tplErrorOrigin)
		return
	}
	if obj.Map.Content != fg1 {
		t.Errorf(tplErrorContent, obj.Map.Content, fg1)
		return
	}
}

// Тестирование загрузки и разбора простой одноуровневой карты условий сложной фильтрации.
func TestParseMapSimple(t *testing.T) {
	const (
		fg1, fg2        = "group1", "group2"
		vMap            = "map=" + fg1 + ":or:" + fg2 + ":or:" + fg1 + ":and:" + fg2
		filters         = "&" + fg1 + "=id:eq:1&" + fg2 + "=key:eq:1"
		size            = 37
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorMap     = "ParseMap() != nil, ошибка не ожидалась."
		tplErrorNil     = "Map == nil, ожидался не пустой узел карты."
		tplErrorOrigin  = "Map.Origin != OriginUnknown, ожидался пустой корневой узел содержащий узлы."
		tplErrorSize    = "Map.SumSize == %d, ожидался суммарный размер карты узлов равный %d"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", "http://localhost?"+vMap+filters, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorMap)
		return
	}
	if obj.Map == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if !obj.Map.IsEmpty() && obj.Map.Origin != OriginUnknown {
		t.Errorf(tplErrorOrigin)
		return
	}
	if obj.Map.SumSize() != size {
		t.Errorf(tplErrorSize, obj.Map.SumSize(), size)
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
	const (
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorEmpty   = "ParseMap() == nil, ожидались ошибки."
		tplErrorIs      = "ParseMap() = %q, ожидалось: %q"
		tplErrorPanic   = "ParseRequest() == nil, ожидалась ошибка."
		tplErrorPanicIs = "ParseRequest() = %q, ожидалось: %q"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", strings.ReplaceAll(data, "\n", ""), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) <= 0 {
		t.Errorf(tplErrorEmpty)
		return
	}
	for n = range ers {
		if !Errors().PairedTagNotMatch.Is(ers[n].Ei) {
			t.Errorf(tplErrorIs, ers[n].Ei, Errors().PairedTagNotMatch.Bind())
			return
		}
	}
	// Проверка ошибки паники.
	obj = New().(*impl)
	if err = obj.ParseRequest(rq, func(a []byte, b error) {
		panic(string(a) + b.Error())
	}); err == nil {
		t.Errorf(tplErrorPanic)
		return
	}
	if !Errors().Panic.Is(err) {
		t.Errorf(tplErrorPanicIs, err, Errors().Panic.Bind())
		return
	}
}

// Тестирование экспорта карты условий сложной фильтрации.
func TestExportMap(t *testing.T) {
	const (
		datUrn          = "http://localhost?map=f1:or:f2&f1=id:eq:1&f2=key:eq:2"
		datMap          = "f1:or:f2"
		datF1           = "id:eq:1"
		datF2           = "key:eq:2"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseRequest() = %v, ошибка не ожидалась."
		tplErrorEmpty   = "exportMap() вернул значения: %q = %v, %v, ожидалось: %q, %v, %v"
		tplErrorMap     = "exportMap() вернул значения: %q = %q, ожидалось: %q, %q"
		tplErrorFilters = "exportMap() вернул фильтры: %v, ожидалось: %v"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		key string
		val string
		flt url.Values
	)

	obj = New().(*impl)
	// Пустая фильтрация.
	if key, val, flt = obj.exportMap(); key != "" || val != "" || len(flt) > 0 {
		t.Errorf(
			tplErrorEmpty,
			key, val, flt,
			"", "", flt,
		)
		return
	}
	// Настроенная фильтрация.
	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if err = obj.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	if key, val, flt = obj.exportMap(); key != keyMap || val != datMap {
		t.Errorf(tplErrorMap, key, val, keyMap, datMap)
		return
	}
	if flt.Get("f1") != datF1 || flt.Get("f2") != datF2 {
		t.Errorf(
			tplErrorFilters,
			flt, url.Values{"f1": {datF1}, "f2": {datF2}},
		)
		return
	}
}
