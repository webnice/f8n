package f8n

import (
	"net/http"
	"testing"
)

// Тестирование ошибки получения множественного TIE.
func TestImpl_ParseTie(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", "http://localhost?tie=or&tie=and", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseTie(rq); len(ers) == 0 {
		t.Errorf("ParseTie() = nil, ожидалась ошибка.")
		return
	}
	for n = range ers {
		if !Errors().TieModeThanOne.Is(ers[n].Ei) || Errors().TieModeThanOne.CodeI().Get() != 6 {
			t.Errorf("ParseTie() = %q, ожидалось: %q", ers[n].Ei.Error(), Errors().TieModeThanOne.Bind())
			return
		}
	}
}
