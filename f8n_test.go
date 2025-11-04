package f8n

import (
	"net/http"
	"testing"
)

func TestImpl_IsFiltration(t *testing.T) {
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		ok  bool
	)

	if rq, err = http.NewRequest("GET", "http://localhost?limit=0:10&filter=domain:ke:*wd*", nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ok = obj.IsFiltration(); ok {
		t.Errorf("IsFiltration() = %t, ожидалось: %t.", ok, false)
		return
	}
	if ers = obj.ParseFilter(rq); len(ers) != 0 {
		t.Errorf("ParseFilter() ошибок: %d, ошибка не ожидалась.", len(ers))
		return
	}
	if ok = obj.IsFiltration(); !ok {
		t.Errorf("IsFiltration() = %t, ожидалось: %t.", ok, true)
		return
	}
}
