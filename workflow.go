// Package wf provides Alfred workflow functionality for searching MetaCPAN modules.
package wf

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/oalders/go-metacpan"
)

type userAgentTransport struct {
	agent string
	wrap  http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("User-Agent", t.agent)
	return t.wrap.RoundTrip(r)
}

func init() {
	http.DefaultClient = &http.Client{
		Transport: &userAgentTransport{
			agent: "libwww-perl/6.72",
			wrap:  http.DefaultTransport,
		},
	}
}

// The ModulesXML defines XML struct of list of distributions.
type ModulesXML struct {
	XMLName xml.Name    `xml:"items"`
	Item    []ModuleXML `xml:"item"`
}

// The ModuleXML defines XML struct of list of distributions.
type ModuleXML struct {
	XMLName  xml.Name `xml:"item"`
	Arg      string   `xml:"arg,attr"`
	Title    string   `xml:"title"`
	Subtitle string   `xml:"subtitle"`
}

// SearchModule returns search distribution by query(q) and returns results as XML.
func SearchModule(q string) string {
	suggestions, err := metacpan.SearchAutocompleteSuggest(q)
	if err != nil {
		return errorToXML(err)
	}

	xmlType := ModulesXML{
		XMLName: xml.Name{},
		Item:    []ModuleXML{},
	}

	for _, suggestion := range suggestions {
		xmlType.Item = append(xmlType.Item, ModuleXML{
			XMLName: xml.Name{},
			Arg:     suggestion.Name,
			Title:   suggestion.Name,
			Subtitle: fmt.Sprintf(
				"%s/%s (%s)",
				suggestion.Author,
				suggestion.Release,
				suggestion.Date[0:10],
			),
		})
	}

	xmlBytes, err := xml.Marshal(xmlType)
	if err != nil {
		return errorToXML(err)
	}

	return xml.Header + string(xmlBytes)
}

// errorToXML convert error to XML.
func errorToXML(_ error) string {
	return xml.Header + `
<items>
  <item arg="">
    <title>ERROR</title>
    <subtitle>Something wrong</subtitle>
  </item>
</items>`
}
