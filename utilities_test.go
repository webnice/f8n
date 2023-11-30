package f8n

import "testing"

// Тестирование надстройки над strconv для int64.
func TestParseInt64(t *testing.T) {
	if v, e := parseInt64("", ""); v != 0 || e != nil {
		t.Errorf("parseInt64(``, ``) = %d, %s; want 0, nil", v, e)
	}
	if v, e := parseInt64("test", "-1"); v != -1 || e != nil {
		t.Errorf("parseInt64(`test`, `-1`) = %d, %s; want -1, nil", v, e)
	}
	if v, e := parseInt64("test", "a-1"); v != 0 || e == nil {
		t.Errorf("parseInt64(`test`, `a-1`) = %d, %s; want 0, error", v, e)
	}
}
