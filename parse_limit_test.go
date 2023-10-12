// Package f8n
package f8n

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
)

// Тестирование ошибки получения множественного лимита.
func TestParseLimitReceivedMoreThanOne(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		ero Err
		n   int
		ok  bool
	)

	if rq, err = http.NewRequest("GET", "http://localhost?limit=1&limit=2", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) == 0 {
		t.Errorf("ParseLimit() = nil, ожидалась ошибка.")
		return
	}
	for n = range ers {
		if ero, ok = ers[n].Ei.(Err); !ok {
			t.Errorf("ParseLimit() = %q, не верный тип ошибки.", ers[n])
			return
		}
		if ero.Anchor() != Errors().LimitReceivedMoreThanOne().Anchor() ||
			Errors().LimitReceivedMoreThanOne().Code() != 4 {
			t.Errorf("ParseLimit() = %q, ожидалось: %q", ero.Error(), Errors().LimitReceivedMoreThanOne().Error())
			return
		}
	}
}

// Тестирование ошибки конвертации не верного значения лимита.
func TestParseLimitInvalidValue(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		ero Err
		n   int
		ok  bool
	)

	if rq, err = http.NewRequest("GET", "http://localhost?limit=1abc2", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) == 0 {
		t.Errorf("ParseLimit() = nil, ожидалась ошибка.")
		return
	}
	for n = range ers {
		if ero, ok = ers[n].Ei.(Err); !ok {
			t.Errorf("ParseLimit() = %q, не верный тип ошибки.", ers[n])
			return
		}
		if ero.Anchor() != Errors().
			LimitInvalidValue("1abc2", errors.New(`strconv.ParseInt: parsing "1abc2": invalid syntax`)).
			Anchor() {
			t.Errorf("ParseLimit() = %q, ожидалось: %q",
				ero.Error(),
				Errors().
					LimitInvalidValue("1abc2", errors.New(`strconv.ParseInt: parsing "1abc2": invalid syntax`)).
					Error(),
			)
			return
		}
	}
}

// Тестирование отсутствия ошибки при запросе без лимитов.
func TestParseLimitEmpty(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", "http://localhost", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf("ParseLimit() = %v, ошибка не ожидалась.", ers)
		return
	}
}

// Тестирование лимита переданного с одним значением.
func TestParseLimitOneValue(t *testing.T) {
	const limitP, limitN = 1234543210, -123123
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", "http://localhost?limit="+strconv.FormatInt(limitP, 10), nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf("ParseLimit() = %v, ошибка не ожидалась.", ers)
	}
	if obj.Limit != limitP {
		t.Errorf("ParseLimit() = %d, ожидалаось: %d.", obj.Limit, limitP)
	}
	if rq, err = http.NewRequest("GET", "http://localhost?limit="+strconv.FormatInt(limitN, 10), nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf("ParseLimit() = %v, ошибка не ожидалась.", ers)
	}
	if obj.Limit != 0 {
		t.Errorf("ParseLimit() = %d, ожидалаось: %d.", obj.Limit, 0)
	}
}

// Тестирование лимита переданного с двумя значениями.
func TestParseLimit(t *testing.T) {
	const (
		offsetP, limitP   = 1012, 53
		offsetN1, limitN1 = -1, -1
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		ero Err
		n   int
		ok  bool
	)

	if rq, err = http.NewRequest("GET",
		"http://localhost?limit="+strconv.FormatInt(offsetP, 10)+":"+strconv.FormatInt(limitP, 10),
		nil,
	); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf("ParseLimit() = %v, ошибка не ожидалась.", ers)
	}
	if obj.Offset != offsetP || obj.Limit != limitP {
		t.Errorf("ParseLimit() = %d:%d, ожидалаось: %d:%d.", obj.Offset, obj.Limit, offsetP, limitP)
	}
	///
	if rq, err = http.NewRequest("GET",
		"http://localhost?limit="+strconv.FormatInt(offsetN1, 10)+":"+strconv.FormatInt(limitP, 10),
		nil,
	); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseLimit(rq); len(ers) == 0 {
		t.Errorf("ParseLimit() = nil, ожидалась ошибка.")
		return
	}
	for n = range ers {
		if ero, ok = ers[n].Ei.(Err); !ok {
			t.Errorf("ParseLimit() = %q, не верный тип ошибки.", ers[n])
			return
		}
		if ero.Anchor() != Errors().ValueCannotBeNegative(offsetN1).Anchor() {
			t.Errorf("ParseLimit() = %q, ожидалось: %q",
				ero.Error(),
				Errors().ValueCannotBeNegative(offsetN1).Error(),
			)
			return
		}
	}
	///
	if rq, err = http.NewRequest("GET",
		"http://localhost?limit="+strconv.FormatInt(offsetP, 10)+":"+strconv.FormatInt(limitN1, 10),
		nil,
	); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf("ParseLimit() = %v, ошибка не ожидалась.", ers)
	}
	if obj.Offset != offsetP || obj.Limit != 0 {
		t.Errorf("ParseLimit() = %d:%d, ожидалаось: %d:%d.", obj.Offset, obj.Limit, offsetP, 0)
	}
	///
	if rq, err = http.NewRequest("GET", "http://localhost?limit=a1:b1", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if ers = obj.ParseLimit(rq); len(ers) != 2 {
		t.Errorf("ParseLimit() = %v, ошидалось две ошибки.", ers)
	}
}

// Тестирование экспорта лимита.
func TestExportLimit(t *testing.T) {
	var (
		obj *impl
		key string
		val string
	)

	obj = New().(*impl)
	// Пустой лимит.
	if key, val = obj.exportLimit(); key != "" || val != "" {
		t.Errorf(
			"exportLimit() вернул начения: %q = %q, ожидалось: %q, %q",
			key, val,
			"", "",
		)
		return
	}
	// Настроенный лимит.
	obj.Offset, obj.Limit = 187, 43
	if key, val = obj.exportLimit(); key != keyLimit || val != "187:43" {
		t.Errorf(
			"exportLimit() вернул начения: %q = %q, ожидалось: %q, %q",
			key, val,
			keyLimit, "187:43",
		)
		return
	}
	// Сокращённые варианты лимита.
	obj.Offset, obj.Limit = 0, 43
	if key, val = obj.exportLimit(); key != keyLimit || val != ":43" {
		t.Errorf(
			"exportLimit() вернул начения: %q = %q, ожидалось: %q, %q",
			key, val,
			keyLimit, ":43",
		)
		return
	}
	obj.Offset, obj.Limit = 871, 0
	if key, val = obj.exportLimit(); key != keyLimit || val != "871:" {
		t.Errorf(
			"exportLimit() вернул начения: %q = %q, ожидалось: %q, %q",
			key, val,
			keyLimit, "871:",
		)
		return
	}
}
