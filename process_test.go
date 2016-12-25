package gimage

import "testing"

func TestProcess(t *testing.T) {
	buf, err := ReadFile("fixtures/canyon.jpg")
	if err != nil {
		t.Fail()
	}
	options := &Options{}
	Process(buf, options)
}
