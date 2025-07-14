package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	// Serwowanie pliku HTML na stronie głównej
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "gardON.html")
	})

	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Błąd formularza", http.StatusBadRequest)
			return
		}

		// Pobieranie wartości z formularza
		days := r.Form["days"]
		initialTime := r.FormValue("initialTime")
		endTime := r.FormValue("endTime")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Zaznaczone dni: %s<br>Od: %s<br>Do: %s",
			strings.Join(days, ", "),
			initialTime,
			endTime,
		)

	})

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/gardON.css", fs)
	mux.Handle("/gardON.js", fs)
	mux.Handle("/logoOFF.png", fs)
	mux.Handle("/logoON.png", fs)
	mux.Handle("/logoOFFswitch.png", fs)
	mux.Handle("/logoONswitch.png", fs)

	http.ListenAndServe(":12346", mux)

}
