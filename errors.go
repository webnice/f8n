package f8n

// Все ошибки определены как константы. Коды ошибок приложения:

// Обычные ошибки
const (
	ePanic uint8 = iota + 1
	eRequestIsNil
	eMultipleErrorsFound
	eLimitReceivedMoreThanOne
	eLimitInvalidValue
	eTieModeThanOne
	eTieModeInvalidValue
	eValueCannotBeNegative
	eByFormat
	eByDirectionUnknown
	eByDirectionField
	eFilterFormat
	eFilterUnknownMethod
	eFilterUnknownField
	eFilterValueType
	ePairedTagNotMatch
	eWrongLogicalOperation
	eOperatorBracketEmpty
	eOperatorBracketOneItem
	eNoLogicalOperationBetweenBrackets
	eFilterCalledByNameWasNotFound
	eOrmIsNil
)

// Текстовые значения кодов ошибок на основном языке приложения.
const (
	cPanic                             = "Восстановление после паники.\n%v\n%s"
	cRequestIsNil                      = `Передан http.Request равный nil.`
	cMultipleErrorsFound               = `Найдены множественные ошибки.`
	cLimitReceivedMoreThanOne          = `Получено больше одного значения лимита. Поддерживается только одно значение лимита.`
	cLimitInvalidValue                 = `Получено недопустимое значение лимита %q, ошибка: %s.`
	cTieModeThanOne                    = `Получено больше одного значения режима фильтрации. Поддерживается только одно значение режима фильтрации.`
	cTieModeInvalidValue               = `Получено недопустимое значение режима фильтрации %q, ошибка: %s.`
	cValueCannotBeNegative             = `Передано значение %d, значение не может быть отрицательным.`
	cByFormat                          = `Формат значения сортировки не верный.`
	cByDirectionUnknown                = `Неизвестный порядок сортировки %q, ожидался asc или desc.`
	cByDirectionField                  = `Указана сортировка по не известному полю %q.`
	cFilterFormat                      = `Формат значения фильтрации не верный.`
	cFilterUnknownMethod               = `Неизвестный способ фильтрации %q.`
	cFilterUnknownField                = `Указана фильтрация по не известному полю %q.`
	cFilterValueType                   = `Указано значение %q не соответствующее типу %q, ошибка: %s.`
	cPairedTagNotMatch                 = `Открывающий тег %q не соответствует закрывающему тегу %q.`
	cWrongLogicalOperation             = `Не верная логическая операция %q`
	cOperatorBracketEmpty              = `Указаны пустые операторные скобки.`
	cOperatorBracketOneItem            = `В операторные скобки заключён один элемент, уберите операторные скобки.`
	cNoLogicalOperationBetweenBrackets = `Необходимо указать логическую операцию между операторных скобок.`
	cFilterCalledByNameWasNotFound     = `Значение фильтра с названием %q не было найдено в запросе.`
	cOrmIsNil                          = `Передан ORM объект являющейся nil. Ожидался не nil объект.`
)

// Константы указаны в объектах, адрес которых фиксирован всё время работы приложения.
// Это позволяет сравнивать ошибки между собой используя обычное сравнение "==", но сравнивать необходимо только якорь "Anchor()" объекта ошибки.
var (
	errSingleton                         = &Error{}
	errPanic                             = err{tpl: cPanic, code: ePanic}
	errRequestIsNil                      = err{tpl: cRequestIsNil, code: eRequestIsNil}
	errMultipleErrorsFound               = err{tpl: cMultipleErrorsFound, code: eMultipleErrorsFound}
	errLimitReceivedMoreThanOne          = err{tpl: cLimitReceivedMoreThanOne, code: eLimitReceivedMoreThanOne}
	errLimitInvalidValue                 = err{tpl: cLimitInvalidValue, code: eLimitInvalidValue}
	errTieModeThanOne                    = err{tpl: cTieModeThanOne, code: eTieModeThanOne}
	errTieModeInvalidValue               = err{tpl: cTieModeInvalidValue, code: eTieModeInvalidValue}
	errValueCannotBeNegative             = err{tpl: cValueCannotBeNegative, code: eValueCannotBeNegative}
	errByFormat                          = err{tpl: cByFormat, code: eByFormat}
	errByDirectionUnknown                = err{tpl: cByDirectionUnknown, code: eByDirectionUnknown}
	errByDirectionField                  = err{tpl: cByDirectionField, code: eByDirectionField}
	errFilterFormat                      = err{tpl: cFilterFormat, code: eFilterFormat}
	errFilterUnknownMethod               = err{tpl: cFilterUnknownMethod, code: eFilterUnknownMethod}
	errFilterUnknownField                = err{tpl: cFilterUnknownField, code: eFilterUnknownField}
	errFilterValueType                   = err{tpl: cFilterValueType, code: eFilterValueType}
	errPairedTagNotMatch                 = err{tpl: cPairedTagNotMatch, code: ePairedTagNotMatch}
	errWrongLogicalOperation             = err{tpl: cWrongLogicalOperation, code: eWrongLogicalOperation}
	errOperatorBracketEmpty              = err{tpl: cOperatorBracketEmpty, code: eOperatorBracketEmpty}
	errOperatorBracketOneItem            = err{tpl: cOperatorBracketOneItem, code: eOperatorBracketOneItem}
	errNoLogicalOperationBetweenBrackets = err{tpl: cNoLogicalOperationBetweenBrackets, code: eNoLogicalOperationBetweenBrackets}
	errFilterCalledByNameWasNotFound     = err{tpl: cFilterCalledByNameWasNotFound, code: eFilterCalledByNameWasNotFound}
	errOrmIsNil                          = err{tpl: cOrmIsNil, code: eOrmIsNil}
)

// ERRORS: Реализация ошибок с возможностью сравнения ошибок между собой.

// Panic Восстановление после паники. ... ...
func (e *Error) Panic(err interface{}, stack []byte) Err {
	return newErr(&errPanic, 0, err, string(stack))
}

// RequestIsNil Передан http.Request равный nil.
func (e *Error) RequestIsNil() Err { return newErr(&errRequestIsNil, 0) }

// MultipleErrorsFound Найдены множественные ошибки.
func (e *Error) MultipleErrorsFound() Err { return newErr(&errMultipleErrorsFound, 0) }

// LimitReceivedMoreThanOne Получено больше одного значения лимита. Поддерживается только одно значение лимита.
func (e *Error) LimitReceivedMoreThanOne() Err { return newErr(&errLimitReceivedMoreThanOne, 0) }

// LimitInvalidValue Получено недопустимое значение лимита: ...
func (e *Error) LimitInvalidValue(s string, err error) Err {
	return newErr(&errLimitInvalidValue, 0, s, err)
}

// TieModeThanOne Получено больше одного значения режима фильтрации. Поддерживается только одно значение режима
// фильтрации.
func (e *Error) TieModeThanOne() Err { return newErr(&errTieModeThanOne, 0) }

// TieModeInvalidValue Получено недопустимое значение режима фильтрации ..., ошибка: ...
func (e *Error) TieModeInvalidValue(s string, err error) Err {
	return newErr(&errTieModeInvalidValue, 0, s, err)
}

// ValueCannotBeNegative Передано значение ..., значение не может быть отрицательным.
func (e *Error) ValueCannotBeNegative(v interface{}) Err {
	return newErr(&errValueCannotBeNegative, 0, v)
}

// ByFormat Формат значения сортировки не верный.
func (e *Error) ByFormat() Err { return newErr(&errByFormat, 0) }

// ByDirectionUnknown Неизвестный порядок сортировки ..., ожидался asc или desc.
func (e *Error) ByDirectionUnknown(s string) Err { return newErr(&errByDirectionUnknown, 0, s) }

// ByDirectionField Указана сортировка по не известному полю ...
func (e *Error) ByDirectionField(s string) Err { return newErr(&errByDirectionField, 0, s) }

// FilterFormat Формат значения фильтрации не верный.
func (e *Error) FilterFormat() Err { return newErr(&errFilterFormat, 0) }

// FilterUnknownMethod Неизвестный способ фильтрации ...
func (e *Error) FilterUnknownMethod(s string) Err { return newErr(&errFilterUnknownMethod, 0, s) }

// FilterUnknownField Указана фильтрация по не известному полю ...
func (e *Error) FilterUnknownField(s string) Err { return newErr(&errFilterUnknownField, 0, s) }

// FilterValueType Указано значение ... не соответствующее типу ..., ошибка: ...
func (e *Error) FilterValueType(s string, t FieldType, err error) Err {
	return newErr(&errFilterValueType, 0, s, t, err)
}

// PairedTagNotMatch Открывающий тег ... не соответствует закрывающему тегу ...
func (e *Error) PairedTagNotMatch(beg string, end string) Err {
	return newErr(&errPairedTagNotMatch, 0, beg, end)
}

// WrongLogicalOperation Не верная логическая операция ...
func (e *Error) WrongLogicalOperation(op Origin) Err {
	return newErr(&errWrongLogicalOperation, 0, op)
}

// OperatorBracketEmpty Указаны пустые операторные скобки.
func (e *Error) OperatorBracketEmpty() Err { return newErr(&errOperatorBracketEmpty, 0) }

// OperatorBracketOneItem В операторные скобки заключён один элемент, уберите операторные скобки.
func (e *Error) OperatorBracketOneItem() Err { return newErr(&errOperatorBracketOneItem, 0) }

// NoLogicalOperationBetweenBrackets Необходимо указать логическую операцию между операторных скобок.
func (e *Error) NoLogicalOperationBetweenBrackets() Err {
	return newErr(&errNoLogicalOperationBetweenBrackets, 0)
}

// FilterCalledByNameWasNotFound Значение фильтра с названием ... не было найдено в запросе.
func (e *Error) FilterCalledByNameWasNotFound(name string) Err {
	return newErr(&errFilterCalledByNameWasNotFound, 0, name)
}

// OrmIsNil Передан ORM объект являющейся nil. Ожидался не nil объект.
func (e *Error) OrmIsNil() Err { return newErr(&errOrmIsNil, 0) }
