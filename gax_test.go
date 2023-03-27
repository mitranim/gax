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

	eqs(
		t,
		strings.TrimSpace(`
			<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>
		`),
		bui,
	)
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

	eqs(
		t,
		strings.TrimSpace(`
			<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>
		`),
		bui,
	)
}

func Test_Bui_Attr(t *testing.T) {
	var bui Bui
	bui.Attr(Attr{`class`, `<one>&"</one>`})
	eqs(t, ` class="<one>&amp;&quot;</one>"`, bui)
}

func Test_Bui_Attrs(t *testing.T) {
	var bui Bui

	bui.Attrs(
		Attr{`class`, `<one>&"</one>`},
		Attr{`style`, `<two>&"</two>`},
	)

	eqs(t, ` class="<one>&amp;&quot;</one>" style="<two>&amp;&quot;</two>"`, bui)
}

func Test_Bui_EscString(t *testing.T) {
	var bui Bui
	bui.EscString(`<one>&"</one>`)
	eqs(t, `&lt;one&gt;&amp;"&lt;/one&gt;`, bui)
}

func Test_Bui_EscBytes(t *testing.T) {
	var bui Bui
	bui.EscBytes([]byte(`<one>&"</one>`))
	eqs(t, `&lt;one&gt;&amp;"&lt;/one&gt;`, bui)
}

func Test_Bui_Child_escaping(t *testing.T) {
	test := childTest(t)
	test(`&lt;one&gt;&amp;"&lt;/one&gt;`, `<one>&"</one>`)
	test(`&lt;one&gt;&amp;"&lt;/one&gt;`, []byte(`<one>&"</one>`))
}

func Test_Bui_Child_Ren(t *testing.T) {
	test := childTest(t)

	test(`str`, Str(`str`))
	test(`str`, Bui(`str`))

	test(
		`<one>&"</one>`,
		Str(`<one>&"</one>`),
	)

	test(
		`<one>&"</one>`,
		Bui(`<one>&"</one>`),
	)

	test(
		`<a>one</a><bui>two</bui><c>three</c>`,
		Str(`<a>one</a><bui>two</bui><c>three</c>`),
	)

	test(
		`<a>one</a><bui>two</bui><c>three</c>`,
		Bui(`<a>one</a><bui>two</bui><c>three</c>`),
	)

	test(
		`<one two="three">four</one>`,
		E(`one`, AP(`two`, `three`), `four`),
	)
}

func Test_Bui_Child_slices(t *testing.T) {
	test := childTest(t)

	test(`10str20`, []any{10, nil, "str", []any{nil, 20}})

	test(
		`<one></one>two<three></three>`,
		[]Ren{
			E(`one`, nil),
			nil,
			Str(`two`),
			nil,
			E(`three`, nil),
		},
	)

	test(
		`<one></one><two></two>`,
		[]Elem{
			E(`one`, nil),
			Elem{},
			E(`two`, nil),
		},
	)
}

func Test_Bui_Child_funcs(t *testing.T) {
	test := childTest(t)

	test(``, (*func())(nil))
	test(``, (*func(*Bui))(nil))
	test(`<div></div>`, func(b *Bui) { b.E(`div`, nil) })
	test(`str`, func(b *Bui) { b.T(`str`) })
}

func Test_Bui_Child_misc(t *testing.T) {
	test := childTest(t)

	test(``, nil)
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
}

func childTest(t testing.TB) func(string, any) {
	return func(exp string, val any) {
		var bui Bui
		bui.Child(val)
		eqs(t, exp, bui)
	}
}

func Test_F(t *testing.T) {
	eqs(
		t,
		`<!doctype html><html lang="en"></html>`,
		F(
			Str(Doctype),
			func(b *Bui) { b.E(`html`, AP(`lang`, `en`)) },
		),
	)
}

// Incomplete test; should also verify zero-alloc.
func Test_Bui_Bytes(t *testing.T) {
	eqs(t, `<div>hello world!</div>`, Bui(`<div>hello world!</div>`))
}

// Incomplete test; should also verify zero-alloc.
func Test_Bui_String(t *testing.T) {
	eq(t, `<div>hello world!</div>`, Bui(`<div>hello world!</div>`).String())
	eq(t, `<div>hello world!</div>`, string(Bui(`<div>hello world!</div>`)))
}

func Test_AttrWri_Write(t *testing.T) {
	var wri AttrWri
	tryInt(wri.Write([]byte("A&B\u00a0C\"D<E>F")))
	eqs(t, `A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_AttrWri_WriteString(t *testing.T) {
	var wri AttrWri
	tryInt(wri.WriteString("A&B\u00a0C\"D<E>F"))
	eqs(t, `A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_AttrWri_WriteRune(t *testing.T) {
	var wri AttrWri

	tryInt(wri.WriteRune('A'))
	eqs(t, `A`, wri)

	tryInt(wri.WriteRune('&'))
	eqs(t, `A&amp;`, wri)

	tryInt(wri.WriteRune('B'))
	eqs(t, `A&amp;B`, wri)

	tryInt(wri.WriteRune('\u00a0'))
	eqs(t, `A&amp;B&nbsp;`, wri)

	tryInt(wri.WriteRune('C'))
	eqs(t, `A&amp;B&nbsp;C`, wri)

	tryInt(wri.WriteRune('"'))
	eqs(t, `A&amp;B&nbsp;C&quot;`, wri)

	tryInt(wri.WriteRune('D'))
	eqs(t, `A&amp;B&nbsp;C&quot;D`, wri)

	tryInt(wri.WriteRune('<'))
	eqs(t, `A&amp;B&nbsp;C&quot;D<`, wri)

	tryInt(wri.WriteRune('E'))
	eqs(t, `A&amp;B&nbsp;C&quot;D<E`, wri)

	tryInt(wri.WriteRune('>'))
	eqs(t, `A&amp;B&nbsp;C&quot;D<E>`, wri)

	tryInt(wri.WriteRune('F'))
	eqs(t, `A&amp;B&nbsp;C&quot;D<E>F`, wri)
}

func Test_TextWri_Write(t *testing.T) {
	var wri TextWri
	tryInt(wri.Write([]byte("A&B\u00a0C\"D<E>F")))
	eqs(t, `A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func Test_TextWri_WriteString(t *testing.T) {
	var wri TextWri
	tryInt(wri.WriteString("A&B\u00a0C\"D<E>F"))
	eqs(t, `A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func Test_TextWri_WriteRune(t *testing.T) {
	var wri TextWri

	tryInt(wri.WriteRune('A'))
	eqs(t, `A`, wri)

	tryInt(wri.WriteRune('&'))
	eqs(t, `A&amp;`, wri)

	tryInt(wri.WriteRune('B'))
	eqs(t, `A&amp;B`, wri)

	tryInt(wri.WriteRune('\u00a0'))
	eqs(t, `A&amp;B&nbsp;`, wri)

	tryInt(wri.WriteRune('C'))
	eqs(t, `A&amp;B&nbsp;C`, wri)

	tryInt(wri.WriteRune('"'))
	eqs(t, `A&amp;B&nbsp;C"`, wri)

	tryInt(wri.WriteRune('D'))
	eqs(t, `A&amp;B&nbsp;C"D`, wri)

	tryInt(wri.WriteRune('<'))
	eqs(t, `A&amp;B&nbsp;C"D&lt;`, wri)

	tryInt(wri.WriteRune('E'))
	eqs(t, `A&amp;B&nbsp;C"D&lt;E`, wri)

	tryInt(wri.WriteRune('>'))
	eqs(t, `A&amp;B&nbsp;C"D&lt;E&gt;`, wri)

	tryInt(wri.WriteRune('F'))
	eqs(t, `A&amp;B&nbsp;C"D&lt;E&gt;F`, wri)
}

func TestElem_GoString(t *testing.T) {
	eq(t,
		"E(`div`, AP(`class`, `one`), `two`, 10, `three`)",
		fmt.Sprintf(
			`%#v`,
			E(`div`, AP(`class`, `one`), `two`, []any{10, `three`}),
		),
	)
}

func TestA_GoString(t *testing.T) {
	eq(t,
		"AP(`one`, `two`, `three`, `four`)",
		fmt.Sprintf(
			`%#v`,
			AP(`one`, `two`, `three`, `four`).A(Attr{}).A(Attr{}).A(Attr{}),
		),
	)
}

func TestVac(t *testing.T) {
	eq(t, nil, Vac(nil))
	eq(t, nil, Vac((*string)(nil)))
	eq(t, nil, Vac([]any{}))
	eq(t, nil, Vac([]any{nil}))
	eq(t, nil, Vac([]any{nil, (*string)(nil)}))
	eq(t, nil, Vac([]byte(nil)))
	eq(t, nil, Vac(Bui(nil)))

	eq(t, "", Vac(""))
	eq(t, 0, Vac(0))
	eq(t, []any{""}, Vac([]any{""}))
	eq(t, []any{0}, Vac([]any{0}))
}

func TestNonEscWri_grow(t *testing.T) {
	var wri NonEscWri

	eq(t, 0, len(wri))
	eq(t, 0, cap(wri))

	wri.grow(1)
	eq(t, 0, len(wri))
	eq(t, 1, cap(wri))

	wri.grow(11)
	eq(t, 0, len(wri))
	eq(t, 13, cap(wri))

	wri.grow(7)
	eq(t, 0, len(wri))
	eq(t, 13, cap(wri))

	wri = wri[:1][1:]
	eq(t, 0, len(wri))
	eq(t, 12, cap(wri))

	wri.grow(7)
	eq(t, 0, len(wri))
	eq(t, 12, cap(wri))

	wri.grow(11)
	eq(t, 0, len(wri))
	eq(t, 12, cap(wri))

	wri.grow(13)
	eq(t, 0, len(wri))
	eq(t, 37, cap(wri))
}

func eqs(t testing.TB, exp string, act []byte) {
	eq(t, exp, string(act))
}

func eq(t testing.TB, exp, act any) {
	t.Helper()
	if !r.DeepEqual(exp, act) {
		t.Fatalf(`
expected (detailed):
	%#[1]v
actual (detailed):
	%#[2]v
expected (simple):
	%[1]v
actual (simple):
	%[2]v
`, exp, act)
	}
}

func tryInt(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
