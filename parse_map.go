package f8n

import (
	"net/http"
	"net/url"
)

// ParseMap Загрузка и разбор карты условий сложной фильтрации.
// Если передано несколько не пустых параметров map, тогда DOM деревья условий объединяются
// по принципу логического «ИЛИ», каждое дерево является операндом логической операции.
func (f8n *impl) ParseMap(rq *http.Request) (ret []*ParseError) {
	var (
		mapValue []string
		parts    []*Map
		root     *Map
		part     *Map
		operand  *Map
		ers      []*ParseError
		n, j     int
	)

	if mapValue = rq.URL.Query()[keyMap]; len(mapValue) == 0 {
		return
	}
	// Разбор каждого не пустого параметра map в отдельное DOM дерево условий.
	for n = range mapValue {
		if len(mapValue[n]) == 0 {
			continue
		}
		part = f8n.parseMap([]byte(mapValue[n]), 0)
		// Анализ корректности разбора части.
		if ers = part.Analysis(); len(ers) > 0 {
			// Подстановка разбираемого значения в описание ошибок.
			for j = range ers {
				for n = range ers[j].Ev {
					if ers[j].Ev[n].FieldValue == "" {
						ers[j].Ev[n].FieldValue = mapValue[n]
					}
				}
			}
			ret = append(ret, ers...)
			continue
		}
		parts = append(parts, part)
	}
	// Единственное дерево условий используется без объединения.
	switch len(parts) {
	case 0:
		return
	case 1:
		f8n.Map = parts[0]
	default:
		// Объединение деревьев условий по принципу логического «ИЛИ».
		root = &Map{parent: f8n}
		for n = range parts {
			part = parts[n]
			switch {
			case len(part.Node) == 0:
				// Корень фильтрации без подчинённых узлов, является операндом логической операции.
				operand = part
			case len(part.Node) == 1:
				// Единственный узел дерева является операндом логической операции.
				operand = part.Node[0]
			default:
				// Условия дерева группируются в операторные скобки, группа является операндом логической операции.
				operand = &Map{
					parent:   f8n,
					TagBegin: pairBeg,
					TagEnd:   pairEnd,
					Origin:   OriginOperatorBracket,
					Node:     part.Node,
					Size:     len(pairBeg) + len(pairEnd),
				}
			}
			// Операнды объединяются логической операцией «ИЛИ».
			if len(root.Node) > 0 {
				root.Node = append(root.Node, &Map{
					parent:   f8n,
					TagBegin: keyOr,
					TagEnd:   keyOr,
					Origin:   OriginOr,
					Size:     len(keyOr),
				})
			}
			root.Node = append(root.Node, operand)
		}
		f8n.Map = root
	}
	if len(ret) > 0 {
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
