// Package f8n
package f8n

import (
	"math"
	"net/http"
	"strings"
	"testing"
)

// Тестирование сброса лимита.
func TestResetLimit(t *testing.T) {
	var (
		obj *impl
		ret *impl
	)

	obj = New().(*impl)
	if obj.Offset != 0 || obj.Limit != 0 {
		t.Errorf("Offset или Limit имеют не верные первоначальные значения.")
		return
	}
	ret = obj.ResetLimit().(*impl)
	if obj != ret {
		t.Errorf("ResetLimit() != object, ожидалось то что функция вернёт объект.")
		return
	}
	obj.Offset, obj.Limit = math.MaxUint64, math.MaxUint64
	obj.ResetLimit()
	if obj.Offset != 0 || obj.Limit != 0 {
		t.Errorf("Функция ResetLimit() не выполнила сброс лимитов.")
		return
	}
}

// Тестирование сброса сортировки.
func TestResetBy(t *testing.T) {
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
		ers []*ParseError
	)

	obj = New().(*impl)
	if len(obj.By) != 0 || cap(obj.By) != 0 {
		t.Errorf("By[] имеет не верные первоначальные значения.")
		return
	}
	ret = obj.ResetBy().(*impl)
	if obj != ret {
		t.Errorf("ResetBy() != object, ожидалось то что функция вернёт объект.")
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
	if len(obj.By) != 2 || cap(obj.By) != 2 {
		t.Errorf(
			"len(By[]) = %d, cap(By[]) = %d ожидался срез из двух значений с размером выделенной памяти два значения.",
			len(obj.By), cap(obj.By),
		)
		return
	}
	obj.ResetBy()
	if len(obj.By) != 0 || cap(obj.By) != 2 {
		t.Errorf(
			"len(By[]) = %d, cap(By[]) = %d ожидался пустой срез с размером выделенной памяти два значения.",
			len(obj.By), cap(obj.By),
		)
		return
	}
}

// Тестирование сброса настроек фильтрации.
func TestResetFilter(t *testing.T) {
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
		ers []*ParseError
	)

	obj = New().(*impl)
	if len(obj.Filter) != 0 || cap(obj.Filter) != 0 {
		t.Errorf("Filter[] имеет не верные первоначальные значения.")
		return
	}
	ret = obj.ResetFilter().(*impl)
	if obj != ret {
		t.Errorf("ResetFilter() != object, ожидалось то что функция вернёт объект.")
		return
	}
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
	if len(obj.Filter) != 2 || cap(obj.Filter) != 2 {
		t.Errorf(
			"len(Filter[]) = %d, cap(Filter[]) = %d ожидался срез из двух значений с размером выделенной "+
				"памяти два значения.",
			len(obj.Filter), cap(obj.Filter),
		)
		return
	}
	obj.ResetFilter()
	if len(obj.Filter) != 0 || cap(obj.Filter) != 2 {
		t.Errorf(
			"len(Filter[]) = %d, cap(Filter[]) = %d ожидался пустой срез с размером выделенной памяти два значения.",
			len(obj.Filter), cap(obj.Filter),
		)
		return
	}
}

// Тестирование сброса DOM дерева карты настроек сложной фильтрации.
func TestResetMap(t *testing.T) {
	const defMap = "g1:or:g2:or:(g3:and:g4)"
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
		ers []*ParseError
	)

	obj = New().(*impl)
	if obj.Map != nil {
		t.Errorf("Map имеет не верные первоначальные значения.")
		return
	}
	ret = obj.ResetMap().(*impl)
	if obj != ret {
		t.Errorf("ResetMap() != object, ожидалось то что функция вернёт объект.")
		return
	}
	if rq, err = http.NewRequest(
		"GET",
		"http://localhost?map="+defMap+"&g1=id:eq:1&g2=key:eq:1&g3=uu:le:10&g4=zz:in:1,2,3,4",
		nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseMap(rq); len(ers) > 0 {
		t.Errorf("ParseMap() = %v, ошибка не ожидалась.", ers)
		return
	}
	if obj.Map == nil {
		t.Errorf("Map == nil, ожидалась не пустая карта DOM дерева сложной фильтрации.")
		return
	}
	if obj.Map.SumSize() != len(defMap) {
		t.Errorf("Map.SumSize() == %d, ожидался размер карты DOM дерева равный %d.", obj.Map.SumSize(), len(defMap))
		return
	}
	obj.ResetMap()
	if obj.Map != nil {
		t.Errorf("Map имеет не верное значение после сброса.")
		return
	}
}

// Тестирование сброса переопределения названий полей или замены полей.
func TestResetRedefinition(t *testing.T) {
	var (
		obj *impl
		ret *impl
	)

	obj = New().(*impl)
	if len(obj.Remap) != 0 {
		t.Errorf("Remap имеет не верные первоначальные значения.")
		return
	}
	ret = obj.ResetRedefinition().(*impl)
	if obj != ret {
		t.Errorf("ResetRedefinition() != object, ожидалось то что функция вернёт объект.")
		return
	}
	obj.
		Redefinition("mId", "`method`.`mId`").
		Redefinition("key", "`cursor`.`id`")
	if len(obj.Remap) != 2 {
		t.Errorf("len(Remap) = %d, ожидалась карта из двух значений.", len(obj.Remap))
		return
	}
	obj.ResetRedefinition()
	if len(obj.Remap) != 0 {
		t.Errorf("len(Remap) = %d, ожидалась пустая карта переопределения названий полей.", len(obj.Remap))
		return
	}
}

// Тестирование сброса назначения типа данных для полей.
func TestResetDatatype(t *testing.T) {
	var (
		obj *impl
		ret *impl
	)

	obj = New().(*impl)
	if len(obj.Datatype) != 0 {
		t.Errorf("Datatype имеет не верные первоначальные значения.")
		return
	}
	ret = obj.ResetDatatype().(*impl)
	if obj != ret {
		t.Errorf("ResetDatatype() != object, ожидалось то что функция вернёт объект.")
		return
	}
	obj.
		FieldDatatype("id", TypeUint64).
		FieldDatatype("createAt", TypeTime).
		FieldDatatype("order", TypeFloat64).
		FieldDatatype("normal", TypeBool).
		FieldDatatype("mId", TypeInt64)
	if len(obj.Datatype) != 5 {
		t.Errorf("len(Datatype) = %d, ожидалась карта типов данных равная %d.", len(obj.Datatype), 5)
		return
	}
	obj.ResetDatatype()
	if len(obj.Datatype) != 0 {
		t.Errorf("len(Datatype) = %d, ожидалась пустая карта типов данных полей.", len(obj.Datatype))
		return
	}
}

// Тестирование сброса списка обрабатываемых полей.
func TestResetFieldSet(t *testing.T) {
	var (
		obj *impl
		ret *impl
		dat []string
	)

	obj = New().(*impl)
	if len(obj.FieldOnly) != 0 || cap(obj.FieldOnly) != 0 {
		t.Errorf("FieldOnly[] имеет не верные первоначальные значения.")
		return
	}
	ret = obj.ResetFieldSet().(*impl)
	if obj != ret {
		t.Errorf("ResetFieldSet() != object, ожидалось то что функция вернёт объект.")
		return
	}
	dat = []string{"i`d`", "'c'r`e`a`te'At", `update"A"t`, "o,r,d,e,r", "normal"}
	obj.FieldSet(dat...)
	if len(obj.FieldOnly) != 5 || cap(obj.FieldOnly) != 5 {
		t.Errorf(
			"len(FieldOnly[]) = %d, cap(FieldOnly[]) = %d ожидался срез длинной %d и выделенной памятью "+
				"для %d значений.",
			len(obj.FieldOnly), cap(obj.FieldOnly), len(dat), cap(dat),
		)
		return
	}
	obj.ResetFieldSet()
	if len(obj.FieldOnly) != 0 || cap(obj.FieldOnly) != 5 {
		t.Errorf(
			"len(FieldOnly[]) = %d, cap(FieldOnly[]) = %d ожидался пустой срез с размером выделенной памяти "+
				"для значения %d.",
			len(obj.FieldOnly), cap(obj.FieldOnly), 5,
		)
		return
	}
}

// Тестирование сброса всех данных.
func TestReset(t *testing.T) {
	const data = `http://localhost?
&map=(group1:and:group2:or:group3):or:(group4:or:group5):or:(group4:and:group5)
&map=group4:and:group5
&group1=id:eq:1
&group2=id:ke:2
&group3=name:ne:value3
&group4=section:ke:value4
&group5=updateAt:ke:value5
&by=id:asc&by=name:desc
&limit=100:10
&filter=section:kn:П*й
&filter=id:ge:10
`
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
	)

	obj = New().(*impl)
	ret = obj.Reset().(*impl)
	if obj != ret {
		t.Errorf("Reset() != object, ожидалось то что функция вернёт объект.")
		return
	}
	obj.
		FieldDatatype("id", TypeUint64).
		FieldDatatype("createAt", TypeTime).
		FieldDatatype("order", TypeFloat64).
		FieldDatatype("normal", TypeBool).
		FieldDatatype("mId", TypeInt64).
		Redefinition("mId", "`method`.`mId`").
		Redefinition("key", "`cursor`.`id`").
		FieldSet([]string{"i`d`", "'c'r`e`a`te'At", `update"A"t`, "o,r,d,e,r", "normal", "section", "name"}...)
	if rq, err = http.NewRequest("GET", strings.ReplaceAll(data, "\n", ""), nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if err = obj.ParseRequest(rq); err != nil {
		t.Errorf("Выполнение функции ParseRequest() прервано, ошибкой: %s.", err)
		return
	}
	if obj.Offset == 0 || obj.Limit == 0 {
		t.Errorf("ParseRequest() - Offset или Limit имеют не верные значения.")
		return
	}
	if len(obj.By) == 0 || cap(obj.By) == 0 {
		t.Errorf("ParseRequest() - By[] имеет не верное значение.")
		return
	}
	if len(obj.Filter) == 0 || cap(obj.Filter) == 0 {
		t.Errorf("ParseRequest() - Filter[] имеет не верное значениея.")
		return
	}
	if obj.Map == nil {
		t.Errorf("ParseRequest() - Map имеет не верное значение.")
		return
	}
	if obj.Map.SumSize() != 100 {
		t.Errorf(
			"ParseRequest() - Map.SumSize() == %d, ожидался размер карты DOM дерева равный %d.",
			obj.Map.SumSize(), 100,
		)
		return
	}
	if len(obj.Remap) == 0 {
		t.Errorf("ParseRequest() - Remap имеет не верное значение.")
		return
	}
	if len(obj.Datatype) == 0 {
		t.Errorf("ParseRequest() - Datatype имеет не верное значение.")
		return
	}
	if len(obj.FieldOnly) == 0 || cap(obj.FieldOnly) == 0 {
		t.Errorf("ParseRequest() - FieldOnly[] имеет не верное значение.")
		return
	}
	// Сброс и проверка того что всё сбросилось правильно.
	obj.Reset()
	if obj.Offset != 0 || obj.Limit != 0 {
		t.Errorf(
			"ParseRequest() - Offset == %d, Limit == %d, ожидалось Offset == %d, Limit == %d.",
			obj.Offset, obj.Limit, 0, 0,
		)
		return
	}
	if len(obj.By) != 0 || cap(obj.By) != 2 {
		t.Errorf(
			"ParseRequest() - len(By[]) == %d, cap(By[]) == %d, ожидалось len(By[]) == %d, cap(By[]) == %d.",
			len(obj.By), cap(obj.By), 0, 2,
		)
		return
	}
	if len(obj.Filter) != 0 || cap(obj.Filter) != 2 {
		t.Errorf(
			"ParseRequest() - len(Filter[]) == %d, cap(Filter[]) == %d, ожидалось "+
				"len(Filter[]) == %d, cap(Filter[]) == %d.",
			len(obj.Filter), cap(obj.Filter), 0, 2,
		)
		return
	}
	if obj.Map != nil {
		t.Errorf("ParseRequest() - Map не сбросило значение.")
		return
	}
	if len(obj.Remap) != 0 {
		t.Errorf("ParseRequest() - Remap не сбросило значение.")
		return
	}
	if len(obj.Datatype) != 0 {
		t.Errorf("ParseRequest() - Datatype не сбросило значение.")
		return
	}
	if len(obj.FieldOnly) != 0 || cap(obj.FieldOnly) != 7 {
		t.Errorf(
			"ParseRequest() - len(FieldOnly[]) == %d, cap(FieldOnly[]) == %d, ожидалось "+
				"len(FieldOnly[]) == %d, cap(FieldOnly[]) == %d.",
			len(obj.FieldOnly), cap(obj.FieldOnly), 0, 7,
		)
		return
	}
}
