// Package f8n
package f8n

import (
	"net/http"
	runtimeDebug "runtime/debug"

	"github.com/webnice/kit/v2/module/verify"
)

// ParseRequest Загрузка настроек фильтрации из запроса http.Request.
// Если при разборе запроса возникает ошибка, она возвращается в качестве интерфейса error.
// Если в функцию передан не нулевой объект функции OnErrorFunc, тогда будет вызвана эта функция и в неё
// будут переданы описания возникшей ошибки.
func (f8n *impl) ParseRequest(rq *http.Request, errorFn ...OnErrorFunc) (err error) {
	var (
		vyi  verify.Interface
		ers  []*ParseError
		tmp  []*ParseError
		buf  []byte
		n, j int
	)

	// При вызове внешней функции errorFn, возможна паника.
	defer func() {
		if e := recover(); e != nil {
			err = f8n.Errors().Panic(e, runtimeDebug.Stack())
		}
	}()
	if rq == nil {
		err = f8n.Errors().RequestIsNil()
		return
	}
	// Загрузка лимита.
	if tmp = f8n.ParseLimit(rq); len(tmp) > 0 {
		ers = append(ers, tmp...)
	}
	// Загрузка сортировки.
	if tmp = f8n.ParseBy(rq); len(tmp) > 0 {
		ers = append(ers, tmp...)
	}
	// Загрузка простой фильтрации.
	if tmp = f8n.ParseFilter(rq); len(tmp) > 0 {
		ers = append(ers, tmp...)
	}
	// Загрузка сложной фильтрации.
	if tmp = f8n.ParseMap(rq); len(tmp) > 0 {
		ers = append(ers, tmp...)
	}
	if len(ers) == 0 {
		return
	}
	// Подготовка данных для вызова внешней функции errorFn для передачи ошибок.
	if err = ers[0].Ei; len(ers) > 1 {
		err = f8n.Errors().MultipleErrorsFound()
	}
	vyi = verify.E4xx().Code(-1).Message(err.Error())
	for n = range ers {
		for j = range ers[n].Ev {
			vyi.Add(ers[n].Ev[j])
		}
	}
	buf = vyi.Json()
	for n = range errorFn {
		errorFn[n](buf, err)
	}

	return
}
