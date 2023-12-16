package http

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cameront/go-svelte-sqlite-template/log"
)

// handleStatic serves the static assets (i.e. the entire frontend app)
func InitStatic(ctx context.Context, staticPath string) http.HandlerFunc {
	logger := log.With(nil, "initiator", "static")

	uiPath := http.Dir(staticPath)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		cleaned := filepath.Clean(r.URL.Path)
		trimmed := strings.TrimPrefix(cleaned, "/")
		path := resolvePath(uiPath, trimmed)

		file, err := uiPath.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				logger.Info(fmt.Sprintf("file %s not found: %s", path, err))
				http.NotFound(w, r)
				return
			}
			logger.Info(fmt.Sprintf("file %s cannot be read: %s", path, err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		w.Header().Set("Content-Type", contentType)
		if strings.HasPrefix(path, "static/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}
		stat, err := file.Stat()
		if err == nil && stat.Size() > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		}

		n, _ := io.Copy(w, file)
		logger.Info(fmt.Sprintf("file %s copied %d bytes", path, n))
	}
}

func resolvePath(uiFS http.Dir, path string) string {
	if path == "" || path == "/" {
		return "index.html"
	}

	// anything that looks like a file w/extension
	if strings.Contains(path, ".") {
		return path
	}

	// route anything else to index.html
	_, err := uiFS.Open(path)
	if err != nil && os.IsNotExist(err) {
		return "index.html"
	}

	// else, just return what we got
	return path
}
