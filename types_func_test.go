package f8n

import (
	"encoding/json"
	"math"
	"net"
	"strconv"
	"testing"
	"time"
)

// Тестирование тестирования соответствия значения указанному типу.
//
//goland:noinspection DuplicatedCode
func TestFilterValueTest(t *testing.T) {
	const (
		tplErrorExpected = "FilterValue.Test() = %v, ожидалась ошибка."
		tplErrorNot      = "FilterValue.Test() = %v, ошибка не ожидалась."
		tplErrorString   = "FilterValue.String() = %q, ожидалась %q."
		tplErrorIs       = "FilterValue.Test() = %v, ожидалась ошибка %q."
	)
	var (
		err error
		fvo FilterValue
		tmo time.Time
	)

	// uint64
	fvo = FilterValue{Source: "test", Type: TypeUint64}
	if err = fvo.Test(); err == nil {
		t.Errorf(tplErrorExpected, err)
		return
	}
	if fvo.String() != "test" {
		t.Errorf(tplErrorString, fvo.String(), "test")
		return
	}
	fvo = FilterValue{Source: strconv.FormatUint(math.MaxUint64, 10), Type: TypeUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	if fvo.String() != strconv.FormatUint(math.MaxUint64, 10) {
		t.Errorf(tplErrorString, fvo.String(), strconv.FormatUint(math.MaxUint64, 10))
		return
	}
	// int64
	fvo = FilterValue{Source: "-test18", Type: TypeInt64}
	if err = fvo.Test(); err == nil {
		t.Errorf(tplErrorExpected, err)
		return
	}
	fvo = FilterValue{Source: strconv.FormatInt(math.MinInt64, 10), Type: TypeInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	if fvo.String() != strconv.FormatInt(math.MinInt64, 10) {
		t.Errorf(tplErrorString, fvo.String(), strconv.FormatInt(math.MinInt64, 10))
		return
	}
	fvo = FilterValue{Source: strconv.FormatInt(math.MaxInt64, 10), Type: TypeInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	// float64
	fvo = FilterValue{Source: "1.8.9", Type: TypeFloat64}
	if err = fvo.Test(); err == nil {
		t.Errorf(tplErrorExpected, err)
		return
	}
	fvo = FilterValue{Source: strconv.FormatFloat(1.2345678, 'E', -1, 64), Type: TypeFloat64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: strconv.FormatFloat(0-3.1415, 'E', -1, 64), Type: TypeFloat64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	// bool
	fvo = FilterValue{Source: "Yes", Type: TypeBool}
	if err = fvo.Test(); err == nil {
		t.Errorf(tplErrorExpected, err)
		return
	}
	fvo = FilterValue{Source: "true", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: "false", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: "1", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: "0", Type: TypeBool}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	// time
	fvo = FilterValue{Source: "18:53", Type: TypeTime}
	if err = fvo.Test(); err == nil {
		t.Errorf(tplErrorExpected, err)
		return
	}
	tmo = time.Now()
	fvo = FilterValue{Source: tmo.Format(time.RFC3339Nano), Type: TypeTime}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: tmo.Format(time.RFC3339), Type: TypeTime}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	// SliceString
	fvo = FilterValue{Source: "", Type: TypeSliceString}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: "one,two;three", Type: TypeSliceString}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	if item, ok := fvo.Value().([]string); !ok || len(item) != 2 || item[0] != "one" || item[1] != "two;three" {
		t.Errorf("Value() = %v, ожидалась %v.", item, []string{"one", "two;three"})
		return
	}
	// SliceInt64
	fvo = FilterValue{Source: "", Type: TypeSliceInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: "3,2,1", Type: TypeSliceInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
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
		t.Errorf(tplErrorNot, err)
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
		t.Errorf(tplErrorIs, err, Errors().FilterValueType.Bind())
		return
	}
	// TypeSliceUint64
	fvo = FilterValue{Source: "", Type: TypeSliceUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	fvo = FilterValue{Source: "7,8,9", Type: TypeSliceUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	if item, ok := fvo.Value().([]uint64); !ok || len(item) != 3 || item[0] != 7 || item[1] != 8 || item[2] != 9 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []uint64{3, 2, 1})
		return
	}
	fvo = FilterValue{Source: "," + strconv.FormatUint(math.MaxUint64, 10), Type: TypeSliceUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	if item, ok := fvo.Value().([]uint64); !ok || len(item) != 2 || item[0] != 0 || item[1] != 18446744073709551615 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []uint64{0, 18446744073709551615})
		return
	}
	fvo = FilterValue{Source: "-1", Type: TypeSliceUint64}
	if err = fvo.Test(); !Errors().FilterValueType.Is(err) {
		t.Errorf(tplErrorIs, err, Errors().FilterValueType.Bind())
		return
	}
	// Срезы с пробелами в значениях.
	fvo = FilterValue{Source: " 3, 2 ,1 ", Type: TypeSliceInt64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	if item, ok := fvo.Value().([]int64); !ok || len(item) != 3 || item[0] != 3 || item[1] != 2 || item[2] != 1 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []int64{3, 2, 1})
		return
	}
	fvo = FilterValue{Source: "7, 8 ,9", Type: TypeSliceUint64}
	if err = fvo.Test(); err != nil {
		t.Errorf(tplErrorNot, err)
		return
	}
	if item, ok := fvo.Value().([]uint64); !ok || len(item) != 3 || item[0] != 7 || item[1] != 8 || item[2] != 9 {
		t.Errorf("Value() = %v, ожидалась %v.", item, []uint64{7, 8, 9})
		return
	}
}

// Тестирование разбора и конвертации значений типов данных полей пакета github.com/webnice/nix.
//
//goland:noinspection ALL
func TestFilterValueParse(t *testing.T) {
	const (
		tplErrorTest  = "%s: FilterValue.Test() = %v, ожидалась ошибка."
		tplErrorNot   = "%s: FilterValue.Test() = %v, ошибка не ожидалась."
		tplErrorValue = "%s: FilterValue.Value() = %v (%T), ожидалось: %v."
		tplErrorNil   = "%s: FilterValue.Value() = %v, ожидалось: nil."
	)
	type test struct {
		Name    string
		Source  string
		Type    FieldType
		Want    any
		IsError bool
	}
	var (
		err   error
		fvo   FilterValue
		got   any
		n     int
		tests = []test{
			// int
			{Name: "int", Source: "42", Type: TypeInt, Want: int(42)},
			{Name: "int-negative", Source: "-42", Type: TypeInt, Want: int(-42)},
			{Name: "int-error", Source: "a42", Type: TypeInt, IsError: true},
			// int8
			{Name: "int8", Source: "127", Type: TypeInt8, Want: int8(127)},
			{Name: "int8-error", Source: "128", Type: TypeInt8, IsError: true},
			// int16
			{Name: "int16", Source: "32767", Type: TypeInt16, Want: int16(32767)},
			{Name: "int16-error", Source: "32768", Type: TypeInt16, IsError: true},
			// int32
			{Name: "int32", Source: "2147483647", Type: TypeInt32, Want: int32(2147483647)},
			{Name: "int32-error", Source: "2147483648", Type: TypeInt32, IsError: true},
			// uint
			{Name: "uint", Source: "42", Type: TypeUint, Want: uint(42)},
			{Name: "uint-error", Source: "-1", Type: TypeUint, IsError: true},
			// uint8
			{Name: "uint8", Source: "255", Type: TypeUint8, Want: uint8(255)},
			{Name: "uint8-error", Source: "256", Type: TypeUint8, IsError: true},
			// uint16
			{Name: "uint16", Source: "65535", Type: TypeUint16, Want: uint16(65535)},
			{Name: "uint16-error", Source: "65536", Type: TypeUint16, IsError: true},
			// uint32
			{Name: "uint32", Source: "4294967295", Type: TypeUint32, Want: uint32(4294967295)},
			{Name: "uint32-error", Source: "4294967296", Type: TypeUint32, IsError: true},
			// float32
			{Name: "float32", Source: "1.5", Type: TypeFloat32, Want: float32(1.5)},
			{Name: "float32-error", Source: "1.8.9", Type: TypeFloat32, IsError: true},
			// duration
			{Name: "duration", Source: "1h30m", Type: TypeDuration, Want: 90 * time.Minute},
			{Name: "duration-ms", Source: "300ms", Type: TypeDuration, Want: 300 * time.Millisecond},
			{Name: "duration-error", Source: "18:53", Type: TypeDuration, IsError: true},
			// bytes, base64
			{Name: "bytes", Source: "aGVsbG8=", Type: TypeBytes, Want: []byte("hello")},
			{Name: "bytes-std", Source: "Pj4+", Type: TypeBytes, Want: []byte(">>>")},
			{Name: "bytes-error", Source: "!!!", Type: TypeBytes, IsError: true},
			// json
			{Name: "json", Source: `{"a":1}`, Type: TypeRawMessage, Want: json.RawMessage(`{"a":1}`)},
			{Name: "json-error", Source: "{bad", Type: TypeRawMessage, IsError: true},
			// ip
			{Name: "ip", Source: "192.168.0.1", Type: TypeIP, Want: net.ParseIP("192.168.0.1")},
			{Name: "ip-error", Source: "not-an-ip", Type: TypeIP, IsError: true},
			// bigint
			{Name: "bigint", Source: "123456789012345678901234567890", Type: TypeBigInt, Want: "123456789012345678901234567890"},
			{Name: "bigint-error", Source: "12a", Type: TypeBigInt, IsError: true},
			// uuid
			{Name: "uuid", Source: "550e8400-e29b-41d4-a716-446655440000", Type: TypeUUID, Want: "550e8400-e29b-41d4-a716-446655440000"},
			{Name: "uuid-no-hyphen", Source: "550e8400e29b41d4a716446655440000", Type: TypeUUID, IsError: true},
			{Name: "uuid-bad-length", Source: "550e8400-e29b-41d4-a716-44665544000", Type: TypeUUID, IsError: true},
			{Name: "uuid-bad-hex", Source: "550e8400-e29b-41d4-a716-44665544000g", Type: TypeUUID, IsError: true},
			// complex
			{Name: "complex64", Source: "1.5+2.5i", Type: TypeComplex64, Want: complex64(1.5 + 2.5i)},
			{Name: "complex64-error", Source: "bad", Type: TypeComplex64, IsError: true},
			{Name: "complex128", Source: "1.5+2.5i", Type: TypeComplex128, Want: complex128(1.5 + 2.5i)},
			{Name: "complex128-error", Source: "bad", Type: TypeComplex128, IsError: true},
			// Типы полей допускающие NULL не конвертируются, значение возвращается как строка.
			{Name: "duration-nil-type", Source: "1h30m", Type: TypeDurationNil, Want: "1h30m"},
			{Name: "uuid-nil-type", Source: "550e8400-e29b-41d4-a716-446655440000", Type: TypeUUIDNil, Want: "550e8400-e29b-41d4-a716-446655440000"},
			{Name: "bytes-nil-type", Source: "aGVsbG8=", Type: TypeBytesNil, Want: "aGVsbG8="},
		}
	)

	for n = range tests {
		fvo = FilterValue{Source: tests[n].Source, Type: tests[n].Type}
		// Тестирование соответствия значения типу.
		if err = fvo.Test(); tests[n].IsError && err == nil {
			t.Errorf(tplErrorTest, tests[n].Name, err)
			continue
		}
		if !tests[n].IsError && err != nil {
			t.Errorf(tplErrorNot, tests[n].Name, err)
			continue
		}
		// Тестирование конвертации значения.
		got = fvo.Value()
		if tests[n].IsError {
			if got != nil {
				t.Errorf(tplErrorNil, tests[n].Name, got)
			}
			continue
		}
		if !equalAnyValue(got, tests[n].Want) {
			t.Errorf(tplErrorValue, tests[n].Name, got, got, tests[n].Want)
		}
	}
}

// Тестирование проверки значения NULL для типов данных полей допускающих NULL.
func TestFilterValueIsValueNilNewTypes(t *testing.T) {
	const tplErrorNil = "%s: IsValueNil() = %t, ожидалось: %t."
	var (
		n     int
		tests = []struct {
			Name  string
			Value FilterValue
			IsNil bool
		}{
			{Name: "int-nil", Value: FilterValue{Source: "nil", Type: TypeIntNil}, IsNil: true},
			{Name: "int8-nil", Value: FilterValue{Source: "nil", Type: TypeInt8Nil}, IsNil: true},
			{Name: "int16-nil", Value: FilterValue{Source: "nil", Type: TypeInt16Nil}, IsNil: true},
			{Name: "int32-nil", Value: FilterValue{Source: "nil", Type: TypeInt32Nil}, IsNil: true},
			{Name: "uint-nil", Value: FilterValue{Source: "nil", Type: TypeUintNil}, IsNil: true},
			{Name: "uint8-nil", Value: FilterValue{Source: "nil", Type: TypeUint8Nil}, IsNil: true},
			{Name: "uint16-nil", Value: FilterValue{Source: "nil", Type: TypeUint16Nil}, IsNil: true},
			{Name: "uint32-nil", Value: FilterValue{Source: "nil", Type: TypeUint32Nil}, IsNil: true},
			{Name: "float32-nil", Value: FilterValue{Source: "nil", Type: TypeFloat32Nil}, IsNil: true},
			{Name: "duration-nil", Value: FilterValue{Source: "nil", Type: TypeDurationNil}, IsNil: true},
			{Name: "bytes-nil", Value: FilterValue{Source: "nil", Type: TypeBytesNil}, IsNil: true},
			{Name: "json-nil", Value: FilterValue{Source: "nil", Type: TypeRawMessageNil}, IsNil: true},
			{Name: "ip-nil", Value: FilterValue{Source: "nil", Type: TypeIPNil}, IsNil: true},
			{Name: "bigint-nil", Value: FilterValue{Source: "nil", Type: TypeBigIntNil}, IsNil: true},
			{Name: "uuid-nil", Value: FilterValue{Source: "nil", Type: TypeUUIDNil}, IsNil: true},
			{Name: "complex64-nil", Value: FilterValue{Source: "nil", Type: TypeComplex64Nil}, IsNil: true},
			{Name: "complex128-nil", Value: FilterValue{Source: "nil", Type: TypeComplex128Nil}, IsNil: true},
			{Name: "uuid", Value: FilterValue{Source: "nil", Type: TypeUUID}, IsNil: false},
			{Name: "duration", Value: FilterValue{Source: "nil", Type: TypeDuration}, IsNil: false},
		}
	)

	for n = range tests {
		if ret := tests[n].Value.IsValueNil(); ret != tests[n].IsNil {
			t.Errorf(tplErrorNil, tests[n].Name, ret, tests[n].IsNil)
		}
	}
}
