package gax_test

import (
	"fmt"

	"github.com/mitranim/gax"
)

func ExampleBui() {
	type A = gax.A

	type Dat struct {
		Title string
		Posts []string
	}

	dat := Dat{
		Title: `Posts`,
		Posts: []string{`Post0`, `Post1`},
	}

	bui := gax.Bui(gax.Doctype)
	E := bui.E

	E(`html`, A{{`lang`, `en`}}, func() {
		E(`head`, nil, func() {
			E(`meta`, A{{`charset`, `utf-8`}})
			E(`link`, A{{`rel`, `icon`}, {`href`, `data:;base64,=`}})

			// Use normal Go conditionals.
			if dat.Title != "" {
				E(`title`, nil, dat.Title)
			} else {
				E(`title`, nil, `test markup`)
			}
		})

		E(`body`, nil, func() {
			E(`h1`, A{{`class`, `title`}}, `Posts`)

			// Use normal Go loops.
			for _, post := range dat.Posts {
				E(`h2`, nil, post)
			}
		})
	})

	fmt.Println(bui)
	// Output:
	// <!doctype html><html lang="en"><head><meta charset="utf-8"><link rel="icon" href="data:;base64,="><title>Posts</title></head><body><h1 class="title">Posts</h1><h2>Post0</h2><h2>Post1</h2></body></html>
}

func ExampleBui_With() {
	type A = gax.A
	type E = gax.E

	fmt.Println(gax.Bui(gax.Doctype).With(func(E E) {
		E(`html`, A{{`lang`, `en`}})
	}))
	// Output:
	// <!doctype html><html lang="en"></html>
}

func ExampleEbui() {
	type A = gax.A
	type E = gax.E

	fmt.Println(gax.Ebui(func(E E) {
		E(`span`, A{{`aria-hidden`, `true`}}, `🔥`)
	}))
	// Output:
	// <span aria-hidden="true">🔥</span>
}

func ExampleDoctype() {
	bui := gax.Bui(gax.Doctype)
	bui.E(`html`, nil)
	fmt.Println(bui)
	// Output:
	// <!doctype html><html></html>
}

func ExampleA() {
	attrs := gax.A{
		{"class", "some-class"},
		{"style", "some: style"},
	}
	fmt.Println(attrs)
	// Output:
	// class="some-class" style="some: style"
}

func ExampleAttr() {
	fmt.Println(gax.Attr{"class", "some-class"})
	// Output:
	// class="some-class"
}

func ExampleE() {
	_ = func(E gax.E) {
		E("div", nil, "...")
	}
}

func ExampleString() {
	var b gax.Bui
	b.E(`div`, nil, gax.String(`<script>alert('hacked!')</script>`))
	fmt.Println(b)
	// Output:
	// <div><script>alert('hacked!')</script></div>
}

func ExampleBytes() {
	var b gax.Bui
	b.E(`div`, nil, gax.Bytes(`<script>alert('hacked!')</script>`))
	fmt.Println(b)
	// Output:
	// <div><script>alert('hacked!')</script></div>
}
