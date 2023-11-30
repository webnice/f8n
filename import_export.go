package f8n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// ExportJson Сохранение настроек фильтрации в формате JSON.
// Экспортируются все разобранные поля запроса, а так же настройки переопределения и типов полей.
// Параметры prefix и ident используются для формирования отформатированного человеко-читаемого JSON, но обычно
// используются пустые значения форматирования.
func (f8n *impl) ExportJson(prefix string, ident string) (ret *bytes.Buffer) {
	var (
		err error
		buf *bytes.Buffer
		enc *json.Encoder
	)

	buf = &bytes.Buffer{}
	enc = json.NewEncoder(buf)
	enc.SetIndent(prefix, ident)
	if err = enc.Encode(f8n); err == nil {
		ret = buf
	}

	return
}

// ImportJson Загрузка настроек фильтрации из переданного в функцию JSON.
// Импортируются все ранее разобранные настройки, а так же настройки переопределения и типов полей.
func (f8n *impl) ImportJson(buf []byte) (err error) {
	f8n.Reset()
	if err = json.Unmarshal(buf, f8n); err != nil {
		return
	}
	// Для карты DOM объектов сложной фильтрации необходимо установить parent.
	if f8n.Map != nil {
		// Рекурсивная установка parent всем узлам карты.
		f8n.Map.setParent(f8n)
	}

	return
}

// String Интерфейс Stringify, представляет разобранный запрос в виде строки.
// В результат не попадают:
// 1. Список имён обрабатываемых полей.
// 2. Карта типов данных полей.
// 3. Переопределение наименований полей.
func (f8n *impl) String() (ret string) {
	var (
		uvl  url.Values
		keys []string
		vals []string
		key  string
		n, j int
		b    strings.Builder
	)

	uvl = f8n.Query()
	keys = make([]string, 0, len(uvl))
	// Сортируем ключи по алфавиту, часто применяется для расчёта подписи к запросу.
	for key = range uvl {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool { return strings.Compare(keys[i], keys[j]) < 0 })
	// Сборка отсортированного результата.
	for n = range keys {
		vals = uvl[keys[n]]
		for j = range vals {
			if n+j > 0 {
				_ = b.WriteByte('&')
			}
			_, _ = b.WriteString(keys[n])
			_ = b.WriteByte('=')
			_, _ = b.WriteString(vals[j])
		}

	}
	ret = b.String()

	return
}

// Query Возвращаются все разобранные значения в формате http.Request.Query().Values, тип данных url.Values.
// В результат не попадают:
// 1. Список имён обрабатываемых полей.
// 2. Карта типов данных полей.
// 3. Переопределение наименований полей.
func (f8n *impl) Query() (ret url.Values) {
	var (
		key, val string
		tmp      []string
		flt      url.Values
		n        int
	)

	ret = make(url.Values)
	// Простая фильтрация.
	if key, tmp = f8n.exportFilter(); key != "" || len(tmp) > 0 {
		for n = range tmp {
			val = tmp[n]
			ret.Add(key, val)
		}
	}
	// Сложная фильтрация.
	if key, val, flt = f8n.exportMap(); key != "" || val != "" {
		ret.Add(key, val)
		for key, tmp = range flt {
			for n = range tmp {
				val = tmp[n]
				ret.Add(key, val)
			}
		}
	}
	// Лимит.
	if key, val = f8n.exportLimit(); key != "" || val != "" {
		ret.Add(key, val)
	}
	// Сортировка.
	if key, tmp = f8n.exportBy(); key != "" || len(tmp) > 0 {
		for n = range tmp {
			val = tmp[n]
			ret.Add(key, val)
		}
	}

	return
}

// ImportString Импорт строки запроса.
// Строка запроса может быть строкой ранее полученной через вызов функции String(), либо строкой полученной
// из параметров URN запроса.
func (f8n *impl) ImportString(s string) (err error) {
	const localhost = `http://localhost`
	var (
		sb strings.Builder
		rq *http.Request
	)

	_, _ = sb.WriteString(localhost)
	if s != "" {
		_ = sb.WriteByte('?')
		_, _ = sb.WriteString(s)
	}
	if rq, err = http.NewRequest("GET", sb.String(), nil); err != nil {
		err = fmt.Errorf("импорт строки: %q, прерван ошибкой: %s", s, err)
		return
	}
	// Сброс только тех значений которые экспортируются через функцию String()
	f8n.ResetLimit().ResetBy().ResetFilter().ResetMap()
	err = f8n.ParseRequest(rq)

	return
}
