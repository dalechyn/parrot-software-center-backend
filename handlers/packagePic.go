package handlers

import (
	"net/http"
)

func AssetsServer(assetsDir string) http.Handler {
	return http.FileServer(http.Dir(assetsDir))
}
