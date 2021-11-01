package gax

import (
	"fmt"
	"reflect"
)

/*
Short for "fragment" or "document fragment". Shortcut for making `Bui` with
these children.
*/
func F(vals ...interface{}) (bui Bui) {
	bui.F(vals...)
	return
}

/*
Short for "builder" or "builder for UI". Has methods for generating HTML/XML
markup, declarative but efficient. See `E`, `F`, and `Bui.E` for 99% of the API
you will use.

When used as a child (see `Bui.E`, `Bui.F`, `Bui.Child`), this also indicates
pre-escaped markup, appending itself to another `Bui` without HTML/XML
escaping. For strings, see `Str`.
*/
type Bui []byte

/*
One of the primary APIs. Counterpart to the function `E`. Short for "element"
or "HTML element". Writes an HTML/XML tag, with attributes and inner content.

For a runnable example, see the definition of `Bui`.

Special rules for children:

	* `nil` is ignored.
	* `[]interface{}` is recursively traversed.
	* `func()`, `func(*Bui)`, or `Ren.Render` is called for side effects.
	* Other values are stringified and escaped via `TextWri`.

To write text without escaping, use `Str` for strings and `Bui` for byte
slices.
*/
func (self *Bui) E(tag string, attrs Attrs, children ...interface{}) {
	self.Begin(tag, attrs)
	self.F(children...)
	self.End(tag)
}

/*
Mostly for internal use. Writes the beginning of an HTML/XML element, with
optional attrs. Supports HTML special cases; see `Bui.Attrs`. Sanity-checks the
tag. Using an invalid tag causes a panic.
*/
func (self *Bui) Begin(tag string, attrs Attrs) {
	validTag(tag)

	self.NonEscString(`<`)
	self.NonEscString(tag)
	self.Attrs(attrs...)
	self.NonEscString(`>`)
}

/*
Mostly for internal use. Writes the end of an HTML/XML element. Supports HTML
void elements, also known as self-closing tags: if `Void.Has(tag)`, this method
is a nop. Sanity-checks the tag. Using an invalid tag causes a panic.
*/
func (self *Bui) End(tag string) {
	validTag(tag)

	if !Void.Has(tag) {
		self.NonEscString(`</`)
		self.NonEscString(tag)
		self.NonEscString(`>`)
	}
}

/*
Mostly for internal use. Writes HTML/XML attributes. Supports HTML special
cases; see `Bui.Attr`.
*/
func (self *Bui) Attrs(vals ...Attr) { *self = Bui(Attrs(vals).Append(*self)) }

/*
Mostly for internal use. Writes an HTML/XML attribute, preceded with a space.
Supports HTML bool attrs: if `Bool.Has(key)`, the attribute value may be
adjusted for spec compliance. Automatically escapes the attribute value.

Sanity-checks the attribute name. Using an invalid name causes a panic.
*/
func (self *Bui) Attr(val Attr) { *self = Bui(val.Append(*self)) }

// Writes multiple children via `Bui.Child`. Like the "tail part" of `Bui.E`.
// Counterpart to the function `F`.
func (self *Bui) F(vals ...interface{}) {
	for _, val := range vals {
		self.Child(val)
	}
}

/*
Mostly for internal use. Writes an arbitrary child. See `Bui.E` for the list of
special rules.
*/
func (self *Bui) Child(val interface{}) {
	switch val := val.(type) {
	case nil:
	case []interface{}:
		self.F(val...)
	case string:
		self.EscString(val)
	case []byte:
		self.EscBytes(val)
	case func():
		if val != nil {
			val()
		}
	case func(*Bui):
		if val != nil {
			val(self)
		}
	case func() interface{}:
		if val != nil {
			self.Child(val())
		}
	case Ren:
		if val != nil {
			val.Render(self)
		}
	default:
		self.Unknown(val)
	}
}

// Shorter alias for `Bui.Child`.
func (self *Bui) C(val interface{}) { self.Child(val) }

/*
Mostly for internal use. Writes regular text without escaping. For writing
`string`, see `Bui.NonEscString`. For escaping, see `Bui.EscBytes`.
*/
func (self *Bui) NonEscBytes(val []byte) {
	*self = append(*self, val...)
}

/*
Mostly for internal use. Writes regular text without escaping. For writing
`[]byte`, see `Bui.NonEscBytes`. For escaping, see `Bui.EscString`.
*/
func (self *Bui) NonEscString(val string) {
	*self = append(*self, val...)
}

/*
Writes regular text, escaping if necessary. For writing `string`, see
`Bui.EscBytes`.
*/
func (self *Bui) EscBytes(val []byte) {
	_, _ = (*TextWri)(self).Write(val)
}

/*
Writes regular text, escaping if necessary. For writing `[]byte`, see
`Bui.EscBytes`.
*/
func (self *Bui) EscString(val string) {
	_, _ = (*TextWri)(self).WriteString(val)
}

// Shorter alias for `Bui.EscString`.
func (self *Bui) T(val string) { self.EscString(val) }

/*
Mostly for internal use. If the provided value is not nil, it's printed via
`fmt.Fprint` and escaped via `TextWri`. Bypasses other special rules for child
encoding. Use `Bui.F` instead.
*/
func (self *Bui) Unknown(val interface{}) {
	if val == nil {
		return
	}

	rval := reflect.ValueOf(val)
	if isRvalNil(rval) {
		return
	}

	switch rval.Kind() {
	case reflect.Invalid, reflect.Func:
		panic(fmt.Errorf(`[gax] can't render %T`, val))
	}

	fmt.Fprint((*TextWri)(self), val)
}

// Implement `Ren`. Appends itself without HTML/XML escaping.
func (self Bui) Render(bui *Bui) { bui.NonEscBytes(self.Bytes()) }

// Free cast to `[]byte`.
func (self Bui) Bytes() []byte { return self }

// Free cast to `string`.
func (self Bui) String() string { return bytesToMutableString(self) }
