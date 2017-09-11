package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewTestServer allocates a testing server that serves up content
// to try to extract links from
func NewTestServer() *httptest.Server {
	m := http.NewServeMux()
	m.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {
		body := []byte(`
<html>
  <body>
    <a>This shouldn't work</a>
    <a href="#bad"></a>
    <a href="/rel-endpoint"></a>
    <a href="http://youtube.com/external-link"></a>
  </body>
</html>`)
		w.Write(body)
	})

	return httptest.NewServer(m)
}

func TestFetchAndWriteHrefAttrs(t *testing.T) {
	s := NewTestServer()

	cases := []struct {
		endpoint   string
		selector   string
		fetchError error
		stats      *Stats
	}{
		{"/1", "a", nil, &Stats{Elements: 4, WithHref: 3, ValidUrl: 2, Duplicates: 1}},
	}

	for i, c := range cases {
		w := &bytes.Buffer{}

		root := fmt.Sprintf("%s%s", s.URL, c.endpoint)
		stats, err := FetchAndWriteHrefAttrs(root, "a", w)
		if err != c.fetchError {
			t.Errorf("case %d fetch error mismatch: %s != %s", err, c.fetchError)
		}

		if err := CompareStats(c.stats, stats); err != nil {
			t.Errorf("case %d stats error: %s", i, err.Error())
			continue
		}
	}
}

func CompareStats(a, b *Stats) error {
	if a.Elements != b.Elements {
		return fmt.Errorf("element count mismatch: %d != %d", a.Elements, b.Elements)
	}
	if a.WithHref != b.WithHref {
		return fmt.Errorf("WithHref count mismatch: %d != %d", a.WithHref, b.WithHref)
	}
	if a.Duplicates != b.Duplicates {
		return fmt.Errorf("Duplicate count mismatch: %d != %d", a.Duplicates, b.Duplicates)
	}
	if a.ValidUrl != b.ValidUrl {
		return fmt.Errorf("valid url count mismatch: %d != %d", a.ValidUrl, b.ValidUrl)
	}
	return nil
}
