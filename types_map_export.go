package f8n

import (
	"net/url"
	"strings"
)

// Экспорт карты сложной фильтрации.
// Возвращаются значения:
// 1. Тело карты DOM сложной фильтрации.
// 2. Значения фильтров используемых в карте сложной фильтрации в формате name=value.
func (mp *Map) exportAsString() (ret string, uvl url.Values) {
	uvl = make(url.Values)
	mp.exportFilters(uvl)
	ret = mp.recursivelyAsString()

	return
}

// Рекурсивный экспорт карты сложной фильтрации.
func (mp *Map) recursivelyAsString() (ret string) {
	var (
		tmp []string
		n   int
	)

	switch mp.Origin {
	case OriginFiltration:
		ret = mp.Content
	case OriginAnd, OriginOr:
		ret = mp.Origin.String()
	case OriginUnknown:
		if len(mp.Node) == 0 {
			return
		}
		fallthrough
	case OriginOperatorBracket:
		tmp = make([]string, 0, len(mp.Node)+2)
		tmp = append(tmp, mp.TagBegin)
		for n = range mp.Node {
			tmp = append(tmp, mp.Node[n].recursivelyAsString())
		}
		tmp = append(tmp, mp.TagEnd)
		ret = strings.Join(tmp, "")
	}

	return
}

// Рекурсивный обход карты сложной фильтрации и экспорт в url.Values всех фильтров используемых в карте.
func (mp *Map) exportFilters(uvl url.Values) {
	var n int

	if mp.Origin == OriginFiltration && mp.Filter != nil {
		uvl.Set(mp.Content, mp.Filter.exportAsString())
	}
	for n = range mp.Node {
		mp.Node[n].exportFilters(uvl)
	}
}
