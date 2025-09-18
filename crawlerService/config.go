package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	crw "github.com/RatNexus/CrawlerGoLib"
)

// TODO: change the default path to be more sane

type ServiceConfig struct {
	LogCfg  *LoggingConfig `json:"loggingConfig"`
	CrwCfg  *crw.Config    `json:"crawlerConfig"`
	toClose []io.Closer
}

func (cfg *ServiceConfig) String() string {
	data, err := json.Marshal(cfg)
	if err != nil {
		return ""
	}
	return string(data)
}

func (cfg *ServiceConfig) Close() error {
	for _, closer := range cfg.toClose {
		err := closer.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func MakeDefaultServiceConfig() (*ServiceConfig, error) {

	lc := MakeDefaultLogConfig()
	writer, file, err := lc.GetLogWriter()
	if err != nil {
		return nil, err
	}
	toClose := make([]io.Closer, 1)
	toClose[0] = file

	prefix := ""
	flag := log.Ldate | log.Ltime | log.Lshortfile
	logger := log.New(writer, prefix, flag)

	ccfg := &crw.Config{
		BaseURL:       "",
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

	sc := &ServiceConfig{
		LogCfg:  lc,
		CrwCfg:  ccfg,
		toClose: toClose,
	}
	return sc, nil
}

func GetConfigPath() (string, error) {
	key := "WENGINE_CRAWLER_CONFIG"
	defaultPath := "./crawler.conf"

	path := os.Getenv(key)
	if path == "" {
		path = defaultPath
		os.Setenv(key, path)
	}

	return path, nil
}

func (cfg *ServiceConfig) Save(path string) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0664)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig(path string) (*ServiceConfig, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if string(data) == "" {
		return nil, fmt.Errorf("config file is empty")
	}

	cfg := &ServiceConfig{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
