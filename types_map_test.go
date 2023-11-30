package f8n

import (
	"net/http"
	"strings"
	"testing"
)

// Тестирование отсутствия именованных фильтров.
func TestMapLoadAndParseFilterCalledByNameWasNotFound(t *testing.T) {
	const data = `http://localhost
?map=filter1:or:filter2
&filter1=id:eq:1
`
	var (
		err error
		rq  *http.Request
		obj *impl
		ers []*ParseError
		ero Err
		n   int
		ok  bool
	)

	if rq, err = http.NewRequest("GET", strings.ReplaceAll(data, "\n", ""), nil); err != nil {
		t.Errorf("Ошибка создания запроса: %s", err)
	}
	obj = New().(*impl)
	if ers = obj.ParseMap(rq); len(ers) <= 0 {
		t.Errorf("ParseMap() == nil, ожидались ошибки.")
		return
	}
	for n = range ers {
		if ero, ok = ers[n].Ei.(Err); !ok {
			t.Errorf("ParseMap() = %q, не верный тип ошибки.", ers[n])
			return
		}
		if ero.Anchor() != Errors().FilterCalledByNameWasNotFound("").Anchor() {
			t.Errorf("ParseMap() = %q, ожидалось: %q",
				ero.Error(),
				Errors().FilterCalledByNameWasNotFound("").Error(),
			)
			return
		}
	}
}

// Тестирование экспорта пустой простой фильтрации.
func TestExportAsStringEmptyFilter(t *testing.T) {
	var (
		obj Filter
		tmp string
	)

	if tmp = obj.exportAsString(); tmp != "" {
		t.Errorf("Filter = %q, ожидались пустая простая фильтрация.", tmp)
		return
	}
}

// Тестирование экспорта пустой сортировки.
func TestExportAsStringEmptyDirection(t *testing.T) {
	var (
		obj Direction
		tmp string
	)

	if tmp = obj.exportAsString(); tmp != "" {
		t.Errorf("Direction = %q, ожидались пустая сортировка.", tmp)
		return
	}
}

// Тестирование экспорта пустой карты сложной фильтрации.
func TestExportAsStringEmptyMap(t *testing.T) {
	var (
		obj *Map
		tmp string
	)

	obj = new(Map)
	if tmp, _ = obj.exportAsString(); tmp != "" {
		t.Errorf("Map = %q, ожидались пустая карта сложной фильтрации.", tmp)
		return
	}
}
