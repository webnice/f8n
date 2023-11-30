package f8n

import (
	"net/http"
	"strings"
	"testing"
)

// Тестирование ошибок анализатора карты сложной фильтрации.
func TestAnalysis(t *testing.T) {
	type (
		test struct {
			Data  string
			Error Err
		}
	)
	var (
		err   error
		obj   *impl
		tests []test
		n     int
		rq    *http.Request
		ero   Err
		ok    bool
	)

	obj = New().(*impl)
	tests = []test{
		{
			Error: Errors().PairedTagNotMatch("", ""),
			Data:  `http://host?map=g1:or:g2):or:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().PairedTagNotMatch("", ""),
			Data:  `http://host?map=)):or:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().PairedTagNotMatch("", ""),
			Data:  `http://host?map=(g1:and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().PairedTagNotMatch("", ""),
			Data:  `http://host?map=g1:and:g2:or:(g1:and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().OperatorBracketEmpty(),
			Data:  `http://host?map=g1:or:g2:or:g3:or:()&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().OperatorBracketOneItem(),
			Data:  `http://host?map=(g1)&g1=f1:eq:v1`,
		},
		{
			Error: Errors().NoLogicalOperationBetweenBrackets(),
			Data:  `http://host?map=(g1:or:g2:or:g3)(g2:and:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().NoLogicalOperationBetweenBrackets(),
			Data:  `http://host?map=g1:or:(g1:or:g2:or:g3)(g2:and:g3)&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
			Data:  `http://host?map=(g1:or:):and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
			Data:  `http://host?map=(:or:g1):and:g2&g1=f1:eq:v1&g2=f2:ke:v2&g3=f3:ke:v3`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
			Data:  `http://host?map=:and::or:`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
			Data:  `http://host?map=:or::and:`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
			Data:  `http://host?map=g1:or::and::or::or:&g1=f1:eq:v1`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
			Data:  `http://host?map=:or::and::or::or:g1&g1=f1:eq:v1`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
			Data:  `http://host?map=(g1:or:g1)g1&g1=f1:eq:v1&g2=f2:ke:v2`,
		},
		{
			Error: Errors().WrongLogicalOperation(OriginUnknown),
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
		if ero, ok = err.(Err); !ok {
			t.Errorf("ParseRequest() = %q, не верный тип ошибки.", err)
			return
		}
		if ero.Anchor() != tests[n].Error.Anchor() {
			t.Logf("данные: %s", tests[n].Data)
			t.Errorf("ParseRequest() = %q, ожидалось: %q", ero.Error(), tests[n].Error.Error())
			return
		}
	}
}
