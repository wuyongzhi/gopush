package jpush

import "testing"




func Test_Request(t *testing.T) {
	m := NewRequest()
	m.AppKey("app_key")
	m.SendNo(1)
	m.ReceiverType(ReceiverTypeBoardcast)
//	m.ReceiverValue("value1","value2")
	m.Platform("android")
	m.Sign("master_secret")
	m.MsgType(MsgTypeCustom)
	m.MsgContent(0, "", "hello,world", "")

	response, err := m.Send()

	t.Log(response, err)
}
