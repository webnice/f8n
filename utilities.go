package f8n

import (
	"fmt"
	"strconv"
	"strings"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// CleanNameField Очистка названия полей от возможных не корректных символов.
func (f8n *impl) CleanNameField() {
	var (
		badWord []string
		n, j    int
	)

	badWord = []string{"`", "'", `"`, ","}
	for n = range f8n.FieldOnly {
		for j = range badWord {
			f8n.FieldOnly[n] = strings.ReplaceAll(f8n.FieldOnly[n], badWord[j], "")
		}
	}

	return
}

// Выполнение замены наименования поля базы данных.
func (f8n *impl) fieldName(name string) (ret string) {
	var ok bool

	if ret, ok = f8n.Remap[name]; !ok {
		ret = fmt.Sprintf("`%s`", name)
	}

	return
}

// Конвертация строки в число int64 без реакции на пустое значение.
func parseInt64(field string, s string) (ret int64, ero *ParseError) {
	var err error

	if s == "" {
		return
	}
	if ret, err = strconv.ParseInt(s, 10, 64); err != nil {
		ero = &ParseError{Ei: Errors().LimitInvalidValue(s, err)}
		ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
			Field:      field,
			FieldValue: s,
			Message:    ero.Ei.Error(),
		})
		return
	}

	return
}
