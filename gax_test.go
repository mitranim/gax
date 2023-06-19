package gax

import (
	"fmt"
	r "reflect"
	"strings"
	"testing"
)

func Test_E_and_F(t *testing.T) {
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

	eqs(t, bui, strings.TrimSpace(`
		<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>
	`))
}

func Test_Bui_E(t *testing.T) {
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

	eqs(t, bui, strings.TrimSpace(`
		<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>
	`))
}

func Test_Bui_Attr(t *testing.T) {
	var bui Bui
	bui.Attr(Attr{`class`, `<one>&"</one>`})
	eqs(t, bui, ` class="<one>&amp;&quot;</one>"`)
}

func Test_Bui_Attrs(t *testing.T) {
	var bui Bui

	bui.Attrs(
		Attr{`class`, `<one>&"</one>`},
		Attr{`style`, `<two>&"</two>`},
	)

	eqs(t, bui, ` class="<one>&amp;&quot;</one>" style="<two>&amp;&quot;</two>"`)
}

func Test_Bui_EscString(t *testing.T) {
	var bui Bui
	bui.EscString(`<one>&"</one>`)
	eqs(t, bui, `&lt;one&gt;&amp;"&lt;/one&gt;`)
}

func Test_Bui_EscBytes(t *testing.T) {
	var bui Bui
	bui.EscBytes([]byte(`<one>&"</one>`))
	eqs(t, bui, `&lt;one&gt;&amp;"&lt;/one&gt;`)
}

func Test_Bui_Child_escaping(t *testing.T) {
	test := childTest(t)
	test(`<one>&"</one>`, `&lt;one&gt;&amp;"&lt;/one&gt;`)
	test([]byte(`<one>&"</one>`), `&lt;one&gt;&amp;"&lt;/one&gt;`)
}

func Test_Bui_Child_Ren(t *testing.T) {
	test := childTest(t)

	test(Str(`str`), `str`)
	test(Bui(`str`), `str`)

	test(
		Str(`<one>&"</one>`),
		`<one>&"</one>`,
	)

	test(
		Bui(`<one>&"</one>`),
		`<one>&"</one>`,
	)

	test(
		Str(`<a>one</a><bui>two</bui><c>three</c>`),
		`<a>one</a><bui>two</bui><c>three</c>`,
	)

	test(
		Bui(`<a>one</a><bui>two</bui><c>three</c>`),
		`<a>one</a><bui>two</bui><c>three</c>`,
	)

	test(
		E(`one`, AP(`two`, `three`), `four`),
		`<one two="three">four</one>`,
	)
}

func Test_Bui_Child_slices(t *testing.T) {
	test := childTest(t)

	test([]any{10, nil, "str", []any{nil, 20}}, `10str20`)

	test(
		[]Ren{
			E(`one`, nil),
			nil,
			Str(`two`),
			nil,
			E(`three`, nil),
		},
		`<one></one>two<three></three>`,
	)

	test(
		[]Elem{
			E(`one`, nil),
			Elem{},
			E(`two`, nil),
		},
		`<one></one><two></two>`,
	)
}

func Test_Bui_Child_funcs(t *testing.T) {
	test := childTest(t)

	test((*func())(nil), ``)
	test((*func(*Bui))(nil), ``)
	test(func(b *Bui) { b.E(`div`, nil) }, `<div></div>`)
	test(func(b *Bui) { b.T(`str`) }, `str`)
}

func Test_Bui_Child_misc(t *testing.T) {
	test := childTest(t)

	test(nil, ``)
	test("str", `str`)
	test([]byte("str"), `str`)
	test(false, `false`)
	test(true, `true`)
	test(0, `0`)
	test(1, `1`)
	test(int32(12), `12`)
	test(int64(12), `12`)
	test(int32(-12), `-12`)
	test(int64(-12), `-12`)
	test(float32(12.34), `12.34`)
	test(float64(12.34), `12.34`)
	test(float32(-12.34), `-12.34`)
	test(float64(-12.34), `-12.34`)
	test([]int{10, 20, 30}, `[10 20 30]`)
}

func childTest(t testing.TB) func(any, string) {
	return func(val any, exp string) {
		var bui Bui
		bui.Child(val)
		eqs(t, bui, exp)
	}
}

func Test_F(t *testing.T) {
	eqs(
		t,
		F(
			Str(Doctype),
			func(b *Bui) { b.E(`html`, AP(`lang`, `en`)) },
		),
		`<!doctype html><html lang="en"></html>`,
	)
}

// Incomplete test; should also verify zero-alloc.
func Test_Bui_Bytes(t *testing.T) {
	eqs(t, Bui(`<div>hello world!</div>`), `<div>hello world!</div>`)
}

// Incomplete test; should also verify zero-alloc.
func Test_Bui_String(t *testing.T) {
	eq(t, Bui(`<div>hello world!</div>`).String(), `<div>hello world!</div>`)
	eq(t, string(Bui(`<div>hello world!</div>`)), `<div>hello world!</div>`)
}

func Test_AttrWri_Write(t *testing.T) {
	var wri AttrWri
	tryInt(wri.Write([]byte("A&B\u00a0C\"D<E>F")))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;D<E>F`)
}

func Test_AttrWri_WriteString(t *testing.T) {
	var wri AttrWri
	tryInt(wri.WriteString("A&B\u00a0C\"D<E>F"))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;D<E>F`)
}

func Test_AttrWri_WriteRune(t *testing.T) {
	var wri AttrWri

	tryInt(wri.WriteRune('A'))
	eqs(t, wri, `A`)

	tryInt(wri.WriteRune('&'))
	eqs(t, wri, `A&amp;`)

	tryInt(wri.WriteRune('B'))
	eqs(t, wri, `A&amp;B`)

	tryInt(wri.WriteRune('\u00a0'))
	eqs(t, wri, `A&amp;B&nbsp;`)

	tryInt(wri.WriteRune('C'))
	eqs(t, wri, `A&amp;B&nbsp;C`)

	tryInt(wri.WriteRune('"'))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;`)

	tryInt(wri.WriteRune('D'))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;D`)

	tryInt(wri.WriteRune('<'))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;D<`)

	tryInt(wri.WriteRune('E'))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;D<E`)

	tryInt(wri.WriteRune('>'))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;D<E>`)

	tryInt(wri.WriteRune('F'))
	eqs(t, wri, `A&amp;B&nbsp;C&quot;D<E>F`)
}

func Test_TextWri_Write(t *testing.T) {
	var wri TextWri
	tryInt(wri.Write([]byte("A&B\u00a0C\"D<E>F")))
	eqs(t, wri, `A&amp;B&nbsp;C"D&lt;E&gt;F`)
}

func Test_TextWri_WriteString(t *testing.T) {
	var wri TextWri
	tryInt(wri.WriteString("A&B\u00a0C\"D<E>F"))
	eqs(t, wri, `A&amp;B&nbsp;C"D&lt;E&gt;F`)
}

func Test_TextWri_WriteRune(t *testing.T) {
	var wri TextWri

	tryInt(wri.WriteRune('A'))
	eqs(t, wri, `A`)

	tryInt(wri.WriteRune('&'))
	eqs(t, wri, `A&amp;`)

	tryInt(wri.WriteRune('B'))
	eqs(t, wri, `A&amp;B`)

	tryInt(wri.WriteRune('\u00a0'))
	eqs(t, wri, `A&amp;B&nbsp;`)

	tryInt(wri.WriteRune('C'))
	eqs(t, wri, `A&amp;B&nbsp;C`)

	tryInt(wri.WriteRune('"'))
	eqs(t, wri, `A&amp;B&nbsp;C"`)

	tryInt(wri.WriteRune('D'))
	eqs(t, wri, `A&amp;B&nbsp;C"D`)

	tryInt(wri.WriteRune('<'))
	eqs(t, wri, `A&amp;B&nbsp;C"D&lt;`)

	tryInt(wri.WriteRune('E'))
	eqs(t, wri, `A&amp;B&nbsp;C"D&lt;E`)

	tryInt(wri.WriteRune('>'))
	eqs(t, wri, `A&amp;B&nbsp;C"D&lt;E&gt;`)

	tryInt(wri.WriteRune('F'))
	eqs(t, wri, `A&amp;B&nbsp;C"D&lt;E&gt;F`)
}

func TestAttr_SetName(t *testing.T) {
	eq(t, Attr{``, ``}.SetName(``), Attr{``, ``})
	eq(t, Attr{`one`, ``}.SetName(``), Attr{``, ``})
	eq(t, Attr{``, `one`}.SetName(``), Attr{``, `one`})
	eq(t, Attr{`one`, `two`}.SetName(``), Attr{``, `two`})

	eq(t, Attr{``, ``}.SetName(`one`), Attr{`one`, ``})
	eq(t, Attr{`one`, ``}.SetName(`two`), Attr{`two`, ``})
	eq(t, Attr{``, `one`}.SetName(`two`), Attr{`two`, `one`})
	eq(t, Attr{`one`, `two`}.SetName(`three`), Attr{`three`, `two`})
}

func TestAttr_Set(t *testing.T) {
	eq(t, Attr{``, ``}.Set(``), Attr{``, ``})
	eq(t, Attr{`one`, ``}.Set(``), Attr{`one`, ``})
	eq(t, Attr{``, `one`}.Set(``), Attr{``, ``})
	eq(t, Attr{`one`, `two`}.Set(``), Attr{`one`, ``})

	eq(t, Attr{``, ``}.Set(`one`), Attr{``, `one`})
	eq(t, Attr{`one`, ``}.Set(`two`), Attr{`one`, `two`})
	eq(t, Attr{``, `one`}.Set(`two`), Attr{``, `two`})
	eq(t, Attr{`one`, `two`}.Set(`three`), Attr{`one`, `three`})
}

func TestAttr_Add(t *testing.T) {
	eq(t, Attr{``, ``}.Add(``), Attr{``, ``})
	eq(t, Attr{`one`, ``}.Add(``), Attr{`one`, ``})
	eq(t, Attr{``, `one`}.Add(``), Attr{``, `one`})
	eq(t, Attr{`one`, `two`}.Add(``), Attr{`one`, `two`})

	eq(t, Attr{``, ``}.Add(`one`), Attr{``, `one`})
	eq(t, Attr{`one`, ``}.Add(`two`), Attr{`one`, `two`})
	eq(t, Attr{``, `one`}.Add(`two`), Attr{``, `one two`})
	eq(t, Attr{`one`, `two`}.Add(`three`), Attr{`one`, `two three`})
}

func TestAttrs_Set(t *testing.T) {
	eq(t, AP().Set(``, ``), nil)
	eq(t, AP().Set(``, `one`), nil)

	eq(t, AP(`one`, `two`).Set(``, ``), AP(`one`, `two`))
	eq(t, AP(`one`, `two`).Set(``, `one`), AP(`one`, `two`))

	eq(
		t,
		AP().Set(`one`, `two`),
		AP(`one`, `two`),
	)

	eq(
		t,
		AP(`one`, `two`).Set(`one`, `three`),
		AP(`one`, `three`),
	)

	eq(
		t,
		AP(`one`, `two`).Set(`three`, `four`),
		AP(`one`, `two`, `three`, `four`),
	)

	eq(
		t,
		AP(`one`, `two`, `one`, `three`).Set(`one`, `four`),
		AP(`one`, `four`, `one`, `four`),
	)

	eq(
		t,
		AP(`one`, `two`, `three`, `four`).Set(`three`, `five`),
		AP(`one`, `two`, `three`, `five`),
	)
}

func TestAttrs_Add(t *testing.T) {
	eq(t, AP().Add(``, ``), nil)
	eq(t, AP().Add(``, `one`), nil)

	eq(t, AP(`one`, `two`).Add(``, ``), AP(`one`, `two`))
	eq(t, AP(`one`, `two`).Add(``, `one`), AP(`one`, `two`))

	eq(
		t,
		AP().Add(`one`, `two`),
		AP(`one`, `two`),
	)

	eq(
		t,
		AP(`one`, `two`).Add(`one`, `three`),
		AP(`one`, `two three`),
	)

	eq(
		t,
		AP(`one`, `two`).Add(`three`, `four`),
		AP(`one`, `two`, `three`, `four`),
	)

	eq(
		t,
		AP(`one`, `two`, `one`, `three`).Add(`one`, `four`),
		AP(`one`, `two four`, `one`, `three four`),
	)

	eq(
		t,
		AP(`one`, `two`, `three`, `four`).Add(`three`, `five`),
		AP(`one`, `two`, `three`, `four five`),
	)
}

func TestElem_AttrSet(t *testing.T) {
	eq(t,
		E(`div`, nil).AttrSet(``, ``),
		Elem{`div`, nil, nil},
	)

	eq(t,
		E(`div`, nil).AttrSet(``, `one`),
		Elem{`div`, nil, nil},
	)

	eq(t,
		E(`div`, nil).AttrSet(`one`, ``),
		Elem{`div`, AP(`one`, ``), nil},
	)

	eq(t,
		E(`div`, nil).AttrSet(`one`, `two`),
		Elem{`div`, AP(`one`, `two`), nil},
	)

	eq(t,
		E(`div`, nil).AttrSet(`one`, `two`).AttrSet(`one`, `three`),
		Elem{`div`, AP(`one`, `three`), nil},
	)

	eq(t,
		E(`div`, nil, `chi`).AttrSet(`one`, `two`).AttrSet(`one`, `three`),
		Elem{`div`, AP(`one`, `three`), []any{`chi`}},
	)

	eq(t,
		E(`div`, nil, `chi`).
			AttrSet(`one`, `two`).
			AttrSet(`three`, `four`).
			AttrSet(`one`, `five`),
		Elem{`div`, AP(`one`, `five`, `three`, `four`), []any{`chi`}},
	)

	eq(t,
		E(`div`, nil, `chi`).
			AttrSet(`one`, `two`).
			AttrSet(`three`, `four`).
			AttrSet(`three`, `five`),
		Elem{`div`, AP(`one`, `two`, `three`, `five`), []any{`chi`}},
	)
}

func TestElem_AttrAdd(t *testing.T) {
	eq(t,
		E(`div`, nil).AttrAdd(``, ``),
		Elem{`div`, nil, nil},
	)

	eq(t,
		E(`div`, nil).AttrAdd(``, `one`),
		Elem{`div`, nil, nil},
	)

	eq(t,
		E(`div`, nil).AttrAdd(`one`, ``),
		Elem{`div`, AP(`one`, ``), nil},
	)

	eq(t,
		E(`div`, nil).AttrAdd(`one`, `two`),
		Elem{`div`, AP(`one`, `two`), nil},
	)

	eq(t,
		E(`div`, nil).AttrAdd(`one`, `two`).AttrAdd(`one`, `three`),
		Elem{`div`, AP(`one`, `two three`), nil},
	)

	eq(t,
		E(`div`, nil, `chi`).AttrAdd(`one`, `two`).AttrAdd(`one`, `three`),
		Elem{`div`, AP(`one`, `two three`), []any{`chi`}},
	)

	eq(t,
		E(`div`, nil, `chi`).
			AttrAdd(`one`, `two`).
			AttrAdd(`three`, `four`).
			AttrAdd(`one`, `five`),
		Elem{`div`, AP(`one`, `two five`, `three`, `four`), []any{`chi`}},
	)

	eq(t,
		E(`div`, nil, `chi`).
			AttrAdd(`one`, `two`).
			AttrAdd(`three`, `four`).
			AttrAdd(`three`, `five`),
		Elem{`div`, AP(`one`, `two`, `three`, `four five`), []any{`chi`}},
	)
}

func TestElem_GoString(t *testing.T) {
	eq(t,
		fmt.Sprintf(
			`%#v`,
			E(`div`, AP(`class`, `one`), `two`, []any{10, `three`}),
		),
		"E(`div`, AP(`class`, `one`), `two`, 10, `three`)",
	)
}

func TestA_GoString(t *testing.T) {
	eq(t,
		fmt.Sprintf(
			`%#v`,
			AP(`one`, `two`, `three`, `four`).A(Attr{}).A(Attr{}).A(Attr{}),
		),
		"AP(`one`, `two`, `three`, `four`)",
	)
}

func TestVac(t *testing.T) {
	eq(t, Vac(nil), nil)
	eq(t, Vac((*string)(nil)), nil)
	eq(t, Vac([]any{}), nil)
	eq(t, Vac([]any{nil}), nil)
	eq(t, Vac([]any{nil, (*string)(nil)}), nil)
	eq(t, Vac([]byte(nil)), nil)
	eq(t, Vac(Bui(nil)), nil)

	eq(t, Vac(""), "")
	eq(t, Vac(0), 0)
	eq(t, any(Vac([]any{""})), any([]any{""}))
	eq(t, any(Vac([]any{0})), any([]any{0}))
}

func TestNonEscWri_grow(t *testing.T) {
	var wri NonEscWri

	eq(t, len(wri), 0)
	eq(t, cap(wri), 0)

	wri.grow(1)
	eq(t, len(wri), 0)
	eq(t, cap(wri), 1)

	wri.grow(11)
	eq(t, len(wri), 0)
	eq(t, cap(wri), 13)

	wri.grow(7)
	eq(t, len(wri), 0)
	eq(t, cap(wri), 13)

	wri = wri[:1][1:]
	eq(t, len(wri), 0)
	eq(t, cap(wri), 12)

	wri.grow(7)
	eq(t, len(wri), 0)
	eq(t, cap(wri), 12)

	wri.grow(11)
	eq(t, len(wri), 0)
	eq(t, cap(wri), 12)

	wri.grow(13)
	eq(t, len(wri), 0)
	eq(t, cap(wri), 37)
}

func eqs[A fmt.Stringer](t testing.TB, act A, exp string) {
	eq(t, act.String(), exp)
}

func eq[A any](t testing.TB, act, exp A) {
	t.Helper()

	if !r.DeepEqual(act, exp) {
		t.Fatalf(`
actual (detailed):
	%#[1]v
expected (detailed):
	%#[2]v
actual (simple):
	%[1]v
expected (simple):
	%[2]v
`, act, exp)
	}
}

func tryInt(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
