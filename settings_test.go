package f8n

import "testing"

// Тестирование перечисления обрабатываемых полей с заменой недопустимых символов в названии полей.
func TestFieldSetGet(t *testing.T) {
	const (
		tplErrorFieldSet = "FieldSet() = %v, ожидалось %v."
		tplErrorFieldGet = "FieldGet() = %v, ожидалось %v."
	)
	var (
		dat []string
		obj *impl
		get []string
	)

	dat = []string{"i`d`", "'c'r`e`a`te'At", `update"A"t`, "o,r,d,e,r", "normal"}
	obj = New().(*impl)
	obj.FieldSet(dat...)
	if len(obj.FieldOnly) != len(dat) {
		t.Errorf(tplErrorFieldSet, obj.FieldOnly, dat)
	}
	if get = obj.FieldGet(); get[4] != dat[4] ||
		get[0] == dat[0] ||
		get[1] == dat[1] ||
		get[2] == dat[2] ||
		get[3] == dat[3] {
		t.Errorf(tplErrorFieldGet, get, dat)
	}
}

// Тестирование назначения типа данных загружаемым значениям полей.
func TestFieldDatatype(t *testing.T) {
	const (
		tplErrorDatatypeNil  = "Datatype map = nil, ожидалась инициализация map[string]Datatype."
		tplErrorDatatypeLen  = "Datatype map = %d, ожидалось %d."
		tplErrorDatatypeInit = "Datatype map не пустой, ожидалась пустая карта."
	)
	var (
		obj *impl
	)

	obj = New().(*impl)
	if obj.Datatype == nil {
		t.Errorf(tplErrorDatatypeNil)
	}
	if len(obj.Datatype) != 0 {
		t.Errorf(tplErrorDatatypeInit)
	}
	obj.
		FieldDatatype("id", TypeUint64).
		FieldDatatype("createAt", TypeTime).
		FieldDatatype("deleteAt", TypeTimeNil).
		FieldDatatype("order", TypeFloat64).
		FieldDatatype("normal", TypeBool).
		FieldDatatype("mId", TypeInt64)
	if len(obj.Datatype) != 6 {
		t.Errorf(tplErrorDatatypeLen, len(obj.Datatype), 6)
	}
}

// Тестирование переопределения названий полей.
func TestRedefinition(t *testing.T) {
	const (
		tplErrorRemapNil  = "Remap map = nil, ожидалась инициализация map[string]string."
		tplErrorRemapLen  = "Remap map = %d, ожидалось %d."
		tplErrorRemapInit = "Remap map не пустой, ожидалась пустая карта."
	)
	var (
		obj *impl
	)

	obj = New().(*impl)
	if obj.Remap == nil {
		t.Errorf(tplErrorRemapNil)
	}
	if len(obj.Remap) != 0 {
		t.Errorf(tplErrorRemapInit)
	}
	obj.Redefinition("mId", "`method`.`mId`")
	if len(obj.Remap) != 1 {
		t.Errorf(tplErrorRemapLen, len(obj.Remap), 1)
	}
}
