package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalOutputText(t *testing.T) {
	var color ColorHint = "mycolor"
	style := TextStyle{color, true, false, true}
	text := Text{"This is the text", 2, style}

	b, err := json.Marshal(text)

	if err != nil {
		t.Error(err)
	}

	assert.JSONEq(t, `{"kind":"text", "data":{"text":"This is the text","indent":2,"style":{"color":"mycolor","isItalic":true,"isBold":true}}}`, string(b))
}

func TestMarshalOutputHorizontalRule(t *testing.T) {
	var color ColorHint = "mycolor"
	style := TextStyle{color, true, false, true}
	text := HorizontalRule{RuleChar: "a", Style: style}

	b, err := json.Marshal(text)

	if err != nil {
		t.Error(err)
	}

	assert.JSONEq(t, `{"kind":"horizontalRule", "data":{"ruleChar":"a","style":{"color":"mycolor","isItalic":true,"isBold":true}}}`, string(b))
}
