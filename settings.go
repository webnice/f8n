// Package f8n
package f8n

import "strings"

// FieldSet Назначение списка обрабатываемых полей.
// Если список не установлен или установлен пустым значением, обрабатываются все полня, переданные в запросе.
// Если список установлен, но в запросе присутствует поле, отсутствующее в списке, будет возвращена ошибка.
func (f8n *impl) FieldSet(field ...string) Interface {
	var n int

	if len(f8n.FieldOnly) == 0 && len(field) > 0 {
		f8n.FieldOnly = make([]string, 0, len(field))
	}
	for n = range field {
		f8n.FieldOnly = append(f8n.FieldOnly, strings.TrimSpace(field[n]))
	}
	f8n.CleanNameField()

	return f8n
}

// FieldGet Возвращается список полей.
// Если список полей был установлен, тогда возвращаются установленные поля.
// Если список полей не был установлен, возвращаются поля, найденные в запросе.
func (f8n *impl) FieldGet() (ret []string) { return f8n.FieldOnly }

// FieldDatatype Назначение типа данных для поля.
// Все поля по умолчанию имеют тип данных "строка".
func (f8n *impl) FieldDatatype(fieldName string, t FieldType) Interface {
	fieldName = strings.TrimSpace(fieldName)
	f8n.Datatype[fieldName] = t

	return f8n
}

// Redefinition Назначение переопределения названий полей или замены полей.
// Назначенные или полученных из запроса поля, подменяются названиями полей используемых в базе данных.
func (f8n *impl) Redefinition(fieldName string, newFieldName string) Interface {
	fieldName = strings.TrimSpace(fieldName)
	f8n.Remap[fieldName] = strings.TrimSpace(newFieldName)

	return f8n
}
