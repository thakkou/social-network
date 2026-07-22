package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"01social/pkg/utilities"
)

// GetSwaggerJSON serves the generated swagger.json file from disk.
func GetSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	// Look for swagger.json relative to the working directory (backend/)
	paths := []string{
		filepath.Join("docs", "swagger", "swagger.json"),
		filepath.Join("..", "docs", "swagger", "swagger.json"),
	}

	var data []byte
	var err error
	for _, p := range paths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}

	if err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "swagger.json not found", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

// GetSwaggerUI serves the Swagger UI HTML page that loads swagger.json from the API.
func GetSwaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Social Network API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
        body { margin: 0; padding: 0; background: #f5f5f5; }
        .topbar { display: none !important; }
        .swagger-ui .info { margin: 20px 0; }
        .swagger-ui .info .title { font-size: 28px; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
    <script>
        SwaggerUIBundle({
            url: "/api/docs/swagger.json",
            dom_id: "#swagger-ui",
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.SwaggerUIStandalonePreset
            ],
            layout: "BaseLayout",
            docExpansion: "list",
            defaultModelsExpandDepth: 2,
            filter: true,
            tagsSorter: "alpha",
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
