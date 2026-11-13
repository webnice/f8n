package f8n

import (
	"reflect"
	"testing"
	"time"
)

// Тестирование формирования SQL запроса и параметров простой фильтрации для ORM GORM.
//
//goland:noinspection ALL
func TestFilterQueryGorm(t *testing.T) {
	type (
		test struct {
			Name   string
			Filter Filter
			Query  string
			Values []any
		}
	)
	var (
		tmo   time.Time
		n     int
		tests []test
	)

	tmo = time.Date(2018, 8, 3, 3, 0, 0, 0, time.FixedZone("UTC+3", 3*60*60))
	tests = []test{
		// Равно, обычное значение.
		{
			Name:   "eq-string",
			Filter: Filter{Field: "f", Method: filterEquivalent, Value: FilterValue{Source: "v", Type: TypeString}},
			Query:  " = ?", Values: []any{"v"},
		},
		// Равно, тип допускающий nil, значение не nil.
		{
			Name:   "eq-nil-type-value",
			Filter: Filter{Field: "f", Method: filterEquivalent, Value: FilterValue{Source: "10", Type: TypeUint64Nil}},
			Query:  " = ?", Values: []any{"10"},
		},
		// Равно, тип допускающий nil, значение nil.
		{
			Name:   "eq-nil-value",
			Filter: Filter{Field: "f", Method: filterEquivalent, Value: FilterValue{Source: "nil", Type: TypeInt64Nil}},
			Query:  " IS NULL",
		},
		// Меньше, значение nil - условие не формируется.
		{
			Name:   "lt-nil-value",
			Filter: Filter{Field: "f", Method: filterLessThan, Value: FilterValue{Source: "NIL", Type: TypeFloat64Nil}},
			Query:  "",
		},
		// Меньше, тип допускающий nil, значение не nil.
		{
			Name:   "lt-nil-type",
			Filter: Filter{Field: "f", Method: filterLessThan, Value: FilterValue{Source: "-5", Type: TypeInt64Nil}},
			Query:  " < ?", Values: []any{"-5"},
		},
		// Меньше.
		{
			Name:   "lt",
			Filter: Filter{Field: "f", Method: filterLessThan, Value: FilterValue{Source: "1.5", Type: TypeFloat64}},
			Query:  " < ?", Values: []any{1.5},
		},
		// Меньше или равно, значение nil - условие не формируется.
		{
			Name:   "le-nil-value",
			Filter: Filter{Field: "f", Method: filterLessThanOrEquivalent, Value: FilterValue{Source: "nil", Type: TypeBoolNil}},
			Query:  "",
		},
		// Меньше или равно, тип допускающий nil, значение не nil.
		{
			Name:   "le-nil-type",
			Filter: Filter{Field: "f", Method: filterLessThanOrEquivalent, Value: FilterValue{Source: "7", Type: TypeUint64Nil}},
			Query:  " <= ?", Values: []any{"7"},
		},
		// Меньше или равно.
		{
			Name:   "le",
			Filter: Filter{Field: "f", Method: filterLessThanOrEquivalent, Value: FilterValue{Source: "-7", Type: TypeInt64}},
			Query:  " <= ?", Values: []any{int64(-7)},
		},
		// Больше, значение nil - условие не формируется.
		{
			Name:   "gt-nil-value",
			Filter: Filter{Field: "f", Method: filterGreaterThan, Value: FilterValue{Source: "nil", Type: TypeTimeNil}},
			Query:  "",
		},
		// Больше, тип допускающий nil, значение не nil.
		{
			Name:   "gt-nil-type",
			Filter: Filter{Field: "f", Method: filterGreaterThan, Value: FilterValue{Source: "1.5", Type: TypeFloat64Nil}},
			Query:  " > ?", Values: []any{"1.5"},
		},
		// Больше.
		{
			Name:   "gt",
			Filter: Filter{Field: "f", Method: filterGreaterThan, Value: FilterValue{Source: "7", Type: TypeUint64}},
			Query:  " > ?", Values: []any{uint64(7)},
		},
		// Больше или равно, значение nil - условие не формируется.
		{
			Name:   "ge-nil-value",
			Filter: Filter{Field: "f", Method: filterGreaterThanOrEquivalent, Value: FilterValue{Source: "nil", Type: TypeStringNil}},
			Query:  "",
		},
		// Больше или равно, тип допускающий nil, значение не nil.
		{
			Name:   "ge-nil-type",
			Filter: Filter{Field: "f", Method: filterGreaterThanOrEquivalent, Value: FilterValue{Source: "true", Type: TypeBoolNil}},
			Query:  " >= ?", Values: []any{"true"},
		},
		// Больше или равно.
		{
			Name:   "ge",
			Filter: Filter{Field: "f", Method: filterGreaterThanOrEquivalent, Value: FilterValue{Source: tmo.Format(time.RFC3339), Type: TypeTime}},
			Query:  " >= ?", Values: []any{tmo},
		},
		// Не равно, значение nil.
		{
			Name:   "ne-nil-value",
			Filter: Filter{Field: "f", Method: filterNotEquivalent, Value: FilterValue{Source: "nil", Type: TypeStringNil}},
			Query:  " IS NOT NULL",
		},
		// Не равно, тип допускающий nil, значение не nil.
		{
			Name:   "ne-nil-type",
			Filter: Filter{Field: "f", Method: filterNotEquivalent, Value: FilterValue{Source: "v", Type: TypeStringNil}},
			Query:  " != ?", Values: []any{"v"},
		},
		// Не равно.
		{
			Name:   "ne",
			Filter: Filter{Field: "f", Method: filterNotEquivalent, Value: FilterValue{Source: "false", Type: TypeBool}},
			Query:  " != ?", Values: []any{false},
		},
		// Похоже, значение nil - условие не формируется.
		{
			Name:   "ke-nil-value",
			Filter: Filter{Field: "f", Method: filterLikeThan, Value: FilterValue{Source: "nil", Type: TypeStringNil}},
			Query:  "",
		},
		// Похоже, тип допускающий nil, значение не nil.
		{
			Name:   "ke-nil-type",
			Filter: Filter{Field: "f", Method: filterLikeThan, Value: FilterValue{Source: "a?b*", Type: TypeStringNil}},
			Query:  " LIKE ?", Values: []any{"a_b%"},
		},
		// Похоже.
		{
			Name:   "ke",
			Filter: Filter{Field: "f", Method: filterLikeThan, Value: FilterValue{Source: "П?рв*й", Type: TypeString}},
			Query:  " LIKE ?", Values: []any{"П_рв%й"},
		},
		// Не похоже, значение nil - условие не формируется.
		{
			Name:   "kn-nil-value",
			Filter: Filter{Field: "f", Method: filterNotLikeThan, Value: FilterValue{Source: "nil", Type: TypeStringNil}},
			Query:  "",
		},
		// Не похоже, тип допускающий nil, значение не nil.
		{
			Name:   "kn-nil-type",
			Filter: Filter{Field: "f", Method: filterNotLikeThan, Value: FilterValue{Source: "a?b*", Type: TypeStringNil}},
			Query:  " NOT LIKE ?", Values: []any{"a_b%"},
		},
		// Не похоже.
		{
			Name:   "kn",
			Filter: Filter{Field: "f", Method: filterNotLikeThan, Value: FilterValue{Source: "a?c*", Type: TypeString}},
			Query:  " NOT LIKE ?", Values: []any{"a_c%"},
		},
		// В списке, срез строк.
		{
			Name:   "in-string",
			Filter: Filter{Field: "f", Method: filterIn, Value: FilterValue{Source: "a, b", Type: TypeSliceString}},
			Query:  " IN (?)", Values: []any{[]string{"a", "b"}},
		},
		// В списке, срез int64.
		{
			Name:   "in-int64",
			Filter: Filter{Field: "f", Method: filterIn, Value: FilterValue{Source: "3,2,1", Type: TypeSliceInt64}},
			Query:  " IN (?)", Values: []any{[]int64{3, 2, 1}},
		},
		// В списке, срез uint64.
		{
			Name:   "in-uint64",
			Filter: Filter{Field: "f", Method: filterIn, Value: FilterValue{Source: "7,8", Type: TypeSliceUint64}},
			Query:  " IN (?)", Values: []any{[]uint64{7, 8}},
		},
		// В списке, тип не срез - условие не формируется.
		{
			Name:   "in-not-slice",
			Filter: Filter{Field: "f", Method: filterIn, Value: FilterValue{Source: "7,8", Type: TypeString}},
			Query:  "",
		},
		// Не в списке, срез строк.
		{
			Name:   "ni-string",
			Filter: Filter{Field: "f", Method: filterNotIn, Value: FilterValue{Source: "a,b", Type: TypeSliceString}},
			Query:  " NOT IN (?)", Values: []any{[]string{"a", "b"}},
		},
		// Не в списке, срез int64.
		{
			Name:   "ni-int64",
			Filter: Filter{Field: "f", Method: filterNotIn, Value: FilterValue{Source: "3,2", Type: TypeSliceInt64}},
			Query:  " NOT IN (?)", Values: []any{[]int64{3, 2}},
		},
		// Не в списке, срез uint64.
		{
			Name:   "ni-uint64",
			Filter: Filter{Field: "f", Method: filterNotIn, Value: FilterValue{Source: "7", Type: TypeSliceUint64}},
			Query:  " NOT IN (?)", Values: []any{[]uint64{7}},
		},
		// Не в списке, тип не срез - условие не формируется.
		{
			Name:   "ni-not-slice",
			Filter: Filter{Field: "f", Method: filterNotIn, Value: FilterValue{Source: "7", Type: TypeUint64}},
			Query:  "",
		},
		// Неизвестный метод фильтрации.
		{
			Name:   "unknown-method",
			Filter: Filter{Field: "f", Method: filterUnknown, Value: FilterValue{Source: "v", Type: TypeString}},
			Query:  "Метод sql запроса `` не реализован.",
		},
	}
	for n = range tests {
		var (
			query  string
			values []any
		)
		query, values = tests[n].Filter.queryGorm()
		if query != tests[n].Query {
			t.Errorf("%s: queryGorm() = %q, ожидалось: %q.", tests[n].Name, query, tests[n].Query)
			continue
		}
		if !equalAnyValues(values, tests[n].Values) {
			t.Errorf("%s: queryGorm() значения = %v, ожидалось: %v.", tests[n].Name, values, tests[n].Values)
		}
	}
}

// Сравнение срезов значений с учётом сравнения времени через Equal.
func equalAnyValues(a []any, b []any) (ret bool) {
	var n int

	if len(a) != len(b) {
		return
	}
	for n = range a {
		if !equalAnyValue(a[n], b[n]) {
			return
		}
	}
	ret = true

	return
}

// Сравнение значений с учётом сравнения времени через Equal.
func equalAnyValue(a any, b any) (ret bool) {
	var (
		ta, tb time.Time
		okA    bool
	)

	if ret = reflect.DeepEqual(a, b); ret {
		return
	}
	if ta, okA = a.(time.Time); !okA {
		return
	}
	if tb, okA = b.(time.Time); !okA {
		return
	}
	ret = ta.Equal(tb)

	return
}

// Тестирование замены символов в значении фильтрации для запросов SQL LIKE.
func TestFilterQueryGormValueLike(t *testing.T) {
	const (
		src  = "П?рв*й ?*"
		want = "П_рв%й _%"
	)
	var (
		filter Filter
		ret    string
	)

	filter = Filter{Value: FilterValue{Source: src, Type: TypeString}}
	if ret = filter.queryGormValueLike(); ret != want {
		t.Errorf("queryGormValueLike() = %q, ожидалось: %q.", ret, want)
		return
	}
}

// Тестирование проверки значения фильтрации на nil.
func TestFilterValueIsValueNil(t *testing.T) {
	type (
		test struct {
			Name  string
			Value FilterValue
			IsNil bool
		}
	)
	var (
		n     int
		tests []test
	)

	tests = []test{
		// Типы допускающие nil, значение nil в разном регистре.
		{Name: "string-nil", Value: FilterValue{Source: "nil", Type: TypeStringNil}, IsNil: true},
		{Name: "string-NIL", Value: FilterValue{Source: "NIL", Type: TypeStringNil}, IsNil: true},
		{Name: "uint64-nil", Value: FilterValue{Source: "nil", Type: TypeUint64Nil}, IsNil: true},
		{Name: "int64-nil", Value: FilterValue{Source: "nil", Type: TypeInt64Nil}, IsNil: true},
		{Name: "float64-nil", Value: FilterValue{Source: "nil", Type: TypeFloat64Nil}, IsNil: true},
		{Name: "bool-nil", Value: FilterValue{Source: "nil", Type: TypeBoolNil}, IsNil: true},
		{Name: "time-nil", Value: FilterValue{Source: "nil", Type: TypeTimeNil}, IsNil: true},
		// Типы допускающие nil, значение не nil.
		{Name: "string-nil-value", Value: FilterValue{Source: "тест", Type: TypeStringNil}, IsNil: false},
		{Name: "uint64-nil-value", Value: FilterValue{Source: "10", Type: TypeUint64Nil}, IsNil: false},
		// Типы не допускающие nil.
		{Name: "string", Value: FilterValue{Source: "nil", Type: TypeString}, IsNil: false},
		{Name: "uint64", Value: FilterValue{Source: "nil", Type: TypeUint64}, IsNil: false},
		{Name: "int64", Value: FilterValue{Source: "nil", Type: TypeInt64}, IsNil: false},
		{Name: "float64", Value: FilterValue{Source: "nil", Type: TypeFloat64}, IsNil: false},
		{Name: "bool", Value: FilterValue{Source: "nil", Type: TypeBool}, IsNil: false},
		{Name: "time", Value: FilterValue{Source: "nil", Type: TypeTime}, IsNil: false},
		{Name: "slice-string", Value: FilterValue{Source: "nil", Type: TypeSliceString}, IsNil: false},
		{Name: "slice-int64", Value: FilterValue{Source: "nil", Type: TypeSliceInt64}, IsNil: false},
		{Name: "slice-uint64", Value: FilterValue{Source: "nil", Type: TypeSliceUint64}, IsNil: false},
	}
	for n = range tests {
		if ret := tests[n].Value.IsValueNil(); ret != tests[n].IsNil {
			t.Errorf("%s: IsValueNil() = %t, ожидалось: %t.", tests[n].Name, ret, tests[n].IsNil)
		}
	}
}

// Тестирование возвращения значения фильтрации в типе указанном у значения.
//
//goland:noinspection ALL
func TestFilterValueValue(t *testing.T) {
	const (
		srcUint    = "18446744073709551615"
		srcInt     = "-9223372036854775808"
		srcFloat   = "-3.1415"
		srcBool    = "true"
		srcTime    = "2018-08-03T03:00:00+03:00"
		srcBadUint = "bad"
		srcBadInt  = "bad"
		srcBadBool = "bad"
		srcBadTime = "bad"
		srcSlices  = "3,2,1"
		srcString  = "произвольная строка"
	)
	var (
		tmo  time.Time
		err  error
		fvo  FilterValue
		item any
	)

	tmo, err = time.Parse(time.RFC3339, srcTime)
	if err != nil {
		t.Errorf("Ошибка подготовки тестовых данных: %s.", err)
		return
	}
	// Строка.
	fvo = FilterValue{Source: srcString, Type: TypeString}
	if item = fvo.Value(); item != srcString {
		t.Errorf("Value() = %v (%T), ожидалось: %q.", item, item, srcString)
		return
	}
	// Тип по умолчанию, для не известного типа значение возвращается как строка.
	fvo = FilterValue{Source: srcString, Type: TypeStringNil}
	if item = fvo.Value(); item != srcString {
		t.Errorf("Value() = %v (%T), ожидалось: %q.", item, item, srcString)
		return
	}
	// uint64.
	fvo = FilterValue{Source: srcUint, Type: TypeUint64}
	if item = fvo.Value(); item != uint64(18446744073709551615) {
		t.Errorf("Value() = %v (%T), ожидалось: %d.", item, item, uint64(18446744073709551615))
		return
	}
	// Ошибка конвертации uint64, ожидается nil.
	fvo = FilterValue{Source: srcBadUint, Type: TypeUint64}
	if item = fvo.Value(); item != nil {
		t.Errorf("Value() = %v, ожидалось: nil.", item)
		return
	}
	// int64.
	fvo = FilterValue{Source: srcInt, Type: TypeInt64}
	if item = fvo.Value(); item != int64(-9223372036854775808) {
		t.Errorf("Value() = %v (%T), ожидалось: %d.", item, item, -9223372036854775808)
		return
	}
	// Ошибка конвертации int64, ожидается nil.
	fvo = FilterValue{Source: srcBadInt, Type: TypeInt64}
	if item = fvo.Value(); item != nil {
		t.Errorf("Value() = %v, ожидалось: nil.", item)
		return
	}
	// float64.
	fvo = FilterValue{Source: srcFloat, Type: TypeFloat64}
	if item = fvo.Value(); item != -3.1415 {
		t.Errorf("Value() = %v (%T), ожидалось: %f.", item, item, -3.1415)
		return
	}
	// bool.
	fvo = FilterValue{Source: srcBool, Type: TypeBool}
	if item = fvo.Value(); item != true {
		t.Errorf("Value() = %v (%T), ожидалось: %t.", item, item, true)
		return
	}
	// Ошибка конвертации bool, ожидается nil.
	fvo = FilterValue{Source: srcBadBool, Type: TypeBool}
	if item = fvo.Value(); item != nil {
		t.Errorf("Value() = %v, ожидалось: nil.", item)
		return
	}
	// time.
	fvo = FilterValue{Source: srcTime, Type: TypeTime}
	if item = fvo.Value(); !item.(time.Time).Equal(tmo) {
		t.Errorf("Value() = %v, ожидалось: %s.", item, tmo)
		return
	}
	// Ошибка конвертации time, ожидается nil.
	fvo = FilterValue{Source: srcBadTime, Type: TypeTime}
	if item = fvo.Value(); item != nil {
		t.Errorf("Value() = %v, ожидалось: nil.", item)
		return
	}
	// Срез строк.
	fvo = FilterValue{Source: srcSlices, Type: TypeSliceString}
	if item = fvo.Value(); !reflect.DeepEqual(item, []string{"3", "2", "1"}) {
		t.Errorf("Value() = %v, ожидалось: %v.", item, []string{"3", "2", "1"})
		return
	}
	// Срез int64.
	fvo = FilterValue{Source: srcSlices, Type: TypeSliceInt64}
	if item = fvo.Value(); !reflect.DeepEqual(item, []int64{3, 2, 1}) {
		t.Errorf("Value() = %v, ожидалось: %v.", item, []int64{3, 2, 1})
		return
	}
	// Срез uint64.
	fvo = FilterValue{Source: srcSlices, Type: TypeSliceUint64}
	if item = fvo.Value(); !reflect.DeepEqual(item, []uint64{3, 2, 1}) {
		t.Errorf("Value() = %v, ожидалось: %v.", item, []uint64{3, 2, 1})
		return
	}
}
