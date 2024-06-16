//+build dev
//go:build dev
// +build dev

package main

import (
	"fmt"
	"net/http"
	"os"
)

// This allows tailwind to rebuild the css without rebuilding the entire binary

func public() http.Handler {
	fmt.Println("building static files for development")
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public")))
}
