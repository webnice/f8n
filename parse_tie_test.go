package f8n

import (
	"net/http"
	"testing"
)

// Тестирование загрузки устаревшего режима фильтрации tie, режим "И".
func TestImpl_ParseTieAnd(t *testing.T) {
	const (
		datUrn          = "http://localhost?tie=AND"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseTie() ошибок: %d, ошибка не ожидалась."
		tplErrorTie     = "Tie.String() = %s, ожидалась: %q."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		tmp string
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseTie(rq); len(ers) != 0 {
		t.Errorf(tplErrorParse, len(ers))
		return
	}
	if tmp = obj.Tie.String(); tmp != "and" {
		t.Errorf(tplErrorTie, tmp, "and")
		return
	}
}

// Тестирование загрузки устаревшего режима фильтрации tie, режим "ИЛИ".
func TestImpl_ParseTieOr(t *testing.T) {
	const (
		datUrn          = "http://localhost?tie=OR"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorParse   = "ParseTie() ошибок: %d, ошибка не ожидалась."
		tplErrorTie     = "Tie.String() = %s, ожидалась: %q."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		tmp string
	)

	if rq, err = http.NewRequest("GET", datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseTie(rq); len(ers) != 0 {
		t.Errorf(tplErrorParse, len(ers))
		return
	}
	if tmp = obj.Tie.String(); tmp != "or" {
		t.Errorf(tplErrorTie, tmp, "or")
		return
	}
}

// Тестирование ошибки получения множественного TIE.
func TestImpl_ParseTieThanOne(t *testing.T) {
	const (
		datUrn          = "http://localhost?tie=or&tie=and"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorEmpty   = "ParseTie() ошибок: %d, ожидалась ошибка."
		tplErrorTie     = "ParseTie() = %q, ожидалось: %q"
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
	if ers = obj.ParseTie(rq); len(ers) == 0 {
		t.Errorf(tplErrorEmpty, len(ers))
		return
	}
	for n = range ers {
		if !Errors().TieModeThanOne.Is(ers[n].Ei) || Errors().TieModeThanOne.CodeI().Get() != 6 {
			t.Errorf(tplErrorTie, ers[n].Ei.Error(), Errors().TieModeThanOne.Bind())
			return
		}
	}
}

// Тестирование ошибки недопустимого значения режима фильтрации tie.
func TestImpl_ParseTieInvalidValue(t *testing.T) {
	const (
		datUrn          = "http://localhost?tie=die"
		tplErrorRequest = "Ошибка создания запроса: %s."
		tplErrorEmpty   = "ParseTie() ошибок: %d, ожидалась ошибка."
		tplErrorTie     = "ParseTie() = %q, ожидалось: %q"
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
	if ers = obj.ParseTie(rq); len(ers) == 0 {
		t.Errorf(tplErrorEmpty, len(ers))
		return
	}
	for n = range ers {
		if !Errors().TieModeInvalidValue.Is(ers[n].Ei) || Errors().TieModeInvalidValue.CodeI().Get() != 7 {
			t.Errorf(tplErrorTie, ers[n].Ei.Error(), Errors().TieModeThanOne.Bind())
			return
		}
	}
}
