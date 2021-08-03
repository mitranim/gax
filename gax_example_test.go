package gax_test

import (
	"fmt"

	x "github.com/mitranim/gax"
)

func ExampleBui() {
	var (
		E  = x.E
		AP = x.AP
	)

	type Dat struct {
		Title string
		Posts []string
	}

	dat := Dat{
		Title: `Posts`,
		Posts: []string{`Post0`, `Post1`},
	}

	bui := x.F(
		x.Str(x.Doctype),
		E(`html`, AP(`lang`, `en`),
			E(`head`, nil,
				E(`meta`, AP(`charset`, `utf-8`)),
				E(`link`, AP(`rel`, `icon`, `href`, `data:;base64,=`)),

				// Use normal Go conditionals.
				func(b *x.Bui) {
					if dat.Title != "" {
						b.E(`title`, nil, dat.Title)
					} else {
						b.E(`title`, nil, `test markup`)
					}
				},
			),

			E(`body`, nil,
				E(`h1`, AP(`class`, `title`), `Posts`),

				// Use normal Go loops.
				func(b *x.Bui) {
					for _, post := range dat.Posts {
						b.E(`h2`, nil, post)
					}
				},
			),
		),
	)

	fmt.Println(bui)
	// Output:
	// <!doctype html><html lang="en"><head><meta charset="utf-8"><link rel="icon" href="data:;base64,="><title>Posts</title></head><body><h1 class="title">Posts</h1><h2>Post0</h2><h2>Post1</h2></body></html>
}

func ExampleE() {
	var (
		E  = x.E
		AP = x.AP
	)

	fmt.Println(
		E(`span`, AP(`aria-hidden`, `true`), `🔥`),
	)
	// Output:
	// <span aria-hidden="true">🔥</span>
}

func ExampleF() {
	var doc = x.F(
		x.Str(x.Doctype),
		x.E(`html`, nil),
	)

	fmt.Println(doc)
	// Output:
	// <!doctype html><html></html>
}

func ExampleDoctype() {
	bui := x.Bui(x.Doctype)
	bui.E(`html`, nil)

	fmt.Println(bui)
	// Output:
	// <!doctype html><html></html>
}

func ExampleAP() {
	fmt.Println(
		x.AP(
			`href`, `/`,
			`aria-current`, `page`,
			`class`, `some-class`,
		),
	)
	// Output:
	// href="/" aria-current="page" class="some-class"
}

func ExampleA() {
	attrs := x.A(
		x.Attr{`class`, `some-class`},
		x.Attr{`style`, `some: style`},
	)
	fmt.Println(attrs)
	// Output:
	// class="some-class" style="some: style"
}

func ExampleAttrs() {
	attrs := x.Attrs{
		{`class`, `some-class`},
		{`style`, `some: style`},
	}
	fmt.Println(attrs)
	// Output:
	// class="some-class" style="some: style"
}

func ExampleAttrs_A() {
	cur := func() x.Attr { return x.Attr{`aria-current`, `page`} }
	bg := func() x.Attr { return x.Attr{`style`, `background-image: url(...)`} }

	fmt.Println(
		x.AP(`class`, `some-class`).A(cur(), bg()),
	)
	// class="some-class" aria-current="page" style="background-image: url(...)"
}

func ExampleAttrs_AP() {
	fmt.Println(
		x.AP(`class`, `some-class`).AP(`href`, `/`),
	)
	// class="some-class" href="/"
}

func ExampleAttr() {
	fmt.Println(x.Attr{`class`, `some-class`})
	// Output:
	// class="some-class"
}

func ExampleStr_Render() {
	var bui x.Bui
	bui.E(`div`, nil, x.Str(`<script>alert('hacked!')</script>`))

	fmt.Println(bui)
	// Output:
	// <div><script>alert('hacked!')</script></div>
}

func ExampleBui_Render() {
	var bui x.Bui
	bui.E(`div`, nil, x.Bui(`<script>alert('hacked!')</script>`))

	fmt.Println(bui)
	// Output:
	// <div><script>alert('hacked!')</script></div>
}
