package f8n

import (
	"fmt"
	"strconv"
	"strings"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// CleanNameField Очистка названия полей от возможных не корректных символов.
func (f8n *impl) CleanNameField() {
	var n int

	for n = range f8n.FieldOnly {
		f8n.FieldOnly[n] = cleanNameField(f8n.FieldOnly[n])
	}

	return
}

// Очистка названия поля от символов опасных при формировании SQL запроса.
// Удаляются символы: ` ' " , - позволяющие выйти за пределы идентификатора поля.
func cleanNameField(name string) (ret string) {
	const dangerous = "`'\","
	var n int

	ret = name
	for n = 0; n < len(dangerous); n++ {
		ret = strings.ReplaceAll(ret, dangerous[n:n+1], "")
	}

	return
}

// Выполнение замены наименования поля базы данных.
// Название поля очищается от символов опасных при формировании SQL запроса.
// Название полученное из переопределения полей используется без изменений, является доверенным вводом приложения.
func (f8n *impl) fieldName(name string) (ret string) {
	var ok bool

	if ret, ok = f8n.Remap[name]; !ok {
		ret = fmt.Sprintf("`%s`", cleanNameField(name))
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
		ero = &ParseError{Ei: Errors().LimitInvalidValue.Bind(s, err)}
		ero.Ev = append(ero.Ev, kitModuleAns.RestErrorField{
			Field:      field,
			FieldValue: s,
			Message:    ero.Ei.Error(),
		})
		return
	}

	return
}
