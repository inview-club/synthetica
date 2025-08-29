package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"uuid": func() string {
			return uuid.New().String()
		},
		"timestamp": func() int64 {
			return time.Now().UnixNano() / int64(time.Millisecond)
		},
		"date": func() string {
			return time.Now().Format(time.RFC3339Nano)
		},
		"oneOf": func(options ...string) string {
			if len(options) == 0 {
				return ""
			}
			return options[rand.Intn(len(options))]
		},
	}

}

func RenderBodyBuffer(templateText string) (*[]byte, error) {
	flattenedTemplate := strings.ReplaceAll(string(templateText), "\n", "")
	tmpl := template.Must(template.New("body").Funcs(templateFuncs()).Parse(flattenedTemplate))

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, nil); err != nil {
		return nil, fmt.Errorf("ошибка рендера: %v", err)
	}

	// Удаляем лишние пробелы и символы новой строки
	renderedStr := strings.TrimSpace(rendered.String())

	// Проверяем валидность JSON
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(renderedStr), &jsonObj); err != nil {
		return nil, fmt.Errorf("неверный JSON: %v", err)
	}

	// Сериализуем обратно в JSON (в компактной форме)
	compactJSON, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации JSON: %v", err)
	}

	return &compactJSON, nil
}
