// Package wf provides Alfred workflow functionality for searching MetaCPAN modules.
package wf

import (
	"encoding/json"
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
			agent: "libwww-perl/6.81",
			wrap:  http.DefaultTransport,
		},
	}
}

type modulesJSON struct {
	Items []moduleJSON `json:"items"`
}

type moduleJSON struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
}

// SearchModule searches for a distribution by query and returns results as JSON.
func SearchModule(q string) string {
	suggestions, err := metacpan.SearchAutocompleteSuggest(q)
	if err != nil {
		return errorToJSON(err)
	}

	result := modulesJSON{Items: []moduleJSON{}}

	for _, suggestion := range suggestions {
		result.Items = append(result.Items, moduleJSON{
			Arg:   suggestion.Name,
			Title: suggestion.Name,
			Subtitle: fmt.Sprintf(
				"%s/%s (%s)",
				suggestion.Author,
				suggestion.Release,
				suggestion.Date[0:10],
			),
		})
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return errorToJSON(err)
	}

	return string(jsonBytes)
}

// errorToJSON converts an error to an Alfred JSON response.
func errorToJSON(err error) string {
	result := modulesJSON{
		Items: []moduleJSON{
			{
				Title:    "ERROR",
				Subtitle: err.Error(),
				Arg:      "",
			},
		},
	}
	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}
