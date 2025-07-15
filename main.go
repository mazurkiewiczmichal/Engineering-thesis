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
	waterLevel1 := false
	waterLevel2 := false
	waterLevel3 := true
	valveSwitch := true
	pumpSwitch := true

	mux.HandleFunc("/pumpOn", func(w http.ResponseWriter, r *http.Request) {
		pumpSwitch = true
	})
	mux.HandleFunc("/pumpOff", func(w http.ResponseWriter, r *http.Request) {
		pumpSwitch = false
	})

	mux.HandleFunc("/valveOn", func(w http.ResponseWriter, r *http.Request) {
		valveSwitch = true
	})
	mux.HandleFunc("/valveOff", func(w http.ResponseWriter, r *http.Request) {
		valveSwitch = false
	})

	_ = valveSwitch
	_ = pumpSwitch

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

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := struct {
			WaterLevel1 bool `json:"waterLevel1"`
			WaterLevel2 bool `json:"waterLevel2"`
			WaterLevel3 bool `json:"waterLevel3"`
			ValveSwitch bool `json:"valveSwitch"`
			PumpSwitch  bool `json:"pumpSwitch"`
		}{
			WaterLevel1: waterLevel1,
			WaterLevel2: waterLevel2,
			WaterLevel3: waterLevel3,
			ValveSwitch: valveSwitch,
			PumpSwitch:  pumpSwitch,
		}

		json.NewEncoder(w).Encode(status)
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
