package f8n

import (
	"net/http"
	"strings"
	"testing"
)

// Тестирование вызова функции с ошибкой в аргументах.
func TestParseRequestWithNil(t *testing.T) {
	const (
		tplErrorNil = "ParseRequest() == nil, ожидалась ошибка."
		tplErrorIs  = "ParseRequest() = %q, ожидалось: %q"
	)
	var (
		err error
		obj *impl
	)

	// Проверка ошибки паники.
	obj = New().(*impl)
	if err = obj.ParseRequest(nil); err == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if !Errors().RequestIsNil.Is(err) {
		t.Errorf(tplErrorIs, err, Errors().RequestIsNil.Bind())
		return
	}
}

// Тестирование множественных ошибок.
func TestMultipleErrorsFound(t *testing.T) {
	const data = `http://localhost
?map=(group4:and:group5):and:(group1:and:group2:or:group3):or:(group4:or:group5):or:(group4:or:group5)
&map=group4:and:group5:and:
&map=(group5:or:(group1:or:():and:group5):and:group5):and:group4
&group1=field1:eq:value1
&group2=field2:ke:value2
&group3=field3:ke:value3
&group4=field1:ke:value4
&group5=field2:ke:value5
&limit=1&limit=2
&by=zzzz-1
&filter=aaa:aaa:aaa
`
	const (
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorNil     = "ParseRequest() == nil, ожидалась ошибка."
		tplErrorIs      = "ParseRequest() = %q, ожидалось: %q"
	)
	var (
		err      error
		rq       *http.Request
		obj      *impl
		errCount uint64
	)

	// Проверка ошибки паники.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", strings.ReplaceAll(data, "\n", ""), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if err = obj.ParseRequest(rq, func(a []byte, b error) { errCount++ }); err == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if !Errors().MultipleErrorsFound.Is(err) {
		t.Errorf(tplErrorIs, err, Errors().MultipleErrorsFound.Bind())
		return
	}
}

// Тестирование одиночной ошибки разбора устаревшего режима фильтрации tie.
func TestParseRequestTieError(t *testing.T) {
	const (
		datUrn          = `http://localhost?tie=die`
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorNil     = "ParseRequest() == nil, ожидалась ошибка."
		tplErrorIs      = "ParseRequest() = %q, ожидалось: %q"
		tplErrorCount   = "ParseRequest() количество вызовов errorFn = %d, ожидался: %d."
	)
	var (
		err      error
		rq       *http.Request
		obj      *impl
		errCount uint64
	)

	// Проверка одиночной ошибки режима фильтрации.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	if err = obj.ParseRequest(rq, func(a []byte, b error) { errCount++ }); err == nil {
		t.Errorf(tplErrorNil)
		return
	}
	if !Errors().TieModeInvalidValue.Is(err) {
		t.Errorf(tplErrorIs, err, Errors().TieModeInvalidValue.Bind())
		return
	}
	if errCount != 1 {
		t.Errorf(tplErrorCount, errCount, 1)
		return
	}
}
