package f8n

import "testing"

// Тестирование надстройки над strconv для int64.
func TestParseInt64(t *testing.T) {
	const (
		tplErrorResult = "parseInt64(`%s`, `%s`) = %d, %v; ожидалось: %d, %v."
		datField       = "test"
		datEmpty       = ""
		datNegative    = "-1"
		datInvalid     = "a-1"
	)

	if v, e := parseInt64(datEmpty, datEmpty); v != 0 || e != nil {
		t.Errorf(tplErrorResult, datField, datEmpty, v, e, 0, nil)
	}
	if v, e := parseInt64(datField, datNegative); v != -1 || e != nil {
		t.Errorf(tplErrorResult, datField, datNegative, v, e, -1, nil)
	}
	if v, e := parseInt64(datField, datInvalid); v != 0 || e == nil {
		t.Errorf(tplErrorResult, datField, datInvalid, v, e, 0, "error")
	}
}

// Тестирование защиты наименования поля базы данных от символов опасных при формировании SQL запроса.
func TestFieldNameCleanInjection(t *testing.T) {
	const (
		tplErrorResult = "fieldName(%q) = %q, ожидалось: %q."
		datRemap       = "`cursor`.`id`"
	)
	var (
		n     int
		obj   *impl
		tests = []struct {
			Source string
			Want   string
		}{
			{Source: "id", Want: "`id`"},
			{Source: "name", Want: "`name`"},
			// Попытки выхода за пределы идентификатора поля.
			{Source: "`", Want: "``"},
			{Source: "` IS NULL--", Want: "` IS NULL--`"},
			{Source: "`=1 OR `1", Want: "`=1 OR 1`"},
			{Source: "a`, `b", Want: "`a b`"},
			{Source: "a' OR '1'='1", Want: "`a OR 1=1`"},
			{Source: `a" OR "1"="1`, Want: "`a OR 1=1`"},
			{Source: "name` = `x` OR `1", Want: "`name = x OR 1`"},
		}
	)

	obj = New().(*impl)
	for n = range tests {
		if ret := obj.fieldName(tests[n].Source); ret != tests[n].Want {
			t.Errorf(tplErrorResult, tests[n].Source, ret, tests[n].Want)
			return
		}
	}
	// Переопределение поля является доверенным вводом и используется без изменений.
	obj.Redefinition("id", datRemap)
	if ret := obj.fieldName("id"); ret != datRemap {
		t.Errorf(tplErrorResult, "id", ret, datRemap)
		return
	}
}
