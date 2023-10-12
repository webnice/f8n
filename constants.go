// Package f8n
package f8n

const (
	keyLimit     = "limit"  // Наименование параметра с лимитом.
	keyBy        = "by"     // Наименование параметра с сортировкой.
	keyFilter    = "filter" // Наименование параметра с фильтром.
	keyMap       = "map"    // Наименование параметра с логической картой.
	keyDelimiter = ":"      // Разделитель и ограничитель.
	keyOr        = ":or:"   // Логическое "ИЛИ".
	keyAnd       = ":and:"  // Логическое "И".
	pairBeg      = "("      // Начало операторных скобок.
	pairEnd      = ")"      // Окончание операторных скобок.
)

const (
	// Неизвестный метод фильтрации.
	filterUnknown = FilterMethod(``)

	// Равно.
	filterEquivalent = FilterMethod(`eq`)

	// Меньше.
	filterLessThan = FilterMethod(`lt`)

	// Меньше или равно.
	filterLessThanOrEquivalent = FilterMethod(`le`)

	// Больше.
	filterGreaterThan = FilterMethod(`gt`)

	// Больше или равно.
	filterGreaterThanOrEquivalent = FilterMethod(`ge`)

	// Не равно.
	filterNotEquivalent = FilterMethod(`ne`)

	// Похоже.
	filterLikeThan = FilterMethod(`ke`)

	// Не похоже.
	filterNotLikeThan = FilterMethod(`kn`)

	// Идентификаторы из списка.
	filterIn = FilterMethod(`in`)

	// Идентификаторы не из списка.
	filterNotIn = FilterMethod(`ni`)
)

const (
	// Неизвестный порядок сортировки.
	byUnknown = ByDirection(``)

	// Прямой порядок сортировки.
	byAsc = ByDirection(`asc`)

	// Обратный порядок сортировки.
	byDesc = ByDirection(`desc`)
)

const (
	// TypeString Тип поля string, по умолчанию для всех полей у которых не указан иноф тип данных.
	TypeString = FieldType("string")

	// TypeUint64 Тип поля uint64.
	TypeUint64 = FieldType("uint64")

	// TypeInt64 Тип поля int64.
	TypeInt64 = FieldType("int64")

	// TypeFloat64 Тип поля float64.
	TypeFloat64 = FieldType("float64")

	// TypeBool Тип поля bool.
	TypeBool = FieldType("bool")

	// TypeTime Тип поля time.
	TypeTime = FieldType("time")

	// TypeArrayInt64 Тип поля int64 для SQL IN.
	//TypeArrayInt64 = FieldType("array_int64")

	// TypeArrayUint64 Тип поля uint64 для SQL IN.
	//TypeArrayUint64 = FieldType("array_uint64")
)

const (
	// OriginUnknown Неизвестный тег.
	OriginUnknown = Origin("")

	// OriginFiltration Наименование фильтрации.
	OriginFiltration = Origin("filtration")

	// OriginOperatorBracket Тег операторных скобок.
	OriginOperatorBracket = Origin("operator_bracket")

	// OriginAnd Тег логического "и".
	OriginAnd = Origin(keyAnd)

	// OriginOr Тег логического "или".
	OriginOr = Origin(keyOr)
)
