package testdb

import "testing"

// Тестирование тестового окружения базы данных.
func TestTestdb(t *testing.T) {
	var err error

	// Проверка соединения с тестовой базой данных.
	if DB == nil {
		t.Errorf("DB == nil, ожидалось соединение с тестовой базой данных.")
		return
	}
	if err = DB.Exec("SELECT 1").Error; err != nil {
		t.Errorf("Ошибка выполнения тестового запроса: %s.", err)
		return
	}
	// Проверка повторного пересоздания таблиц.
	Migrations()
	if !DB.Migrator().HasTable(&Parent{}) || !DB.Migrator().HasTable(&Child{}) {
		t.Errorf("Таблицы Parent и Child не найдены после миграции.")
		return
	}
}
