package f8n

import (
	"net/http"
	"testing"
)

// Тестирование информации о наличии фильтрации в запросе.
func TestImpl_IsFiltration(t *testing.T) {
	const (
		methodGet          = "GET"
		datUrn             = "http://localhost?limit=0:10&filter=domain:ke:*wd*"
		tplErrorRequest    = "Ошибка создания запроса: %s."
		tplErrorParse      = "ParseFilter() ошибок: %d, ошибка не ожидалась."
		tplErrorFilterTrue = "IsFiltration() = %t, ожидалось: %t."
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		ok  bool
	)

	if rq, err = http.NewRequest(methodGet, datUrn, nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ok = obj.IsFiltration(); ok {
		t.Errorf(tplErrorFilterTrue, ok, false)
		return
	}
	if ers = obj.ParseFilter(rq); len(ers) != 0 {
		t.Errorf(tplErrorParse, len(ers))
		return
	}
	if ok = obj.IsFiltration(); !ok {
		t.Errorf(tplErrorFilterTrue, ok, true)
		return
	}
}
