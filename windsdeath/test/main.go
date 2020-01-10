package main

import (
        "net/http"
)

func main() {
        http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
        http.ListenAndServe(":3000", nil)
}