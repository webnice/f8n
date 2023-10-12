# Filtration

---
[![GoDoc](https://godoc.org/gopkg.in/webnice/f8n.v1/f8n?status.svg)](https://godoc.org/gopkg.in/webnice/f8n.v1/f8n)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/webnice/f8n.v1)](https://goreportcard.com/report/gopkg.in/webnice/f8n.v1)
[![Coverage Status](https://coveralls.io/repos/github/webnice/f8n/badge.svg?branch=v1)](https://coveralls.io/github/webnice/f8n?branch=v1)

### filtration -> f[iltratio]n -> f[8]n -> f8n

Библиотека golang.

---

### Фильтрация, сортировка, пагинация, лимиты.

Библиотека для работы с реляционными базами данных, реализующая формированием SQL запроса с фильтрацией, сортировкой, лимитами и постраничным выводом, используя данные полученные через параметры URN запроса.

---

## Подключение пакета.

```bash
  go get github.com/webnice/f8n
```

## Использование пакета в приложении.

```go
package main

import (
  "github.com/webnice/f8n"
)
```
