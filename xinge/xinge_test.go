package xinge

import (
	"testing"
	"log"
)



func Test_request (t *testing.T ) {
	r := NewRequest()

	resp, err := r.PushTagsAnd("room1")
	if err != nil {
		t.Error(err)
	}
//	t.Log(err, resp)
	log.Println(resp)
}
