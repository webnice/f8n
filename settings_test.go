package f8n

import "testing"

// Тестирование перечисление обрабатываемых полей с заменой недопустимых символов в названии полей.
func TestFieldSetGet(t *testing.T) {
	var (
		dat []string
		obj *impl
		get []string
	)

	dat = []string{"i`d`", "'c'r`e`a`te'At", `update"A"t`, "o,r,d,e,r", "normal"}
	obj = New().(*impl)
	obj.FieldSet(dat...)
	if len(obj.FieldOnly) != len(dat) {
		t.Errorf("FieldSet() = %v, ожидалаось %v.", obj.FieldOnly, dat)
	}
	if get = obj.FieldGet(); get[4] != dat[4] ||
		get[0] == dat[0] ||
		get[1] == dat[1] ||
		get[2] == dat[2] ||
		get[3] == dat[3] {
		t.Errorf("FieldGet() = %v, ожидалаось %v.", get, dat)
	}
}

// Тестирование назначения типа данных загружаемым значениям полей.
func TestFieldDatatype(t *testing.T) {
	var (
		obj *impl
	)

	obj = New().(*impl)
	if obj.Datatype == nil {
		t.Errorf("Datatype map = nil, ожидалась инициализация map[string]Datatype.")
	}
	if len(obj.Datatype) != 0 {
		t.Errorf("Datatype map не пустой, ожидалась пустая карта.")
	}
	obj.
		FieldDatatype("id", TypeUint64).
		FieldDatatype("createAt", TypeTime).
		FieldDatatype("order", TypeFloat64).
		FieldDatatype("normal", TypeBool).
		FieldDatatype("mId", TypeInt64)
	if len(obj.Datatype) != 5 {
		t.Errorf("Datatype map = %d, ожидалось %d.", len(obj.Datatype), 5)
	}
}

// Тестирование переопределения названий полей.
func TestRedefinition(t *testing.T) {
	var (
		obj *impl
	)

	obj = New().(*impl)
	if obj.Remap == nil {
		t.Errorf("Remap map = nil, ожидалась инициализация map[string]string.")
	}
	if len(obj.Remap) != 0 {
		t.Errorf("Remap map не пустой, ожидалась пустая карта.")
	}
	obj.Redefinition("mId", "`method`.`mId`")
	if len(obj.Remap) != 1 {
		t.Errorf("Remap map = %d, ожидалось %d.", len(obj.Remap), 1)
	}
}
