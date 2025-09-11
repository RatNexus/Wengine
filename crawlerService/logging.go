package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	crw "github.com/RatNexus/CrawlerGoLib"
)

type LoggingConfig struct {
	crw.LoggingOptions `json:"crawlerLoggingOptions"`

	LogsFolder string `json:"logsFolder"`
	LogName    string `json:"logName"`
	DateSuffix string `json:"dateSuffix"`

	LogToFile   bool `json:"logToFile"`
	LogToScreen bool `json:"logToScreen"`
}

func (lc *LoggingConfig) SetDefultLogFileOptions() {
	lc.LogsFolder = "/tmp/protoCrawler/logs"
	lc.LogName = "crawler"
	lc.DateSuffix = "2006-01-02_15:04:05"
}

func (lc *LoggingConfig) GetLogWriter() (io.Writer, error) {
	if lc.LogToScreen && lc.LogToFile {
		return io.Discard, nil
	}

	if lc.LogToScreen && !lc.LogToFile {
		return os.Stdout, nil
	}

	if !lc.LogToScreen && lc.LogToFile {
		file, err := lc.GetLogFile()
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return io.Discard, nil
}

func (lc *LoggingConfig) GetLogFile() (Lf *os.File, err error) {
	var LogFileName string
	if lc.DateSuffix != "" {
		CurrentTime := time.Now()
		Tf := CurrentTime.Format(lc.DateSuffix)

		LogFileName = fmt.Sprintf("%s_%s.log", lc.LogName, Tf)
	} else {
		LogFileName = fmt.Sprintf("%s.log", lc.LogName)
	}

	lc.LogsFolder = strings.TrimRight(lc.LogsFolder, "/")
	err = os.MkdirAll(lc.LogsFolder, 0755)
	if err != nil {
		return nil, err
	}

	LogPath := fmt.Sprintf("%s/%s", lc.LogsFolder, LogFileName)
	LogFile, err := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return LogFile, nil
}
