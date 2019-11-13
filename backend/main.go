package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/stuffbin"

	flag "github.com/spf13/pflag"
)

var (
	version = "1.0.1"
	ko      = koanf.New(".")
)

// App contains all the global components which
// are injected into HTTP request handlers.
type App struct {
	fs     stuffbin.FileSystem
	logger *log.Logger
}

func wrap(app *App, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "app", app)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func initLogger() *log.Logger {
	return log.New(os.Stdout, "newsletter: ", log.Ldate|log.Ltime|log.Llongfile)
}

func initConfig() {
	// Register --help handler.
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.StringSlice("config", []string{"config.toml"},
		"Path to one or more TOML config files to load in order")
	f.StringSlice("prov", []string{"smtp.prov"},
		"Path to a provider plugin. Can specify multiple values.")
	f.Bool("version", false, "Show build version")
	f.Parse(os.Args[1:])

	// Display version.
	if ok, _ := f.GetBool("version"); ok {
		fmt.Println(version)
		os.Exit(0)
	}

	// Read the config files.
	cFiles, _ := f.GetStringSlice("config")
	for _, f := range cFiles {
		log.Printf("reading config: %s", f)
		if err := ko.Load(file.Provider(f), toml.Parser()); err != nil {
			log.Printf("error reading config: %v", err)
		}
	}
	// Load environment variables and merge into the loaded config.
	if err := ko.Load(env.Provider("NEWSLETTER_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "NEWSLETTER_")), "__", ".", -1)
	}), nil); err != nil {
		log.Printf("error loading env config: %v", err)
	}

	ko.Load(posflag.Provider(f, ".", ko), nil)
}

// initFileSystem initializes the stuffbin FileSystem to provide
// access to bunded static assets to the app.
func initFileSystem(binPath string) (stuffbin.FileSystem, error) {
	fs, err := stuffbin.UnStuff(os.Args[0])
	if err != nil {
		return nil, err
	}
	fmt.Println("loaded files", fs.List())
	return fs, nil
}

func main() {
	app := &App{}
	initConfig()
	app.logger = initLogger()
	app.logger.Printf("Booting program...")

	// Initialize the static file system into which all
	// required static assets (.css, .js files etc.) are loaded.
	fs, err := initFileSystem(os.Args[0])
	if err != nil {
		app.logger.Fatalf("error reading stuffed binary: %v", err)
	}
	app.fs = fs

	// Register handles.
	r := chi.NewRouter()

	r.Get("/api/", wrap(app, handleAPIRoot))
	r.Get("/api/health", wrap(app, handleHealthCheck))
	r.Post("/api/create", wrap(app, handleNewSubscription))
	// TODO: Confirm email endpoint.

	r.Get("/", wrap(app, handleIndex))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		app.fs.FileServer().ServeHTTP(w, r)
	})

	// Start a web server
	srv := &http.Server{
		Addr:         ko.String("server.address"),
		ReadTimeout:  ko.Duration("server.read_timeout") * time.Second,
		WriteTimeout: ko.Duration("server.write_timeout") * time.Second,
		IdleTimeout:  ko.Duration("server.keepalive_timeout") * time.Second,
		Handler:      r,
	}

	app.logger.Printf("starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		app.logger.Fatalf("couldn't start server: %v", err)
	}
}
