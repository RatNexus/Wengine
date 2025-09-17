package main

import (
	"fmt"
	crw "github.com/RatNexus/CrawlerGoLib"
	"log"
	"os"
	"sync"
)

// TODO: add config save and load functionality for crw and main packages
// TODO: add cli handling
// TODO: add database handling
// TODO: get rid of placeholders here
// TODO: make sure the logger flags are sane

func logReport(pages map[string]string, baseURL string, logger *log.Logger) {
	logger.Println("=============================")
	logger.Printf("  REPORT for %s\n", baseURL)
	logger.Println("=============================")

	for url, n := range pages {
		pages[url] = n
	}

	for url := range pages {
		logger.Printf("Found %s", url)
	}
}

func main() {
	pages := make(map[string]string)
	mu := sync.Mutex{}

	storePage := func(url string, html string, siteLinks []string, imgLinks []string) error {
		mu.Lock()

		pages[url] = html

		mu.Unlock()
		return nil
	}

	isPageStored := func(url string) (bool, error) {
		mu.Lock()
		_, exists := pages[url]
		mu.Unlock()
		return exists, nil
	}

	lc := MakeDefaultLogConfig()
	writer, file, err := lc.GetLogWriter()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if file != nil {
		defer file.Close()
	}

	prefix := ""
	flag := log.Ldate | log.Ltime | log.Lshortfile
	logger := log.New(writer, prefix, flag)

	page := "https://crawler-test.com/"
	ccfg := crw.Config{
		BaseURL:       page,
		MaxDepth:      1,
		MaxPages:      10,
		MaxGoroutines: 10, // Todo: make 10 the default in crawler lib

		LoggingOptions: &crw.LoggingOptions{
			DoLogging:    true,
			DoStart:      true,
			DoEnd:        true,
			DoPageAbyss:  false,
			DoDepthAbyss: false,
			DoDepth:      true,
			DoWidth:      true,
			DoErrors:     true,
			DoPages:      true,
		},
		Logger: logger,
	}

	crawler, err := ccfg.MakeCrawler(storePage, isPageStored)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	crawler.CrawlPage(page)

	logReport(pages, page, logger)
}
