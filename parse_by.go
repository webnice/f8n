package f8n

import (
	"net/http"
	"strings"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// ParseBy Загрузка сортировки.
func (f8n *impl) ParseBy(rq *http.Request) (ret []*ParseError) {
	const nSection = 2
	var (
		ero      *ParseError
		src, tmp []string
		n, j     int
		dir      ByDirection
		found    bool
	)

	if src = rq.URL.Query()[keyBy]; len(src) == 0 {
		return
	}
	for n = range src {
		// Проверка формата.
		if tmp = strings.SplitN(src[n], keyDelimiter, nSection); len(tmp) != nSection {
			ero = &ParseError{Ei: f8n.Errors().ByFormat.Bind()}
			ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
				Field:      keyBy,
				FieldValue: src[n],
				Message:    ero.Ei.Error(),
			})
			ret = append(ret, ero)
			continue
		}
		// Проверка ключевого слова сортировки.
		if dir = parseByDirection(tmp[1]); dir == byUnknown {
			ero = &ParseError{Ei: f8n.Errors().ByDirectionUnknown.Bind(tmp[1])}
			ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
				Field:      keyBy,
				FieldValue: src[n],
				Message:    ero.Ei.Error(),
			})
			ret = append(ret, ero)
			continue
		}
		f8n.By = append(f8n.By, Direction{Field: tmp[0], By: dir})
	}
	// Если заданы поля, проверка имени поля сортировки.
	if len(f8n.FieldOnly) > 0 {
		for n = range f8n.By {
			found = false
			for j = range f8n.FieldOnly {
				if f8n.FieldOnly[j] == f8n.By[n].Field {
					found = true
				}
			}
			if !found {
				ero = &ParseError{Ei: f8n.Errors().ByDirectionField.Bind(f8n.By[n].Field)}
				ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
					Field:      keyBy,
					FieldValue: src[n],
					Message:    ero.Ei.Error(),
				})
				ret = append(ret, ero)
			}
		}
	}

	return
}

// Экспорт сортировки.
func (f8n *impl) exportBy() (key string, val []string) {
	var n int

	if len(f8n.By) == 0 {
		return
	}
	key, val = keyBy, make([]string, 0, len(f8n.By))
	for n = range f8n.By {
		val = append(val, f8n.By[n].exportAsString())
	}

	return
}
