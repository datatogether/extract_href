# extract_href

<!-- Repo Badges for: Github Project, Slack, License-->

[![GitHub](https://img.shields.io/badge/project-Data_Together-487b57.svg?style=flat-square)](http://github.com/datatogether)
[![Slack](https://img.shields.io/badge/slack-Archivers-b44e88.svg?style=flat-square)](https://archivers-slack.herokuapp.com/)
[![License](https://img.shields.io/github/license/datatogether/extract_href.svg?style=flat-square)](./LICENSE) 

extract_href is a command line tool for extracting urls from a HTML web page,
writing each url on a new line. Each matched url is:

* absolute (referenced by source url)
* unique - (no duplicates are added to the list)
* refers to a separate resource - (no url fragments)

It uses a jquery-style selector to search the HTML document for elements that
have an href attribute to construct a de-duplicated list of href attributes It
has three major command line options:

* **u** - (required) url of html document
* **o** - path to output file
* **s** - (required, default: "a") - jquery-style selector to match html elements

## License & Copyright

Copyright (C) 2017 Data Together  
This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, version 3.0.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.

See the [`LICENSE`](./LICENSE) file for details.

## Getting Involved

We would love involvement from more people! If you notice any errors or would 
like to submit changes, please see our [Contributing Guidelines](./.github/CONTRIBUTING.md). 

We use GitHub issues for [tracking bugs and feature requests](https://github.com/datatogether/extract_href/issues) 
and Pull Requests (PRs) for [submitting changes](https://github.com/datatogether/extract_href/pulls)

### Installation

You will need [go](http://golang.org) installed to build `extract_href` from 
source, then run:

```bash
go get -u github.com/datatogether/extract_href
```

test that it's working by running a raw `extract_href`, which should output help text.

### Usage

```bash
extract_href -u https://www.epa.gov/endangered-species/biological-evaluation-chapters-chlorpyrifos-esa-assessment -s '.main-column.clearfix a'
``` 
This will fetch the epa.gov url, select all "a" tags in the document that
are a decendant of any element with the classes "main-column" and "clearfix" and
build a deduplicated list of absolute urls using the `href` attribute of all
found anchor tags. Run that same command adding `-o urls.txt` to save the
results to a file and see output stats instead. Picking the right jQuery
selector is a bit of an art. The goal is to isolate the most general part of the
page that contains all of the links that you're after. For more information on
jquery selectors and how they work, have a look [at the jQuery selection docs](https://learn.jquery.com/using-jquery-core/selecting-elements/). When in doubt,
it's often fine to leave the default "a" selector, which will generate lots of
links you may not want, and manually remove those from the output file.

## Development

Coming soon!
