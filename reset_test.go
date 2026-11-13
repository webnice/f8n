package f8n

import (
	"math"
	"net/http"
	"strings"
	"testing"
)

// Тестирование сброса лимита.
func TestResetLimit(t *testing.T) {
	const (
		tplErrorInit  = "Offset или Limit имеют не верные первоначальные значения."
		tplErrorRet   = "ResetLimit() != object, ожидалось то что функция вернёт объект."
		tplErrorReset = "Функция ResetLimit() не выполнила сброс лимитов."
	)
	var (
		obj *impl
		ret *impl
	)

	obj = New().(*impl)
	if obj.Offset != 0 || obj.Limit != 0 {
		t.Errorf(tplErrorInit)
		return
	}
	ret = obj.ResetLimit().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
		return
	}
	obj.Offset, obj.Limit = math.MaxUint64, math.MaxUint64
	obj.ResetLimit()
	if obj.Offset != 0 || obj.Limit != 0 {
		t.Errorf(tplErrorReset)
		return
	}
}

// Тестирование сброса сортировки.
func TestResetBy(t *testing.T) {
	const (
		datUrn           = "http://localhost?by=id:asc&by=name:desc"
		tplErrorRequest  = "Ошибка создания запроса: %s."
		tplErrorInit     = "By[] имеет не верные первоначальные значения."
		tplErrorRet      = "ResetBy() != object, ожидалось то что функция вернёт объект."
		tplErrorParse    = "ParseBy() = %v, ошибка не ожидалась."
		tplErrorLen      = "len(By[]) = %d, cap(By[]) = %d ожидался срез из двух значений с размером выделенной памяти два значения."
		tplErrorLenReset = "len(By[]) = %d, cap(By[]) = %d ожидался пустой срез с размером выделенной памяти два значения."
	)
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
		ers []*ParseError
	)

	obj = New().(*impl)
	if len(obj.By) != 0 || cap(obj.By) != 0 {
		t.Errorf(tplErrorInit)
		return
	}
	ret = obj.ResetBy().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
		return
	}
	// Два значения сортировки.
	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseBy(rq); len(ers) > 0 {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if len(obj.By) != 2 || cap(obj.By) != 2 {
		t.Errorf(tplErrorLen, len(obj.By), cap(obj.By))
		return
	}
	obj.ResetBy()
	if len(obj.By) != 0 || cap(obj.By) != 2 {
		t.Errorf(tplErrorLenReset, len(obj.By), cap(obj.By))
		return
	}
}

// Тестирование сброса настроек фильтрации.
func TestResetFilter(t *testing.T) {
	const (
		datUrn           = "http://localhost?filter=aaa:eq:asc&filter=bbb:in:desc"
		tplErrorRequest  = "Ошибка создания запроса: %s."
		tplErrorInit     = "Filter[] имеет не верные первоначальные значения."
		tplErrorRet      = "ResetFilter() != object, ожидалось то что функция вернёт объект."
		tplErrorParse    = "ParseFilter() = %v, ошибка не ожидалась."
		tplErrorLen      = "len(Filter[]) = %d, cap(Filter[]) = %d ожидался срез из двух значений с размером выделенной памяти два значения."
		tplErrorLenReset = "len(Filter[]) = %d, cap(Filter[]) = %d ожидался пустой срез с размером выделенной памяти два значения."
	)
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
		ers []*ParseError
	)

	obj = New().(*impl)
	if len(obj.Filter) != 0 || cap(obj.Filter) != 0 {
		t.Errorf(tplErrorInit)
		return
	}
	ret = obj.ResetFilter().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
		return
	}
	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseFilter(rq); len(ers) > 0 {
		t.Errorf(tplErrorParse, ers)
		return
	}
	if len(obj.Filter) != 2 || cap(obj.Filter) != 2 {
		t.Errorf(tplErrorLen, len(obj.Filter), cap(obj.Filter))
		return
	}
	obj.ResetFilter()
	if len(obj.Filter) != 0 || cap(obj.Filter) != 2 {
		t.Errorf(tplErrorLenReset, len(obj.Filter), cap(obj.Filter))
		return
	}
}

// Тестирование сброса DOM дерева карты настроек сложной фильтрации.
func TestResetMap(t *testing.T) {
	const (
		defMap          = "g1:or:g2:or:(g3:and:g4)"
		datUrn          = "http://localhost?map=" + defMap + "&g1=id:eq:1&g2=key:eq:1&g3=uu:le:10&g4=zz:in:1,2,3,4"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorInit    = "Map имеет не верные первоначальные значения."
		tplErrorRet     = "ResetMap() != object, ожидалось то что функция вернёт объект."
		tplErrorParse   = "ParseMap() = %v, ошибка не ожидалась."
		tplErrorNil     = "Map == nil, ожидалась не пустая карта DOM дерева сложной фильтрации."
		tplErrorSize    = "Map.SumSize() == %d, ожидался размер карты DOM дерева равный %d."
		tplErrorReset   = "Map имеет не верное значение после сброса."
	)
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
		ers []*ParseError
	)

	obj = New().(*impl)
	if obj.Map != nil {
		t.Errorf(tplErrorInit)
		return
	}
	ret = obj.ResetMap().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
		return
	}
	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
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
	if obj.Map.SumSize() != len(defMap) {
		t.Errorf(tplErrorSize, obj.Map.SumSize(), len(defMap))
		return
	}
	obj.ResetMap()
	if obj.Map != nil {
		t.Errorf(tplErrorReset)
		return
	}
}

// Тестирование сброса переопределения названий полей или замены полей.
func TestResetRedefinition(t *testing.T) {
	const (
		tplErrorInit  = "Remap имеет не верные первоначальные значения."
		tplErrorRet   = "ResetRedefinition() != object, ожидалось то что функция вернёт объект."
		tplErrorLen   = "len(Remap) = %d, ожидалась карта из двух значений."
		tplErrorReset = "len(Remap) = %d, ожидалась пустая карта переопределения названий полей."
	)
	var (
		obj *impl
		ret *impl
	)

	obj = New().(*impl)
	if len(obj.Remap) != 0 {
		t.Errorf(tplErrorInit)
		return
	}
	ret = obj.ResetRedefinition().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
		return
	}
	obj.
		Redefinition("mId", "`method`.`mId`").
		Redefinition("key", "`cursor`.`id`")
	if len(obj.Remap) != 2 {
		t.Errorf(tplErrorLen, len(obj.Remap))
		return
	}
	obj.ResetRedefinition()
	if len(obj.Remap) != 0 {
		t.Errorf(tplErrorReset, len(obj.Remap))
		return
	}
}

// Тестирование сброса назначения типа данных для полей.
func TestResetDatatype(t *testing.T) {
	const (
		tplErrorInit  = "Datatype имеет не верные первоначальные значения."
		tplErrorRet   = "ResetDatatype() != object, ожидалось то что функция вернёт объект."
		tplErrorLen   = "len(Datatype) = %d, ожидалась карта типов данных равная %d."
		tplErrorReset = "len(Datatype) = %d, ожидалась пустая карта типов данных полей."
	)
	var (
		obj *impl
		ret *impl
	)

	obj = New().(*impl)
	if len(obj.Datatype) != 0 {
		t.Errorf(tplErrorInit)
		return
	}
	ret = obj.ResetDatatype().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
		return
	}
	obj.
		FieldDatatype("id", TypeUint64).
		FieldDatatype("createAt", TypeTime).
		FieldDatatype("order", TypeFloat64).
		FieldDatatype("normal", TypeBool).
		FieldDatatype("mId", TypeInt64)
	if len(obj.Datatype) != 5 {
		t.Errorf(tplErrorLen, len(obj.Datatype), 5)
		return
	}
	obj.ResetDatatype()
	if len(obj.Datatype) != 0 {
		t.Errorf(tplErrorReset, len(obj.Datatype))
		return
	}
}

// Тестирование сброса списка обрабатываемых полей.
func TestResetFieldSet(t *testing.T) {
	const (
		tplErrorInit     = "FieldOnly[] имеет не верные первоначальные значения."
		tplErrorRet      = "ResetFieldSet() != object, ожидалось то что функция вернёт объект."
		tplErrorLen      = "len(FieldOnly[]) = %d, cap(FieldOnly[]) = %d ожидался срез длинной %d и выделенной памятью для %d значений."
		tplErrorLenReset = "len(FieldOnly[]) = %d, cap(FieldOnly[]) = %d ожидался пустой срез с размером выделенной памяти для значения %d."
	)
	var (
		obj *impl
		ret *impl
		dat []string
	)

	obj = New().(*impl)
	if len(obj.FieldOnly) != 0 || cap(obj.FieldOnly) != 0 {
		t.Errorf(tplErrorInit)
		return
	}
	ret = obj.ResetFieldSet().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
		return
	}
	dat = []string{"i`d`", "'c'r`e`a`te'At", `update"A"t`, "o,r,d,e,r", "normal"}
	obj.FieldSet(dat...)
	if len(obj.FieldOnly) != 5 || cap(obj.FieldOnly) != 5 {
		t.Errorf(tplErrorLen, len(obj.FieldOnly), cap(obj.FieldOnly), len(dat), cap(dat))
		return
	}
	obj.ResetFieldSet()
	if len(obj.FieldOnly) != 0 || cap(obj.FieldOnly) != 5 {
		t.Errorf(tplErrorLenReset, len(obj.FieldOnly), cap(obj.FieldOnly), 5)
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
	const (
		tplErrorRequest   = "Ошибка создания запроса: %s."
		tplErrorRet       = "Reset() != object, ожидалось то что функция вернёт объект."
		tplErrorParse     = "Выполнение функции ParseRequest() прервано, ошибкой: %s."
		tplErrorLimit     = "ParseRequest() - Offset или Limit имеют не верные значения."
		tplErrorBy        = "ParseRequest() - By[] имеет не верное значение."
		tplErrorFilter    = "ParseRequest() - Filter[] имеет не верное значениея."
		tplErrorMap       = "ParseRequest() - Map имеет не верное значение."
		tplErrorMapSize   = "ParseRequest() - Map.SumSize() == %d, ожидался размер карты DOM дерева равный %d."
		tplErrorRemap     = "ParseRequest() - Remap имеет не верное значение."
		tplErrorDatatype  = "ParseRequest() - Datatype имеет не верное значение."
		tplErrorFieldOnly = "ParseRequest() - FieldOnly[] имеет не верное значение."
		tplErrorResetLim  = "ParseRequest() - Offset == %d, Limit == %d, ожидалось Offset == %d, Limit == %d."
		tplErrorResetBy   = "ParseRequest() - len(By[]) == %d, cap(By[]) == %d, ожидалось len(By[]) == %d, cap(By[]) == %d."
		tplErrorResetFlt  = "ParseRequest() - len(Filter[]) == %d, cap(Filter[]) == %d, ожидалось len(Filter[]) == %d, cap(Filter[]) == %d."
		tplErrorResetMap  = "ParseRequest() - Map не сбросило значение."
		tplErrorResetRmp  = "ParseRequest() - Remap не сбросило значение."
		tplErrorResetDtp  = "ParseRequest() - Datatype не сбросило значение."
		tplErrorResetFld  = "ParseRequest() - len(FieldOnly[]) == %d, cap(FieldOnly[]) == %d, ожидалось len(FieldOnly[]) == %d, cap(FieldOnly[]) == %d."
	)
	var (
		err error
		obj *impl
		ret *impl
		rq  *http.Request
	)

	obj = New().(*impl)
	ret = obj.Reset().(*impl)
	if obj != ret {
		t.Errorf(tplErrorRet)
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
		t.Errorf(tplErrorRequest, err)
	}
	if err = obj.ParseRequest(rq); err != nil {
		t.Errorf(tplErrorParse, err)
		return
	}
	if obj.Offset == 0 || obj.Limit == 0 {
		t.Errorf(tplErrorLimit)
		return
	}
	if len(obj.By) == 0 || cap(obj.By) == 0 {
		t.Errorf(tplErrorBy)
		return
	}
	if len(obj.Filter) == 0 || cap(obj.Filter) == 0 {
		t.Errorf(tplErrorFilter)
		return
	}
	if obj.Map == nil {
		t.Errorf(tplErrorMap)
		return
	}
	if obj.Map.SumSize() != 100 {
		t.Errorf(tplErrorMapSize, obj.Map.SumSize(), 100)
		return
	}
	if len(obj.Remap) == 0 {
		t.Errorf(tplErrorRemap)
		return
	}
	if len(obj.Datatype) == 0 {
		t.Errorf(tplErrorDatatype)
		return
	}
	if len(obj.FieldOnly) == 0 || cap(obj.FieldOnly) == 0 {
		t.Errorf(tplErrorFieldOnly)
		return
	}
	// Сброс и проверка того что всё сбросилось правильно.
	obj.Reset()
	if obj.Offset != 0 || obj.Limit != 0 {
		t.Errorf(tplErrorResetLim, obj.Offset, obj.Limit, 0, 0)
		return
	}
	if len(obj.By) != 0 || cap(obj.By) != 2 {
		t.Errorf(tplErrorResetBy, len(obj.By), cap(obj.By), 0, 2)
		return
	}
	if len(obj.Filter) != 0 || cap(obj.Filter) != 2 {
		t.Errorf(tplErrorResetFlt, len(obj.Filter), cap(obj.Filter), 0, 2)
		return
	}
	if obj.Map != nil {
		t.Errorf(tplErrorResetMap)
		return
	}
	if len(obj.Remap) != 0 {
		t.Errorf(tplErrorResetRmp)
		return
	}
	if len(obj.Datatype) != 0 {
		t.Errorf(tplErrorResetDtp)
		return
	}
	if len(obj.FieldOnly) != 0 || cap(obj.FieldOnly) != 7 {
		t.Errorf(tplErrorResetFld, len(obj.FieldOnly), cap(obj.FieldOnly), 0, 7)
		return
	}
}
