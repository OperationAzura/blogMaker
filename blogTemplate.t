{{define "blogTemplate.t"}}
{
  "title": "{{.Title}}",
  "description": "{{.Description}}",
  "image": "{{.Image}}",
  "video": "{{.Video}}",
  "date": "{{.Date}}",
  "tags": [{{range $Index, $Tag := .TagSlice}}"{{$Tag}}"{{ if gt (len $.TagSlice) (NextIndex $Index)}},{{end}}{{end}}],
  "categories": [{{range $Index, $Category := .CatSlice}}"{{$Category}}"{{ if gt (len $.CatSlice) (NextIndex $Index) }},{{end}}{{end}}],
  "draft": "{{.Draft}}"
}
{{if .Video}}<iframe width="560" height="315" src={{.Video}} frameborder="0" allowfullscreen></iframe>{{end}}

{{.Body}}
{{end}}
