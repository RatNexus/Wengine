package main

import (
	"fmt"
	crw "github.com/RatNexus/CrawlerGoLib"
	"log"
	"os"
	"sync"
)

// TODO: add a build script
// TODO: add config save and load functionality for crw and main packages
// TODO: add cli handling
// TODO: add database handling
// TODO: get rid of placeholders here

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
	var lc *LoggingConfig

	lc = &LoggingConfig{}
	lc.LoggingOptions = crw.LoggingOptions{
		DoLogging:    true,
		DoStart:      true,
		DoEnd:        true,
		DoPageAbyss:  false,
		DoDepthAbyss: false,
		DoDepth:      true,
		DoWidth:      true,
		DoErrors:     true,
		DoPages:      true,
	}
	lc.SetDefultLogFileOptions()
	lc.LogToFile = false
	lc.LogToScreen = true

	prefix := ""
	flag := log.Ldate | log.Ltime | log.Lshortfile
	writer, err := lc.GetLogWriter()
	if err != nil {
		fmt.Print(err)
	}
	logger := log.New(writer, prefix, flag)

	// Why not use lo directly? So it can be saved and retrieved as JSON more easily.
	// This is a setup for an abstraction to be added later.
	page := "https://crawler-test.com/"
	cfg, err := crw.ConfigSetup(1, 10, 10, page, logger, &lc.LoggingOptions,
		storePage, isPageStored)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	cfg.CrawlPage(page)

	logReport(pages, page, logger)
}
