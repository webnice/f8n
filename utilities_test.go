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
