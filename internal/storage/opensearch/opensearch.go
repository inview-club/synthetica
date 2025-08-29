package opens

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/inview-club/synthetica/internal/template"
	"github.com/opensearch-project/opensearch-go"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Index       string      `json:"index"`
	Nodes       []string    `json:"nodes"`
	Credentials Credentials `json:"credentials"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Worker(id int, cfg Config, documentsAmount int, iterations int, templateText string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Infof("Worker %d starting for index %s\n", id, cfg.Index)

	cert, err := tls.LoadX509KeyPair("./elkcer.crt", "./elkcer.key")
	if err != nil {
		log.Errorf("failed to get certs: %s", err.Error())
		os.Exit(1)
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			},
		},
		Addresses: cfg.Nodes,
		Username:  cfg.Credentials.Username,
		Password:  cfg.Credentials.Password,
	})

	if err != nil {
		log.Errorf("failed to create client: %s", err.Error())
		os.Exit(1)
	}

	log.Info(client.Info())

	bufferPool := &sync.Pool{
		New: func() interface{} {
			data, err := template.RenderBodyBuffer(templateText)
			if err != nil {
				log.Errorf("Ошибка генерации буфера: %v", err)
				return nil
			}
			return data
		},
	}

	for i := 0; i < iterations; i++ {
		body := &bytes.Buffer{}

		for j := 0; j < documentsAmount; j++ {
			meta := `{"create":{"_index":"` + cfg.Index + `"}}`
			body.WriteString(meta)
			body.WriteString("\n")
			data := bufferPool.Get().(*[]byte)
			body.Write(*data)
			body.WriteString("\n")
		}
		fmt.Printf(body.String())
		log.Infof("Worker %d send bulk to index %s", id, cfg.Index)
		bulk, err := client.Bulk(body)
		if err != nil {
			log.Errorf("failed to bulk: %s", err.Error())
			continue
		}
		statusCode := bulk.StatusCode
		body.Reset()
		bulk.Body.Close()
		log.Infof("Worker %d get status %d for index %s", id, statusCode, cfg.Index)

		time.Sleep(1 * time.Second)
	}

	log.Infof("Worker %d stop work for index %s\n", id, cfg.Index)
}
