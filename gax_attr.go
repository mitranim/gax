package gax

import "fmt"

/*
Short for "attributes". Same as the `Attrs{}` constructor, but uses parentheses,
which is sometimes more convenient. Symmetric with `Attrs.A`.
*/
func A(vals ...Attr) Attrs { return Attrs(vals) }

/*
Short for "attributes from pairs". Recommended way to write attributes, due to
its brevity. Symmetric with `Attrs.AP`.
*/
func AP(pairs ...string) Attrs {
	return make(Attrs, 0, len(pairs)/2).AP(pairs...)
}

/*
Short for "attributes". List of arbitrary HTML/XML attributes. Used by `Elem`.
Usually passed to `E` or `Bui.E`.
*/
type Attrs []Attr

/*
Shortcut for appending more attributes. Useful when combining attributes from
hardcoded pairs (via `AP`) with attributes created as `Attr`. For example, you
can write a function that generates a specific attribute, and use this to
append the result.
*/
func (self Attrs) A(vals ...Attr) Attrs { return append(self, vals...) }

/*
Shortcut for appending more attributes from pairs, as if by calling `AP`.
Panics if the argument count is not even.
*/
func (self Attrs) AP(pairs ...string) Attrs {
	if len(pairs)%2 != 0 {
		panic(fmt.Errorf(`[gax] AP expects an even amount of args, got %#v`, pairs))
	}

	ind := 0
	for ind < len(pairs) {
		key := pairs[ind]
		ind++
		val := pairs[ind]
		ind++
		self = append(self, Attr{key, val})
	}
	return self
}

// Mostly for internal use.
func (self Attrs) AppendTo(buf []byte) []byte {
	for _, val := range self {
		buf = val.AppendTo(buf)
	}
	return buf
}

// Implement `fmt.Stringer` for debug purposes. Not used by builder methods.
func (self Attrs) String() string {
	return NonEscWri(self.AppendTo(nil)).String()
}

/*
Implement `fmt.GoStringer` for debug purposes. Not used by builder methods.
Represents itself as a call to `AP`, which is the recommended way to write
this.
*/
func (self Attrs) GoString() string {
	if self == nil {
		return `nil`
	}

	var buf NonEscWri
	_, _ = buf.WriteString(`AP(`)

	found := false
	for _, val := range self {
		if val == (Attr{}) {
			continue
		}

		if !found {
			found = true
		} else {
			_, _ = buf.WriteString(`, `)
		}

		buf = NonEscWri(appendQuote(buf, val.Name()))
		_, _ = buf.WriteString(`, `)
		buf = NonEscWri(appendQuote(buf, val.Value()))
	}

	_, _ = buf.WriteString(`)`)
	return buf.String()
}

/*
Represents an arbitrary HTML/XML attribute. Usually part of `Attrs{}`. An
empty/zero attr (equal to `Attr{}`) is ignored during encoding.
*/
type Attr [2]string

/*
Attribute name. If the attr is not equal to `Attr{}`, the name is validated
during encoding. Using an invalid name causes a panic.
*/
func (self Attr) Name() string { return self[0] }

/*
Attribute value. Automatically escaped and quoted when encoding the attr. For
known HTML boolean attrs, listed in `Bool`, the value may be tweaked for better
spec compliance, or the attr may be omitted entirely.
*/
func (self Attr) Value() string { return self[1] }

// Mostly for internal use.
func (self Attr) AppendTo(buf []byte) []byte {
	if self == (Attr{}) {
		return buf
	}

	key, val := self.Name(), self.Value()
	validAttr(key)

	if Bool.Has(key) {
		// Dumb hack. Should revise.
		if val == "false" {
			return buf
		}
		val = ""
	}

	buf = append(buf, ` `...)
	buf = append(buf, key...)
	buf = append(buf, `="`...)
	_, _ = (*AttrWri)(&buf).WriteString(val)
	buf = append(buf, `"`...)
	return buf
}

// Implement `fmt.Stringer` for debug purposes. Not used by builder methods.
func (self Attr) String() string {
	return NonEscWri(self.AppendTo(nil)).String()
}
