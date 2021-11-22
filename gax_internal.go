package gax

import (
	"fmt"
	r "reflect"
	"strconv"
	"strings"
	"unsafe"
)

var typeRen = r.TypeOf((*Ren)(nil)).Elem()

func newStringSet(vals ...string) stringSet {
	set := make(stringSet, len(vals))
	for _, val := range vals {
		set.Add(val)
	}
	return set
}

type stringSet map[string]struct{}

func (self stringSet) Has(val string) bool { _, ok := self[val]; return ok }
func (self stringSet) Add(val string)      { self[val] = struct{}{} }
func (self stringSet) Del(val string)      { delete(self, val) }

/*
Extremely permissive. Should prevent weird gotchas without interfering with
non-ASCII XML.

Reference for HTML:

	https://www.w3.org/TR/html52/syntax.html#tag-name
	https://www.w3.org/TR/html52/infrastructure.html#alphanumeric-ascii-characters

Also see for attrs, unused:

	https://www.w3.org/TR/html52/syntax.html#elements-attributes
*/
func validTag(val string) {
	if invalidTagOrAttr(val) {
		panic(fmt.Errorf(`[gax] invalid tag name %q`, val))
	}
}

/*
Extremely permissive. Intended only to prevent weird gotchas without interfering
with non-ASCII XML.
*/
func validAttr(val string) {
	if invalidTagOrAttr(val) {
		panic(fmt.Errorf(`[gax] invalid attribute name %q`, val))
	}
}

func invalidTagOrAttr(val string) bool {
	return strings.ContainsAny(val, " \t\n\r\v<>\"=")
}

func isNil(val interface{}) bool {
	return val == nil || isRvalNil(r.ValueOf(val))
}

func isRvalNil(rval r.Value) bool {
	return !rval.IsValid() || isRkindNilable(rval.Kind()) && rval.IsNil()
}

func isRkindNilable(kind r.Kind) bool {
	switch kind {
	case r.Chan, r.Func, r.Interface, r.Map, r.Ptr, r.Slice:
		return true
	default:
		return false
	}
}

/*
Allocation-free conversion. Reinterprets a byte slice as a string. Borrowed from
the standard library. Reasonably safe.
*/
func bytesString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func appendQuote(buf []byte, val string) []byte {
	if strconv.CanBackquote(val) {
		buf = append(buf, "`"...)
		buf = append(buf, val...)
		buf = append(buf, "`"...)
		return buf
	}
	return strconv.AppendQuote(buf, val)
}

func grow(prev []byte, size int) []byte {
	len, cap := len(prev), cap(prev)
	if cap-len >= size {
		return prev
	}

	next := make([]byte, len, 2*cap+size)
	copy(next, prev)
	return next
}

func iter(count int) []struct{} { return make([]struct{}, count) }
