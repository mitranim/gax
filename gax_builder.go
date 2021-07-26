package gax

import "fmt"

/*
Short for "builder". Has methods for generating HTML/XML markup, declarative but
efficient. See `Bui.E` for 99% of the API you will use.
*/
type Bui []byte

/*
99% of the API of this library. Short for "element" or "HTML element". Writes an
HTML/XML tag, with attributes and inner content.

For a runnable example, see the definition of `Bui`.

Special rules for children:

	* `nil` is ignored.
	* `[]interface{}` is recursively traversed.
	* `string` or `[]bytes` is escaped via `TextWri`.
	* `String` or `Bytes` is written as-is, without escaping.
	* `func()` or `func(*Bui)` is called for side effects.
	* Other values are stringified.
*/
func (self *Bui) E(tag string, attrs A, children ...interface{}) {
	self.Begin(tag, attrs)
	self.Children(children...)
	self.End(tag)
}

/*
Mostly for internal use. Writes the beginning of an HTML/XML element, with
optional attrs. Supports HTML special cases; see `Bui.Attrs`. Sanity-checks the
tag. Using an invalid tag causes a panic.
*/
func (self *Bui) Begin(tag string, attrs A) {
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
func (self *Bui) Attrs(attrs ...Attr) {
	*self = Bui(A(attrs).Append([]byte(*self)))
}

/*
Mostly for internal use. Writes an HTML/XML attribute, preceded with a space.
Supports HTML bool attrs: if `Bool.Has(key)`, the attribute value may be
adjusted for spec compliance. Automatically escapes the attribute value.

Sanity-checks the attribute name. Using an invalid name causes a panic.
*/
func (self *Bui) Attr(attr Attr) {
	*self = Bui(attr.Append([]byte(*self)))
}

// Mostly for internal use. Writes children via `Bui.Child`.
func (self *Bui) Children(vals ...interface{}) {
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
		self.Children(val...)

	case String:
		self.NonEscString(string(val))

	case string:
		self.EscString(val)

	case Bytes:
		self.NonEscBytes(val)

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

	default:
		self.Unknown(val)
	}
}

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

/*
Mostly for internal use. If the provided value is not nil, it's printed via
`fmt.Fprint` and escaped via `TextWri`. Bypasses other special rules for child
encoding. Use `Bui.Child` instead.
*/
func (self *Bui) Unknown(val interface{}) {
	if isNil(val) {
		return
	}
	fmt.Fprint((*TextWri)(self), val)
}

/*
Shortcut for calling a function with `Bui.E` and returning the same `Bui`
instance.
*/
func (self Bui) With(fun func(E)) Bui {
	fun(self.E)
	return self
}

// Free cast to `[]byte`.
func (self Bui) Bytes() []byte { return self }

// Free cast to `string`.
func (self Bui) String() string { return bytesToMutableString(self) }

// Same as `Bui{}.With(fun)` but marginally shorter/cleaner.
func Ebui(fun func(E E)) (bui Bui) {
	fun(bui.E)
	return bui
}
