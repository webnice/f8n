package f8n

import (
	"net/http"
	"strings"
	"testing"
)

// Тестирование вызова функции с ошибкой в аргументах.
func TestParseRequestWithNil(t *testing.T) {
	var (
		err error
		obj *impl
	)

	// Проверка ошибки паники.
	obj = New().(*impl)
	if err = obj.ParseRequest(nil); err == nil {
		t.Errorf("ParseRequest() == nil, ожидалась ошибка.")
		return
	}
	if !Errors().RequestIsNil.Is(err) {
		t.Errorf("ParseRequest() = %q, ожидалось: %q", err, Errors().RequestIsNil.Bind())
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
	var (
		err      error
		rq       *http.Request
		obj      *impl
		errCount uint64
	)

	// Проверка ошибки паники.
	obj = New().(*impl)
	if rq, err = http.NewRequest("GET", strings.ReplaceAll(data, "\n", ""), nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	if err = obj.ParseRequest(rq, func(a []byte, b error) { errCount++ }); err == nil {
		t.Errorf("ParseRequest() == nil, ожидалась ошибка.")
		return
	}
	if !Errors().MultipleErrorsFound.Is(err) {
		t.Errorf("ParseRequest() = %q, ожидалось: %q", err, Errors().MultipleErrorsFound.Bind())
		return
	}
}
