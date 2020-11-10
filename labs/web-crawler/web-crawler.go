// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.

type Site struct{
	depth int
	link string
}

var tokens = make(chan struct{}, 20)

var depth = flag.Int("depth", 1, "crawler depth")
var name = flag.String("results","results.txt", "Result file")


func crawl(url Site) []Site {
	f, err := os.OpenFile(*name, os.O_APPEND| os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print(err)
	}

	_, err = fmt.Fprintln(f, url.link)
	if err != nil {
		log.Print(err)
		f.Close()
	}

	err = f.Close()
	if err != nil {
		fmt.Print(err)
	}

	if url.depth < *depth{

		tokens <- struct{}{} // acquire a token
		list, err := links.Extract(url.link)
		if err != nil{
			fmt.Println(err)
		}
		sites := make([]Site,0)
		for _, link := range list{
			sites = append(sites, Site{link: link,depth: url.depth +1})
		}
		<-tokens // release the token
		
		return sites

	}

	

	
	return []Site{}
}

//!-sema

//!+
func main() {
	worklist := make(chan []Site)
	var n int // number of pending sends to worklist

	flag.Parse()
	

	// Start with the command-line arguments.
	if len(os.Args) != 4{
		fmt.Println("Sugested usage: ./web-crawler -depth=3 -results=results.txt https://google.com")
		os.Exit(1)
	}

	n++
	file, err := os.Create(*name)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	err = file.Close()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}



	go func() {
		urls := make([]Site, 0)
		urls = append(urls, Site{link: (os.Args[3:])[0],depth: 0})
		worklist <- urls
	}()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link.link] {
				seen[link.link] = true
				n++
				go func(link Site) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

//!-