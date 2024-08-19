package f8n

import (
	"regexp"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

var (
	// Регулярное выражение описывающее все возможные теги для разбора карты условий сложной фильтрации.
	rexTag = regexp.MustCompile(`(?mi)(\(|\)|:and:|:or:)`)

	// OriginAll Справочник всех значений сущностей тегов сложной фильтрации.
	OriginAll = []Origin{
		OriginFiltration, OriginOperatorBracket, OriginAnd, OriginOr,
	}
)

// Объект сущности, реализующий интерфейс Interface.
type impl struct {
	FieldOnly []string    `json:"field"`         // Имена обрабатываемых полей. Пусто - обрабатываются все поля.
	Datatype  Datatype    `json:"field_type"`    // Карта типов данных полей.
	Remap     Remap       `json:"remap"`         // Переопределение наименований полей.
	Offset    uint64      `json:"offset"`        // Лимит - позиция выборки.
	Limit     uint64      `json:"limit"`         // Лимит - размер выборки.
	By        []Direction `json:"by"`            // Опции сортировки результата выборки.
	Filter    []Filter    `json:"filter"`        // Простая фильтрация.
	Tie       TieMode     `json:"tie,omitempty"` // Устаревший режим простой фильтрации.
	Map       *Map        `json:"map"`           // Сложная фильтрация разобранная в карту DOM объектов.
}

// ParseError Ошибки возникшие в результате разбора входящих параметров.
type ParseError struct {
	Ei error
	Ev []kitModuleAns.RestErrorField
}

// OnErrorFunc Функция, вызываемая при возникновении ошибки при разборе запроса.
// Функция получит:
// []byte - Готовые данные с описанием ошибки.
// error  - Интерфейс возникшей ошибки.
type OnErrorFunc func([]byte, error)

// FieldType Тип данных для поля.
type FieldType string

// Datatype Тип данных для карты типов данных полей.
type Datatype map[string]FieldType

// Remap Тип данных для переопределения наименований полей.
type Remap map[string]string

// SqlDialect Тип диалекта базы данных.
type SqlDialect string

// Direction Сортировка.
type Direction struct {
	Field string      `json:"field"`
	By    ByDirection `json:"direction"`
}

// ByDirection Тип направления сортировки.
type ByDirection string

// Filter Структура настроек фильтрации.
type Filter struct {
	Field  string       `json:"field"`
	Method FilterMethod `json:"method"`
	Value  FilterValue  `json:"value"`
}

// FilterMethod Тип способа сравнения.
type FilterMethod string

// FilterValue Значение фильтрации.
type FilterValue struct {
	Source string    `json:"source"` // Исходное значение фильтрации.
	Type   FieldType `json:"type"`   // Тип значения.
}

// TieMode Устаревший режим работы простой фильтрации.
type TieMode string
