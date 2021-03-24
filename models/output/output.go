package output

import (
	"encoding/json"
	"strings"
)

type Helper interface {
	json.Marshaler
}

type ColorHint string

const (
	Body    ColorHint = "body"
	Subtle  ColorHint = "subtle"
	Primary ColorHint = "primary"
	Success ColorHint = "success"
	Info    ColorHint = "info"
	Warning ColorHint = "warning"
	Error   ColorHint = "error"
)

type TextStyle struct {
	Color       ColorHint `json:"color"`
	IsItalic    bool      `json:"isItalic,omitempty"`
	IsUnderline bool      `json:"isUnderline,omitempty"`
	IsBold      bool      `json:"isBold,omitempty"`
}

func (ts TextStyle) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var DefaultStyle *TextStyle = &TextStyle{Color: Body, IsItalic: false, IsUnderline: false, IsBold: false}
var HeaderStyle *TextStyle = &TextStyle{Color: Body, IsItalic: false, IsUnderline: false, IsBold: true}

type Text struct {
	Text   string    `json:"text"`
	Indent int       `json:"indent"`
	Style  TextStyle `json:"style"`
}

func (t Text) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type FakeText Text // to avoid recursive JSON marshaling
func (t Text) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind string   `json:"kind"`
		Data FakeText `json:"data"`
	}{
		"text",
		FakeText(t),
	})
}

var NormalBulletChar string = "•"

type UnorderedList struct {
	BulletChar string    `json:"bullet"`
	Indent     int       `json:"indent"`
	Style      TextStyle `json:"style"`
	Items      []Text    `json:"items"`
}

func (ul UnorderedList) String() string {
	b, err := json.Marshal(ul)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type FakeUnorderedList UnorderedList // to avoid recursive JSON marshaling
func (ul UnorderedList) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind string            `json:"kind"`
		Data FakeUnorderedList `json:"data"`
	}{
		"unorderedList",
		FakeUnorderedList(ul),
	})
}

type BulletStyle string

const (
	Decimal    BulletStyle = "decimal"
	UpperAlpha BulletStyle = "upper-alpha"
	LowerAlpha BulletStyle = "lower-alpha"
	UpperRoman BulletStyle = "upper-roman"
	LowerRoman BulletStyle = "lower-roman"
)

type OrderedList struct {
	BulletStyle BulletStyle `json:"bulletStyle"`
	Indent      int         `json:"indent"`
	Style       TextStyle   `json:"style"`
	Items       []Text      `json:"items"`
}

func (ol OrderedList) String() string {
	b, err := json.Marshal(ol)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type FakeOrderedList OrderedList // to avoid recursive JSON marshaling
func (ol OrderedList) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind string          `json:"kind"`
		Data FakeOrderedList `json:"data"`
	}{
		"orderedList",
		FakeOrderedList(ol),
	})
}

type HorizontalRule struct {
	RuleChar string    `json:"ruleChar"`
	Style    TextStyle `json:"style"`
}

func (hr HorizontalRule) String() string {
	b, err := json.Marshal(hr)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type FakeHorizontalRule HorizontalRule // to avoid recursive JSON marshaling
func (hr HorizontalRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind string             `json:"kind"`
		Data FakeHorizontalRule `json:"data"`
	}{
		"horizontalRule",
		FakeHorizontalRule(hr),
	})
}

// Functions for making output

type OutputGroup struct {
	items         []Helper
	currentIndent int
	styleStack    []TextStyle
}

func EmptyOutputGroup() *OutputGroup {
	return &OutputGroup{currentIndent: 0, styleStack: []TextStyle{*DefaultStyle}}
}

func (g *OutputGroup) Indent() *OutputGroup {
	(*g).currentIndent++
	return g
}

func (g *OutputGroup) Unindent() *OutputGroup {
	(*g).currentIndent--
	return g
}

func (g *OutputGroup) PushStyle(style TextStyle) *OutputGroup {
	g.styleStack = append(g.styleStack, style)
	return g
}

func (g *OutputGroup) PopStyle() *OutputGroup {
	if len(g.styleStack) > 0 {
		g.styleStack = g.styleStack[:len(g.styleStack)-1]
	}
	return g
}

func (g *OutputGroup) CurrentStyle() TextStyle {
	return g.styleStack[len(g.styleStack)-1]
}

func (g *OutputGroup) Header(text string) *OutputGroup {
	g.items = append(g.items, Text{Text: text, Indent: g.currentIndent, Style: *HeaderStyle})
	return g
}
func (g *OutputGroup) HorizontalRule(ruleChar string) *OutputGroup {
	g.items = append(g.items, HorizontalRule{RuleChar: ruleChar, Style: g.CurrentStyle()})
	return g
}

func (g *OutputGroup) Paragraph(text string) *OutputGroup {
	g.items = append(g.items, Text{Text: text, Indent: g.currentIndent, Style: g.CurrentStyle()})
	return g
}

func (g *OutputGroup) OrderedList(items []string, bulletStyle BulletStyle) *OutputGroup {
	var listItems []Text
	for _, item := range items {
		listItems = append(listItems, Text{Text: item, Indent: g.currentIndent, Style: g.CurrentStyle()})
	}

	g.items = append(g.items, OrderedList{BulletStyle: bulletStyle, Items: listItems, Indent: g.currentIndent, Style: g.CurrentStyle()})
	return g
}

func (g *OutputGroup) UnorderedList(items []string, bulletChar string) *OutputGroup {
	var listItems []Text
	for _, item := range items {
		listItems = append(listItems, Text{Text: item, Indent: g.currentIndent, Style: g.CurrentStyle()})
	}

	g.items = append(g.items, UnorderedList{BulletChar: bulletChar, Items: listItems, Indent: g.currentIndent, Style: g.CurrentStyle()})
	return g
}

func (g *OutputGroup) EmptyLines(numLines int) *OutputGroup {
	g.items = append(g.items, Text{Text: strings.Repeat("\n", numLines), Indent: g.currentIndent, Style: g.CurrentStyle()})
	return g
}

func (g *OutputGroup) ToSlice() []Helper {
	return g.items
}
