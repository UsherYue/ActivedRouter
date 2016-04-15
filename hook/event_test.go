package hook

import (
	"testing"
)

//测试
func Test_event(t *testing.T) {
	ParseHookScript("../config/hook.json")
	DispatchEvent()
}
