package f8n

import (
	"net/http"
	"strings"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// ParseFilter Загрузка простой фильтрации.
func (f8n *impl) ParseFilter(rq *http.Request) (ret []*ParseError) {
	var (
		err     []*ParseError
		filters []string
		filter  Filter
		n       int
	)

	if filters = rq.URL.Query()[keyFilter]; len(filters) == 0 {
		return
	}
	for n = range filters {
		if filter, err = f8n.parseFilterSimple(filters[n]); len(err) > 0 {
			ret = append(ret, err...)
			continue
		}
		f8n.Filter = append(f8n.Filter, filter)
	}

	return
}

// Разбор значения простой фильтрации.
func (f8n *impl) parseFilterSimple(s string) (ret Filter, err []*ParseError) {
	const nSection = 3
	var (
		e     error
		tmp   []string
		ero   *ParseError
		n     int
		found bool
		tpe   FieldType
	)

	// Проверка формата.
	if tmp = strings.SplitN(s, keyDelimiter, nSection); len(tmp) != nSection {
		ero = &ParseError{Ei: f8n.Errors().FilterFormat.Bind()}
		ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
			Field:      keyFilter,
			FieldValue: s,
			Message:    ero.Ei.Error(),
		})
		err = append(err, ero)
		return
	}
	// Проверка константы метода фильтрации.
	if ret.Method = parseFilterMethod(tmp[1]); ret.Method == filterUnknown {
		ero = &ParseError{Ei: f8n.Errors().FilterUnknownMethod.Bind(tmp[1])}
		ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
			Field:      keyFilter,
			FieldValue: tmp[1],
			Message:    ero.Ei.Error(),
		})
		err = append(err, ero)
		return
	}
	// Если заданы поля, проверка имени поля фильтрации.
	if ret.Field = tmp[0]; len(f8n.FieldOnly) > 0 {
		for n = range f8n.FieldOnly {
			if f8n.FieldOnly[n] == ret.Field {
				found = true
			}
		}
		if !found {
			ero = &ParseError{Ei: f8n.Errors().FilterUnknownField.Bind(ret.Field)}
			ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
				Field:      keyFilter,
				FieldValue: ret.Field,
				Message:    ero.Ei.Error(),
			})
			err = append(err, ero)
			return
		}
	}
	// Присвоения значения фильтрации.
	ret.Value = FilterValue{Source: tmp[2], Type: TypeString}
	if tpe, found = f8n.Datatype[ret.Field]; found {
		ret.Value.Type = tpe
	}
	// Тестирование значения.
	if e = ret.Value.Test(); e != nil {
		ero = &ParseError{Ei: f8n.Errors().FilterValueType.Bind(ret.Value.Source, ret.Value.Type, e)}
		ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
			Field:      keyFilter,
			FieldValue: s,
			Message:    ero.Ei.Error(),
		})
		err = append(err, ero)
		return
	}

	return
}

// Экспорт простой фильтрации.
func (f8n *impl) exportFilter() (key string, val []string) {
	var n int

	if len(f8n.Filter) == 0 {
		return
	}
	key, val = keyFilter, make([]string, 0, len(f8n.Filter))
	for n = range f8n.Filter {
		val = append(val, f8n.Filter[n].exportAsString())
	}

	return
}
