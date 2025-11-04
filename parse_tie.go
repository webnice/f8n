package f8n

import (
	"fmt"
	"net/http"
	"strings"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// ParseTie Загрузка устаревшего режима фильтрации - tie.
// Режим tie не совместим с map.
func (f8n *impl) ParseTie(rq *http.Request) (ret []*ParseError) {
	var (
		err   error
		ero   *ParseError
		modes []string
		tie   TieMode
	)

	if modes = rq.URL.Query()[keyTie]; len(modes) > 1 {
		ero = &ParseError{Ei: f8n.Errors().TieModeThanOne.Bind()}
		ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
			Field:      keyLimit,
			FieldValue: strings.Join(modes, ", "),
			Message:    ero.Ei.Error(),
		})
		ret = append(ret, ero)
		return
	}
	if len(modes) <= 0 {
		return
	}
	switch tie = parseTie(modes[0]); tie {
	case tieAnd, tieOr:
		f8n.Tie = tie
	default:
		err = fmt.Errorf("допустимые значения: %q или %q", tieAnd, tieOr)
		ero = &ParseError{Ei: Errors().TieModeInvalidValue.Bind(modes[0], err)}
		ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
			Field:      keyTie,
			FieldValue: modes[0],
			Message:    ero.Ei.Error(),
		})
		ret = append(ret, ero)
		return
	}

	return
}
