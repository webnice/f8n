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
	sepInNi      = ","      // Разделитель значений для "IN" и "NI".
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

	// Значения из списка.
	filterIn = FilterMethod("in")

	// Значения не из списка.
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

	// TypeSliceString Тип поля среза строк для SQL IN.
	TypeSliceString = FieldType("slice_of_strings")

	// TypeSliceInt64 Тип поля среза int64 для SQL IN.
	TypeSliceInt64 = FieldType("slice_of_int64")

	// TypeSliceUint64 Тип поля среза uint64 для SQL IN.
	TypeSliceUint64 = FieldType("slice_of_uint64")
)

// Типы данных полей соответствующие типам пакета github.com/webnice/nix.
const (
	// TypeInt Тип поля int.
	TypeInt = FieldType("int")

	// TypeIntNil Тип поля int, который в базе данных может быть NULL.
	TypeIntNil = FieldType("int-or-nil")

	// TypeInt8 Тип поля int8.
	TypeInt8 = FieldType("int8")

	// TypeInt8Nil Тип поля int8, который в базе данных может быть NULL.
	TypeInt8Nil = FieldType("int8-or-nil")

	// TypeInt16 Тип поля int16.
	TypeInt16 = FieldType("int16")

	// TypeInt16Nil Тип поля int16, который в базе данных может быть NULL.
	TypeInt16Nil = FieldType("int16-or-nil")

	// TypeInt32 Тип поля int32.
	TypeInt32 = FieldType("int32")

	// TypeInt32Nil Тип поля int32, который в базе данных может быть NULL.
	TypeInt32Nil = FieldType("int32-or-nil")

	// TypeUint Тип поля uint.
	TypeUint = FieldType("uint")

	// TypeUintNil Тип поля uint, который в базе данных может быть NULL.
	TypeUintNil = FieldType("uint-or-nil")

	// TypeUint8 Тип поля uint8.
	TypeUint8 = FieldType("uint8")

	// TypeUint8Nil Тип поля uint8, который в базе данных может быть NULL.
	TypeUint8Nil = FieldType("uint8-or-nil")

	// TypeUint16 Тип поля uint16.
	TypeUint16 = FieldType("uint16")

	// TypeUint16Nil Тип поля uint16, который в базе данных может быть NULL.
	TypeUint16Nil = FieldType("uint16-or-nil")

	// TypeUint32 Тип поля uint32.
	TypeUint32 = FieldType("uint32")

	// TypeUint32Nil Тип поля uint32, который в базе данных может быть NULL.
	TypeUint32Nil = FieldType("uint32-or-nil")

	// TypeFloat32 Тип поля float32.
	TypeFloat32 = FieldType("float32")

	// TypeFloat32Nil Тип поля float32, который в базе данных может быть NULL.
	TypeFloat32Nil = FieldType("float32-or-nil")

	// TypeBytes Тип поля среза байт, в запросе значение кодируется в base64.
	TypeBytes = FieldType("bytes")

	// TypeBytesNil Тип поля среза байт, который в базе данных может быть NULL.
	TypeBytesNil = FieldType("bytes-or-nil")

	// TypeDuration Тип поля time.Duration, в запросе значение в формате time.ParseDuration.
	TypeDuration = FieldType("duration")

	// TypeDurationNil Тип поля time.Duration, который в базе данных может быть NULL.
	TypeDurationNil = FieldType("duration-or-nil")

	// TypeRawMessage Тип поля json.RawMessage.
	TypeRawMessage = FieldType("json")

	// TypeRawMessageNil Тип поля json.RawMessage, который в базе данных может быть NULL.
	TypeRawMessageNil = FieldType("json-or-nil")

	// TypeIP Тип поля net.IP.
	TypeIP = FieldType("ip")

	// TypeIPNil Тип поля net.IP, который в базе данных может быть NULL.
	TypeIPNil = FieldType("ip-or-nil")

	// TypeBigInt Тип поля big.Int, десятичная строка.
	TypeBigInt = FieldType("bigint")

	// TypeBigIntNil Тип поля big.Int, который в базе данных может быть NULL.
	TypeBigIntNil = FieldType("bigint-or-nil")

	// TypeUUID Тип поля UUID, строка формата xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
	TypeUUID = FieldType("uuid")

	// TypeUUIDNil Тип поля UUID, который в базе данных может быть NULL.
	TypeUUIDNil = FieldType("uuid-or-nil")

	// TypeComplex64 Тип поля complex64.
	TypeComplex64 = FieldType("complex64")

	// TypeComplex64Nil Тип поля complex64, который в базе данных может быть NULL.
	TypeComplex64Nil = FieldType("complex64-or-nil")

	// TypeComplex128 Тип поля complex128.
	TypeComplex128 = FieldType("complex128")

	// TypeComplex128Nil Тип поля complex128, который в базе данных может быть NULL.
	TypeComplex128Nil = FieldType("complex128-or-nil")
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

const (
	// SqlDialectMySQL Диалект базы данных MySQL.
	SqlDialectMySQL = SqlDialect("mysql")

	// SqlDialectMariaDB Диалект базы данных MariaDB.
	SqlDialectMariaDB = SqlDialect("mariadb")

	// SqlDialectPostgreSQL Диалект базы данных PostgreSQL.
	SqlDialectPostgreSQL = SqlDialect("postgres")

	// SqlDialectSqlite3 Диалект базы данных sqlite3.
	SqlDialectSqlite3 = SqlDialect("sqlite3")
)
