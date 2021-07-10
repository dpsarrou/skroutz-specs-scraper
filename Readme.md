# Skroutz Specs Scraper

A simple specs scraper for [skroutz.gr](https://www.skroutz.gr/) made in [golang](https://golang.org/)
that can crawl product pages and output the results in a different formats, such as CSV.

Originally it was created to help me choose the best aircondition unit to buy for
my apartment. However it is built modular enough with support of custom `Filters`,
`Store Parsers`so that it can crawl and parse other products and shops as well.

In addition it supports local caching of requests in the [cache](./cache) directory
to ensure that it does not burden the website with requests for products that have
already been crawled before. The cache at the time being does not have an expiry 
date, if you want to recrawl the website because you believe something has been
updated / changed in the products, just delete it and it will be recreated on the
next run.
## Usage

1.  Copy the URL of skroutz.gr that lists the products you want to crawl. For example
    https://www.skroutz.gr/c/407/oikiaka-klimatistika-inverter.html

2.  If you have [golang](https://golang.org/) already installed in your system you can
    simply run

        go run . https://www.skroutz.gr/c/407/oikiaka-klimatistika-inverter.html > data.csv

to crawl the products and export them into a csv file named `data.csv`.
In the future I might also compile OS specific binaries so that you don't have to install
golang.

### Attention: Skroutz Cookie

[skroutz.gr](https://www.skroutz.gr/) will attempt to block automated
software such as bots and crawlers from reading data from its website. It will do
so after a number of requests that the program has sent, and will present a page
with a CAPTCHA to be solved by a human. Naturally the crawler cannot solve this on its
own but there is a work around. If you open [skroutz.gr](https://www.skroutz.gr/)
in your browser, (and provided you are savvy enough to use the browser's Dev Tools),
copy the Cookie that [skroutz.gr](https://www.skroutz.gr/) has created and
paste it in the `skroutz_cookie` field at [config.yml](./config.yml).

## Configuration using config.yml

TODO

## Development

### Testing

To ensure the application runs as expected run the automated tests that have been created so far

    go test ./...

And the expected output will be something like:

    ?       skroutz-specs-scraper   [no test files]
    ?       skroutz-specs-scraper/crawler   [no test files]
    ?       skroutz-specs-scraper/exporter  [no test files]
    ok      skroutz-specs-scraper/filter    0.324s
    ok      skroutz-specs-scraper/parser    0.471s

### Adding more shop parsers

TODO

### Adding more filters

TODO
