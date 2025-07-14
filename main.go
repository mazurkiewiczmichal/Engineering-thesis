package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	soilMoisture := 22
	// waterLevel1 := false
	// waterLevel2 := false
	// waterLevel3 := false

	// Serwowanie pliku HTML na stronie głównej
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "gardON.html")
	})

	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Form error", http.StatusBadRequest)
			return
		}

		// Pobieranie wartości z formularza
		days := r.Form["days"]
		initialTime := r.FormValue("initialTime")
		endTime := r.FormValue("endTime")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Selected days: %s<br>Since: %s<br>Until: %s",
			strings.Join(days, ", "),
			initialTime,
			endTime,
		)

	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(soilMoisture)
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
