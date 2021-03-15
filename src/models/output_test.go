package models

import (
	"encoding/json"
	"testing"
)

func TestMarshalOutputText(t *testing.T) {
	var color ColorHint = "mycolor"
	style := TextStyle{color, true, false, true}
	text := Text{"This is the text", 2, &style}

	b, err := json.Marshal(text)

	if err != nil {
		t.Error(err)
	}

	if string(b) != `{"text":"This is the text","indent":2,"style":{"color":"mycolor","isItalic":true,"isBold":true}}` {
		t.Errorf("Got %s but expected %s", b, `{"text":"This is the text","indent":2,"style":{"color":"mycolor","isItalic":true,"isBold":true}}`)
	}
}
