package f8n

import "testing"

// Тестирование создания SQL запроса с фильтрацией данных.
func TestSql(t *testing.T) {
	const dialect = "sqlite"
	var (
		f8n Interface
		ret string
		err error
	)

	f8n = New()
	if ret, err = f8n.Sql(dialect); err != nil {
		t.Errorf("Sql() ошибка: %s, ошибка не ожидалась.", err)
		return
	}
	if ret != "" {
		t.Errorf("Sql() = %q, ожидалась пустая строка.", ret)
		return
	}
}
