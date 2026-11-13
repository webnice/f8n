package f8n

import (
	"net/http"
	"strconv"
	"testing"
)

// Тестирование ошибки получения множественного лимита.
func TestParseLimitReceivedMoreThanOne(t *testing.T) {
	const (
		datUrn          = "http://localhost?limit=1&limit=2"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorEmpty   = "ParseLimit() = nil, ожидалась ошибка."
		tplErrorLimit   = "ParseLimit() = %q, ожидалось: %q"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) == 0 {
		t.Errorf(tplErrorEmpty)
		return
	}
	for n = range ers {
		if !Errors().LimitReceivedMoreThanOne.Is(ers[n].Ei) ||
			Errors().LimitReceivedMoreThanOne.CodeI().Get() != 4 {
			t.Errorf(tplErrorLimit, ers[n].Ei, Errors().LimitReceivedMoreThanOne.Bind())
			return
		}
	}
}

// Тестирование ошибки конвертации не верного значения лимита.
func TestParseLimitInvalidValue(t *testing.T) {
	const (
		datUrn          = "http://localhost?limit=1abc2"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorEmpty   = "ParseLimit() = nil, ожидалась ошибка."
		tplErrorLimit   = "ParseLimit() = %q, ожидалось: %q"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) == 0 {
		t.Errorf(tplErrorEmpty)
		return
	}
	for n = range ers {
		if !Errors().LimitInvalidValue.Is(ers[n].Ei) {
			t.Errorf(tplErrorLimit, ers[n].Ei, Errors().LimitInvalidValue.Bind())
			return
		}
	}
}

// Тестирование отсутствия ошибки при запросе без лимитов.
func TestParseLimitEmpty(t *testing.T) {
	const (
		datUrn          = "http://localhost"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorLimit   = "ParseLimit() = %v, ошибка не ожидалась."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf(tplErrorLimit, ers)
		return
	}
}

// Тестирование лимита переданного с одним значением.
func TestParseLimitOneValue(t *testing.T) {
	const (
		limitP, limitN    = 1234543210, -123123
		tplErrorRequest   = "Ошибка создания запроса: %s."
		tplErrorLimit     = "ParseLimit() = %v, ошибка не ожидалась."
		tplErrorLimitWant = "ParseLimit() = %d, ожидалось: %d."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
	)

	if rq, err = http.NewRequest("GET", "http://localhost?limit="+strconv.FormatInt(limitP, 10), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf(tplErrorLimit, ers)
	}
	if obj.Limit != limitP {
		t.Errorf(tplErrorLimitWant, obj.Limit, limitP)
	}
	if rq, err = http.NewRequest("GET", "http://localhost?limit="+strconv.FormatInt(limitN, 10), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf(tplErrorLimit, ers)
	}
	if obj.Limit != 0 {
		t.Errorf(tplErrorLimitWant, obj.Limit, 0)
	}
}

// Тестирование лимита переданного с двумя значениями.
//
//goland:noinspection DuplicatedCode
func TestParseLimit(t *testing.T) {
	const (
		offsetP, limitP   = 1012, 53
		offsetN1, limitN1 = -1, -1
		tplErrorRequest   = "Ошибка создания запроса: %s."
		tplErrorLimit     = "ParseLimit() = %v, ошибка не ожидалась."
		tplErrorEmpty     = "ParseLimit() = nil, ожидалась ошибка."
		tplErrorLimitWant = "ParseLimit() = %d:%d, ожидалось: %d:%d."
		tplErrorTwo       = "ParseLimit() = %v, ожидалось две ошибки."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET",
		"http://localhost?limit="+strconv.FormatInt(offsetP, 10)+":"+strconv.FormatInt(limitP, 10),
		nil,
	); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf(tplErrorLimit, ers)
	}
	if obj.Offset != offsetP || obj.Limit != limitP {
		t.Errorf(tplErrorLimitWant, obj.Offset, obj.Limit, offsetP, limitP)
	}
	/// Отрицательное значение позиции выборки.
	if rq, err = http.NewRequest("GET",
		"http://localhost?limit="+strconv.FormatInt(offsetN1, 10)+":"+strconv.FormatInt(limitP, 10),
		nil,
	); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseLimit(rq); len(ers) == 0 {
		t.Errorf(tplErrorEmpty)
		return
	}
	for n = range ers {
		if !Errors().ValueCannotBeNegative.Is(ers[n].Ei) {
			t.Errorf("ParseLimit() = %q, ожидалось: %q", ers[n].Ei, Errors().ValueCannotBeNegative.Bind())
			return
		}
	}
	/// Отрицательное значение размера выборки.
	if rq, err = http.NewRequest("GET",
		"http://localhost?limit="+strconv.FormatInt(offsetP, 10)+":"+strconv.FormatInt(limitN1, 10),
		nil,
	); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseLimit(rq); len(ers) > 0 {
		t.Errorf(tplErrorLimit, ers)
	}
	if obj.Offset != offsetP || obj.Limit != 0 {
		t.Errorf(tplErrorLimitWant, obj.Offset, obj.Limit, offsetP, 0)
	}
	/// Ошибки конвертации обоих значений.
	if rq, err = http.NewRequest("GET", "http://localhost?limit=a1:b1", nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if ers = obj.ParseLimit(rq); len(ers) != 2 {
		t.Errorf(tplErrorTwo, ers)
	}
}

// Тестирование экспорта лимита.
func TestExportLimit(t *testing.T) {
	const (
		tplErrorExport = "exportLimit() вернул значения: %q = %q, ожидалось: %q, %q"
	)
	var (
		obj *impl
		key string
		val string
	)

	obj = New().(*impl)
	// Пустой лимит.
	if key, val = obj.exportLimit(); key != "" || val != "" {
		t.Errorf(tplErrorExport, key, val, "", "")
		return
	}
	// Настроенный лимит.
	obj.Offset, obj.Limit = 187, 43
	if key, val = obj.exportLimit(); key != keyLimit || val != "187:43" {
		t.Errorf(tplErrorExport, key, val, keyLimit, "187:43")
		return
	}
	// Сокращённые варианты лимита.
	obj.Offset, obj.Limit = 0, 43
	if key, val = obj.exportLimit(); key != keyLimit || val != ":43" {
		t.Errorf(tplErrorExport, key, val, keyLimit, ":43")
		return
	}
	obj.Offset, obj.Limit = 871, 0
	if key, val = obj.exportLimit(); key != keyLimit || val != "871:" {
		t.Errorf(tplErrorExport, key, val, keyLimit, "871:")
		return
	}
}
