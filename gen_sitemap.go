package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

func isFullURL(uri string) bool {
	return strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://")
}

func addSitemapURL(b *Book, uri string) {
	if !isFullURL(uri) {
		uri = urlJoin(b.BaseURL(), uri)
	}
	b.muSitemapURLS.Lock()
	b.sitemapURLS[uri] = struct{}{}
	b.muSitemapURLS.Unlock()
}

const (
	sitemapTmpl = `User-agent: *
Disallow:

Sitemap: %s
`
)

func writeSitemap(b *Book) {
	// http://www.advancedhtml.co.uk/robots-sitemaps.htm
	sitemapURL := urlJoin(b.BaseURL(), "sitemap.txt")
	robotsTxt := fmt.Sprintf(sitemapTmpl, sitemapURL)
	robotsTxtPath := filepath.Join(b.DirWWW, "robots.txt")
	err := ioutil.WriteFile(robotsTxtPath, []byte(robotsTxt), 0644)
	must(err)

	addSitemapURL(b, "/")
	//addSitemapURL(b, "about")

	var urls []string
	for uri := range b.sitemapURLS {
		urls = append(urls, uri)
	}
	sort.Strings(urls)
	s := strings.Join(urls, "\n")
	sitemapPath := filepath.Join(b.DirWWW, "sitemap.txt")
	err = ioutil.WriteFile(sitemapPath, []byte(s), 0644)
	must(err)
}
