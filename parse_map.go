// Package f8n
package f8n

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// ParseMap Загрузка и разбор карты условий сложной фильтрации.
func (f8n *impl) ParseMap(rq *http.Request) (ret []*ParseError) {
	const tpl1, tpl2 = "(%s)", keyAnd + "(%s)"
	var (
		values   url.Values
		mapValue []string
		buf      *bytes.Buffer
		n, j     int
		tpl, tmp string
	)

	values = rq.URL.Query()
	mapValue = values[keyMap]
	switch len(mapValue) {
	case 0:
		return
	case 1:
		if len(mapValue[0]) == 0 {
			return
		}
		buf = bytes.NewBufferString(mapValue[0])
	default:
		buf = &bytes.Buffer{}
		for n = range mapValue {
			if tpl = tpl1; n > 0 {
				tpl = tpl2
			}
			_, _ = buf.WriteString(fmt.Sprintf(tpl, mapValue[n]))
		}
	}
	f8n.Map = f8n.parseMap(buf.Bytes(), 0)
	// Анализ корректности разбора.
	if ret = f8n.Map.Analysis(); len(ret) > 0 {
		tmp = buf.String()
		for n = range ret {
			for j = range ret[n].Ev {
				if ret[n].Ev[j].FieldValue == "" {
					ret[n].Ev[j].FieldValue = tmp
				}
			}
		}
		return
	}
	// Рекурсивное копирование фильтров из запроса в DOM объект, а так же разбор фильтров в соответствии с синтаксисом.
	if ret = f8n.Map.LoadAndParseFilter(rq); len(ret) > 0 {
		return
	}

	return
}

// Рекурсивная функция разбора карты условий сложной фильтрации.
func (f8n *impl) parseMap(buf []byte, cursor int) (ret *Map) {
	var (
		loc  []int
		next *Map
		tag  string
	)

	ret = &Map{parent: f8n}
	for {
		if len(buf[cursor:]) == 0 {
			break
		}
		// Поиск первого тега.
		loc = rexTag.FindIndex(buf[cursor:])
		// Если ни один тег не найден.
		if len(loc) == 0 && cursor == 0 {
			ret.Content = string(buf[cursor:])
			ret.Size = len(ret.Content)
			cursor += ret.Size
			continue
		}
		// Всё после последнего тега, добавляется в последний узел.
		if len(loc) == 0 && cursor > 0 {
			next = &Map{
				parent:  f8n,
				Content: string(buf[cursor:]),
				Size:    len(buf[cursor:]),
			}
			ret.Node = append(ret.Node, next)
			cursor += next.Size
			continue
		}
		// Всё до первого тега, добавляется в первый узел.
		if loc[0] > 0 {
			next = &Map{
				parent:  f8n,
				Content: string(buf[cursor : cursor+loc[0]]),
				Size:    loc[0],
			}
			ret.Node = append(ret.Node, next)
			cursor += next.Size
			continue
		}
		// На этой позиции первый тег найден.
		tag = string(buf[cursor:][:loc[1]])
		// Далее надо выяснить парный или не парный тег.
		// --- Не парные теги.
		if tag != pairBeg && tag != pairEnd {
			next = &Map{parent: f8n, TagBegin: tag, TagEnd: tag, Size: loc[1]}
			ret.Node = append(ret.Node, next)
			cursor += next.Size
			continue
		}
		loc = rexTag.FindIndex(buf[cursor:])
		// --- Парный тег. Открывающие всех парных тегов.
		if tag == pairBeg {
			next = f8n.parseMap(buf, cursor+loc[1])
			next.TagBegin = tag
			next.Size += loc[1]
			ret.Node = append(ret.Node, next)
			cursor += next.SumSize()
			continue
		}
		// --- Парный тег. Закрывающие всех парных тегов.
		if tag == pairEnd {
			ret.TagEnd = tag
			ret.Size = loc[1]
			break
		}
	}

	return
}

// Экспорт карты условий сложной фильтрации.
func (f8n *impl) exportMap() (key string, val string, uvl url.Values) {
	if f8n.Map == nil {
		return
	}
	key = keyMap
	val, uvl = f8n.Map.exportAsString()

	return
}
