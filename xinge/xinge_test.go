package xinge

import (
	"testing"
	"log"
)



func Test_request (t *testing.T ) {
	r := NewRequest()
	Method = "POST"
	r.SetAccessId(123456789)
	r.SetMessageType(2)
	r.SetTitle("hello")
	r.SetContent("world")
	r.SetSecret("your secret")

	resp, err := r.PushTagsAnd("room1")
	if err != nil {
		t.Error(err)
	}
//	t.Log(err, resp)
	log.Println(resp)
}
