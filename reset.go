// Package f8n
package f8n

// Reset Сброс всех загруженных или импортированных настроек, переопределений и типов полей.
func (f8n *impl) Reset() Interface {
	return f8n.ResetLimit().ResetBy().ResetFilter().ResetMap().ResetRedefinition().ResetDatatype().ResetFieldSet()
}

// ResetLimit Сброс настроек лимитирования и постраничного вывода.
func (f8n *impl) ResetLimit() Interface {
	f8n.Offset, f8n.Limit = 0, 0
	return f8n
}

// ResetBy Сброс настроек сортировки.
func (f8n *impl) ResetBy() Interface {
	f8n.By = f8n.By[:0]
	return f8n
}

// ResetFilter Сброс настроек простой фильтрации.
func (f8n *impl) ResetFilter() Interface {
	f8n.Filter = f8n.Filter[:0]
	return f8n
}

// ResetMap Сброс DOM дерева карты настроек сложной фильтрации.
func (f8n *impl) ResetMap() Interface {
	f8n.Map = nil
	return f8n
}

// ResetRedefinition Сброс переопределения названий полей или замены полей.
func (f8n *impl) ResetRedefinition() Interface {
	f8n.Remap = make(Remap)
	return f8n
}

// ResetDatatype Сброс назначения типа данных для полей.
func (f8n *impl) ResetDatatype() Interface {
	f8n.Datatype = make(Datatype)
	return f8n
}

// ResetFieldSet Сброс списка обрабатываемых полей.
func (f8n *impl) ResetFieldSet() Interface {
	f8n.FieldOnly = f8n.FieldOnly[:0]
	return f8n
}
