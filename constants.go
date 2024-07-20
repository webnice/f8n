package f8n

const (
	keyLimit     = "limit"  // Наименование параметра с лимитом.
	keyBy        = "by"     // Наименование параметра с сортировкой.
	keyFilter    = "filter" // Наименование параметра с фильтром.
	keyTie       = "tie"    // Устаревший режим работы фильтрации. Для совместимости.
	keyMap       = "map"    // Наименование параметра с логической картой.
	keyDelimiter = ":"      // Разделитель и ограничитель.
	keyOr        = ":or:"   // Логическое "ИЛИ".
	keyAnd       = ":and:"  // Логическое "И".
	keyNilValue  = "nil"    // Значение определяющие nil.
	pairBeg      = "("      // Начало операторных скобок.
	pairEnd      = ")"      // Окончание операторных скобок.
)

const (
	// Неизвестный метод фильтрации.
	filterUnknown = FilterMethod("")

	// Равно.
	filterEquivalent = FilterMethod("eq")

	// Меньше.
	filterLessThan = FilterMethod("lt")

	// Меньше или равно.
	filterLessThanOrEquivalent = FilterMethod("le")

	// Больше.
	filterGreaterThan = FilterMethod("gt")

	// Больше или равно.
	filterGreaterThanOrEquivalent = FilterMethod("ge")

	// Не равно.
	filterNotEquivalent = FilterMethod("ne")

	// Похоже.
	filterLikeThan = FilterMethod("ke")

	// Не похоже.
	filterNotLikeThan = FilterMethod("kn")

	// Идентификаторы из списка.
	filterIn = FilterMethod("in")

	// Идентификаторы не из списка.
	filterNotIn = FilterMethod("ni")
)

const (
	// Неизвестный порядок сортировки.
	byUnknown = ByDirection("")

	// Прямой порядок сортировки.
	byAsc = ByDirection("asc")

	// Обратный порядок сортировки.
	byDesc = ByDirection("desc")
)

const (
	// Неизвестный режим.
	tieUnknown = TieMode("")

	// Режим "И" (по умолчанию).
	tieAnd = TieMode("and")

	// Режим "ИЛИ".
	tieOr = TieMode("or")
)

const (
	// TypeString Тип поля string, по умолчанию для всех полей у которых не указан иной тип данных.
	TypeString = FieldType("string")

	// TypeStringNil Тип поля string, который в базе данных может быть NULL.
	TypeStringNil = FieldType("string-or-nil")

	// TypeUint64 Тип поля uint64.
	TypeUint64 = FieldType("uint64")

	// TypeUint64Nil Тип поля uint64, который в базе данных может быть NULL.
	TypeUint64Nil = FieldType("uint64-or-nil")

	// TypeInt64 Тип поля int64.
	TypeInt64 = FieldType("int64")

	// TypeInt64Nil Тип поля int64, который в базе данных может быть NULL.
	TypeInt64Nil = FieldType("int64-or-nil")

	// TypeFloat64 Тип поля float64.
	TypeFloat64 = FieldType("float64")

	// TypeFloat64Nil Тип поля float64, который в базе данных может быть NULL.
	TypeFloat64Nil = FieldType("float64-or-nil")

	// TypeBool Тип поля bool.
	TypeBool = FieldType("bool")

	// TypeBoolNil Тип поля bool, который в базе данных может быть NULL.
	TypeBoolNil = FieldType("bool-or-nil")

	// TypeTime Тип поля time.
	TypeTime = FieldType("time")

	// TypeTimeNil Тип поля time, который в базе данных может быть NULL.
	TypeTimeNil = FieldType("time-or-nil")

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
