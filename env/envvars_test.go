// derived from github.com/kelseyhightower/envconfig
// Original: Copyright(C) Kelsey Hightower, all rights reserved
// derivation: Copyright(C) Larry Rau. all rights reserved
// see LICENSE below.
package env

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type Specification struct {
	Debug                        bool
	Port                         int
	Rate                         float32
	User                         string
	TTL                          uint32
	Timeout                      time.Duration
	AdminUsers                   []string
	MagicNumbers                 []int
	ColorCodes                   map[string]int
	MultiWordVar                 string
	SomePointer                  *string
	SomePointerWithDefault       *string `default:"foo2baz" desc:"foorbar is the word"`
	MultiWordVarWithAlt          string  `alias:"MULTI_WORD_VAR_WITH_ALT" desc:"what alt"`
	MultiWordVarWithLowerCaseAlt string  `alias:"multi_word_var_with_lower_case_alt"`
	NoPrefixWithAlt              string  `alias:"SERVICE_HOST"`
	DefaultVar                   string  `default:"foobar"`
	RequiredVar                  string  `require:"True"`
	NoPrefixDefault              string  `alias:"BROKER" default:"127.0.0.1"`
	RequiredDefault              string  `require:"true" default:"foo2bar"`
	Ignored                      string  `ignore:"true"`
	Datetime                     time.Time
	MapField                     map[string]string `default:"one:two,three:four"`
}

type BadSpecHasStruct struct {
	foo                 string
	bar                 string
	NestedSpecification struct {
		Property            string `alias:"inner"`
		PropertyWithDefault string
	} `alias:"outer"`
	AfterNested string
}

type BadSpecHasEmbedded struct {
	Embedded
	Foo int
}

type Embedded struct {
	Enabled bool `desc:"some embedded value"`
}

func TestLoad(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_DEBUG", "true")
	os.Setenv("EV_PORT", "8080")
	os.Setenv("EV_RATE", "0.5")
	os.Setenv("EV_USER", "user")
	os.Setenv("EV_TIMEOUT", "2m")
	os.Setenv("EV_ADMINUSERS", "user1,user2,user3")
	os.Setenv("EV_MAGICNUMBERS", "5,10,20")
	os.Setenv("EV_COLORCODES", "red:1,green:2,blue:3")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("EV_TTL", "30")
	os.Setenv("EV_REQUIREDVAR", "foo")
	os.Setenv("EV_IGNORED", "was-not-ignored")
	os.Setenv("EV_OUTER_INNER", "iamnested")
	os.Setenv("EV_AFTERNESTED", "after")
	os.Setenv("EV_HONOR", "honor")
	os.Setenv("EV_DATETIME", "2016-08-16T18:57:05Z")

	err := Load("ev", &s)
	if err != nil {
		t.Error(err.Error())
	}
	if s.NoPrefixWithAlt != "127.0.0.1" {
		t.Errorf("expected %v, got %v", "127.0.0.1", s.NoPrefixWithAlt)
	}
	if !s.Debug {
		t.Errorf("expected %v, got %v", true, s.Debug)
	}
	if s.Port != 8080 {
		t.Errorf("expected %d, got %v", 8080, s.Port)
	}
	if s.Rate != 0.5 {
		t.Errorf("expected %f, got %v", 0.5, s.Rate)
	}
	if s.TTL != 30 {
		t.Errorf("expected %d, got %v", 30, s.TTL)
	}
	if s.User != "user" {
		t.Errorf("expected %s, got %s", "user", s.User)
	}
	if s.Timeout != 2*time.Minute {
		t.Errorf("expected %s, got %s", 2*time.Minute, s.Timeout)
	}
	if s.RequiredVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.RequiredVar)
	}
	if len(s.AdminUsers) != 3 ||
		s.AdminUsers[0] != "user1" ||
		s.AdminUsers[1] != "user2" ||
		s.AdminUsers[2] != "user3" {
		t.Errorf("expected %#v, got %#v", []string{"John", "Adam", "Will"}, s.AdminUsers)
	}
	if len(s.MagicNumbers) != 3 ||
		s.MagicNumbers[0] != 5 ||
		s.MagicNumbers[1] != 10 ||
		s.MagicNumbers[2] != 20 {
		t.Errorf("expected %#v, got %#v", []int{5, 10, 20}, s.MagicNumbers)
	}
	if s.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}

	if len(s.ColorCodes) != 3 ||
		s.ColorCodes["red"] != 1 ||
		s.ColorCodes["green"] != 2 ||
		s.ColorCodes["blue"] != 3 {
		t.Errorf(
			"expected %#v, got %#v",
			map[string]int{
				"red":   1,
				"green": 2,
				"blue":  3,
			},
			s.ColorCodes,
		)
	}

	if expected := time.Date(2016, 8, 16, 18, 57, 05, 0, time.UTC); !s.Datetime.Equal(expected) {
		t.Errorf("expected %s, got %s", expected.Format(time.RFC3339), s.Datetime.Format(time.RFC3339))
	}

}

func TestBadSpec(t *testing.T) {
	var s BadSpecHasStruct
	err := Load("ev", &s)
	if err == nil {
		t.Error("expected bad specification error")
	}

}

func TestBadSpecEmbedded(t *testing.T) {
	var s BadSpecHasEmbedded
	err := Load("ev", &s)
	if err == nil {
		t.Error("expected bad specification error")
	}

}

func TestParseErrorBool(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_DEBUG", "string")
	os.Setenv("EV_REQUIREDVAR", "foo")
	err := Load("ev", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Debug" {
		t.Errorf("expected %s, got %v", "Debug", v.FieldName)
	}
	if s.Debug != false {
		t.Errorf("expected %v, got %v", false, s.Debug)
	}
}

func TestParseErrorFloat32(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_RATE", "string")
	os.Setenv("EV_REQUIREDVAR", "foo")
	err := Load("ev", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Rate" {
		t.Errorf("expected %s, got %v", "Rate", v.FieldName)
	}
	if s.Rate != 0 {
		t.Errorf("expected %v, got %v", 0, s.Rate)
	}
}

func TestParseErrorInt(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_PORT", "string")
	os.Setenv("EV_REQUIREDVAR", "foo")
	err := Load("ev", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Port" {
		t.Errorf("expected %s, got %v", "Port", v.FieldName)
	}
	if s.Port != 0 {
		t.Errorf("expected %v, got %v", 0, s.Port)
	}
}

func TestParseErrorUint(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_TTL", "-30")
	err := Load("ev", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "TTL" {
		t.Errorf("expected %s, got %v", "TTL", v.FieldName)
	}
	if s.TTL != 0 {
		t.Errorf("expected %v, got %v", 0, s.TTL)
	}
}

func TestErrBadSpec(t *testing.T) {
	m := make(map[string]string)
	err := Load("ev", &m)
	if err != ErrBadSpec {
		t.Errorf("expected %v, got %v", ErrBadSpec, err)
	}
}

func TestUnsetVars(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("USER", "foo")
	os.Setenv("EV_REQUIREDVAR", "foo")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	// If the var is not defined the non-prefixed version should not be used
	// unless the struct tag says so
	if s.User != "" {
		t.Errorf("expected %q, got %q", "", s.User)
	}
}

func TestAlternateVarNames(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_MULTI_WORD_VAR", "foo")
	os.Setenv("EV_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("EV_MULTI_WORD_VAR_WITH_LOWER_CASE_ALT", "baz")
	os.Setenv("EV_REQUIREDVAR", "foo")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	// Setting the alt version of the var in the environment has no effect if
	// the struct tag is not supplied
	if s.MultiWordVar != "" {
		t.Errorf("expected %q, got %q", "", s.MultiWordVar)
	}

	// Setting the alt version of the var in the environment correctly sets
	// the value if the struct tag IS supplied
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %q, got %q", "bar", s.MultiWordVarWithAlt)
	}

	// Alt value is not case sensitive and is treated as all uppercase
	if s.MultiWordVarWithLowerCaseAlt != "baz" {
		t.Errorf("expected %q, got %q", "baz", s.MultiWordVarWithLowerCaseAlt)
	}
}

func TestRequiredVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "foobar")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.RequiredVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.RequiredVar)
	}
}

func TestRequiredMissing(t *testing.T) {
	var s Specification
	os.Clearenv()

	err := Load("ev", &s)
	if err == nil {
		t.Error("no failure when missing required variable")
	}
}

func TestBlankDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "requiredvalue")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.DefaultVar)
	}

	if *s.SomePointerWithDefault != "foo2baz" {
		t.Errorf("expected %s, got %s", "foo2baz", *s.SomePointerWithDefault)
	}
}

func TestNonBlankDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_DEFAULTVAR", "nondefaultval")
	os.Setenv("EV_REQUIREDVAR", "requiredvalue")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "nondefaultval" {
		t.Errorf("expected %s, got %s", "nondefaultval", s.DefaultVar)
	}
}

func TestExplicitBlankDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_DEFAULTVAR", "")
	os.Setenv("EV_REQUIREDVAR", "")

	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "" {
		t.Errorf("expected %s, got %s", "\"\"", s.DefaultVar)
	}
}

func TestAlternateNameDefaultVar(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("BROKER", "betterbroker")
	os.Setenv("EV_REQUIREDVAR", "foo")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.NoPrefixDefault != "betterbroker" {
		t.Errorf("expected %q, got %q", "betterbroker", s.NoPrefixDefault)
	}

	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "foo")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.NoPrefixDefault != "127.0.0.1" {
		t.Errorf("expected %q, got %q", "127.0.0.1", s.NoPrefixDefault)
	}
}

func TestRequiredDefault(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "foo")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.RequiredDefault != "foo2bar" {
		t.Errorf("expected %q, got %q", "foo2bar", s.RequiredDefault)
	}
}

func TestPointerFieldBlank(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "foo")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.SomePointer != nil {
		t.Errorf("expected <nil>, got %q", *s.SomePointer)
	}
}

func TestEmptyMapFieldOverride(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "foo")
	os.Setenv("EV_MAPFIELD", "")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.MapField == nil {
		t.Error("expected empty map, got <nil>")
	}

	if len(s.MapField) != 0 {
		t.Errorf("expected empty map, got map of size %d", len(s.MapField))
	}
}

func TestEmbeddedStruct(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "required")
	os.Setenv("EV_ENABLED", "true")
	os.Setenv("EV_EMBEDDEDPORT", "1234")
	os.Setenv("EV_MULTIWORDVAR", "foo")
	os.Setenv("EV_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("EV_MULTI_WITH_DIFFERENT_ALT", "baz")
	os.Setenv("EV_EMBEDDED_WITH_ALT", "foobar")
	os.Setenv("EV_SOMEPOINTER", "foobaz")
	os.Setenv("EV_EMBEDDED_IGNORED", "was-not-ignored")
	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if s.MultiWordVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.MultiWordVar)
	}
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %s, got %s", "bar", s.MultiWordVarWithAlt)
	}
	if *s.SomePointer != "foobaz" {
		t.Errorf("expected %s, got %s", "foobaz", *s.SomePointer)
	}
}

func TestNonPointerFailsProperly(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "snap")

	err := Load("ev", s)
	if err != ErrBadSpec {
		t.Errorf("non-pointer should fail with ErrBadSpec, was instead %s", err)
	}
}

func TestCustomValueFields(t *testing.T) {
	var s struct {
		Foo    string
		Bar    bracketed
		Struct setterStruct
	}

	os.Clearenv()
	os.Setenv("EV_FOO", "foo")
	os.Setenv("EV_BAR", "bar")
	os.Setenv("EV_STRUCT", "inner")

	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if want := "foo"; s.Foo != want {
		t.Errorf("foo: got %#q, want %#q", s.Foo, want)
	}

	if want := "[bar]"; s.Bar.String() != want {
		t.Errorf("bar: got %#q, want %#q", s.Bar, want)
	}

	if want := `setterstruct{"inner"}`; s.Struct.Inner != want {
		t.Errorf(`Struct.Inner: got %#q, want %#q`, s.Struct.Inner, want)
	}
}

func TestCustomPointerFields(t *testing.T) {
	var s struct {
		Foo    string
		Bar    *bracketed
		Struct *setterStruct
	}

	// Set would panic when the receiver is nil,
	// so make sure they have initial values to replace.
	s.Bar = new(bracketed)

	os.Clearenv()
	os.Setenv("EV_FOO", "foo")
	os.Setenv("EV_BAR", "bar")
	os.Setenv("EV_STRUCT", "inner")

	if err := Load("ev", &s); err != nil {
		t.Error(err.Error())
	}

	if want := "foo"; s.Foo != want {
		t.Errorf("foo: got %#q, want %#q", s.Foo, want)
	}

	if want := "[bar]"; s.Bar.String() != want {
		t.Errorf("bar: got %#q, want %#q", s.Bar, want)
	}

	if want := `setterstruct{"inner"}`; s.Struct.Inner != want {
		t.Errorf(`Struct.Inner: got %#q, want %#q`, s.Struct.Inner, want)
	}
}

func TestEmptyPrefixUsesFieldNames(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("REQUIREDVAR", "foo")

	err := Load("", &s)
	if err == nil {
		t.Errorf("expected Load() to fail with: a 'prefix must be provided' message.")
	}
}

func TestTextUnmarshalerError(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("EV_REQUIREDVAR", "foo")
	os.Setenv("EV_DATETIME", "I'M NOT A DATE")

	err := Load("ev", &s)

	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Datetime" {
		t.Errorf("expected %s, got %v", "Debug", v.FieldName)
	}

	expectedLowLevelError := time.ParseError{
		Layout:     time.RFC3339,
		Value:      "I'M NOT A DATE",
		LayoutElem: "2006",
		ValueElem:  "I'M NOT A DATE",
	}

	if v.Err.Error() != expectedLowLevelError.Error() {
		t.Errorf("expected %s, got %s", expectedLowLevelError, v.Err)
	}
	if s.Debug != false {
		t.Errorf("expected %v, got %v", false, s.Debug)
	}
}

type bracketed string

func (b *bracketed) Set(value string) error {
	*b = bracketed("[" + value + "]")
	return nil
}

func (b bracketed) String() string {
	return string(b)
}

type setterStruct struct {
	Inner string
}

func (ss *setterStruct) Set(value string) error {
	ss.Inner = fmt.Sprintf("setterstruct{%q}", value)
	return nil
}
