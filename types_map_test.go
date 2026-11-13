package f8n

import (
	"net/http"
	"strings"
	"testing"
)

// Тестирование отсутствия именованных фильтров.
func TestMapLoadAndParseFilterCalledByNameWasNotFound(t *testing.T) {
	const (
		datRequest = `http://localhost
?map=filter1:or:filter2
&filter1=id:eq:1
`
		tplErrorRequest  = "Ошибка создания запроса: %s."
		tplErrorEmpty    = "ParseMap() == nil, ожидались ошибки."
		tplErrorNotFound = "ParseMap() = %q, ожидалось: %q"
	)
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		n   int
	)

	if rq, err = http.NewRequest("GET", strings.ReplaceAll(datRequest, "\n", ""), nil); err != nil {
		t.Errorf(tplErrorRequest, err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) <= 0 {
		t.Errorf(tplErrorEmpty)
		return
	}
	for n = range ers {
		if !Errors().FilterCalledByNameWasNotFound.Is(ers[n].Ei) {
			t.Errorf(tplErrorNotFound, ers[n].Ei, Errors().FilterCalledByNameWasNotFound.Bind())
			return
		}
	}
}

// Тестирование экспорта пустой простой фильтрации.
func TestExportAsStringEmptyFilter(t *testing.T) {
	const tplErrorExport = "Filter = %q, ожидались пустая простая фильтрация."
	var (
		obj Filter
		tmp string
	)

	if tmp = obj.exportAsString(); tmp != "" {
		t.Errorf(tplErrorExport, tmp)
		return
	}
}

// Тестирование экспорта пустой сортировки.
func TestExportAsStringEmptyDirection(t *testing.T) {
	const tplErrorExport = "Direction = %q, ожидались пустая сортировка."
	var (
		obj Direction
		tmp string
	)

	if tmp = obj.exportAsString(); tmp != "" {
		t.Errorf(tplErrorExport, tmp)
		return
	}
}

// Тестирование экспорта пустой карты сложной фильтрации.
func TestExportAsStringEmptyMap(t *testing.T) {
	const tplErrorExport = "Map = %q, ожидались пустая карта сложной фильтрации."
	var (
		obj *Map
		tmp string
	)

	obj = new(Map)
	if tmp, _ = obj.exportAsString(); tmp != "" {
		t.Errorf(tplErrorExport, tmp)
		return
	}
}
