package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

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

	var cfg *ServiceConfig

	path, err := GetConfigPath()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	cfg, err = LoadConfig(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg, err = MakeDefaultServiceConfig()
		fmt.Println("No config at the specified path.")
		fmt.Println("Using Defaults.")
	}
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer cfg.Close()

	writer, file, err := cfg.LogCfg.GetLogWriter()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	cfg.toClose = append(cfg.toClose, file)

	cfg.CrwCfg.Logger = log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile)

	page := "https://crawler-test.com/"
	cfg.CrwCfg.BaseURL = page
	crawler, err := cfg.CrwCfg.MakeCrawler(storePage, isPageStored)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	crawler.CrawlPage(page)

	logReport(pages, page, cfg.CrwCfg.Logger)

	fmt.Println(cfg)
	err = cfg.Save(path)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
