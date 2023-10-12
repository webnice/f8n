// Package f8n
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
