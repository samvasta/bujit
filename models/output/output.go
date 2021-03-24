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

var DefaultStyle *TextStyle = &TextStyle{Color: Body, IsItalic: false, IsUnderline: false, IsBold: false}
var HeaderStyle *TextStyle = &TextStyle{Color: Body, IsItalic: false, IsUnderline: false, IsBold: true}

type Text struct {
	Text   string    `json:"text"`
	Indent int       `json:"indent"`
	Style  TextStyle `json:"style"`
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

type UnorderedList struct {
	BulletChar string    `json:"bullet"`
	Indent     int       `json:"indent"`
	Style      TextStyle `json:"style"`
	Items      []Text    `json:"items"`
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

func (g *OutputGroup) Indent() {
	(*g).currentIndent++
}

func (g *OutputGroup) Unindent() {
	(*g).currentIndent--
}

func (g *OutputGroup) PushStyle(style TextStyle) {
	g.styleStack = append(g.styleStack, style)
}

func (g *OutputGroup) PopStyle() {
	if len(g.styleStack) > 0 {
		g.styleStack = g.styleStack[:len(g.styleStack)-1]
	}
}

func (g *OutputGroup) CurrentStyle() TextStyle {
	return g.styleStack[len(g.styleStack)-1]
}

func (g *OutputGroup) Header(text string) *OutputGroup {
	g.items = append(g.items, &Text{Text: text, Indent: g.currentIndent, Style: *HeaderStyle})
	return g
}
func (g *OutputGroup) HorizontalRule(ruleChar string) *OutputGroup {
	g.items = append(g.items, &HorizontalRule{RuleChar: ruleChar, Style: g.CurrentStyle()})
	return g
}

func (g *OutputGroup) Paragraph(text string) *OutputGroup {
	g.items = append(g.items, &Text{Text: text, Indent: g.currentIndent, Style: g.CurrentStyle()}, g.EmptyLine(1))
	return g
}

func (g *OutputGroup) OrderedList(items []string) *OutputGroup {
	var listItems []Text
	for _, item := range items {
		listItems = append(listItems, Text{Text: item, Indent: g.currentIndent, Style: g.CurrentStyle()})
	}

	g.items = append(g.items, &OrderedList{BulletStyle: Decimal, Items: listItems, Indent: g.currentIndent, Style: g.CurrentStyle()}, g.EmptyLine(1))
	return g
}

func (g *OutputGroup) UnorderedList(items []string) *OutputGroup {
	var listItems []Text
	for _, item := range items {
		listItems = append(listItems, Text{Text: item, Indent: g.currentIndent, Style: g.CurrentStyle()})
	}

	g.items = append(g.items, &UnorderedList{BulletChar: "•", Items: listItems, Indent: g.currentIndent, Style: g.CurrentStyle()}, g.EmptyLine(1))
	return g
}

func (g *OutputGroup) EmptyLine(numLines int) Text {
	return Text{Text: strings.Repeat("\n", numLines), Indent: g.currentIndent, Style: g.CurrentStyle()}
}

func (g *OutputGroup) ToSlice() []Helper {
	return g.items
}
