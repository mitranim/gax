package gax

import (
	"fmt"
)

/*
Primary API. Short for "element" or "HTML element". Expresses an HTML/XML tag,
with attributes and inner content. Creates an instance of `Elem`, which
implements the `Ren` interface. It can render itself as HTML/XML, or be passed
as a child to `F`, `E`, `Bui.E`.

For special rules regarding child encoding, see `Bui.E`.
*/
func E(tag string, attrs Attrs, child ...any) Elem {
	return Elem{tag, attrs, child}
}

/*
Represents an HTML element. Usually created via `E`. Can render itself, or be
passed as a child to `F` or `Bui.E`.
*/
type Elem struct {
	Tag   string
	Attrs Attrs
	Child any
}

var _ = Ren(Elem{})

/*
Implement `Ren`. This allows `Elem` to be passed as a child to the various
rendering functions like `E`, `F`, `Bui.E`. As a special case, `Elem` with
an empty `.Tag` does not render anything.
*/
func (self Elem) Render(b *Bui) {
	if self.Tag != `` {
		b.E(self.Tag, self.Attrs, self.Child)
	}
}

// Implement `fmt.Stringer` for debug purposes. Not used by builder methods.
func (self Elem) String() string { return F(self).String() }

/*
Implement `fmt.GoStringer` for debug purposes. Not used by builder methods.
Represents itself as a call to `E`, which is the recommended way to write
this.
*/
func (self Elem) GoString() string {
	var buf NonEscWri
	_, _ = buf.WriteString(`E(`)
	buf = NonEscWri(appendQuote(buf, self.Tag))
	_, _ = buf.WriteString(`, `)
	buf = append(buf, self.Attrs.GoString()...)
	buf = appendElemChild(buf, self.Child)
	_, _ = buf.WriteString(`)`)
	return buf.String()
}

func appendElemChild(buf NonEscWri, val any) NonEscWri {
	switch val := val.(type) {
	case nil:
	case []any:
		for _, val := range val {
			buf = appendElemChild(buf, val)
		}
	case string:
		_, _ = buf.WriteString(`, `)
		buf = NonEscWri(appendQuote(buf, val))
	default:
		_, _ = buf.WriteString(`, `)
		fmt.Fprintf(&buf, `%#v`, val)
	}
	return buf
}
