package gax

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_static_gax(_ *testing.T) {
	bui := Bui(Doctype)
	renderStatic(bui.E)
	eq(
		`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>`,
		bui.String(),
	)
}

func Test_static_template(_ *testing.T) {
	eq(
		`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>`,
		templateToString(tplStatic, nil),
	)
}

func Test_static_equiv(_ *testing.T) {
	bui := Bui(Doctype)
	renderStatic(bui.E)
	eq(bui.String(), templateToString(tplStatic, nil))
}

func Test_dynamic_gax(_ *testing.T) {
	bui := Bui(Doctype)
	renderDynamic(bui.E, mockDat)
	eq(
		`<!doctype html><html lang="en"><head><link rel="icon" href="data:;base64,="><title>Posts</title><meta property="og:title" content="Posts"><meta name="description" content="Random notes and thoughts"></head><body><nav><a href="/">home</a><a href="/works">works</a><a href="/posts" aria-current="page">posts</a><a href="/demos">demos</a><span>Updated Apr 05 3123</span></nav><div role="main"><h1>Posts</h1><h2><a href="/posts/post1.html">post 1</a></h2><h2><a href="/posts/post2.html">post 2</a></h2><h2><a href="/posts/post3.html">post 3</a></h2></div></body></html>`,
		bui.String(),
	)
}

func Test_dynamic_template(_ *testing.T) {
	eq(
		`<!doctype html><html lang="en"><head><link rel="icon" href="data:;base64,="><title>Posts</title><meta property="og:title" content="Posts"><meta name="description" content="Random notes and thoughts"></head><body><nav><a href="/">home</a><a href="/works">works</a><a href="/posts" aria-current="page">posts</a><a href="/demos">demos</a><span>Updated Apr 05 3123</span></nav><div role="main"><h1>Posts</h1><h2><a href="/posts/post1.html">post 1</a></h2><h2><a href="/posts/post2.html">post 2</a></h2><h2><a href="/posts/post3.html">post 3</a></h2></div></body></html>`,
		templateToString(tplDynamic, mockDat),
	)
}

func Test_dynamic_equiv(_ *testing.T) {
	bui := Bui(Doctype)
	renderDynamic(bui.E, mockDat)
	eq(bui.String(), templateToString(tplDynamic, mockDat))
}

func Benchmark_static_gax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bui := Bui(Doctype)
		renderStatic(bui.E)
	}
}

func Benchmark_static_template(b *testing.B) {
	for i := 0; i < b.N; i++ {
		must(tplStatic.Execute(new(bytes.Buffer), nil))
	}
}

func Benchmark_dynamic_gax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bui := make(Bui, 0, 1024)
		renderDynamic(bui.E, mockDat)
	}
}

func Benchmark_dynamic_template(b *testing.B) {
	for i := 0; i < b.N; i++ {
		must(tplDynamic.Execute(new(bytes.Buffer), mockDat))
	}
}

type MockDat struct {
	Page
	Posts []Post
}

type Page struct {
	Path        string
	Title       string
	Desc        string
	Type        string
	Image       string
	ColorScheme string
	Posts       []Post
}

type Post struct {
	Page
	HtmlBody    []byte
	PublishedAt *time.Time
	UpdatedAt   *time.Time
	IsListed    bool
}

var mockDat = MockDat{
	Page: Page{
		Path:  `posts.html`,
		Title: `Posts`,
		Desc:  `Random notes and thoughts`,
	},
	Posts: []Post{
		Post{Page: Page{Title: `post 0`, Path: `/posts/post0.html`}, IsListed: false},
		Post{Page: Page{Title: `post 1`, Path: `/posts/post1.html`}, IsListed: true},
		Post{Page: Page{Title: `post 2`, Path: `/posts/post2.html`}, IsListed: true},
		Post{Page: Page{Title: `post 3`, Path: `/posts/post3.html`}, IsListed: true},
	},
}

/*
Using a bound method, rather than the builder itself, incurs a significant
slowdown. We're using this in benchmarks because we kinda recommend this for
syntactic usability.
*/
func renderStatic(E E) {
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
}

var tplStatic = templateMake(nil, `
<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="icon" href="data:;base64,=">
		<title>test markup</title>
	</head>
	<body class="stretch-to-viewport">
		<h1 class="title">mock markup</h1>
		<div class="main">hello world!</div>
	</body>
</html>
`)

func renderDynamic(E E, dat MockDat) {
	var scheme Attr
	if dat.ColorScheme != "" {
		scheme = Attr{`class`, dat.ColorScheme}
	}

	E(`html`, A{{`lang`, `en`}, scheme}, func() {
		E(`head`, nil, func() {
			E(`link`, A{{`rel`, `icon`}, {`href`, `data:;base64,=`}})

			if dat.Title != "" {
				E(`title`, nil, dat.Title)
				E(`meta`, A{{`property`, `og:title`}, {`content`, dat.Title}})
			} else {
				E(`title`, nil, `mock site`)
			}

			if dat.Desc != "" {
				E(`meta`, A{{`name`, `description`}, {`content`, dat.Desc}})
			}

			if dat.Image != "" {
				E(`meta`, A{{`property`, `og:image`}, {`content`, `/images/` + dat.Image}})
			}

			if dat.Type != "" {
				E(`meta`, A{{`property`, `og:type`}, {`content`, dat.Type}})
				E(`meta`, A{{`property`, `og:site_name`}, {`content`, `mock site`}})
			}
		})

		E(`body`, nil, func() {
			E(`nav`, nil, func() {
				E(`a`, A{{`href`, "/"}, cur(dat, "index.html")}, `home`)
				E(`a`, A{{`href`, "/works"}, cur(dat, "works.html")}, `works`)
				E(`a`, A{{`href`, "/posts"}, cur(dat, "posts.html")}, `posts`)
				E(`a`, A{{`href`, "/demos"}, cur(dat, "demos.html")}, `demos`)
				E(`span`, nil, `Updated `+inst)
			})

			E(`div`, A{{`role`, `main`}}, func() {
				E(`h1`, nil, `Posts`)

				if len(dat.Posts) > 0 {
					for _, post := range dat.Posts {
						if !post.IsListed {
							continue
						}
						E(`h2`, nil, func() {
							E(`a`, A{{`href`, post.Path}}, post.Title)
						})
						if post.Desc != "" {
							E(`p`, nil, post.Desc)
						}
					}
				} else {
					E(`p`, nil, `Oops! It appears there are no public posts yet.`)
				}
			})
		})
	})
}

var tplDynamic = templateMake(
	template.FuncMap{
		"cur":         curFun,
		"now":         func() string { return inst },
		"listedPosts": listedPosts,
	},
	`
{{define "site-top.html"}}
<!doctype html>
<html lang="en"{{if .ColorScheme}} class="{{.ColorScheme}}"{{end}}>
	<head>
		<link rel="icon" href="data:;base64,=">
		{{if .Title}}
			<title>{{.Title}}</title>
			<meta property="og:title" content="{{.Title}}">
		{{else}}
			<title>mock site</title>
		{{end}}
		{{if .Desc}}
			<meta name="description" content="{{.Desc}}">
		{{end}}
		{{if .Image}}
			<meta property="og:image" content="/images/{{.Image}}">
		{{end}}
		{{if .Type}}
			<meta property="og:type" content="{{.Type}}">
			<meta property="og:site_name" content="about:mitranim">
		{{end}}
	</head>
	<body>
{{end}}

{{define "site-bottom.html"}}
	</body>
</html>
{{end}}

{{define "nav.html"}}
	<nav>
		<a href="/"      {{- cur . "index.html"}}>home</a>
		<a href="/works" {{- cur . "works.html"}}>works</a>
		<a href="/posts" {{- cur . "posts.html"}}>posts</a>
		<a href="/demos" {{- cur . "demos.html"}}>demos</a>
		<span>Updated {{now}}</span>
	</nav>
{{end}}

{{template "site-top.html" .}}
{{template "nav.html" .}}
	<div role="main">
		<h1>Posts</h1>

		{{range $post := listedPosts .Posts}}
			<h2><a href="{{$post.Path}}">{{$post.Title}}</a></h2>
			{{if $post.Desc}}<p>{{$post.Desc}}</p>{{end}}
		{{else}}
			<p>Oops! It appears there are no public posts yet.</p>
		{{end}}
	</div>
{{template "site-bottom.html" .}}
`,
)

func templateMake(funs template.FuncMap, str string) *template.Template {
	tpl := template.New("")
	tpl.Option(`missingkey=error`)
	tpl.Funcs(funs)
	template.Must(tpl.Parse(trimLines(str)))
	return tpl
}

func templateToString(tpl *template.Template, val interface{}) string {
	var buf strings.Builder
	must(tpl.Execute(&buf, val))
	return buf.String()
}

func cur(dat MockDat, path string) Attr {
	if dat.Path == path {
		return Attr{`aria-current`, `page`}
	}
	return Attr{}
}

func curFun(dat MockDat, path string) template.HTMLAttr {
	if dat.Path == path {
		return ` aria-current="page"`
	}
	return ``
}

func trimLines(str string) string {
	return strings.TrimSpace(strings.Join(reLines.Split(str, -1), ""))
}

var reLines = regexp.MustCompile(`\s*(?:\r|\n)\s*`)

var inst = time.Date(3123, 4, 5, 6, 7, 8, 9, time.UTC).Format("Jan 02 2006")

func listedPosts(vals []Post) (out []Post) {
	for _, val := range vals {
		if val.IsListed {
			out = append(out, val)
		}
	}
	return out
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
