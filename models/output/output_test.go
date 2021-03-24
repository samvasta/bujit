package output

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextStyleString(t *testing.T) {
	var color ColorHint = "mycolor"
	style := TextStyle{color, true, false, true}

	str := style.String()
	b, err := json.Marshal(style)

	assert.Nil(t, err, err)

	assert.Equal(t, string(b), str)
}

func TestMarshalOutputText(t *testing.T) {
	var color ColorHint = "mycolor"
	style := TextStyle{color, true, false, true}
	text := Text{"This is the text", 2, style}

	b, err := json.Marshal(text)

	if err != nil {
		t.Error(err)
	}

	assert.JSONEq(t, fmt.Sprintf(`{"kind":"text", "data":{"text":"This is the text","indent":2,"style":%s}}`, style.String()), string(b))
}

func TestMarshalOutputHorizontalRule(t *testing.T) {
	var color ColorHint = "mycolor"
	style := TextStyle{color, true, false, true}
	text := HorizontalRule{RuleChar: "a", Style: style}

	b, err := json.Marshal(text)

	if err != nil {
		t.Error(err)
	}

	assert.JSONEq(t, fmt.Sprintf(`{"kind":"horizontalRule", "data":{"ruleChar":"a","style":%s}}`, style.String()), string(b))
}

func TestMarshalOutputOrderedList(t *testing.T) {
	var color ColorHint = "mycolor"
	var bulletStyle BulletStyle = "bulletstyle"
	style := TextStyle{color, true, false, true}

	items := []Text{
		{Text: "Item 1"},
		{Text: "Item2"},
	}

	ol := OrderedList{BulletStyle: bulletStyle, Indent: 2, Style: style, Items: items}

	b, err := json.Marshal(ol)

	if err != nil {
		t.Error(err)
	}

	expected := fmt.Sprintf(
		`{
			"kind": "orderedList",
			"data":{
				"bulletStyle": "bulletstyle",
				"indent": 2,
				"style": %s,
				"items": [
					%s,
					%s
				]
			}
		}`, style.String(), items[0].String(), items[1].String())

	assert.JSONEq(t, expected, string(b))
}

func TestMarshalOutputUnorderedList(t *testing.T) {
	var color ColorHint = "mycolor"
	style := TextStyle{color, true, false, true}

	items := []Text{
		{Text: "Item 1"},
		{Text: "Item2"},
	}

	ul := UnorderedList{BulletChar: "!", Indent: 2, Style: style, Items: items}

	b, err := json.Marshal(ul)

	if err != nil {
		t.Error(err)
	}

	expected := fmt.Sprintf(
		`{
			"kind": "unorderedList",
			"data":{
				"bullet": "!",
				"indent": 2,
				"style": %s,
				"items": [
					%s,
					%s
				]
			}
		}`, style.String(), items[0].String(), items[1].String())

	assert.JSONEq(t, expected, string(b))
}

func TestEmptyOutputGroup(t *testing.T) {
	group := EmptyOutputGroup()

	slice := group.ToSlice()

	assert.Empty(t, slice)
}

func TestOutputGroupBuilder(t *testing.T) {

	style1 := TextStyle{Color: Primary, IsItalic: false, IsUnderline: true, IsBold: true}

	group := EmptyOutputGroup()

	items := group.Header("Header Text").
		HorizontalRule("=").
		Paragraph("paragraph text").
		PushStyle(style1).
		Paragraph("Styled paragraph").
		Indent().
		OrderedList([]string{"OL item1", "OL item2"}, LowerRoman).
		PopStyle().
		Paragraph("Indented, normal style").
		Unindent().
		HorizontalRule("-").
		EmptyLines(2).
		UnorderedList([]string{"UL item1", "UL item2"}, "f").
		ToSlice()

	header := items[0].(Text)
	assert.Equal(t, "Header Text", header.Text)
	assert.Equal(t, *HeaderStyle, header.Style)
	assert.Equal(t, 0, header.Indent)

	hRule1 := items[1].(HorizontalRule)
	assert.Equal(t, *DefaultStyle, hRule1.Style)
	assert.Equal(t, "=", hRule1.RuleChar)

	para1 := items[2].(Text)
	assert.Equal(t, "paragraph text", para1.Text)
	assert.Equal(t, *DefaultStyle, para1.Style)
	assert.Equal(t, 0, para1.Indent)

	para2 := items[3].(Text)
	assert.Equal(t, "Styled paragraph", para2.Text)
	assert.Equal(t, style1, para2.Style)
	assert.Equal(t, 0, para2.Indent)

	ol := items[4].(OrderedList)
	assert.Equal(t, style1, ol.Style)
	assert.Equal(t, 1, ol.Indent)
	assert.Equal(t, LowerRoman, ol.BulletStyle)
	assert.Len(t, ol.Items, 2)
	ol_item1 := ol.Items[0]
	assert.Equal(t, "OL item1", ol_item1.Text)
	assert.Equal(t, style1, ol_item1.Style)
	assert.Equal(t, 1, ol_item1.Indent)
	ol_item2 := ol.Items[1]
	assert.Equal(t, "OL item2", ol_item2.Text)
	assert.Equal(t, style1, ol_item2.Style)
	assert.Equal(t, 1, ol_item2.Indent)

	para3 := items[5].(Text)
	assert.Equal(t, "Indented, normal style", para3.Text)
	assert.Equal(t, *DefaultStyle, para3.Style)
	assert.Equal(t, 1, para3.Indent)

	hRule2 := items[6].(HorizontalRule)
	assert.Equal(t, *DefaultStyle, hRule2.Style)
	assert.Equal(t, "-", hRule2.RuleChar)

	para4 := items[7].(Text)
	assert.Equal(t, "\n\n", para4.Text)
	assert.Equal(t, *DefaultStyle, para4.Style)
	assert.Equal(t, 0, para4.Indent)

	ul := items[8].(UnorderedList)
	assert.Equal(t, *DefaultStyle, ul.Style)
	assert.Equal(t, 0, ul.Indent)
	assert.Equal(t, "f", ul.BulletChar)
	assert.Len(t, ul.Items, 2)
	ul_item1 := ul.Items[0]
	assert.Equal(t, "UL item1", ul_item1.Text)
	assert.Equal(t, *DefaultStyle, ul_item1.Style)
	assert.Equal(t, 0, ul_item1.Indent)
	ul_item2 := ul.Items[1]
	assert.Equal(t, "UL item2", ul_item2.Text)
	assert.Equal(t, *DefaultStyle, ul_item2.Style)
	assert.Equal(t, 0, ul_item2.Indent)
}
