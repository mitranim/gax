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
	eqs(
		`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>`,
		renderedStatic,
	)
}

func Test_static_template(_ *testing.T) {
	eq(
		`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="icon" href="data:;base64,="><title>test markup</title></head><body class="stretch-to-viewport"><h1 class="title">mock markup</h1><div class="main">hello world!</div></body></html>`,
		templateToString(tplStatic, nil),
	)
}

func Test_static_equiv(_ *testing.T) {
	eqs(
		templateToString(tplStatic, nil),
		renderedStatic,
	)
}

func Test_dynamic_gax(_ *testing.T) {
	eqs(
		`<!doctype html><html lang="en"><head><link rel="icon" href="data:;base64,="><title>Posts</title><meta property="og:title" content="Posts"><meta name="description" content="Random notes and thoughts"></head><body><nav><a href="/">home</a><a href="/works">works</a><a href="/posts" aria-current="page">posts</a><a href="/demos">demos</a><span>Updated Apr 05 3123</span></nav><div role="main"><h1>Posts</h1><h2><a href="/posts/post1.html">post 1</a></h2><h2><a href="/posts/post2.html">post 2</a></h2><h2><a href="/posts/post3.html">post 3</a></h2><h2><a href="/posts/post4.html">post 4</a></h2><h2><a href="/posts/post5.html">post 5</a></h2><h2><a href="/posts/post6.html">post 6</a></h2><h2><a href="/posts/post7.html">post 7</a></h2><h2><a href="/posts/post8.html">post 8</a></h2><h2><a href="/posts/post9.html">post 9</a></h2><h2><a href="/posts/post10.html">post 10</a></h2><h2><a href="/posts/post11.html">post 11</a></h2><h2><a href="/posts/post12.html">post 12</a></h2><h2><a href="/posts/post13.html">post 13</a></h2><h2><a href="/posts/post14.html">post 14</a></h2><h2><a href="/posts/post15.html">post 15</a></h2><h2><a href="/posts/post16.html">post 16</a></h2><h2><a href="/posts/post17.html">post 17</a></h2><h2><a href="/posts/post18.html">post 18</a></h2><h2><a href="/posts/post19.html">post 19</a></h2><h2><a href="/posts/post20.html">post 20</a></h2><h2><a href="/posts/post21.html">post 21</a></h2><h2><a href="/posts/post22.html">post 22</a></h2><h2><a href="/posts/post23.html">post 23</a></h2><h2><a href="/posts/post24.html">post 24</a></h2></div></body></html>`,
		renderDynamic(mockDat),
	)
}

func Test_dynamic_template(_ *testing.T) {
	eq(
		`<!doctype html><html lang="en"><head><link rel="icon" href="data:;base64,="><title>Posts</title><meta property="og:title" content="Posts"><meta name="description" content="Random notes and thoughts"></head><body><nav><a href="/">home</a><a href="/works">works</a><a href="/posts" aria-current="page">posts</a><a href="/demos">demos</a><span>Updated Apr 05 3123</span></nav><div role="main"><h1>Posts</h1><h2><a href="/posts/post1.html">post 1</a></h2><h2><a href="/posts/post2.html">post 2</a></h2><h2><a href="/posts/post3.html">post 3</a></h2><h2><a href="/posts/post4.html">post 4</a></h2><h2><a href="/posts/post5.html">post 5</a></h2><h2><a href="/posts/post6.html">post 6</a></h2><h2><a href="/posts/post7.html">post 7</a></h2><h2><a href="/posts/post8.html">post 8</a></h2><h2><a href="/posts/post9.html">post 9</a></h2><h2><a href="/posts/post10.html">post 10</a></h2><h2><a href="/posts/post11.html">post 11</a></h2><h2><a href="/posts/post12.html">post 12</a></h2><h2><a href="/posts/post13.html">post 13</a></h2><h2><a href="/posts/post14.html">post 14</a></h2><h2><a href="/posts/post15.html">post 15</a></h2><h2><a href="/posts/post16.html">post 16</a></h2><h2><a href="/posts/post17.html">post 17</a></h2><h2><a href="/posts/post18.html">post 18</a></h2><h2><a href="/posts/post19.html">post 19</a></h2><h2><a href="/posts/post20.html">post 20</a></h2><h2><a href="/posts/post21.html">post 21</a></h2><h2><a href="/posts/post22.html">post 22</a></h2><h2><a href="/posts/post23.html">post 23</a></h2><h2><a href="/posts/post24.html">post 24</a></h2></div></body></html>`,
		templateToString(tplDynamic, mockDat),
	)
}

func Test_dynamic_equiv(_ *testing.T) {
	eqs(
		templateToString(tplDynamic, mockDat),
		renderDynamic(mockDat),
	)
}

func Benchmark_static_gax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = new(bytes.Buffer).Write(renderedStatic)
	}
}

func Benchmark_static_template(b *testing.B) {
	for i := 0; i < b.N; i++ {
		must(tplStatic.Execute(new(bytes.Buffer), nil))
	}
}

func Benchmark_dynamic_gax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = renderDynamic(mockDat)
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
		Post{Page: Page{Title: `post 4`, Path: `/posts/post4.html`}, IsListed: true},
		Post{Page: Page{Title: `post 5`, Path: `/posts/post5.html`}, IsListed: true},
		Post{Page: Page{Title: `post 6`, Path: `/posts/post6.html`}, IsListed: true},
		Post{Page: Page{Title: `post 7`, Path: `/posts/post7.html`}, IsListed: true},
		Post{Page: Page{Title: `post 8`, Path: `/posts/post8.html`}, IsListed: true},
		Post{Page: Page{Title: `post 9`, Path: `/posts/post9.html`}, IsListed: true},
		Post{Page: Page{Title: `post 10`, Path: `/posts/post10.html`}, IsListed: true},
		Post{Page: Page{Title: `post 11`, Path: `/posts/post11.html`}, IsListed: true},
		Post{Page: Page{Title: `post 12`, Path: `/posts/post12.html`}, IsListed: true},
		Post{Page: Page{Title: `post 13`, Path: `/posts/post13.html`}, IsListed: true},
		Post{Page: Page{Title: `post 14`, Path: `/posts/post14.html`}, IsListed: true},
		Post{Page: Page{Title: `post 15`, Path: `/posts/post15.html`}, IsListed: true},
		Post{Page: Page{Title: `post 16`, Path: `/posts/post16.html`}, IsListed: true},
		Post{Page: Page{Title: `post 17`, Path: `/posts/post17.html`}, IsListed: true},
		Post{Page: Page{Title: `post 18`, Path: `/posts/post18.html`}, IsListed: true},
		Post{Page: Page{Title: `post 19`, Path: `/posts/post19.html`}, IsListed: true},
		Post{Page: Page{Title: `post 20`, Path: `/posts/post20.html`}, IsListed: true},
		Post{Page: Page{Title: `post 21`, Path: `/posts/post21.html`}, IsListed: true},
		Post{Page: Page{Title: `post 22`, Path: `/posts/post22.html`}, IsListed: true},
		Post{Page: Page{Title: `post 23`, Path: `/posts/post23.html`}, IsListed: true},
		Post{Page: Page{Title: `post 24`, Path: `/posts/post24.html`}, IsListed: true},
	},
}

var renderedStatic = F(
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

func renderDynamic(dat MockDat) Bui {
	var scheme Attr
	if dat.ColorScheme != "" {
		scheme = Attr{`class`, dat.ColorScheme}
	}

	return F(
		Str(Doctype),
		E(`html`, AP(`lang`, `en`).A(scheme),
			E(`head`, nil, func(b *Bui) {
				b.E(`link`, AP(`rel`, `icon`, `href`, `data:;base64,=`))

				if dat.Title != "" {
					b.E(`title`, nil, dat.Title)
					b.E(`meta`, AP(`property`, `og:title`, `content`, dat.Title))
				} else {
					b.E(`title`, nil, `mock site`)
				}

				if dat.Desc != "" {
					b.E(`meta`, AP(`name`, `description`, `content`, dat.Desc))
				}

				if dat.Image != "" {
					b.E(`meta`, AP(`property`, `og:image`, `content`, `/images/`+dat.Image))
				}

				if dat.Type != "" {
					b.E(`meta`, AP(`property`, `og:type`, `content`, dat.Type))
					b.E(`meta`, AP(`property`, `og:site_name`, `content`, `mock site`))
				}
			}),

			E(`body`, nil,
				E(`nav`, nil,
					E(`a`, AP(`href`, `/`).A(cur(dat, `index.html`)), `home`),
					E(`a`, AP(`href`, `/works`).A(cur(dat, `works.html`)), `works`),
					E(`a`, AP(`href`, `/posts`).A(cur(dat, `posts.html`)), `posts`),
					E(`a`, AP(`href`, `/demos`).A(cur(dat, `demos.html`)), `demos`),
					E(`span`, nil, `Updated `+inst),
				),

				E(`div`, AP(`role`, `main`),
					E(`h1`, nil, `Posts`),

					func(b *Bui) {
						if len(dat.Posts) > 0 {
							for _, post := range dat.Posts {
								if !post.IsListed {
									continue
								}
								b.E(`h2`, nil, func() {
									b.E(`a`, AP(`href`, post.Path), post.Title)
								})
								if post.Desc != "" {
									b.E(`p`, nil, post.Desc)
								}
							}
						} else {
							b.E(`p`, nil, `Oops! It appears there are no public posts yet.`)
						}
					},
				),
			),
		),
	)
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
