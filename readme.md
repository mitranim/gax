## Overview

Simple system for writing HTML/XML as Go code. Better-performing replacement for `html/template` and `text/template`. Vaguely inspired by JS library https://github.com/mitranim/prax.

Features / benefits:

  * No weird special language to learn.
  * Use actual Go code.
  * Use normal Go conditionals.
  * Use normal Go loops.
  * Use normal Go functions.
  * Benefit from static typing.
  * Benefit from Go code analysis.
  * Benefit from Go performance.
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

  x "github.com/mitranim/gax"
)

var (
  E  = x.E
  AP = x.AP
)

func main() {
  fmt.Println(Page(mockDat))
  // <!doctype html><html lang="en"><head><meta charset="utf-8"><link rel="icon" href="data:;base64,="><title>Posts</title></head><body><h1 class="title">Posts</h1><h2>Post0</h2><h2>Post1</h2></body></html>
}

func Page(dat Dat) x.Bui {
  return x.F(
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

Gax easily beats `text/template` and `html/template`. The more dynamic a template is, the better it gets.

In the static benchmark, Gax renders the markup just once into a global variable, benchmarking the writing performance of `bytes.Buffer` for "fairness" with templating. This is recommended for all completely static markup.

The dynamic benchmark is intentionally naive, avoiding some Gax optimizations such as static prerender, to mimic simple user code.

```sh
go test -bench . -benchmem
```

```
cpu: Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz
Benchmark_static_gax-12           12817749       88.78 ns/op       384 B/op        1 allocs/op
Benchmark_static_template-12       6270380       193.4 ns/op       480 B/op        3 allocs/op
Benchmark_dynamic_gax-12             68360       16726 ns/op      9320 B/op      140 allocs/op
Benchmark_dynamic_template-12         9432      130345 ns/op     61847 B/op     1376 allocs/op
```

## Changelog

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
