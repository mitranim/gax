package gax

/*
Short for "attributes". List of arbitrary HTML/XML attributes. Usually passed to
`Bui.E`.
*/
type A []Attr

// Mostly for internal use.
func (self A) WriteTo(buf *[]byte) {
	for _, val := range self {
		val.WriteTo(buf)
	}
}

// Implement `fmt.Stringer` for debug purposes. Not used by builder methods.
func (self A) String() string {
	var buf []byte
	self.WriteTo(&buf)
	return NonEscWri(buf).String()
}

/*
Represents an arbitrary HTML/XML attribute. Usually part of `A{}`. An empty/zero
attr (equal to `Attr{}`) is ignored during encoding.
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
func (self Attr) WriteTo(buf *[]byte) {
	if self == (Attr{}) {
		return
	}

	key, val := self.Name(), self.Value()
	validAttr(key)

	if Bool.Has(key) {
		// Dumb hack. Should revise.
		if val == "false" {
			return
		}
		val = ""
	}

	*buf = append(*buf, ` `...)
	*buf = append(*buf, key...)
	*buf = append(*buf, `="`...)
	_, _ = (*AttrWri)(buf).WriteString(val)
	*buf = append(*buf, `"`...)
}

// Implement `fmt.Stringer` for debug purposes. Not used by builder methods.
func (self Attr) String() string {
	var buf []byte
	self.WriteTo(&buf)
	return NonEscWri(buf).String()
}
