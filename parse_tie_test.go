package f8n

import (
	"net/http"
	"testing"
)

func TestImpl_ParseTieAnd(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		tmp string
	)

	if rq, err = http.NewRequest("GET", "http://localhost?tie=AND", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseTie(rq); len(ers) != 0 {
		t.Errorf("ParseTie() ошибок: %d, ошибка не ожидалась.", len(ers))
		return
	}
	if tmp = obj.Tie.String(); tmp != "and" {
		t.Errorf("Tie.String() = %s, ожидалась: %q.", tmp, "and")
		return
	}
}

func TestImpl_ParseTieOr(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		tmp string
	)

	if rq, err = http.NewRequest("GET", "http://localhost?tie=OR", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseTie(rq); len(ers) != 0 {
		t.Errorf("ParseTie() ошибок: %d, ошибка не ожидалась.", len(ers))
		return
	}
	if tmp = obj.Tie.String(); tmp != "or" {
		t.Errorf("Tie.String() = %s, ожидалась: %q.", tmp, "or")
		return
	}
}

// Тестирование ошибки получения множественного TIE.
func TestImpl_ParseTieThanOne(t *testing.T) {
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
		t.Errorf("ParseTie() ошибок: %d, ожидалась ошибка.", len(ers))

		return
	}
	for n = range ers {
		if !Errors().TieModeThanOne.Is(ers[n].Ei) || Errors().TieModeThanOne.CodeI().Get() != 6 {
			t.Errorf("ParseTie() = %q, ожидалось: %q", ers[n].Ei.Error(), Errors().TieModeThanOne.Bind())
			return
		}
	}
}

func TestImpl_ParseTieInvalidValue(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", "http://localhost?tie=die", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseTie(rq); len(ers) == 0 {
		t.Errorf("ParseTie() ошибок: %d, ожидалась ошибка.", len(ers))
		return
	}
	for n = range ers {
		if !Errors().TieModeInvalidValue.Is(ers[n].Ei) || Errors().TieModeInvalidValue.CodeI().Get() != 7 {
			t.Errorf("ParseTie() = %q, ожидалось: %q", ers[n].Ei.Error(), Errors().TieModeThanOne.Bind())
			return
		}
	}
}
