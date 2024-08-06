package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network port")

	// of note, adding a flag like this lets you add -help and see it listed as an option.
	// currently running `go run ./cmd/web -help` will show the addr flag as an option
	// Another positive, if you set flags up, you can pass env variables into the command.
	// ex. `go run ./cmd/web -addr=$SNIPPETBOX_ADDR`

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

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

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So in this code, that means the addr variable
	// is actually a pointer, and we need to dereference it (i.e. prefix it with
	// the * symbol) before using it. Note that we're using the log.Printf()
	// function to interpolate the address with the log message.
	log.Printf("starting server on %s", *addr)

	// TODO: add way to stop disable directory listings
	// And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
