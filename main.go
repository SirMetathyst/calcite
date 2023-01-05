package main

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HelloWorld struct {
	List []string
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.HandleFunc("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("."))).ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = RenderHelloWorld(w, HelloWorld{
			List: []string{"Apple", "Orange", "Pineapple"},
		})
	})

	r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		_ = RenderHelloWorldList(w, []string{"Apple", "Orange", "Pineapple"})
	})

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("err: %s\n", err)
		}
	}()

	log.Println("Server Started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("err: %s\n", err)
	}

	log.Println("Server Stopped")
}

func RenderHTML(w http.ResponseWriter, fn func(wr io.Writer) error) error {

	buf := new(bytes.Buffer)

	if err := fn(buf); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")

	_, err := w.Write(buf.Bytes())
	return err
}

func RenderHelloWorld(w http.ResponseWriter, n HelloWorld) error {
	return RenderHTML(w, func(wr io.Writer) error {
		return template.Must(template.ParseFiles("tmpl_layout.html", "tmpl_hello_world.html")).Execute(w, n)
	})
}

func RenderHelloWorldList(w http.ResponseWriter, n []string) error {
	return RenderHTML(w, func(wr io.Writer) error {
		return template.Must(template.ParseFiles("tmpl_layout.html", "tmpl_hello_world.html")).ExecuteTemplate(w, "list", n)
	})
}
