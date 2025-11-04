package f8n

import "github.com/webnice/dic"

// Коды ошибок.
const (
	ePanic                             uint8 = iota + 1 // 001
	eRequestIsNil                                       // 002
	eMultipleErrorsFound                                // 003
	eLimitReceivedMoreThanOne                           // 004
	eLimitInvalidValue                                  // 005
	eTieModeThanOne                                     // 006
	eTieModeInvalidValue                                // 007
	eValueCannotBeNegative                              // 008
	eByFormat                                           // 009
	eByDirectionUnknown                                 // 010
	eByDirectionField                                   // 011
	eFilterFormat                                       // 012
	eFilterUnknownMethod                                // 013
	eFilterUnknownField                                 // 014
	eFilterValueType                                    // 015
	ePairedTagNotMatch                                  // 016
	eWrongLogicalOperation                              // 017
	eOperatorBracketEmpty                               // 018
	eOperatorBracketOneItem                             // 019
	eNoLogicalOperationBetweenBrackets                  // 020
	eFilterCalledByNameWasNotFound                      // 021
	eOrmIsNil                                           // 022
)

// Константы ошибок.
const (
	cPanic                             = "Восстановление после паники.\n%v\n%s"
	cRequestIsNil                      = "Передан http.Request равный nil."
	cMultipleErrorsFound               = "Найдены множественные ошибки."
	cLimitReceivedMoreThanOne          = "Получено больше одного значения лимита. Поддерживается только одно значение лимита."
	cLimitInvalidValue                 = "Получено недопустимое значение лимита %q, ошибка: %s."
	cTieModeThanOne                    = "Получено больше одного значения режима фильтрации. Поддерживается только одно значение режима фильтрации."
	cTieModeInvalidValue               = "Получено недопустимое значение режима фильтрации %q, ошибка: %s."
	cValueCannotBeNegative             = "Передано значение %d, значение не может быть отрицательным."
	cByFormat                          = "Формат значения сортировки не верный."
	cByDirectionUnknown                = "Неизвестный порядок сортировки %q, ожидался asc или desc."
	cByDirectionField                  = "Указана сортировка по не известному полю %q."
	cFilterFormat                      = "Формат значения фильтрации не верный."
	cFilterUnknownMethod               = "Неизвестный способ фильтрации %q."
	cFilterUnknownField                = "Указана фильтрация по не известному полю %q."
	cFilterValueType                   = "Указано значение %q не соответствующее типу %q, ошибка: %s."
	cPairedTagNotMatch                 = "Открывающий тег %q не соответствует закрывающему тегу %q."
	cWrongLogicalOperation             = "Не верная логическая операция %q"
	cOperatorBracketEmpty              = "Указаны пустые операторные скобки."
	cOperatorBracketOneItem            = "В операторные скобки заключён один элемент, уберите операторные скобки."
	cNoLogicalOperationBetweenBrackets = "Необходимо указать логическую операцию между операторных скобок."
	cFilterCalledByNameWasNotFound     = "Значение фильтра с названием %q не было найдено в запросе."
	cOrmIsNil                          = "Передан ORM объект являющейся nil. Ожидался не nil объект."
)

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors

	// Panic Восстановление после паники. ...
	Panic dic.IError

	// RequestIsNil Передан http.Request равный nil.
	RequestIsNil dic.IError

	// MultipleErrorsFound Найдены множественные ошибки.
	MultipleErrorsFound dic.IError

	// LimitReceivedMoreThanOne Получено больше одного значения лимита. Поддерживается только одно значение лимита.
	LimitReceivedMoreThanOne dic.IError

	// LimitInvalidValue Получено недопустимое значение лимита ..., ошибка: ...
	LimitInvalidValue dic.IError

	// TieModeThanOne Получено больше одного значения режима фильтрации.
	// Поддерживается только одно значение режима фильтрации.
	TieModeThanOne dic.IError

	// TieModeInvalidValue Получено недопустимое значение режима фильтрации ..., ошибка: ...
	TieModeInvalidValue dic.IError

	// ValueCannotBeNegative Передано значение ..., значение не может быть отрицательным.
	ValueCannotBeNegative dic.IError

	// ByFormat Формат значения сортировки не верный.
	ByFormat dic.IError

	// ByDirectionUnknown Неизвестный порядок сортировки ..., ожидался asc или desc.
	ByDirectionUnknown dic.IError

	// ByDirectionField Указана сортировка по не известному полю ...
	ByDirectionField dic.IError

	// FilterFormat Формат значения фильтрации не верный.
	FilterFormat dic.IError

	// FilterUnknownMethod Неизвестный способ фильтрации ...
	FilterUnknownMethod dic.IError

	// FilterUnknownField Указана фильтрация по не известному полю ...
	FilterUnknownField dic.IError

	// FilterValueType Указано значение ... не соответствующее типу ..., ошибка: ...
	FilterValueType dic.IError

	// PairedTagNotMatch Открывающий тег ... не соответствует закрывающему тегу ...
	PairedTagNotMatch dic.IError

	// WrongLogicalOperation Не верная логическая операция ...
	WrongLogicalOperation dic.IError

	// OperatorBracketEmpty Указаны пустые операторные скобки.
	OperatorBracketEmpty dic.IError

	// OperatorBracketOneItem В операторные скобки заключён один элемент, уберите операторные скобки.
	OperatorBracketOneItem dic.IError

	// NoLogicalOperationBetweenBrackets Необходимо указать логическую операцию между операторных скобок.
	NoLogicalOperationBetweenBrackets dic.IError

	// FilterCalledByNameWasNotFound Значение фильтра с названием ... не было найдено в запросе.
	FilterCalledByNameWasNotFound dic.IError

	// OrmIsNil Передан ORM объект являющейся nil. Ожидался не nil объект.
	OrmIsNil dic.IError
}

var errSingleton = &Error{
	Errors:                            dic.Error(),
	Panic:                             dic.NewError(cPanic, "паника", "стек").CodeU8().Set(ePanic),
	RequestIsNil:                      dic.NewError(cRequestIsNil).CodeU8().Set(eRequestIsNil),
	MultipleErrorsFound:               dic.NewError(cMultipleErrorsFound).CodeU8().Set(eMultipleErrorsFound),
	LimitReceivedMoreThanOne:          dic.NewError(cLimitReceivedMoreThanOne).CodeU8().Set(eLimitReceivedMoreThanOne),
	LimitInvalidValue:                 dic.NewError(cLimitInvalidValue, "значение", "ошибка").CodeU8().Set(eLimitInvalidValue),
	TieModeThanOne:                    dic.NewError(cTieModeThanOne).CodeU8().Set(eTieModeThanOne),
	TieModeInvalidValue:               dic.NewError(cTieModeInvalidValue, "значение", "ошибка").CodeU8().Set(eTieModeInvalidValue),
	ValueCannotBeNegative:             dic.NewError(cValueCannotBeNegative, "значение").CodeU8().Set(eValueCannotBeNegative),
	ByFormat:                          dic.NewError(cByFormat).CodeU8().Set(eByFormat),
	ByDirectionUnknown:                dic.NewError(cByDirectionUnknown, "значение").CodeU8().Set(eByDirectionUnknown),
	ByDirectionField:                  dic.NewError(cByDirectionField, "название поля").CodeU8().Set(eByDirectionField),
	FilterFormat:                      dic.NewError(cFilterFormat).CodeU8().Set(eFilterFormat),
	FilterUnknownMethod:               dic.NewError(cFilterUnknownMethod, "значение").CodeU8().Set(eFilterUnknownMethod),
	FilterUnknownField:                dic.NewError(cFilterUnknownField, "название поля").CodeU8().Set(eFilterUnknownField),
	FilterValueType:                   dic.NewError(cFilterValueType, "значение", "тип", "ошибка").CodeU8().Set(eFilterValueType),
	PairedTagNotMatch:                 dic.NewError(cPairedTagNotMatch, "тег", "тег").CodeU8().Set(ePairedTagNotMatch),
	WrongLogicalOperation:             dic.NewError(cWrongLogicalOperation, "операция").CodeU8().Set(eWrongLogicalOperation),
	OperatorBracketEmpty:              dic.NewError(cOperatorBracketEmpty).CodeU8().Set(eOperatorBracketEmpty),
	OperatorBracketOneItem:            dic.NewError(cOperatorBracketOneItem).CodeU8().Set(eOperatorBracketOneItem),
	NoLogicalOperationBetweenBrackets: dic.NewError(cNoLogicalOperationBetweenBrackets).CodeU8().Set(eNoLogicalOperationBetweenBrackets),
	FilterCalledByNameWasNotFound:     dic.NewError(cFilterCalledByNameWasNotFound, "название").CodeU8().Set(eFilterCalledByNameWasNotFound),
	OrmIsNil:                          dic.NewError(cOrmIsNil).CodeU8().Set(eOrmIsNil),
}

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }
