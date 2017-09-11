// extract_href is a command line tool for extracting urls from a HTML web page.
// writing each url on a new line.
// Each matched url is:
//  * absolute (referenced by source url)
//  * unique - (no duplicates are added to the list)
//  * refers to a separate resource - (no url fragments)
// it uses a jquery-style selector to search the HTML document for elements that have an href attribute
// to construct a de-duplicated list of href attributes
// It has three major command line options:
//  * u - (required) url of html document
//  * o - path to output file
//  * s - (required, default: "a") - jquery-style selector to match html elements
// example use:
// ```
// 		./extract_href -u https://www.epa.gov/endangered-species/biological-evaluation-chapters-chlorpyrifos-esa-assessment -s '.main-column.clearfix a'
// ```
// this will fetch the epa.gov url, select all "a" tags in the document that are a decendant of any element with the classes "main-column" and "clearfix"
// and build a deduplicated list of absolute urls using the `href` attribute of all found anchor tags. run that same command adding `-o urls.txt` to save
// the results to a file and see output stats instead.
// Picking the right jquery selector is a bit of an art, the goal is to isolate the most general part of the page that contains all of the links
// that you're after. For more information on jquery selectors and how they work, have a look here: https://learn.jquery.com/using-jquery-core/selecting-elements/
// When in doubt, it's often fine to leave the default "a" selector, which will generate lots of links you may not want,
// and manually remove them from the output file.
package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	showHelp bool
	outFile  string
	rootUrl  string
	selector string
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "print help text")
	flag.StringVar(&outFile, "o", "", "path to write file to")
	flag.StringVar(&rootUrl, "u", "", "url to fetch links from")
	flag.StringVar(&selector, "s", "a", "jquery-style selector to scope url search to, default is 'a'")
}

func main() {
	// parse flags, grabbing values from the command line
	flag.Parse()

	if len(os.Args) == 1 || showHelp {
		PrintHelpText()
		return
	}

	// allocate a new results writer
	w, err := NewResultsWriter(outFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stats, err := FetchAndWriteHrefAttrs(rootUrl, selector, w)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// if stdout isn't being used for output, write stats to stdout
	if w != os.Stderr {
		fmt.Println(stats)
	}

	// check to see if our writer implements the closer interface,
	// call close if so
	if closer, ok := w.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

// Stats tracks state for extracting hrefs from a given document
type Stats struct {
	// elements matched by selector
	Elements int
	// elements that have an "href" attribute
	WithHref int
	// elements who's absolute url was a duplicate
	Duplicates int
	// elements with a valid url
	ValidUrl int
}

// stats implements the stringer interface
func (s *Stats) String() string {
	return fmt.Sprintf("%d matched HTML elements\n%d had a href attribute.\n%d were duplicates\n%d were valid\n", s.Elements, s.WithHref, s.Duplicates, s.ValidUrl)
}

// NewResultsWriter writes to either a file or stderr if no path is provided
func NewResultsWriter(path string) (io.Writer, error) {
	if path != "" {
		return os.Create(path)
	}
	return os.Stderr, nil
}

// FetchAndWriteHrefAttrs fetches a given url, and uses the provided jquery-style selector to grab
// all of the "href" attributes for a given url HTML document, writing a line-delimited list of
// deduplicated absolute urls to w
func FetchAndWriteHrefAttrs(rootUrl, selector string, w io.Writer) (*Stats, error) {
	// check for required params
	if rootUrl == "" {
		return nil, fmt.Errorf("url is required")
	}

	root, err := url.Parse(rootUrl)
	if err != nil {
		return nil, err
	}

	doc, err := fetchGoqueryUrl(rootUrl)
	if err != nil {
		return nil, err
	}

	// find the selected HTML elements
	elements := doc.Find(selector)
	// create a stats object with the total number of matched element
	stats := &Stats{Elements: elements.Length()}
	// added is a list of urls that have been added already
	added := map[string]bool{}

	// iterate through elements
	for i := range elements.Nodes {
		el := elements.Eq(i)

		if href, exists := el.Attr("href"); exists {
			stats.WithHref++

			// Reslove any relative url references by parsing the href in relation
			// to the root url
			abs, err := root.Parse(href)
			if err != nil {
				// TODO - handle error here
				continue
			}

			// remove fragement from url element, this will make
			abs.Fragment = ""

			absStr := abs.String()

			if absStr != rootUrl && added[absStr] == false {
				added[absStr] = true
				stats.ValidUrl++
				// write the url as a line to the writer
				w.Write([]byte(fmt.Sprintf("%s\n", abs.String())))
			} else {
				stats.Duplicates++
			}
		}
	}

	return stats, nil
}

// fetchGoqueryUrl performs a GET to the passed in url, returning a parsed goquery document
func fetchGoqueryUrl(urlstr string) (*goquery.Document, error) {
	resp, err := http.Get(urlstr)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(resp.Body)
}

// PrintHelpText outputs instructions for using this program to os.Stdout
func PrintHelpText() {
	fmt.Println(`
extract_href is a command line tool for extracting urls from a HTML web page.
writing each url on a new line.
Each matched url is:
 * absolute (referenced by source url)
 * unique - (no duplicates are added to the list)
 * refers to a separate resource - (no url fragments)

extract_url uses a jquery-style selector to search the HTML document for elements that have an href attribute
to construct a de-duplicated list of href attributes.

options:`)
	flag.PrintDefaults()
}
