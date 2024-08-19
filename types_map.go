package f8n

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// Origin Структура константы описывающей происхождение тега.
type Origin string

// Map Структура настроек сложной фильтрации.
type Map struct {
	parent   *impl   ``               // Ссылка на объект пакета.
	TagBegin string  `json:"beg"`     // Для парных тегов - открывающий тег, для одиночных весь тег целиком.
	TagEnd   string  `json:"end"`     // Для парных тегов - закрывающий тег, для одиночных весь тег целиком.
	Origin   Origin  `json:"origin"`  // Константа описывающая сущность тега.
	Content  string  `json:"content"` // Контент тега - наименование фильтрации.
	Filter   *Filter `json:"filter"`  // Контент фильтра, разобранный в объект описания простой фильтрации.
	Size     int     `json:"size"`    // Общий размер тега в байтах.
	Node     []*Map  `json:"node"`    // Массив узлов, в порядке следования, сложенных в текущий узел.
}

func init() {
	// Сортировка констант справочника тегов в обратном порядке по длине тега.
	// Чем длиннее тег, тем он более релевантен при совпадении поиска, по отношению к короткому тегу.
	// Пример: list и list-item, при поиске совпадения первым должен сравниваться list-item.
	sort.SliceStable(OriginAll, func(i int, j int) bool {
		return len(OriginAll[i].String()) > len(OriginAll[j].String())
	})
}

// String Реализация интерфейса Stringify.
func (origin Origin) String() string { return string(origin) }

// Экспорт значения настроек сортировки.
func (direction Direction) exportAsString() (ret string) {
	if direction.Field == "" {
		return
	}
	ret = fmt.Sprintf("%s%s%s", direction.Field, keyDelimiter, direction.By.String())

	return
}

// Рекурсивная установка parent для всех узлов карты.
// Применяется после импорта настроек фильтрации.
func (mp *Map) setParent(f8n *impl) {
	var n int

	mp.parent = f8n
	for n = range mp.Node {
		mp.Node[n].setParent(f8n)
	}
}

// SumSize Выполняет суммирование всех размеров узла и всех подчинённых узлов.
// Возвращается размер в байтах.
func (mp *Map) SumSize() (ret int) {
	var (
		size int
		n    int
	)

	for n = range mp.Node {
		size += mp.Node[n].SumSize()
	}

	return mp.Size + size
}

// IsEmpty Возвращает истину если узел карты пустой - не является тегом и не содержит подчинённых узлов.
func (mp *Map) IsEmpty() bool {
	return mp.TagBegin == "" && mp.TagEnd == "" && mp.Content == "" && len(mp.Node) == 0
}

// Analysis Выполнение анализатора карты рекурсивно.
func (mp *Map) Analysis() (ret []*ParseError) {
	var (
		n   int
		tmp []*ParseError
	)

	for n = range mp.Node {
		if tmp = mp.Node[n].Analysis(); len(tmp) > 0 {
			ret = append(ret, tmp...)
		}
	}
	if tmp = mp.parent.analysis(mp); len(tmp) > 0 {
		ret = append(ret, tmp...)
	}

	return
}

// LoadAndParseFilter Рекурсивное копирование фильтров из запроса в DOM объект, а так же разбор фильтров
// в соответствии с синтаксисом.
func (mp *Map) LoadAndParseFilter(rq *http.Request) (ret []*ParseError) {
	var (
		errs      []*ParseError
		ero       *ParseError
		filter    Filter
		n         int
		ok        bool
		keyFilter string
		filters   []string
		tmp       string
	)

	// Обработка текущего узла.
	if mp.Origin == OriginFiltration && mp.Content != "" {
		keyFilter = strings.TrimSpace(mp.Content)
		if filters, ok = rq.URL.Query()[keyFilter]; !ok {
			ero = &ParseError{Ei: mp.parent.Errors().FilterCalledByNameWasNotFound(keyFilter)}
			ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
				Field:      keyMap,
				FieldValue: keyFilter,
				Message:    ero.Ei.Error(),
			})
			ret, filters = append(ret, ero), filters[:0]
		}
		if len(filters) > 0 {
			tmp = strings.Join(filters, keyAnd)
			if filter, errs = mp.parent.parseFilterSimple(tmp); len(errs) == 0 {
				mp.Filter = &filter
			}
			ret = append(ret, errs...)
		}
	}
	// Обработка всех подчинённых узлов.
	for n = range mp.Node {
		if errs = mp.Node[n].LoadAndParseFilter(rq); len(errs) > 0 {
			ret = append(ret, errs...)
		}
	}

	return
}
