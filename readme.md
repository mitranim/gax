## Overview

Simple system for writing HTML/XML as Go code. Better-performing replacement for `html/template` and `text/template`. Vaguely inspired by JS library https://github.com/mitranim/prax, but uses a different design.

API docs: https://pkg.go.dev/github.com/mitranim/gax.

Features / benefits:

  * No weird special language to learn.
  * Use actual Go code.
  * Use normal Go conditionals.
  * Use normal Go loops.
  * Use normal Go functions.
  * Benefit from static typing.
  * Benefit from Go code analysis.
  * Benefit from Go performance.
  * Tiny and dependency-free.

## Usage

```golang
package main

import (
  "fmt"
  "github.com/mitranim/gax"
)

type A = gax.A

func main() {
  bui := gax.Bui(gax.Doctype)

  render(bui.E, mockDat)

  fmt.Println(bui)
  // <!doctype html><html lang="en"><head><meta charset="utf-8"><link rel="icon" href="data:;base64,="><title>Posts</title></head><body><h1 class="title">Posts</h1><h2>Post0</h2><h2>Post1</h2></body></html>
}

func render(E gax.E, dat Dat) {
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

For mostly-static templates, Gax loses to `html/template` but remains more than fast enough. For anything dynamic, Gax seems to perform several times better. The more complicated a template is, the better it gets.

The following output is from `gax_bench_test.go`, ran via:

```sh
go test -bench . -benchmem
```

```
cpu: Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz
Benchmark_static_gax-12           327079     3289 ns/op     1544 B/op     25 allocs/op
Benchmark_static_template-12     6301851    193.9 ns/op      480 B/op      3 allocs/op
Benchmark_dynamic_gax-12          203464     5842 ns/op     2608 B/op     49 allocs/op
Benchmark_dynamic_template-12      39686    30669 ns/op    12862 B/op    297 allocs/op
```

## License

https://unlicense.org

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
