package f8n

import (
	"math"
	"strconv"
	"testing"
	"time"
)

// Тестирование тестирования соответствия значения указанному типу.
//
//goland:noinspection DuplicatedCode
func TestFilterValueTest(t *testing.T) {
	var (
		err error
		fvo FilterValue
		tmo time.Time
	)

	// uint64
	fvo = FilterValue{Source: "test", Type: TypeUint64}
	if err = fvo.Test(); err == nil {
		t.Errorf("FilterValue.Test() = %v, ожидалась ошибка.", err)
		return
	}
	if fvo.String() != "test" {
		t.Errorf("FilterValue.String() = %q, ожидалась %q.", fvo.String(), "test")
		return
	}
	fvo = FilterValue{Source: strconv.FormatUint(math.MaxUint64, 10), Type: TypeUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	if fvo.String() != strconv.FormatUint(math.MaxUint64, 10) {
		t.Errorf("FilterValue.String() = %q, ожидалась %q.", fvo.String(), strconv.FormatUint(math.MaxUint64, 10))
		return
	}
	// int64
	fvo = FilterValue{Source: "-test18", Type: TypeInt64}
	if err = fvo.Test(); err == nil {
		t.Errorf("FilterValue.Test() = %v, ожидалась ошибка.", err)
		return
	}
	fvo = FilterValue{Source: strconv.FormatInt(math.MinInt64, 10), Type: TypeInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	if fvo.String() != strconv.FormatInt(math.MinInt64, 10) {
		t.Errorf("FilterValue.String() = %q, ожидалась %q.", fvo.String(), strconv.FormatInt(math.MinInt64, 10))
		return
	}
	fvo = FilterValue{Source: strconv.FormatInt(math.MaxInt64, 10), Type: TypeInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	// float64
	fvo = FilterValue{Source: "1.8.9", Type: TypeFloat64}
	if err = fvo.Test(); err == nil {
		t.Errorf("FilterValue.Test() = %v, ожидалась ошибка.", err)
		return
	}
	fvo = FilterValue{Source: strconv.FormatFloat(1.2345678, 'E', -1, 64), Type: TypeFloat64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: strconv.FormatFloat(0-3.1415, 'E', -1, 64), Type: TypeFloat64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	// bool
	fvo = FilterValue{Source: "Yes", Type: TypeBool}
	if err = fvo.Test(); err == nil {
		t.Errorf("FilterValue.Test() = %v, ожидалась ошибка.", err)
		return
	}
	fvo = FilterValue{Source: "true", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: "false", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: "1", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: "0", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	// time
	fvo = FilterValue{Source: "18:53", Type: TypeTime}
	if err = fvo.Test(); err == nil {
		t.Errorf("FilterValue.Test() = %v, ожидалась ошибка.", err)
		return
	}
	tmo = time.Now()
	fvo = FilterValue{Source: tmo.Format(time.RFC3339Nano), Type: TypeTime}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: tmo.Format(time.RFC3339), Type: TypeTime}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	// SliceString
	fvo = FilterValue{Source: "", Type: TypeSliceString}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: "one,two;three", Type: TypeSliceString}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	if item, ok := fvo.Value().([]string); !ok || len(item) != 2 || item[0] != "one" || item[1] != "two;three" {
		t.Errorf("Value() = %v, ожидалась %v.", item, []string{"one", "two;three"})
		return
	}
	// SliceInt64
	fvo = FilterValue{Source: "", Type: TypeSliceInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: "3,2,1", Type: TypeSliceInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	if item, ok := fvo.Value().([]int64); !ok || len(item) != 3 || item[0] != 3 || item[1] != 2 || item[2] != 1 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []int64{3, 2, 1})
		return
	}
	fvo = FilterValue{
		Source: "," + strconv.FormatInt(math.MaxInt64, 10) + "," + strconv.FormatInt(math.MinInt64, 10),
		Type:   TypeSliceInt64,
	}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	if item, ok := fvo.Value().([]int64); !ok ||
		len(item) != 3 ||
		item[0] != 0 ||
		item[1] != 9223372036854775807 ||
		item[2] != -9223372036854775808 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []int64{0, 9223372036854775807, -9223372036854775808})
		return
	}
	fvo = FilterValue{Source: "-1,zero", Type: TypeSliceInt64}
	if err = fvo.Test(); !Errors().FilterValueType.Is(err) {
		t.Errorf("FilterValue.Test() = %v, ожидалась ошибка %q.", err, Errors().FilterValueType.Bind())
		return
	}
	// TypeSliceUint64
	fvo = FilterValue{Source: "", Type: TypeSliceUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	fvo = FilterValue{Source: "7,8,9", Type: TypeSliceUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	if item, ok := fvo.Value().([]uint64); !ok || len(item) != 3 || item[0] != 7 || item[1] != 8 || item[2] != 9 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []uint64{3, 2, 1})
		return
	}
	fvo = FilterValue{Source: "," + strconv.FormatUint(math.MaxUint64, 10), Type: TypeSliceUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf("FilterValue.Test() = %v, ошибка не ожидалась.", err)
		return
	}
	if item, ok := fvo.Value().([]uint64); !ok || len(item) != 2 || item[0] != 0 || item[1] != 18446744073709551615 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []uint64{0, 18446744073709551615})
		return
	}
	fvo = FilterValue{Source: "-1", Type: TypeSliceUint64}
	if err = fvo.Test(); !Errors().FilterValueType.Is(err) {
		t.Errorf("FilterValue.Test() = %v, ожидалась ошибка %q.", err, Errors().FilterValueType.Bind())
		return
	}
}
