package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network port")

	// of note, adding a flag like this lets you add -help and see it listed as an option.
	// currently running `go run ./cmd/web -help` will show the addr flag as an option
	// you could even call with env var: `go run ./cmd/web -addr=$SNIPPETBOX_ADDR`

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	// Apparently there's also a JSON handler instead of text.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// filename and line number
		AddSource: true,
		// Info is the default, added for easier changing later
		Level: slog.LevelInfo,
	}))

	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes as normal..
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snipped/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// The value returned from the flag.String() function is a pointer, dereference
	// with a leading *
	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("starting server", "addr", *addr)

	// TODO: add way to stop disable directory listings
	// And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())

	os.Exit(1)
}
