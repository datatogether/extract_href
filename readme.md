# extract_href
--
extract_href is a command line tool for extracting urls from a HTML web page.
writing each url on a new line. Each matched url is:

    * absolute (referenced by source url)
    * unique - (no duplicates are added to the list)
    * refers to a separate resource - (no url fragments)

it uses a jquery-style selector to search the HTML document for elements that
have an href attribute to construct a de-duplicated list of href attributes It
has three major command line options:

    * u - (required) url of html document
    * o - path to output file
    * s - (required, default: "a") - jquery-style selector to match html elements

example use: ```

    ./extract_href -u https://www.epa.gov/endangered-species/biological-evaluation-chapters-chlorpyrifos-esa-assessment -s '.main-column.clearfix a'

``` this will fetch the epa.gov url, select all "a" tags in the document that
are a decendant of any element with the classes "main-column" and "clearfix" and
build a deduplicated list of absolute urls using the `href` attribute of all
found anchor tags. run that same command adding `-o urls.txt` to save the
results to a file and see output stats instead. Picking the right jquery
selector is a bit of an art, the goal is to isolate the most general part of the
page that contains all of the links that you're after. For more information on
jquery selectors and how they work, have a look here:
https://learn.jquery.com/using-jquery-core/selecting-elements/ When in doubt,
it's often fine to leave the default "a" selector, which will generate lots of
links you may not want, and manually remove them from the output file.
