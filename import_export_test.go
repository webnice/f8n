package f8n

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// Тестовый запрос.
func testData1(t *testing.T) string {
	const data = `http://localhost
?map=f1:or:f2:or:(f5:and:f3)
&map=f3:and:f4
&f1=id:ge:10
&f2=key:eq:2
&f3=is_done:eq:true
&f4=devicedisplayName:ke:Солнечная система
&f5=approximately:ge:1.1
&limit=%d:%d`
	var (
		tmp string
		lim []uint64
		dir []Direction
		flt []Filter
		n   int
	)

	lim, dir, flt =
		makeTestContent(t, 0, []uint64{}),
		makeTestContent(t, 0, []Direction{}),
		makeTestContent(t, 0, []Filter{})
	tmp = strings.ReplaceAll(data, "\n", "")
	for n = range dir {
		tmp += fmt.Sprintf("&by=%s:%s", dir[n].Field, dir[n].By)
	}
	for n = range flt {
		tmp += fmt.Sprintf("&filter=%s:%s:%s", flt[n].Field, flt[n].Method.String(), flt[n].Value.String())
	}

	return fmt.Sprintf(tmp, lim[0], lim[1])
}

// Тестовый JSON соответствующий запросу.
func testData2(_ *testing.T) string {
	return "" +
		`{"field":["id","key","is_done","devicedisplayName","is_deleted","approximately"],"field_type":{"approximate` +
		`ly":"float64","id":"uint64","is_deleted":"bool","is_done":"bool","key":"int64"},"remap":{"is_deleted":"isDe` +
		`leted","is_done":"isDone"},"offset":101,"limit":3,"by":[{"field":"devicedisplayName","direction":"asc"},{"f` +
		`ield":"id","direction":"desc"}],"filter":[{"field":"is_deleted","method":"eq","value":{"source":"false","ty` +
		`pe":"bool"}}],"map":{"beg":"","end":"","origin":"","content":"","filter":null,"size":0,"node":[{"beg":"(","` +
		`end":")","origin":"operator_bracket","content":"","filter":null,"size":2,"node":[{"beg":"","end":"","origin` +
		`":"filtration","content":"f1","filter":{"field":"id","method":"ge","value":{"source":"10","type":"uint64"}}` +
		`,"size":2,"node":null},{"beg":":or:","end":":or:","origin":":or:","content":"","filter":null,"size":4,"node` +
		`":null},{"beg":"","end":"","origin":"filtration","content":"f2","filter":{"field":"key","method":"eq","valu` +
		`e":{"source":"2","type":"int64"}},"size":2,"node":null},{"beg":":or:","end":":or:","origin":":or:","content` +
		`":"","filter":null,"size":4,"node":null},{"beg":"(","end":")","origin":"operator_bracket","content":"","fil` +
		`ter":null,"size":2,"node":[{"beg":"","end":"","origin":"filtration","content":"f5","filter":{"field":"appro` +
		`ximately","method":"ge","value":{"source":"1.1","type":"float64"}},"size":2,"node":null},{"beg":":and:","en` +
		`d":":and:","origin":":and:","content":"","filter":null,"size":5,"node":null},{"beg":"","end":"","origin":"f` +
		`iltration","content":"f3","filter":{"field":"is_done","method":"eq","value":{"source":"true","type":"bool"}` +
		`},"size":2,"node":null}]}]},{"beg":":and:","end":":and:","origin":":and:","content":"","filter":null,"size"` +
		`:5,"node":null},{"beg":"(","end":")","origin":"operator_bracket","content":"","filter":null,"size":2,"node"` +
		`:[{"beg":"","end":"","origin":"filtration","content":"f3","filter":{"field":"is_done","method":"eq","value"` +
		`:{"source":"true","type":"bool"}},"size":2,"node":null},{"beg":":and:","end":":and:","origin":":and:","cont` +
		`ent":"","filter":null,"size":5,"node":null},{"beg":"","end":"","origin":"filtration","content":"f4","filter` +
		`":{"field":"devicedisplayName","method":"ke","value":{"source":"Солнечная система","type":"string"}},"size"` +
		`:2,"node":null}]}]}}`
}

// Тип данных для дженерика.
type testContentModel interface {
	string | []string | []uint64 | Remap | Datatype | []Direction | []Filter
}

// Дженерик, возвращает тестовые данные соответствующие запросу и JSON данным.
func makeTestContent[T testContentModel](t *testing.T, dataID int, data T) (ret T) {
	switch typ := any(data).(type) {
	case string:
		switch dataID {
		case 1:
			ret = any(testData1(t)).(T)
		case 2:
			ret = any(testData2(t)).(T)
		case 3:
			ret = any("(f1:or:f2:or:(f5:and:f3)):and:(f3:and:f4)").(T)
		default:
			t.Errorf("Не верный вызов функции makeTestContent[T], необходимо указать dataID > 0")
		}
	case []string:
		ret = any(append(typ, []string{
			"id", "key", "is_done", "devicedisplayName", "is_deleted", "approximately",
		}...)).(T)
	case []uint64:
		ret = any(append(typ, []uint64{101, 3}...)).(T)
	case Remap:
		typ["is_done"], typ["is_deleted"] = "isDone", "isDeleted"
		ret = any(typ).(T)
	case Datatype:
		typ["id"], typ["key"], typ["is_done"], typ["is_deleted"], typ["approximately"] =
			TypeUint64, TypeInt64, TypeBool, TypeBool, TypeFloat64
		ret = any(typ).(T)
	case []Direction:
		typ = append(typ, Direction{Field: "devicedisplayName", By: "asc"})
		typ = append(typ, Direction{Field: "id", By: "desc"})
		ret = any(typ).(T)
	case []Filter:
		typ = append(typ, Filter{Field: "is_deleted", Method: "eq", Value: FilterValue{Source: "false", Type: "bool"}})
		ret = any(typ).(T)
	default:
		t.Errorf("Не верный вызов функции makeTestContent[T], тип: %q", reflect.TypeOf(typ).String())
	}

	return
}

func makeTestObject(t *testing.T) (ret Interface, err error) {
	var rq *http.Request

	if rq, err = http.NewRequest("GET", makeTestContent(t, 1, ""), nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	ret = New().FieldSet(makeTestContent(t, 0, []string{})...)
	for key, val := range makeTestContent(t, 0, Datatype{}) {
		ret = ret.FieldDatatype(key, val)
	}
	for key, val := range makeTestContent(t, 0, Remap{}) {
		ret = ret.Redefinition(key, val)
	}
	err = ret.ParseRequest(rq, func(b []byte, e error) { t.Logf("ошибка: %s", string(b)) })

	return
}

// Тестирование экспорта разобранных значений фильтрации в формат JSON.
func TestExportJson(t *testing.T) {
	var (
		err error
		obj Interface
		exp string
	)

	if obj, err = makeTestObject(t); err != nil {
		t.Errorf("ParseRequest() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	// Выполнение экспорта.
	if exp = obj.ExportJson("", "").
		String(); !strings.EqualFold(strings.TrimSpace(exp), makeTestContent(t, 2, "")) {
		t.Logf("вернулось значение: %s", strings.TrimSpace(exp))
		t.Logf("ожидалось значение: %s", makeTestContent(t, 2, ""))
		t.Errorf("ExportJson() вернул не корректный результат.")
		return
	}
}

// Тестирование импорта JSON.
func TestImportJson(t *testing.T) {
	var (
		err error
		obj Interface
		f8n *impl
		lim []uint64
		cft *Map
	)

	if obj, err = makeTestObject(t); err != nil {
		t.Errorf("ParseRequest() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	cft = obj.(*impl).Map
	obj = New()
	if err = obj.ImportJson([]byte("}}{")); err == nil {
		t.Errorf("ImportJson() == nil, ожидались ошибка.")
		return
	}
	obj.Reset()
	if err = obj.ImportJson([]byte(makeTestContent(t, 2, ""))); err != nil {
		t.Errorf("ImportJson() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	// Проверка импортированных данных.
	f8n = obj.(*impl)
	// Сравнение списка обрабатываемых полей.
	if !reflect.DeepEqual(f8n.FieldOnly, makeTestContent(t, 0, []string{})) {
		t.Errorf(
			"ImportJson() -> FieldSet() не совпадает. Импортировано: [%v] ожидалось: [%v]",
			f8n.FieldOnly,
			makeTestContent(t, 0, []string{}),
		)
		return
	}
	// Сравнение типа данных для поля.
	if !reflect.DeepEqual(f8n.Datatype, makeTestContent(t, 0, Datatype{})) {
		t.Errorf(
			"ImportJson() -> FieldDatatype() не совпадает. Импортировано: [%v] ожидалось: [%v]",
			f8n.Datatype,
			makeTestContent(t, 0, Datatype{}),
		)
		return
	}
	// Сравнение переопределения названий полей или замены полей.
	if !reflect.DeepEqual(f8n.Remap, makeTestContent(t, 0, Remap{})) {
		t.Errorf(
			"ImportJson() -> Redefinition() не совпадает. Импортировано: [%v] ожидалось: [%v]",
			f8n.Remap,
			makeTestContent(t, 0, Remap{}),
		)
		return
	}
	// Сравнение лимитов.
	if lim = makeTestContent(t, 0, []uint64{}); f8n.Offset != lim[0] || f8n.Limit != lim[1] {
		t.Errorf(
			"ImportJson() -> Limit[] не совпадает. Импортировано: [%v] ожидалось: [%v]",
			[]uint64{f8n.Offset, f8n.Limit},
			lim,
		)
		return
	}
	// Сравнение сортировки результата выборки.
	if !reflect.DeepEqual(f8n.By, makeTestContent(t, 0, []Direction{})) {
		t.Errorf(
			"ImportJson() -> By[] не совпадает. Импортировано: [%v] ожидалось: [%v]",
			f8n.By,
			makeTestContent(t, 0, []Direction{}),
		)
		return
	}
	// Сравнение простой фильтрации.
	if !reflect.DeepEqual(f8n.Filter, makeTestContent(t, 0, []Filter{})) {
		t.Errorf(
			"ImportJson() -> Filter[] не совпадает. Импортировано: [%v] ожидалось: [%v]",
			f8n.Filter,
			makeTestContent(t, 0, []Filter{}),
		)
		return
	}
	// Сравнение сложной фильтрации.
	cft.setParent(f8n) // Надо поменять для сравнения, так как объект другой.
	if !reflect.DeepEqual(f8n.Map, cft) {
		t.Errorf(
			"ImportJson() -> Map[] не совпадает. Импортировано: [%v] ожидалось: [%v]",
			f8n.Map,
			cft,
		)
		return
	}
}

// Тестирование экспорта настроек фильтрации в объект http.Request.Query().Values, тип данных url.Values.
func TestQuery(t *testing.T) {
	var (
		err error
		obj Interface
		val url.Values
		tmp string
	)

	if obj, err = makeTestObject(t); err != nil {
		t.Errorf("ParseRequest() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	val = obj.Query()
	// Проверка значений получившегося экспорта.
	if !val.Has(keyMap) {
		t.Errorf("Query() map == \"\". Ошибка значения экспорта.")
		return
	}
	if tmp = val.Get(keyMap); tmp != makeTestContent(t, 3, "") {
		t.Errorf("Query() map == %q. Ожидалось значение map = %q", tmp, makeTestContent(t, 3, ""))
		return
	}
}

// Тестирование экспорта настроек фильтрации в строку.
func TestString(t *testing.T) {
	const localhost = `http://localhost?`
	var (
		err error
		obj Interface
		o1  *impl
		o2  *impl
		rq  *http.Request
		tmp string
	)

	if obj, err = makeTestObject(t); err != nil {
		t.Errorf("ParseRequest() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	if tmp = obj.String(); tmp == "" {
		t.Errorf("String() == %q, ожидался не пустой результат.", tmp)
		return
	}
	o1 = obj.(*impl)
	// Разбор полученного экспорта.
	if rq, err = http.NewRequest("GET", localhost+tmp, nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().FieldSet(makeTestContent(t, 0, []string{})...)
	for key, val := range makeTestContent(t, 0, Datatype{}) {
		obj = obj.FieldDatatype(key, val)
	}
	for key, val := range makeTestContent(t, 0, Remap{}) {
		obj = obj.Redefinition(key, val)
	}
	if err = obj.ParseRequest(rq, func(b []byte, e error) { t.Logf("ошибка: %s", string(b)) }); err != nil {
		t.Errorf("ParseRequest() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	o2 = obj.(*impl)
	// Если экспорт и импорт работают корректно, тогда эти объекты должны совпасть.
	if !reflect.DeepEqual(o1, o2) {
		t.Errorf("Ошибка функции String(). Получен объект: %v, ожидался: %v\n",
			strings.TrimSpace(o2.ExportJson("", "").String()),
			strings.TrimSpace(o1.ExportJson("", "").String()),
		)
		return
	}
}

// Тестирование импорта строки запроса.
func TestImportString(t *testing.T) {
	var (
		err error
		obj Interface
		o1  *impl
		o2  *impl
		tmp string
	)

	if obj, err = makeTestObject(t); err != nil {
		t.Errorf("ParseRequest() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	if tmp = obj.String(); tmp == "" {
		t.Errorf("String() == %q, ожидался не пустой результат.", tmp)
		return
	}
	o1 = obj.(*impl)
	// Разбор полученного экспорта c помощью функции ImportString.
	obj = New().FieldSet(makeTestContent(t, 0, []string{})...)
	for key, val := range makeTestContent(t, 0, Datatype{}) {
		obj = obj.FieldDatatype(key, val)
	}
	for key, val := range makeTestContent(t, 0, Remap{}) {
		obj = obj.Redefinition(key, val)
	}
	if err = obj.ImportString(tmp); err != nil {
		t.Errorf("ImportString() != nil, ошибка не ожидались. Ошибка: %s", err)
		return
	}
	o2 = obj.(*impl)
	// Если экспорт и импорт работают корректно, тогда эти объекты должны совпасть.
	if !reflect.DeepEqual(o1, o2) {
		t.Errorf("Ошибка функции ImportString(). Получен объект: %v, ожидался: %v\n",
			strings.TrimSpace(o2.ExportJson("", "").String()),
			strings.TrimSpace(o1.ExportJson("", "").String()),
		)
		return
	}
}

// Тестирование импорта строки запроса содержащей ошибку.
func TestImportStringWithError(t *testing.T) {
	var (
		err error
		obj Interface
	)

	obj = New()
	if err = obj.ImportString(string([]byte{0, 1, 2, 3, 4, 5})); err == nil {
		t.Errorf("ImportString() == nil, ожидалась ошибка.")
		return
	}
}
