{{define "blogTemplate.t"}}
{
  "title": "{{.Title}}",
  "description": "{{.Description}}",
  "image": "{{.Image}}",
  "video": "{{.Video}}",
  "date": "{{.Date}}",
  "tags": "{{.Tags}}",
  "categories": "{{.Categories}}",
  "draft": {{.Draft}}
}
{{if .Video}}<iframe width="560" height="315" src="{{.Video}}" frameborder="0" allowfullscreen></iframe>{{end}}

{{.Body}}
{{end}}
