## Overview

Simple system for writing HTML/XML as Go code. Better-performing replacement for `html/template` and `text/template`. Vaguely inspired by JS library https://github.com/mitranim/prax.

Advantages over string-based templating:

  * No weird special language to learn.
  * Normal Go code.
  * Normal Go conditionals.
  * Normal Go loops.
  * Normal Go functions.
  * Normal Go static typing.
  * Normal Go code analysis.
  * Much better performance.

Other features / benefits:

  * Tiny and dependency-free (only stdlib).

## TOC

* [Usage](#usage)
* [Performance](#performance)
* [Changelog](#changelog)
* [License](#license)
* [Misc](#misc)

## Usage

API docs: https://pkg.go.dev/github.com/mitranim/gax.

```golang
package main

import (
  "fmt"

  gax "github.com/mitranim/gax"
)

var (
  E  = gax.E
  AP = gax.AP
)

func main() {
  fmt.Println(Page(mockDat))
  // <!doctype html><html lang="en"><head><meta charset="utf-8"><link rel="icon" href="data:;base64,="><title>Posts</title></head><body><h1 class="title">Posts</h1><h2>Post0</h2><h2>Post1</h2></body></html>
}

func Page(dat Dat) gax.Bui {
  return gax.F(
    gax.Str(gax.Doctype),
    E(`html`, AP(`lang`, `en`),
      E(`head`, nil,
        E(`meta`, AP(`charset`, `utf-8`)),
        E(`link`, AP(`rel`, `icon`, `href`, `data:;base64,=`)),

        // Use normal Go conditionals.
        func(bui *gax.Bui) {
          if dat.Title != `` {
            bui.E(`title`, nil, dat.Title)
          } else {
            bui.E(`title`, nil, `test markup`)
          }
        },
      ),

      E(`body`, nil,
        E(`h1`, AP(`class`, `title`), `Posts`),

        // Use normal Go loops.
        func(bui *gax.Bui) {
          for _, post := range dat.Posts {
            bui.E(`h2`, nil, post)
          }
        },
      ),
    ),
  )
}

var mockDat = Dat{
  Title: `Posts`,
  Posts: []string{`Post0`, `Post1`},
}

type Dat struct {
  Title string
  Posts []string
}
```

## Performance

Gax easily beats `text/template` and `html/template`. The more complex a template is, the better it gets.

The static benchmark is "unfair" because the Gax version renders just once into a global variable. This is recommended for all completely static markup. Prerendering is also possible with `text/template` and `html/template`, but syntactically inconvenient and usually avoided. With Gax it's syntactically convenient and easily done, and the benchmark reflects that.

The dynamic benchmark is intentionally naive, avoiding some Gax optimizations such as static prerender, to mimic simple user code.

```sh
go test -bench . -benchmem
```

```
cpu: Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz
Benchmark_template_static-12   17002192    67.31 ns/op     48 B/op     1 allocs/op
Benchmark_gax_static-12        640101200   1.845 ns/op      0 B/op     0 allocs/op
Benchmark_template_dynamic-12      9205   130812 ns/op  51811 B/op  1370 allocs/op
Benchmark_gax_dynamic-12          70465    17090 ns/op  10376 B/op   169 allocs/op
```

## Changelog

### `v0.3.0`

* Renamed `.Append` to `.AppendTo` for consistency with other libraries.
* `Elem` with empty `.Tag` no longer renders anything. As a result, zero value of `Elem` is the same as nil. This can be convenient for functions that return `Elem`.
* Added `LinkBlank`.
* Require Go 1.20.

### `v0.2.1`

Child rendering now supports walking `[]Ren` and `[]T where T is Ren`.

### `v0.2.0`

API revision. Now supports both the append-only style via `Bui.E` and the expression style via free `E`. Mix and match for simpler code.

### `v0.1.4`

`Bui.Child` also supports `func(E)`.

### `v0.1.3`

Added `Bui.With` and `Ebui`.

### `v0.1.2`

Minor syntactic bugfix.

### `v0.1.1`

Converted methods `.WriteTo(*[]byte)` methods to `.Append([]byte) []byte` for better compliance with established interfaces.

### `v0.1.0`

Init.

## License

https://unlicense.org

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
