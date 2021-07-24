package gax

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_Bui_E(_ *testing.T) {
	bui := Bui(Doctype)
	E := bui.E

	E(`html`, A{{`lang`, `en`}}, func() {
		E(`head`, nil, func() {
			E(`meta`, A{{`charset`, `utf-8`}})
			E(`meta`, A{{`http-equiv`, `X-UA-Compatible`}, {`content`, `IE=edge`}})
			E(`meta`, A{{`name`, `viewport`}, {`content`, `width=device-width, initial-scale=1`}})
			E(`link`, A{{`rel`, `icon`}, {`href`, `data:;base64,=`}})
			E(`title`, nil, `test markup`)
		})
		E(`body`, A{{`class`, `stretch-to-viewport`}}, func() {
			E(`h1`, A{{`class`, `title`}}, `mock markup`)
			E(`div`, A{{`class`, `main`}}, `hello world!`)
		})
	})

	eqs(strings.TrimSpace(`
		<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>
	`), bui)
}

func Test_Bui_Attr(_ *testing.T) {
	var bui Bui
	bui.Attr(Attr{`class`, `<one>&"</one>`})
	eqs(` class="<one>&amp;&quot;</one>"`, bui)
}

func Test_Bui_Attrs(_ *testing.T) {
	var bui Bui

	bui.Attrs(
		Attr{`class`, `<one>&"</one>`},
		Attr{`style`, `<two>&"</two>`},
	)

	eqs(` class="<one>&amp;&quot;</one>" style="<two>&amp;&quot;</two>"`, bui)
}

func Test_Bui_EscString(_ *testing.T) {
	var bui Bui
	bui.EscString(`<one>&"</one>`)
	eqs(`&lt;one&gt;&amp;"&lt;/one&gt;`, bui)
}

func Test_Bui_EscBytes(_ *testing.T) {
	var bui Bui
	bui.EscBytes([]byte(`<one>&"</one>`))
	eqs(`&lt;one&gt;&amp;"&lt;/one&gt;`, bui)
}

func Test_Bui_Child(t *testing.T) {
	test := func(exp string, val interface{}) {
		var bui Bui
		bui.Child(val)
		eqs(exp, bui)
	}

	test(``, nil)
	test(``, (*func())(nil))
	test(`0`, 0)
	test(`1`, 1)
	test(`false`, false)
	test(`true`, true)
	test(`str`, "str")
	test(`[10 20 30]`, []int{10, 20, 30})
	test(`10str20`, []interface{}{10, nil, "str", []interface{}{nil, 20}})
	test(`&lt;one&gt;&amp;"&lt;/one&gt;`, `<one>&"</one>`)
	test(`&lt;one&gt;&amp;"&lt;/one&gt;`, []byte(`<one>&"</one>`))

	t.Run("do_not_escape_special_type", func(_ *testing.T) {
		test(`<one>&"</one>`, String(`<one>&"</one>`))
		test(`<one>&"</one>`, Bytes(`<one>&"</one>`))
		test(`<a>one</a><bui>two</bui><c>three</c>`, String(`<a>one</a><bui>two</bui><c>three</c>`))
		test(`<a>one</a><bui>two</bui><c>three</c>`, Bytes(`<a>one</a><bui>two</bui><c>three</c>`))
	})
}

// Incomplete test; should also verify zero-alloc.
func Test_Bui_Bytes(_ *testing.T) {
	eqs(`<div>hello world!</div>`, Bui(`<div>hello world!</div>`))
}

// Incomplete test; should also verify zero-alloc.
func Test_Bui_String(_ *testing.T) {
	eq(`<div>hello world!</div>`, Bui(`<div>hello world!</div>`).String())
	eq(`<div>hello world!</div>`, string(Bui(`<div>hello world!</div>`)))
}

func Test_AttrWri_Write(_ *testing.T) {
	var wri AttrWri
	wri.Write([]byte("A&B\u00a0C\"D<E>F"))
	eqs(`A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_AttrWri_WriteString(_ *testing.T) {
	var wri AttrWri
	wri.WriteString("A&B\u00a0C\"D<E>F")
	eqs(`A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_AttrWri_WriteRune(_ *testing.T) {
	var wri AttrWri

	wri.WriteRune('A')
	eqs(`A`, wri)

	wri.WriteRune('&')
	eqs(`A&amp;`, wri)

	wri.WriteRune('B')
	eqs(`A&amp;B`, wri)

	wri.WriteRune('\u00a0')
	eqs(`A&amp;B&nbsp;`, wri)

	wri.WriteRune('C')
	eqs(`A&amp;B&nbsp;C`, wri)

	wri.WriteRune('"')
	eqs(`A&amp;B&nbsp;C&quot;`, wri)

	wri.WriteRune('D')
	eqs(`A&amp;B&nbsp;C&quot;D`, wri)

	wri.WriteRune('<')
	eqs(`A&amp;B&nbsp;C&quot;D<`, wri)

	wri.WriteRune('E')
	eqs(`A&amp;B&nbsp;C&quot;D<E`, wri)

	wri.WriteRune('>')
	eqs(`A&amp;B&nbsp;C&quot;D<E>`, wri)

	wri.WriteRune('F')
	eqs(`A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_TextWri_Write(_ *testing.T) {
	var wri TextWri
	wri.Write([]byte("A&B\u00a0C\"D<E>F"))
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func Test_TextWri_WriteString(_ *testing.T) {
	var wri TextWri
	wri.WriteString("A&B\u00a0C\"D<E>F")
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func Test_TextWri_WriteRune(_ *testing.T) {
	var wri TextWri

	wri.WriteRune('A')
	eqs(`A`, wri)

	wri.WriteRune('&')
	eqs(`A&amp;`, wri)

	wri.WriteRune('B')
	eqs(`A&amp;B`, wri)

	wri.WriteRune('\u00a0')
	eqs(`A&amp;B&nbsp;`, wri)

	wri.WriteRune('C')
	eqs(`A&amp;B&nbsp;C`, wri)

	wri.WriteRune('"')
	eqs(`A&amp;B&nbsp;C"`, wri)

	wri.WriteRune('D')
	eqs(`A&amp;B&nbsp;C"D`, wri)

	wri.WriteRune('<')
	eqs(`A&amp;B&nbsp;C"D&lt;`, wri)

	wri.WriteRune('E')
	eqs(`A&amp;B&nbsp;C"D&lt;E`, wri)

	wri.WriteRune('>')
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;`, wri)

	wri.WriteRune('F')
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func eqs(exp string, act []byte) {
	eq(exp, string(act))
}

func eq(exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		panic(fmt.Errorf("expected:\n%#v\ngot:\n%#v\n", exp, act))
	}
}
