package gax

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// Match the `fmt` default for consistency.
const floatVerb = 'g'

func newStringSet(vals ...string) stringSet {
	set := make(stringSet, len(vals))
	for _, val := range vals {
		set.Add(val)
	}
	return set
}

type stringSet map[string]struct{}

func (self stringSet) Has(val string) bool {
	_, ok := self[val]
	return ok
}

func (self stringSet) Add(val string) { self[val] = struct{}{} }

func (self stringSet) Del(val string) { delete(self, val) }

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
	return val == nil || isRvalNil(reflect.ValueOf(val))
}

func isRvalNil(rval reflect.Value) bool {
	return !rval.IsValid() || isRkindNilable(rval.Kind()) && rval.IsNil()
}

func isRkindNilable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	default:
		return false
	}
}

/*
Allocation-free conversion. Reinterprets a byte slice as a string. Borrowed from
the standard library. Reasonably safe.
*/
func bytesToMutableString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
