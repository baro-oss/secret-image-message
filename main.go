package main

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"net/http"
	"os/signal"
	"secret-images/template"
	"syscall"

	"secret-images/src"
)

func main() {

	http.HandleFunc("POST /encode", src.HandleEncodeImage)
	http.HandleFunc("POST /decode", src.HandleDecodeImage)
	http.HandleFunc("POST /image-cap", src.HandleGetMaxCapacity)

	http.Handle("GET /test-template", templ.Handler(template.Home()))

	go func() {
		println("Server is listening on port 8080")
		err := http.ListenAndServe(":8080", http.DefaultServeMux)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()
	<-ctx.Done()
}
