package main

import (
	"log"
	"mime"
	"net/http"
	"os"
)

func main() {
	wd := os.Args[1]
	port := os.Args[2]

	mime.AddExtensionType(".wasm", "application/wasm")

	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.Dir(wd)))
	log.Fatal(http.ListenAndServe(":"+port, m))
}
