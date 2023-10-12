// Package f8n
package f8n

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/webnice/kit/v2/module/verify"
)

// ParseLimit Загрузка лимита.
func (f8n *impl) ParseLimit(rq *http.Request) (ret []*ParseError) {
	var (
		ero        *ParseError
		src        []string
		tmp        []string
		i64o, i64l int64
	)

	if src = rq.URL.Query()[keyLimit]; len(src) > 1 {
		ero = &ParseError{Ei: f8n.Errors().LimitReceivedMoreThanOne()}
		ero.Ev = append(ero.Ev, verify.Error{
			Field:      keyLimit,
			FieldValue: strings.Join(src, ", "),
			Message:    ero.Ei.Error(),
		})
		ret = append(ret, ero)
		return
	}
	if len(src) <= 0 {
		return
	}
	// Разбор значения лимита.
	switch tmp = strings.SplitN(src[0], keyDelimiter, 2); len(tmp) {
	case 1:
		if i64l, ero = parseInt64(keyLimit, tmp[0]); ero != nil {
			ret = append(ret, ero)
			return
		}
		if i64l < 0 {
			i64l = 0
		}
	case 2:
		if i64o, ero = parseInt64(keyLimit, tmp[0]); ero != nil {
			ret = append(ret, ero)
		}
		if i64l, ero = parseInt64(keyLimit, tmp[1]); ero != nil {
			ret = append(ret, ero)
		}
		if i64o < 0 {
			ero = &ParseError{Ei: f8n.Errors().ValueCannotBeNegative(i64o)}
			ero.Ev = append(ero.Ev, verify.Error{
				Field:      keyLimit,
				FieldValue: tmp[0],
				Message:    ero.Ei.Error(),
			})
			ret = append(ret, ero)
			return
		}
		if i64l < 0 {
			i64l = 0
		}
	}
	f8n.Offset, f8n.Limit = uint64(i64o), uint64(i64l)

	return
}

// Экспорт лимита.
func (f8n *impl) exportLimit() (key string, val string) {
	if f8n.Offset == 0 && f8n.Limit == 0 {
		return
	}
	key = keyLimit
	if f8n.Offset != 0 {
		val += strconv.FormatUint(f8n.Offset, 10)
	}
	val += keyDelimiter
	if f8n.Limit != 0 {
		val += strconv.FormatUint(f8n.Limit, 10)
	}

	return
}
