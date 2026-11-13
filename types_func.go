package f8n

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"
)

// String Реализация интерфейса Stringify.
func (fm FilterMethod) String() string { return string(fm) }

// String Реализация интерфейса Stringify.
func (sd ByDirection) String() string { return string(sd) }

// String Реализация интерфейса Stringify.
func (fto FieldType) String() string { return string(fto) }

// String Реализация интерфейса Stringify.
func (fvo FilterValue) String() string { return fvo.Source }

// IsValueNil Возвращает "Истина" если поле может принимать значения NIL и получено значение фильтрации nil.
func (fvo FilterValue) IsValueNil() (ret bool) {
	return isFieldTypeNil(fvo.Type) && strings.EqualFold(fvo.Source, keyNilValue)
}

// isFieldTypeNil Возвращает истину если тип данных поля допускает хранение значения NULL.
func isFieldTypeNil(t FieldType) (ret bool) {
	switch t {
	case TypeStringNil,
		TypeIntNil, TypeInt8Nil, TypeInt16Nil, TypeInt32Nil, TypeInt64Nil,
		TypeUintNil, TypeUint8Nil, TypeUint16Nil, TypeUint32Nil, TypeUint64Nil,
		TypeFloat32Nil, TypeFloat64Nil,
		TypeBoolNil, TypeTimeNil, TypeDurationNil, TypeBytesNil,
		TypeRawMessageNil, TypeIPNil, TypeBigIntNil, TypeUUIDNil,
		TypeComplex64Nil, TypeComplex128Nil:
		ret = true
	}

	return
}

// Value Возвращение значения в типе указанном у значения.
// При ошибке конвертации значения возвращается nil.
func (fvo FilterValue) Value() (ret any) {
	var err error

	if ret, err = fvo.parse(); err != nil {
		ret = nil
	}

	return
}

// Test Тестирование соответствия значения указанному типу.
func (fvo FilterValue) Test() (err error) {
	_, err = fvo.parse()

	return
}

// Разбор и конвертация значения фильтрации в тип данных указанный у значения.
// Типы полей допускающие NULL не конвертируются, значение возвращается как строка.
// Типы полей хранящиеся в базе данных строкой возвращаются строкой после проверки формата значения.
func (fvo FilterValue) parse() (ret any, err error) {
	var (
		bigValue big.Int
		i64      int64
		u64      uint64
		f64      float64
		c128     complex128
		ipValue  net.IP
	)

	switch fvo.Type {
	case TypeString:
		ret = fvo.Source
	case TypeInt:
		if i64, err = strconv.ParseInt(fvo.Source, 10, 64); err == nil {
			ret = int(i64)
		}
	case TypeInt8:
		if i64, err = strconv.ParseInt(fvo.Source, 10, 8); err == nil {
			ret = int8(i64)
		}
	case TypeInt16:
		if i64, err = strconv.ParseInt(fvo.Source, 10, 16); err == nil {
			ret = int16(i64)
		}
	case TypeInt32:
		if i64, err = strconv.ParseInt(fvo.Source, 10, 32); err == nil {
			ret = int32(i64)
		}
	case TypeInt64:
		if i64, err = strconv.ParseInt(fvo.Source, 10, 64); err == nil {
			ret = i64
		}
	case TypeUint:
		if u64, err = strconv.ParseUint(fvo.Source, 10, 64); err == nil {
			ret = uint(u64)
		}
	case TypeUint8:
		if u64, err = strconv.ParseUint(fvo.Source, 10, 8); err == nil {
			ret = uint8(u64)
		}
	case TypeUint16:
		if u64, err = strconv.ParseUint(fvo.Source, 10, 16); err == nil {
			ret = uint16(u64)
		}
	case TypeUint32:
		if u64, err = strconv.ParseUint(fvo.Source, 10, 32); err == nil {
			ret = uint32(u64)
		}
	case TypeUint64:
		if u64, err = strconv.ParseUint(fvo.Source, 10, 64); err == nil {
			ret = u64
		}
	case TypeFloat32:
		if f64, err = strconv.ParseFloat(fvo.Source, 32); err == nil {
			ret = float32(f64)
		}
	case TypeFloat64:
		if f64, err = strconv.ParseFloat(fvo.Source, 64); err == nil {
			ret = f64
		}
	case TypeBool:
		if ret, err = strconv.ParseBool(fvo.Source); err != nil {
			ret = nil
		}
	case TypeTime:
		if ret, err = time.Parse(time.RFC3339, fvo.Source); err != nil {
			ret = nil
		}
	case TypeDuration:
		if ret, err = time.ParseDuration(fvo.Source); err != nil {
			ret = nil
		}
	case TypeBytes:
		if ret, err = decodeBase64(fvo.Source); err != nil {
			ret = nil
		}
	case TypeRawMessage:
		if !json.Valid([]byte(fvo.Source)) {
			err = strconv.ErrSyntax
			break
		}
		ret = json.RawMessage(fvo.Source)
	case TypeIP:
		if ipValue = net.ParseIP(fvo.Source); ipValue == nil {
			err = strconv.ErrSyntax
			break
		}
		ret = ipValue
	case TypeBigInt:
		if _, ok := bigValue.SetString(fvo.Source, 10); !ok {
			err = strconv.ErrSyntax
			break
		}
		ret = fvo.Source
	case TypeUUID:
		if !isUUIDFormat(fvo.Source) {
			err = strconv.ErrSyntax
			break
		}
		ret = fvo.Source
	case TypeComplex64:
		if c128, err = strconv.ParseComplex(fvo.Source, 64); err == nil {
			ret = complex64(c128)
		}
	case TypeComplex128:
		if ret, err = strconv.ParseComplex(fvo.Source, 128); err != nil {
			ret = nil
		}
	case TypeSliceString:
		ret, err = fvo.SliceString()
	case TypeSliceInt64:
		ret, err = fvo.SliceInt64()
	case TypeSliceUint64:
		ret, err = fvo.SliceUint64()
	default:
		// Типы полей допускающие NULL и не известные типы, значение возвращается как строка.
		ret = fvo.Source
	}

	return
}

// SliceString Представление значения в виде среза строк.
func (fvo FilterValue) SliceString() (ret []string, err error) {
	const sep = ","
	var n int

	ret = strings.Split(fvo.Source, sep)
	for n = range ret {
		ret[n] = strings.TrimSpace(ret[n])
	}

	return
}

// SliceInt64 Представления значения в виде среза int64.
func (fvo FilterValue) SliceInt64() (ret []int64, err error) {
	var (
		n   int
		str []string
		tmp int64
	)

	str = strings.Split(fvo.Source, sepInNi)
	ret = make([]int64, 0, len(str))
	for n = range str {
		str[n] = strings.TrimSpace(str[n])
		if tmp = 0; str[n] != "" {
			if tmp, err = strconv.ParseInt(str[n], 10, 64); err != nil {
				err = Errors().FilterValueType.Bind(str[n], fvo.Type.String(), err)
				return
			}
		}
		ret = append(ret, tmp)
	}

	return
}

// SliceUint64 Представления значения в виде среза uint64.
func (fvo FilterValue) SliceUint64() (ret []uint64, err error) {
	var (
		n   int
		str []string
		tmp uint64
	)

	str = strings.Split(fvo.Source, sepInNi)
	ret = make([]uint64, 0, len(str))
	for n = range str {
		str[n] = strings.TrimSpace(str[n])
		if tmp = 0; str[n] != "" {
			if tmp, err = strconv.ParseUint(str[n], 10, 64); err != nil {
				err = Errors().FilterValueType.Bind(str[n], fvo.Type.String(), err)
				return
			}
		}
		ret = append(ret, tmp)
	}

	return
}

// Декодирование base64 значения, поддерживаются стандартная и URL безопасная кодировки.
func decodeBase64(s string) (ret []byte, err error) {
	if ret, err = base64.URLEncoding.DecodeString(s); err == nil {
		return
	}
	ret, err = base64.StdEncoding.DecodeString(s)

	return
}

// Проверка формата UUID, строка формата xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func isUUIDFormat(s string) (ret bool) {
	const (
		tplHyphen                    = "-"
		seg1, seg2, seg3, seg4, seg5 = 8, 4, 4, 4, 12
	)
	var (
		seg []string
		n   int
	)

	seg = strings.Split(s, tplHyphen)
	if len(seg) != 5 {
		return
	}
	if len(seg[0]) != seg1 || len(seg[1]) != seg2 || len(seg[2]) != seg3 ||
		len(seg[3]) != seg4 || len(seg[4]) != seg5 {
		return
	}
	for n = range seg {
		if !isHexString(seg[n]) {
			return
		}
	}
	ret = true

	return
}

// Проверка строки, все символы строки должны быть шестнадцатеричными числами.
func isHexString(s string) (ret bool) {
	var n int

	ret = true
	for n = 0; n < len(s); n++ {
		switch c := s[n]; {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		case c >= 'A' && c <= 'F':
		default:
			ret = false
			return
		}
	}

	return
}

// Разбор строки в тип направления сортировки.
func parseByDirection(s string) (ret ByDirection) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case byAsc.String():
		ret = byAsc
	case byDesc.String():
		ret = byDesc
	default:
		ret = byUnknown
	}

	return
}

// Разбор строки в тип фильтрации.
func parseFilterMethod(s string) (ret FilterMethod) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case filterEquivalent.String():
		ret = filterEquivalent
	case filterLessThan.String():
		ret = filterLessThan
	case filterLessThanOrEquivalent.String():
		ret = filterLessThanOrEquivalent
	case filterGreaterThan.String():
		ret = filterGreaterThan
	case filterGreaterThanOrEquivalent.String():
		ret = filterGreaterThanOrEquivalent
	case filterNotEquivalent.String():
		ret = filterNotEquivalent
	case filterLikeThan.String():
		ret = filterLikeThan
	case filterNotLikeThan.String():
		ret = filterNotLikeThan
	case filterIn.String():
		ret = filterIn
	case filterNotIn.String():
		ret = filterNotIn
	default:
		ret = filterUnknown
	}

	return
}

// String Реализация интерфейса Stringify.
func (tmo TieMode) String() string { return string(tmo) }

// Разбор строки в тип режима фильтрации.
func parseTie(s string) (ret TieMode) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case tieAnd.String():
		ret = tieAnd
	case tieOr.String():
		ret = tieOr
	default:
		ret = tieUnknown
	}

	return
}
