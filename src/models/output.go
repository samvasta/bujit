package models

type OutputSerializer interface {
	Serialize() string
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

var DefaultStyle *TextStyle = &TextStyle{Body, false, false, false}

type Text struct {
	Text   string     `json:"text"`
	Indent int        `json:"indent"`
	Style  *TextStyle `json:"style"`
}

type UnorderedList struct {
	BulletChar rune      `json:"bullet"`
	Indent     int       `json:"indent"`
	Style      TextStyle `json:"style"`
	Items      []Text    `json:"items"`
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

type HorizontalRule struct {
	RuleChar rune      `json:"ruleChar"`
	Style    TextStyle `json:"style"`
}
