package gax

import "unicode/utf8"

/*
Mostly for internal use. Similar to `bytes.Buffer` or `strings.Builder`, but
simpler and more flexible, being just a byte slice.
*/
type NonEscWri []byte

// Similar to `strings.Builder.Write`.
func (self *NonEscWri) Write(val []byte) (int, error) {
	*self = append(*self, val...)
	return len(val), nil
}

// Similar to `strings.Builder.WriteString`.
func (self *NonEscWri) WriteString(val string) (int, error) {
	*self = append(*self, val...)
	return len(val), nil
}

// Similar to `strings.Builder.WriteRune`.
func (self *NonEscWri) WriteRune(val rune) (int, error) {
	if uint32(val) < utf8.RuneSelf {
		*self = append(*self, byte(val))
		return 1, nil
	}

	lenPrev := len(*self)
	if cap(*self)-lenPrev < utf8.UTFMax {
		self.grow(utf8.UTFMax)
	}

	wid := utf8.EncodeRune((*self)[lenPrev:lenPrev+utf8.UTFMax], val)
	*self = (*self)[:lenPrev+wid]
	return wid, nil
}

// Similar to `strings.Builder.String`. Free cast with no allocation.
func (self NonEscWri) String() string { return bytesToMutableString(self) }

func (self *NonEscWri) grow(size int) {
	buf := make([]byte, len(*self), 2*cap(*self)+size)
	copy(buf, *self)
	*self = buf
}

/*
Mostly for internal use. Writes text as if it were inside an HTML/XML attribute,
without enclosing quotes, escaping as necessary. For escaping rules, see:

	https://www.w3.org/TR/html52/syntax.html#escaping-a-string
*/
type AttrWri []byte

/*
Similar to `strings.Builder.Write`, but escapes special chars. Technically not
compliant with `io.Writer`: the returned count of written bytes may exceed the
size of the provided chunk.
*/
func (self *AttrWri) Write(val []byte) (int, error) {
	return self.WriteString(bytesToMutableString(val))
}

// Similar to `strings.Builder.WriteString`, but escapes special chars.
func (self *AttrWri) WriteString(val string) (size int, _ error) {
	for _, char := range val {
		delta, _ := self.WriteRune(char)
		size += delta
	}
	return
}

// Similar to `strings.Builder.WriteRune`, but escapes special chars.
func (self *AttrWri) WriteRune(val rune) (int, error) {
	wri := (*NonEscWri)(self)

	switch val {
	case '&':
		return wri.WriteString(`&amp;`)
	case '\u00a0':
		return wri.WriteString(`&nbsp;`)
	case '"':
		return wri.WriteString(`&quot;`)
	default:
		return wri.WriteRune(val)
	}
}

// Similar to `strings.Builder.String`. Free cast with no allocation.
func (self AttrWri) String() string { return bytesToMutableString(self) }

/*
Mostly for internal use. Writes text as if it were inside an HTML/XML element,
escaping as necessary. For escaping rules, see:

	https://www.w3.org/TR/html52/syntax.html#escaping-a-string
*/
type TextWri []byte

/*
Similar to `strings.Builder.Write`, but escapes special chars. Technically not
compliant with `io.Writer`: the returned count of written bytes may exceed the
size of the provided chunk.
*/
func (self *TextWri) Write(val []byte) (int, error) {
	return self.WriteString(bytesToMutableString(val))
}

// Similar to `strings.Builder.WriteString`, but escapes special chars.
func (self *TextWri) WriteString(val string) (size int, _ error) {
	for _, char := range val {
		delta, _ := self.WriteRune(char)
		size += delta
	}
	return
}

// Similar to `strings.Builder.WriteRune`, but escapes special chars.
func (self *TextWri) WriteRune(val rune) (int, error) {
	wri := (*NonEscWri)(self)

	switch val {
	case '&':
		return wri.WriteString(`&amp;`)
	case '\u00a0':
		return wri.WriteString(`&nbsp;`)
	case '<':
		return wri.WriteString(`&lt;`)
	case '>':
		return wri.WriteString(`&gt;`)
	default:
		return wri.WriteRune(val)
	}
}

// Similar to `strings.Builder.String`. Free cast with no allocation.
func (self TextWri) String() string { return bytesToMutableString(self) }