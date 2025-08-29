package main

import (
	"encoding/json"
	"os"
	"sync"

	opens "github.com/inview-club/synthetica/internal/storage/opensearch"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	WorkersAmount    int          `json:"workers"`
	DocumentsAmount  int          `json:"documents_amount"`
	Iterations       int          `json:"iterations"`
	BodyTemplateFile string       `json:"body_template_file"`
	Opensearch       opens.Config `json:"opensearch"`
}

func main() {
	f, err := os.ReadFile("config.json")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	var cfg Config
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	templateContent, err := os.ReadFile(cfg.BodyTemplateFile)
	if err != nil {
		log.Errorf("Failed to read template: %s. Error: %v", cfg.BodyTemplateFile, err)
		os.Exit(1)
	}

	for i := 0; i < cfg.WorkersAmount; i++ {
		wg.Add(1)
		go opens.Worker(i, cfg.Opensearch, cfg.DocumentsAmount, cfg.Iterations, string(templateContent), &wg)
	}

	wg.Wait()
}
