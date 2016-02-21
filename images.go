package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var searchImageTemplate = template.Must(template.New("searchImage").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<title>Dispel - Image Database</title>
		<link rel="stylesheet" href="/static/css/milligram.min.css">
		<link rel="stylesheet" href="/static/css/images.css">
		<script src="/static/js/images.js"></script>
	</head>
	<body>
		<header>
			Header
		</header>
		<div style="margin: 0 1.5% 24px 1.5%;">
			<input id="searchbar" type="search" placeholder="yeb guac" value="{{ .Search }}" />
		</div>
		<div class="imagelist">
			{{ if .Images }}
				{{ range .Images }}
					<a href="/images/show/{{.}}">
						<span class="thumb">
							<img class="preview" src="/static/images/{{.}}" />
						</span>
					</a>
				{{ end }}
			{{ else }}
				<span>No results!</span><br/><br/>
			{{ end }}
		</div>
		<footer>
			Footer
		</footer>
	</body>
</html>
`))

var showImageTemplate = template.Must(template.New("showImage").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<title>Dispel - {{.}}</title>
		<link rel="stylesheet" href="/static/css/milligram.min.css">
		<link rel="stylesheet" href="/static/css/images.css">
	</head>
	<body>
		<header>
			Header
		</header>
		<div class="flex">
			<div class="sidebar">
				Sidebar
			</div>
			<div class="content">
				<img style="max-width: 100%;" src="/static/images/{{.}}" />
			</div>
		</div>
		<footer>
			Footer
		</footer>
	</body>
</html>
`))

func parseTags(tagQuery string) (include, exclude []string) {
	for _, tag := range strings.Split(tagQuery, " ") {
		if strings.TrimPrefix(tag, "-") == "" {
			continue
		}
		if strings.HasPrefix(tag, "-") {
			exclude = append(exclude, tag[1:])
		} else {
			include = append(include, tag)
		}
	}
	return
}

// imageSearchHandler is the handler for the /images route. If
func (db *imageDB) imageSearchHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	searchTags := req.FormValue("t")
	urls, err := db.lookupByTags(parseTags(searchTags))
	if err != nil {
		http.Error(w, "Lookup failed", http.StatusInternalServerError)
		return
	}
	searchImageTemplate.Execute(w, struct {
		Search string
		Images []string
	}{searchTags, urls})
}

func imageShowHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// TODO: look up img
	showImageTemplate.Execute(w, ps.ByName("img"))
}
