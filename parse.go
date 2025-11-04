package f8n

import (
	"bytes"
	"net/http"
	runtimeDebug "runtime/debug"

	"github.com/webnice/dic"
	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// ParseRequest Загрузка настроек фильтрации из запроса http.Request.
// Если при разборе запроса возникает ошибка, она возвращается в качестве интерфейса error.
// Если в функцию передан не нулевой объект функции OnErrorFunc, тогда будет вызвана эта функция и в неё
// будут переданы описания возникшей ошибки.
func (f8n *impl) ParseRequest(rq *http.Request, errorFn ...OnErrorFunc) (err error) {
	var (
		rei  kitModuleAns.RestErrorInterface
		ers  []*ParseError
		tmp  []*ParseError
		buf  *bytes.Buffer
		n, j int
	)

	// При вызове внешней функции errorFn, возможна паника.
	defer func() {
		if e := recover(); e != nil {
			err = f8n.Errors().Panic.Bind(e, runtimeDebug.Stack())
		}
	}()
	if rq == nil {
		err = f8n.Errors().RequestIsNil.Bind()
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
	// Если нет сложной фильтрации, ДЛЯ СОВМЕСТИМОСТИ - загрузка устаревшего режима фильтрации - tie.
	// Переключение режима фильтрации TIE не совместимо с MAP.
	if f8n.Map == nil {
		// Загрузка устаревшего режима фильтрации - tie.
		if tmp = f8n.ParseTie(rq); len(tmp) > 0 {
			ers = append(ers, tmp...)
		}
	}
	if len(ers) == 0 {
		return
	}
	// Подготовка данных для вызова внешней функции errorFn для передачи ошибок.
	if err = ers[0].Ei; len(ers) > 1 {
		err = f8n.Errors().MultipleErrorsFound.Bind()
	}
	rei = kitModuleAns.New(nil).
		NewRestError(dic.Status().BadRequest, err).
		CodeSet(-1)
	for n = range ers {
		for j = range ers[n].Ev {
			rei.AddWithKey(ers[n].Ev[j].Field, ers[n].Ev[j].FieldValue, ers[n].Ev[j].Message, ers[n].Ev[j].I18nKey)
		}
	}
	buf, _ = rei.JsonBytes()
	for n = range errorFn {
		errorFn[n](buf.Bytes(), err)
	}

	return
}
