#!/bin/zsh
# short script to gist-ize the output of href_extracter.
# relies on the presence of the ruby gem "gist" https://github.com/defunkt/gist

# set constants
# file containing URLS
OUTPUT="urls.txt"
# json file
JSON="collection.json"

# process options
# we need to provide "name", "url" and "description", as well as selector
# so we'll use -n, -u, -d,-s, respectively

while getopts ":u:s:n:d:h" opt; do
 case $opt in
   h)
     echo "USAGE: href_gister -u URL [-s SELECTOR ] -d NAME -d \"DESCRIPTION\"
     -u, -n and -d are mandatory. -s defaults to 'a'"
     ;;
   u)
     url=$OPTARG
     ;;
   s)
     selector=$OPTARG
     ;;
   n)
     name=$OPTARG
     ;;
   d)
     description=$OPTARG
     ;;
   
   \?)
     echo "Invalid option: -$OPTARG" >&2
     exit 1
     ;;
   :)
     echo "Option -$OPTARG requires an argument." >&2
     exit 1
     ;;
 esac
done

echo "-u $url -o $OUTPUT -s $selector -n $name -d $description"


if [[ ! ($url || $name ) ]]; then
  echo "-u URL, -n NAME and -d DESCRIPTION are all required."
  exit 1
fi

if [[ ! selector ]]; then
  selector=a
fi

if [[ ! description ]]; then
  description=""
fi

# extract the urls
echo "extract_href -u $url -o $OUTPUT -s $selector  "
extract_href -u $url -o $OUTPUT -s $selector  

# write the config file
echo '{ "name" : "'$name'",
"url" : "'$url'",
"description" : "'$description'"
}' > $JSON
  
  
# call gist
mygist=$(gist collection.json urls.txt -d "extract_href output from $url, \"$name\"")
echo "$mygist,$url" >> gists.txt
echo $mygist
