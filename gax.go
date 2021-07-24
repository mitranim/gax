/*
Simple system for writing HTML as Go code. Better-performing replacement for
`html/template` and `text/template`; see benchmarks in readme.

Vaguely inspired by JS library https://github.com/mitranim/prax,
but uses a different design.

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
is `Bui` and `Bui.E`. See the `Bui` example below.
*/
package gax

/*
Shortcut for prepending HTML doctype. Use `gax.Bui(gax.Doctype)` to create a
document-level HTML builder.
*/
const Doctype = `<!doctype html>`

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
Indicates pre-escaped markup. When using `Bui.E`, values of type `Bytes` are
written as-is without additional HTML/XML escaping. For strings, see `String`.
*/
type Bytes []byte

/*
Indicates pre-escaped markup. When using `Bui.E`, Values of type `String` are
written as-is without additional HTML/XML escaping. For bytes, see `Bytes`.
*/
type String string

// Signature of the method `Bui.E`. Makes it convenient to pass around.
type E = func(string, A, ...interface{})
