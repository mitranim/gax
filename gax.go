/*
Simple system for writing HTML as Go code. Better-performing replacement for
`html/template` and `text/template`; see benchmarks in readme.

Vaguely inspired by JS library https://github.com/mitranim/prax.

Features / benefits:

	* No weird special language to learn.
	* Use actual Go code.
	* Use normal Go conditionals.
	* Use normal Go loops.
	* Use normal Go functions.
	* Benefit from static typing.
	* Benefit from Go code analysis.
	* Benefit from Go performance.
	* Tiny and dependency-free.

The API is bloated with "just in case" public exports, but 99% of what you want
is `E`, `F`, `Bui`, and `Bui.E`. See the `Bui` example below.
*/
package gax

/*
Shortcut for prepending HTML doctype. Use `Bui(Doctype)` to create a
document-level HTML builder, or `Str(Doctype)` to prepend this in `F`.
*/
const Doctype = `<!doctype html>`

/*
Short for "renderer". On children implementing this interface, the `Render`
method is called for side effects, instead of stringifying the child.
*/
type Ren interface{ Render(*Bui) }

/*
Indicates pre-escaped markup. Children of this type are written as-is without
additional HTML/XML escaping. For bytes, use `Bui`.
*/
type Str string

// Implement `Ren`. Appends itself without HTML/XML escaping.
func (self Str) Render(bui *Bui) { bui.NonEscString(string(self)) }

/*
Set of known HTML boolean attributes. Can be modified via `Bool.Add` and
`Bool.Del`. The specification postulates the concept, but where's the standard
list? Taken from non-authoritative sources. Reference:

	https://www.w3.org/TR/html52/infrastructure.html#boolean-attribute
*/
var Bool = newStringSet(
	"allowfullscreen", "allowpaymentrequest", "async", "autofocus", "autoplay",
	"checked", "controls", "default", "disabled", "formnovalidate", "hidden",
	"ismap", "itemscope", "loop", "multiple", "muted", "nomodule", "novalidate",
	"open", "playsinline", "readonly", "required", "reversed", "selected",
	"truespeed",
)

/*
Set of known HTML void elements, also known as self-closing tags. Can be
modified via `Void.Add` and `Void.Del`. Reference:

	https://www.w3.org/TR/html52/
	https://www.w3.org/TR/html52/syntax.html#writing-html-documents-elements
*/
var Void = newStringSet(
	"area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta",
	"param", "source", "track", "wbr",
)

/*
Short for "vacate", "vacuum", "vacuous". Takes a "child" intended for `E` or
`F`. If the child is empty, returns `nil`, otherwise returns the child as-is.
Empty is defined as containing only nils. Just like `E` and `F`, this
recursively traverses `[]interface{}`.
*/
func Vac(val interface{}) interface{} {
	inout := val

	switch val := val.(type) {
	case []interface{}:
		for _, val := range val {
			if Vac(val) != nil {
				return inout
			}
		}
		return nil

	default:
		if !isNil(val) {
			return inout
		}
		return nil
	}
}
