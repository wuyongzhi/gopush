package jpush

import "testing"




func Test_Message(t *testing.T) {
	m := Message{}
	response, err := m.Send(nil)

	t.Log(response, err)
}
