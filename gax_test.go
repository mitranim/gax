package gax

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_E_and_F(_ *testing.T) {
	bui := F(
		Str(Doctype),
		E(`html`, AP(`lang`, `en`),
			E(`head`, nil,
				E(`meta`, AP(`charset`, `utf-8`)),
				E(`meta`, AP(`http-equiv`, `X-UA-Compatible`, `content`, `IE=edge`)),
				E(`meta`, AP(`name`, `viewport`, `content`, `width=device-width, initial-scale=1`)),
				E(`link`, AP(`rel`, `icon`, `href`, `data:;base64,=`)),
				E(`title`, nil, `test markup`),
			),
			E(`body`, AP(`class`, `stretch-to-viewport`),
				E(`h1`, AP(`class`, `title`), `mock markup`),
				E(`div`, AP(`class`, `main`), `hello world!`),
			),
		),
	)

	eqs(
		strings.TrimSpace(`
			<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>
		`),
		bui,
	)
}

func Test_Bui_E(_ *testing.T) {
	bui := Bui(Doctype)
	E := bui.E

	E(`html`, AP(`lang`, `en`), func() {
		E(`head`, nil, func() {
			E(`meta`, AP(`charset`, `utf-8`))
			E(`meta`, AP(`http-equiv`, `X-UA-Compatible`, `content`, `IE=edge`))
			E(`meta`, AP(`name`, `viewport`, `content`, `width=device-width, initial-scale=1`))
			E(`link`, AP(`rel`, `icon`, `href`, `data:;base64,=`))
			E(`title`, nil, `test markup`)
		})
		E(`body`, AP(`class`, `stretch-to-viewport`), func() {
			E(`h1`, AP(`class`, `title`), `mock markup`)
			E(`div`, AP(`class`, `main`), `hello world!`)
		})
	})

	eqs(
		strings.TrimSpace(`
			<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>
		`),
		bui,
	)
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
	test(`str`, "str")
	test(`str`, []byte("str"))
	test(`false`, false)
	test(`true`, true)
	test(`0`, 0)
	test(`1`, 1)
	test(`12`, int32(12))
	test(`12`, int64(12))
	test(`-12`, int32(-12))
	test(`-12`, int64(-12))
	test(`12.34`, float32(12.34))
	test(`12.34`, float64(12.34))
	test(`-12.34`, float32(-12.34))
	test(`-12.34`, float64(-12.34))
	test(`[10 20 30]`, []int{10, 20, 30})
	test(`10str20`, []interface{}{10, nil, "str", []interface{}{nil, 20}})
	test(`&lt;one&gt;&amp;"&lt;/one&gt;`, `<one>&"</one>`)
	test(`&lt;one&gt;&amp;"&lt;/one&gt;`, []byte(`<one>&"</one>`))

	test(`<div></div>`, func(b *Bui) { b.E(`div`, nil) })
	test(`str`, func(b *Bui) { b.T(`str`) })

	t.Run("do_not_escape_special_type", func(_ *testing.T) {
		test(`<one>&"</one>`, Str(`<one>&"</one>`))
		test(`<one>&"</one>`, Bui(`<one>&"</one>`))
		test(`<a>one</a><bui>two</bui><c>three</c>`, Str(`<a>one</a><bui>two</bui><c>three</c>`))
		test(`<a>one</a><bui>two</bui><c>three</c>`, Bui(`<a>one</a><bui>two</bui><c>three</c>`))
	})
}

func Test_F(_ *testing.T) {
	fun := func(b *Bui) { b.E(`html`, AP(`lang`, `en`)) }
	eqs(
		`<!doctype html><html lang="en"></html>`,
		F(Str(Doctype), fun),
	)
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
	tryInt(wri.Write([]byte("A&B\u00a0C\"D<E>F")))
	eqs(`A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_AttrWri_WriteString(_ *testing.T) {
	var wri AttrWri
	tryInt(wri.WriteString("A&B\u00a0C\"D<E>F"))
	eqs(`A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_AttrWri_WriteRune(_ *testing.T) {
	var wri AttrWri

	tryInt(wri.WriteRune('A'))
	eqs(`A`, wri)

	tryInt(wri.WriteRune('&'))
	eqs(`A&amp;`, wri)

	tryInt(wri.WriteRune('B'))
	eqs(`A&amp;B`, wri)

	tryInt(wri.WriteRune('\u00a0'))
	eqs(`A&amp;B&nbsp;`, wri)

	tryInt(wri.WriteRune('C'))
	eqs(`A&amp;B&nbsp;C`, wri)

	tryInt(wri.WriteRune('"'))
	eqs(`A&amp;B&nbsp;C&quot;`, wri)

	tryInt(wri.WriteRune('D'))
	eqs(`A&amp;B&nbsp;C&quot;D`, wri)

	tryInt(wri.WriteRune('<'))
	eqs(`A&amp;B&nbsp;C&quot;D<`, wri)

	tryInt(wri.WriteRune('E'))
	eqs(`A&amp;B&nbsp;C&quot;D<E`, wri)

	tryInt(wri.WriteRune('>'))
	eqs(`A&amp;B&nbsp;C&quot;D<E>`, wri)

	tryInt(wri.WriteRune('F'))
	eqs(`A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_TextWri_Write(_ *testing.T) {
	var wri TextWri
	tryInt(wri.Write([]byte("A&B\u00a0C\"D<E>F")))
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func Test_TextWri_WriteString(_ *testing.T) {
	var wri TextWri
	tryInt(wri.WriteString("A&B\u00a0C\"D<E>F"))
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func Test_TextWri_WriteRune(_ *testing.T) {
	var wri TextWri

	tryInt(wri.WriteRune('A'))
	eqs(`A`, wri)

	tryInt(wri.WriteRune('&'))
	eqs(`A&amp;`, wri)

	tryInt(wri.WriteRune('B'))
	eqs(`A&amp;B`, wri)

	tryInt(wri.WriteRune('\u00a0'))
	eqs(`A&amp;B&nbsp;`, wri)

	tryInt(wri.WriteRune('C'))
	eqs(`A&amp;B&nbsp;C`, wri)

	tryInt(wri.WriteRune('"'))
	eqs(`A&amp;B&nbsp;C"`, wri)

	tryInt(wri.WriteRune('D'))
	eqs(`A&amp;B&nbsp;C"D`, wri)

	tryInt(wri.WriteRune('<'))
	eqs(`A&amp;B&nbsp;C"D&lt;`, wri)

	tryInt(wri.WriteRune('E'))
	eqs(`A&amp;B&nbsp;C"D&lt;E`, wri)

	tryInt(wri.WriteRune('>'))
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;`, wri)

	tryInt(wri.WriteRune('F'))
	eqs(`A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func TestElem_GoString(t *testing.T) {
	eq(
		"E(`div`, AP(`class`, `one`), `two`, 10, `three`)",
		fmt.Sprintf(
			`%#v`,
			E(`div`, AP(`class`, `one`), `two`, []interface{}{10, `three`}),
		),
	)
}

func TestA_GoString(t *testing.T) {
	eq(
		"AP(`one`, `two`, `three`, `four`)",
		fmt.Sprintf(
			`%#v`,
			AP(`one`, `two`, `three`, `four`).A(Attr{}).A(Attr{}).A(Attr{}),
		),
	)
}

func TestVac(t *testing.T) {
	eq(nil, Vac(nil))
	eq(nil, Vac((*string)(nil)))
	eq(nil, Vac([]interface{}{}))
	eq(nil, Vac([]interface{}{nil}))
	eq(nil, Vac([]interface{}{nil, (*string)(nil)}))
	eq(nil, Vac([]byte(nil)))
	eq(nil, Vac(Bui(nil)))

	eq("", Vac(""))
	eq(0, Vac(0))
	eq([]interface{}{""}, Vac([]interface{}{""}))
	eq([]interface{}{0}, Vac([]interface{}{0}))
}

func eqs(exp string, act []byte) {
	eq(exp, string(act))
}

func eq(exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		panic(fmt.Errorf("expected:\n%+v\ngot:\n%+v\n", exp, act))
	}
}

func tryInt(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
