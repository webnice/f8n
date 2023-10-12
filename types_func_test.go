// Package f8n
package f8n

import (
	"math"
	"strconv"
	"testing"
	"time"
)

// Тестирование тестирования соответствия значения указанному типу.
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
}
