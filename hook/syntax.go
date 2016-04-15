package hook

import (
	"strconv"
)

//语法关键字
var DefaultScriptSyntaxKeywords = []string{"eventtarget", //事件
	"gt",       //大于
	"lt",       //小于
	"attr",     //属性
	"callback", //回调cmd
}

//表达式数据
type ExprData struct {
	bGt   bool
	bLt   bool
	gtVal float64
	ltVal float64
}

//关键字分类用于执行语意
//var KeywordsClass = map[string][]string{"expression": []string{"gt", "lt", "eq", "condition"},
//	"targetname": []string{"cpu", "mem", "disk", "status", "load"},
//	"callback":   []string{"callback"},
//	"attr":       []string{"use_percent", "free_percent", "use", "free"},
//}

type Syntaxer interface {
	CheckSyntakKeyWords(keywords string) bool
	GetExpt(event *Event) *ExprData
	IsExpr(keywords string) bool
	CheckFloadValue(expr *ExprData, value float64) bool
}

//new  syntax
func NewDefaultSyntax() Syntaxer {
	return &ScriptSyntax{}
}

type ScriptSyntax struct {
}

//是否是表达式
func (self *ScriptSyntax) IsExpr(keywords string) bool {
	return true
}

//解析condition
func (self *ScriptSyntax) ParseCondition(condition string) bool {
	return true
}

//检查关键字
func (self *ScriptSyntax) CheckSyntakKeyWords(keywords string) bool {
	for _, v := range DefaultScriptSyntaxKeywords {
		if v == keywords {
			return true
		}
	}
	return false
}

//check fload
func (self *ScriptSyntax) CheckFloadValue(expr *ExprData, value float64) bool {
	if expr.bGt && expr.bLt {
		return (value >= expr.gtVal && value <= expr.ltVal)
	} else if expr.bGt && !expr.bLt {
		return (value >= expr.gtVal)
	} else if !expr.bGt && expr.bLt {
		return (value <= expr.ltVal)
	}
	return true
}
func (self *ScriptSyntax) GetExpt(event *Event) *ExprData {
	expr := &ExprData{}
	for exprName, exprValue := range event.EventCondition {
		switch exprName {
		case "gt":
			{
				expr.bGt = true
				expr.gtVal, _ = strconv.ParseFloat(exprValue, 64)
			}
		case "lt":
			{
				expr.bLt = true
				expr.ltVal, _ = strconv.ParseFloat(exprValue, 64)
			}
		}
	}
	return expr
}
