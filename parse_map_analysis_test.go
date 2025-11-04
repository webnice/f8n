package f8n

import (
	"net/http"
	"strings"
	"testing"

	"github.com/webnice/dic"
)

// Тестирование ошибок анализатора карты сложной фильтрации.
//
//goland:noinspection HttpUrlsUsage
func TestAnalysis(t *testing.T) {
	type (
		test struct {
			Data  string
			Error dic.IError
		}
	)
	var (
		err   error
		obj   *impl
		tests []test
		n     int
		rq    *http.Request
	)

	obj = New().(*impl)
	tests = []test{
		{
			Error: Errors().PairedTagNotMatch,
			Data:  `http://host?map=g1:or:g2):or:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().PairedTagNotMatch,
			Data:  `http://host?map=)):or:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().PairedTagNotMatch,
			Data:  `http://host?map=(g1:and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().PairedTagNotMatch,
			Data:  `http://host?map=g1:and:g2:or:(g1:and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().OperatorBracketEmpty,
			Data:  `http://host?map=g1:or:g2:or:g3:or:()&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().OperatorBracketOneItem,
			Data:  `http://host?map=(g1)&g1=f1:eq:v1`,
		},
		{
			Error: Errors().NoLogicalOperationBetweenBrackets,
			Data:  `http://host?map=(g1:or:g2:or:g3)(g2:and:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().NoLogicalOperationBetweenBrackets,
			Data:  `http://host?map=g1:or:(g1:or:g2:or:g3)(g2:and:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=(g1:or:):and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=(:or:g1):and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=:and::or:`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=:or::and:`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=g1:or::and::or::or:&g1=f1:eq:v1`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=:or::and::or::or:g1&g1=f1:eq:v1`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=(g1:or:g1)g1&g1=f1:eq:v1&g2=f2:ke:v2`,
		},
		{
			Error: Errors().WrongLogicalOperation,
			Data:  `http://host?map=g1(g1:or:g2)&g1=f1:eq:v1&g2=f2:ke:v2`,
		},
	}
	for n = range tests {
		obj.Reset()
		if rq, err = http.NewRequest("GET", strings.ReplaceAll(tests[n].Data, "\n", ""), nil); err != nil {
			t.Errorf("Ошибка создания запроса: %s", err)
		}
		if err = obj.ParseRequest(rq); err == nil {
			t.Logf("данные: %s", tests[n].Data)
			t.Errorf("ParseRequest() == nil, ожидалась ошибка.")
			return
		}
		if !tests[n].Error.Is(err) {
			t.Logf("данные: %s", tests[n].Data)
			t.Errorf("ParseRequest() = %q, ожидалось: %q", err, tests[n].Error.Bind())
			return
		}
	}
}
