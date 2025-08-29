# Synthetica

Synthetica is a Go-based tool for generating and bulk-inserting **synthetic log data** into databases and search engines.  

## 🗄️ Supported Databases

Currently supported:

- ![Opensearch](https://img.shields.io/badge/opensearch_2.19-005EB8.svg?style=for-the-badge&logo=OpenSearch&logoColor=white)

## 🚀 Features

- Customizable log templates (using Go `text/template`).
- Multi-worker parallel execution.
- Easily extensible to support new databases.

## ⚙️ Configuration

Synthetica is configured via a `config.json` file.  

Example:

```json
{
  "workers": 2,
  "documents_amount": 5,
  "iterations": 3,
  "body_template_file": "body.json.tmpl",
  "opensearch": {
    "index": "synthetic-logs",
    "nodes": ["https://localhost:9200"],
    "credentials": {
      "username": "admin",
      "password": "admin"
    }
  }
}
```

Global parameters:

- workers → Number of concurrent workers;
- documents_amount → Number of documents per request;
- iterations → Number of requests per worker;
- body_template_file → Path to the JSON log template.

Database-specific parameters:

For OpenSearch:

- opensearch.index → Target index;
- opensearch.nodes → List of OpenSearch nodes
- opensearch.credentials → Username & password.

## 📝 Templating

Synthetica uses Go’s text/template engine with custom helpers.

Available functions:

- uuid → Generates a random UUID.
- timestamp → Current timestamp in milliseconds.
- date → Current timestamp in RFC3339Nano.
- oneOf "a" "b" "c" → Picks a random option from the provided list.

### Example

```
{
  "timestamp": "{{date}}",
  "level": "{{oneOf "INFO" "WARN" "ERROR" "DEBUG"}}",
  "service": "{{oneOf "auth-service" "payment-service" "user-service" "order-service"}}",
  "message": "{{oneOf "User login successful" "User login failed" "Payment processed" "Payment declined" "Order created" "Order cancelled" "Database connection error" "Cache miss detected"}}",
  "host": "{{oneOf "server-1" "server-2" "server-3" "server-4"}}",
  "trace_id": "{{uuid}}"
}
```
