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
		ero Err
		n   int
		ok  bool
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
		if ero, ok = ers[n].Ei.(Err); !ok {
			t.Errorf("ParseTie() = %q, не верный тип ошибки.", ers[n])
			return
		}
		if ero.Anchor() != Errors().TieModeThanOne().Anchor() ||
			Errors().TieModeThanOne().Code() != 6 {
			t.Errorf("ParseTie() = %q, ожидалось: %q", ero.Error(), Errors().TieModeThanOne().Error())
			return
		}
	}
}
