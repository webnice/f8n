package f8n

import (
	"strings"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
)

// Анализатор, выделение сущности, проверка.
func (f8n *impl) analysis(node *Map) (ret []*ParseError) {
	var (
		fxp func([]*ParseError, Err) []*ParseError
		n   int
	)

	// Функция сборки ошибки.
	fxp = func(r []*ParseError, e Err) []*ParseError {
		var eo = &ParseError{Ei: e}
		eo.Ev = append(eo.Ev, kitModuleAns.RestErrorField{
			Field:   keyFilter,
			Message: eo.Ei.Error(),
		})
		return append(r, eo)
	}
	// Присвоение сущности по умолчанию.
	if node.TagBegin == OriginUnknown.String() && node.TagEnd == OriginUnknown.String() && node.Content != "" {
		node.Origin = OriginFiltration
		return
	}
	for n = range OriginAll {
		if strings.EqualFold(OriginAll[n].String(), strings.ToLower(node.TagBegin)) && node.Origin == OriginUnknown {
			node.Origin = OriginAll[n]
		}
		if strings.EqualFold(node.TagBegin, pairBeg) && node.Origin == OriginUnknown {
			node.Origin = OriginOperatorBracket
		}
	}
	defer func() {
		if len(ret) > 0 {
			return
		}
		switch node.Origin {
		case OriginUnknown:
			if node.TagBegin == pairBeg || node.TagEnd == pairEnd {
				ret = fxp(ret, f8n.Errors().PairedTagNotMatch(node.TagBegin, node.TagEnd))
				return
			}
		case OriginOperatorBracket:
			if node.TagBegin != pairBeg || node.TagEnd != pairEnd {
				ret = fxp(ret, f8n.Errors().PairedTagNotMatch(node.TagBegin, node.TagEnd))
				return
			}
		}
	}()
	// Проверка операторных скобок.
	if node.Origin == OriginOperatorBracket && len(node.Node) == 0 {
		ret = fxp(ret, f8n.Errors().OperatorBracketEmpty())
		return
	}
	if node.Origin == OriginOperatorBracket && len(node.Node) == 1 {
		ret = fxp(ret, f8n.Errors().OperatorBracketOneItem())
		return
	}
	// Между операторными скобками должна быть логическая операция.
	for n = range node.Node {
		if node.Node[n].Origin == OriginOperatorBracket {
			switch n {
			case 0:
				if len(node.Node) > 1 && node.Node[n+1].Origin == OriginOperatorBracket {
					ret = fxp(ret, f8n.Errors().NoLogicalOperationBetweenBrackets())
					return
				}
			default:
				if node.Node[n-1].Origin == OriginOperatorBracket {
					ret = fxp(ret, f8n.Errors().NoLogicalOperationBetweenBrackets())
					return
				}
			}
		}
	}
	// Проверка корректности применения логических операций.
	for n = range node.Node {
		if node.Node[n].Origin == OriginOr || node.Node[n].Origin == OriginAnd {
			// Не указан второй аргумент для логической операции.
			if n == 0 || n == len(node.Node)-1 {
				ret = fxp(ret, f8n.Errors().WrongLogicalOperation(node.Node[n].Origin))
				return
			}
			// Логическая операция применяется к логической операции (справа).
			if node.Node[n+1].Origin == OriginOr || node.Node[n+1].Origin == OriginAnd {
				ret = fxp(ret, f8n.Errors().WrongLogicalOperation(node.Node[n].Origin))
				return
			}
		}
		// С одной из сторон от наименования фильтрации должна быть логическая операция.
		if node.Node[n].Origin == OriginFiltration {
			switch n {
			case 0:
				if len(node.Node) > 1 && (node.Node[n+1].Origin != OriginOr && node.Node[n+1].Origin != OriginAnd) {
					ret = fxp(ret, f8n.Errors().WrongLogicalOperation(node.Node[n].Origin))
					return
				}
			default:
				if node.Node[n-1].Origin != OriginOr && node.Node[n-1].Origin != OriginAnd {
					ret = fxp(ret, f8n.Errors().WrongLogicalOperation(node.Node[n].Origin))
					return
				}
			}
		}
	}

	return
}
