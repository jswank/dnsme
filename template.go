package main

import (
	"io"
	"strings"
	"text/template"
)

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

var usageTemplate = `usage: dnsme command [arguments]

'dnsme' is a command-line interface to the DNS Made Easy REST API.  In
order to use it, an API key pair must be configured for your account:
this can be trivially created by visiting your account information page
at https://cp.dnsmadeeasy.com/account/info.

The following environment variables should be set:

    DNSME_API_URL = http://api.dnsmadeeasy.com/V1.2
    DNSME_API_KEY = API key
    DNSME_SECRET_KEY = Secret key

Available commands are:
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-16s"}} {{.Short}}{{end}}{{end}}

Use "dnsme help [command]" for more information about a command.

Global flags:

The -d flag can be used to print raw HTTP requests and responses to
stderr.

The flag "-o" specifies the output type.  Available output types are
"csv", "json", or the default text-based "std".

`

var helpTemplate = `{{if .Runnable}}usage: dnsme {{.UsageLine}}

{{end}}{{.Long | trim}}

`

var domainListTemplate = `{{range .List}}{{printf "%s\n" .}}{{end}}`

var domainInfoTemplate = `{{range .NameServers}}{{printf "Nameserver: %s\n" .}}{{end}}{{range .VanityNameServers}}{{printf "Vanity NS: %s\n" .}}{{end}}GTD Enabled: {{.GtdEnabled}}
`
var domainInfoTemplateCSV = `{{.Name}},{{range .NameServers}}{{.}} {{end}},{{range .VanityNameServers}}{{.}} {{end}},{{.GtdEnabled}}
`

var secondaryTemplate = `{{range .IP}}{{printf "IP: %s\n" .}}{{end}}`
var secondaryTemplateCSV = `{{.Name}},{{range .IP}}{{.}} {{end}}`

// looks like zone file entries
var recordTemplate = `{{if .Name}}{{printf "%-20s" .Name}}{{else}}{{printf "%-20s" "@"}}{{end}} {{printf "%-6d" .TTL}} {{printf "%-5s" .Type}} {{if .Data}}{{printf "%-32s" .Data}}{{else}}{{printf "%-32s" "@"}}{{end}} ; id={{printf "%-7d" .ID}}, gtd={{.GtdLocation}}
`
