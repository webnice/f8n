package f8n

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
func New() Interface {
	var f8n = &impl{
		Datatype: make(Datatype),
		Remap:    make(Remap),
	}

	return f8n
}

// Errors Справочник ошибок.
func (f8n *impl) Errors() *Error { return Errors() }

// IsFiltration Информация о наличии фильтрации в запросе.
// Истина - Есть хотя бы один фильтр.
// Ложь - фильтры отсутствуют.
func (f8n *impl) IsFiltration() bool {
	return len(f8n.Filter) > 0
}
