package f8n

import (
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
	switch fvo.Type {
	case TypeStringNil, TypeUint64Nil, TypeInt64Nil, TypeFloat64Nil, TypeBoolNil, TypeTimeNil:
		ret = strings.ToLower(fvo.Source) == keyNilValue
	}

	return
}

// Value Возвращения значения в типе указанном у значения.
func (fvo FilterValue) Value() (ret any) {
	var err error

	switch fvo.Type {
	case TypeString:
		ret = fvo.Source
	case TypeUint64:
		ret, err = strconv.ParseUint(fvo.Source, 10, 64)
	case TypeInt64:
		ret, err = strconv.ParseInt(fvo.Source, 10, 64)
	case TypeFloat64:
		ret, err = strconv.ParseFloat(fvo.Source, 64)
	case TypeBool:
		ret, err = strconv.ParseBool(fvo.Source)
	case TypeTime:
		ret, err = time.Parse(time.RFC3339, fvo.Source)
	case TypeSliceString:
		ret, err = fvo.SliceString()
	case TypeSliceInt64:
		ret, err = fvo.SliceInt64()
	case TypeSliceUint64:
		ret, err = fvo.SliceUint64()
	default:
		ret = fvo.Source
	}
	if err != nil {
		ret = nil
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
		if tmp = 0; str[n] != "" {
			if tmp, err = strconv.ParseInt(str[n], 10, 64); err != nil {
				err = Errors().FilterValueType.Bind(str[n], "int64", err)
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
		if tmp = 0; str[n] != "" {
			if tmp, err = strconv.ParseUint(str[n], 10, 64); err != nil {
				err = Errors().FilterValueType.Bind(str[n], "uint64", err)
				return
			}
		}
		ret = append(ret, tmp)
	}

	return
}

// Test Тестирование соответствия значения указанному типу.
func (fvo FilterValue) Test() (err error) {
	switch fvo.Type {
	case TypeUint64:
		_, err = strconv.ParseUint(fvo.Source, 10, 64)
	case TypeInt64:
		_, err = strconv.ParseInt(fvo.Source, 10, 64)
	case TypeFloat64:
		_, err = strconv.ParseFloat(fvo.Source, 64)
	case TypeBool:
		_, err = strconv.ParseBool(fvo.Source)
	case TypeTime:
		_, err = time.Parse(time.RFC3339, fvo.Source)
	case TypeSliceString:
		_, err = fvo.SliceString()
	case TypeSliceInt64:
		_, err = fvo.SliceInt64()
	case TypeSliceUint64:
		_, err = fvo.SliceUint64()
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
