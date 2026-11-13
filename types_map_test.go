package f8n

import (
	"net/http"
	"strings"
	"testing"
)

// Тестирование отсутствия именованных фильтров.
func TestMapLoadAndParseFilterCalledByNameWasNotFound(t *testing.T) {
	const (
		datRequest = `http://localhost
?map=filter1:or:filter2
&filter1=id:eq:1
`
		tplErrorRequest  = "Ошибка создания запроса: %s."
		tplErrorEmpty    = "ParseMap() == nil, ожидались ошибки."
		tplErrorNotFound = "ParseMap() = %q, ожидалось: %q"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", strings.ReplaceAll(datRequest, "\n", ""), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) <= 0 {
		t.Errorf(tplErrorEmpty)
		return
	}
	for n = range ers {
		if !Errors().FilterCalledByNameWasNotFound.Is(ers[n].Ei) {
			t.Errorf(tplErrorNotFound, ers[n].Ei, Errors().FilterCalledByNameWasNotFound.Bind())
			return
		}
	}
}

// Тестирование экспорта пустой простой фильтрации.
func TestExportAsStringEmptyFilter(t *testing.T) {
	const tplErrorExport = "Filter = %q, ожидались пустая простая фильтрация."
	var (
		obj Filter
		tmp string
	)

	if tmp = obj.exportAsString(); tmp != "" {
		t.Errorf(tplErrorExport, tmp)
		return
	}
}

// Тестирование экспорта пустой сортировки.
func TestExportAsStringEmptyDirection(t *testing.T) {
	const tplErrorExport = "Direction = %q, ожидались пустая сортировка."
	var (
		obj Direction
		tmp string
	)

	if tmp = obj.exportAsString(); tmp != "" {
		t.Errorf(tplErrorExport, tmp)
		return
	}
}

// Тестирование экспорта пустой карты сложной фильтрации.
func TestExportAsStringEmptyMap(t *testing.T) {
	const tplErrorExport = "Map = %q, ожидались пустая карта сложной фильтрации."
	var (
		obj *Map
		tmp string
	)

	obj = new(Map)
	if tmp, _ = obj.exportAsString(); tmp != "" {
		t.Errorf(tplErrorExport, tmp)
		return
	}
}

// Тестирование формирования SQL запроса и параметров из DOM дерева карты сложной фильтрации.
//
//goland:noinspection ALL
func TestMapQueryGorm(t *testing.T) {
	const (
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseMap() = %v, ошибка не ожидалась."
		tplErrorNil     = "Map == nil, ожидалась не пустая карта."
		tplErrorQuery   = "%s: queryGorm() = %q, ожидалось: %q."
		tplErrorValues  = "%s: queryGorm() значения = %v, ожидалось: %v."
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
			Remap  Remap
			Query  string
			Values []any
		}
	)
	var (
		err    error
		rq     *http.Request
		obj    *impl
		ers    []*ParseError
		n, j   int
		query  string
		values []any
		tests  = []gs{
			{
				Name:   "single-filtration",
				Urn:    "http://localhost?map=f1&f1=id:eq:1",
				Fields: []ft{{Field: "id", Type: TypeUint64}},
				Query:  "`id` = ?",
				Values: []any{uint64(1)},
			},
			{
				Name:   "or-chain",
				Urn:    "http://localhost?map=f1:or:f2&f1=id:eq:1&f2=key:eq:2",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "`id` = ? OR `key` = ?",
				Values: []any{uint64(1), int64(2)},
			},
			{
				Name:   "bracket-and-or",
				Urn:    "http://localhost?map=(f1:and:f2):or:f3&f1=id:eq:1&f2=key:eq:2&f3=name:eq:вася",
				Fields: []ft{{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64}},
				Query:  "(`id` = ? AND `key` = ?) OR `name` = ?",
				Values: []any{uint64(1), int64(2), "вася"},
			},
			{
				Name: "nested-brackets",
				Urn: "http://localhost?map=((f1:or:f2):and:(f3:or:f4))" +
					"&f1=id:eq:1&f2=key:eq:2&f3=name:eq:вася&f4=age:ge:18",
				Fields: []ft{
					{Field: "id", Type: TypeUint64}, {Field: "key", Type: TypeInt64},
					{Field: "age", Type: TypeUint64},
				},
				Query:  "((`id` = ? OR `key` = ?) AND (`name` = ? OR `age` >= ?))",
				Values: []any{uint64(1), int64(2), "вася", uint64(18)},
			},
			{
				Name:   "like",
				Urn:    "http://localhost?map=f1&f1=section:ke:П?рв*й",
				Query:  "`section` LIKE ?",
				Values: []any{"П_рв%й"},
			},
			{
				Name:   "in",
				Urn:    "http://localhost?map=f1&f1=id:in:1,3,5",
				Fields: []ft{{Field: "id", Type: TypeSliceUint64}},
				Query:  "`id` IN (?)",
				Values: []any{[]uint64{1, 3, 5}},
			},
			{
				Name:   "nil-value",
				Urn:    "http://localhost?map=f1:and:f2&f1=deletedAt:eq:nil&f2=key:eq:2",
				Fields: []ft{{Field: "deletedAt", Type: TypeTimeNil}, {Field: "key", Type: TypeInt64}},
				Query:  "`deletedAt` IS NULL AND `key` = ?",
				Values: []any{int64(2)},
			},
			{
				Name:   "redefinition",
				Urn:    "http://localhost?map=f1&f1=id:eq:1",
				Fields: []ft{{Field: "id", Type: TypeUint64}},
				Remap:  Remap{"id": "`cursor`.`id`"},
				Query:  "`cursor`.`id` = ?",
				Values: []any{uint64(1)},
			},
		}
	)

	for n = range tests {
		if !t.Run(tests[n].Name, func(t *testing.T) {
			// Подготовка объекта фильтрации.
			obj = New().(*impl)
			for j = range tests[n].Fields {
				obj.FieldDatatype(tests[n].Fields[j].Field, tests[n].Fields[j].Type)
			}
			for key, val := range tests[n].Remap {
				obj.Redefinition(key, val)
			}
			// Разбор запроса.
			if rq, err = http.NewRequest("GET", tests[n].Urn, nil); err != nil {
				t.Errorf(tplErrorRequest, err)
			}
			if ers = obj.ParseMap(rq); len(ers) > 0 {
				t.Errorf(tplErrorParse, ers)
				return
			}
			if obj.Map == nil {
				t.Errorf(tplErrorNil)
				return
			}
			// Проверка сформированного запроса и параметров.
			query, values = obj.Map.queryGorm()
			if query != tests[n].Query {
				t.Errorf(tplErrorQuery, tests[n].Name, query, tests[n].Query)
				return
			}
			if !equalAnyValues(values, tests[n].Values) {
				t.Errorf(tplErrorValues, tests[n].Name, values, tests[n].Values)
			}
		}) {
			return
		}
	}
}

// Тестирование пропуска узлов карты сложной фильтрации не содержащих фильтров.
func TestMapQueryGormEmptyNodes(t *testing.T) {
	const (
		tplErrorQuery = "queryGorm() = %q, ожидалась пустая строка."
		tplErrorWant  = "queryGorm() = %q, ожидалось: %q."
	)
	var (
		obj   *impl
		mp    *Map
		query string
	)

	obj = New().(*impl)
	// Узел фильтрации без разобранного фильтра пропускается.
	mp = &Map{
		parent: obj,
		Origin: OriginUnknown,
		Node: []*Map{
			{parent: obj, Origin: OriginFiltration, Content: "f1"},
		},
	}
	if query, _ = mp.queryGorm(); query != "" {
		t.Errorf(tplErrorQuery, query)
		return
	}
	// Операторная скобка, группа условий которой пуста, пропускается.
	mp = &Map{
		parent: obj,
		Origin: OriginUnknown,
		Node: []*Map{
			{
				parent: obj,
				Origin: OriginOperatorBracket,
				Node: []*Map{
					{parent: obj, Origin: OriginFiltration, Content: "f1"},
				},
			},
		},
	}
	if query, _ = mp.queryGorm(); query != "" {
		t.Errorf(tplErrorQuery, query)
		return
	}
	// Пропущенный операнд между логическими операциями не оставляет двойных разделителей.
	mp = &Map{
		parent: obj,
		Origin: OriginUnknown,
		Node: []*Map{
			{parent: obj, Origin: OriginFiltration, Content: "f1"},
			{parent: obj, Origin: OriginAnd, TagBegin: keyAnd, TagEnd: keyAnd},
			{parent: obj, Origin: OriginFiltration, Content: "f2"},
			{parent: obj, Origin: OriginOr, TagBegin: keyOr, TagEnd: keyOr},
			{parent: obj, Origin: OriginFiltration, Content: "f3", Filter: &Filter{
				Field: "id", Method: filterEquivalent, Value: FilterValue{Source: "1", Type: TypeUint64},
			}},
		},
	}
	if query, _ = mp.queryGorm(); query != "`id` = ?" {
		t.Errorf(tplErrorWant, query, "`id` = ?")
		return
	}
}

// Тестирование повторной передачи значения именованного фильтра карты сложной фильтрации.
func TestMapLoadAndParseFilterRepeatedValue(t *testing.T) {
	const (
		datUrn          = `http://localhost?map=f1&f1=id:eq:1&f1=key:eq:2`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseMap() = %v, ошибка не ожидалась."
		tplErrorFirst   = "Filter.Field = %q, ожидалось первое значение %q."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorParse, ers)
		return
	}
	// При повторной передаче фильтра с одним именем используется первое значение.
	if obj.Map == nil || obj.Map.Filter == nil {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if obj.Map.Filter.Field != "id" || obj.Map.Filter.Value.Source != "1" {
		t.Errorf(tplErrorFirst, obj.Map.Filter.Field, "id")
		return
	}
}

// Тестирование объединения нескольких параметров map по принципу логического «ИЛИ».
func TestParseMapMultipleOr(t *testing.T) {
	const (
		datUrn          = `http://localhost?map=f1:and:f2&map=f3:and:f4&f1=id:eq:1&f2=key:eq:2&f3=u:eq:3&f4=v:eq:4`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseMap() = %v, ошибка не ожидалась."
		tplErrorExport  = "Map = %q, ожидалось объединение параметров map по принципу «ИЛИ»: %q."
		datWant         = "(f1:and:f2):or:(f3:and:f4)"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		tmp string
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if obj.Map == nil {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if tmp, _ = obj.Map.exportAsString(); tmp != datWant {
		t.Errorf(tplErrorExport, tmp, datWant)
		return
	}
}

// Тестирование объединения одиночных параметров map по принципу логического «ИЛИ».
func TestParseMapMultipleOrSingle(t *testing.T) {
	const (
		datUrn          = `http://localhost?map=f1&map=f2&f1=id:eq:1&f2=key:eq:2`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseMap() = %v, ошибка не ожидалась."
		tplErrorExport  = "Map = %q, ожидалось объединение параметров map по принципу «ИЛИ»: %q."
		datWant         = "f1:or:f2"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		tmp string
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if obj.Map == nil {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if tmp, _ = obj.Map.exportAsString(); tmp != datWant {
		t.Errorf(tplErrorExport, tmp, datWant)
		return
	}
}

// Тестирование передачи только пустых значений карты сложной фильтрации.
func TestParseMapEmptyValues(t *testing.T) {
	const (
		datUrn          = `http://localhost?map=&map=&f1=id:eq:1`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseMap() = %v, ошибка не ожидалась."
		tplErrorNotNil  = "Map != nil, ожидалось отсутствие карты сложной фильтрации."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if obj.Map != nil {
		t.Errorf(tplErrorNotNil)
		return
	}
}
