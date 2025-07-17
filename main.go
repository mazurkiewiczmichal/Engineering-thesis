package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	daysWeekday  []time.Weekday = []time.Weekday{}
	days         []string       = []string{}
	soilMoisture                = 22
	waterLevel1                 = false
	pinLevel1                   = rpio.Pin(4)
	waterLevel2                 = false
	pinLevel2                   = rpio.Pin(17)
	waterLevel3                 = true
	pinLevel3                   = rpio.Pin(27)
	valveSwitch                 = true
	valvePin                    = rpio.Pin(22)
	pumpSwitch                  = true
	pumpPin                     = rpio.Pin(10)
	pouring                     = false
	initialTime  string
	endTime      string
)

func main() {
	err := rpio.Open()
	if err != nil {
		log.Fatal(err)
	}

	pumpPin.Output()
	valvePin.Output()

	pinLevel1.Input()
	pinLevel2.Input()
	pinLevel3.Input()

	pinLevel1.PullUp()
	pinLevel2.PullUp()
	pinLevel3.PullUp()

	mux := http.NewServeMux()

	mux.HandleFunc("/pumpOn", func(w http.ResponseWriter, r *http.Request) {
		pumpSwitch = true
		pumpPin.High()
	})
	mux.HandleFunc("/pumpOff", func(w http.ResponseWriter, r *http.Request) {
		pumpSwitch = false
		pumpPin.Low()
	})

	mux.HandleFunc("/valveOn", func(w http.ResponseWriter, r *http.Request) {
		valveSwitch = true
		valvePin.High()
	})
	mux.HandleFunc("/valveOff", func(w http.ResponseWriter, r *http.Request) {
		valveSwitch = false
		valvePin.Low()
	})

	_ = valveSwitch
	_ = pumpSwitch

	// Serwowanie pliku HTML na stronie głównej
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "gardON.html")
	})

	// mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
	// 	if err := r.ParseForm(); err != nil {
	// 		http.Error(w, "Form error", http.StatusBadRequest)
	// 		return
	// 	}

	// 	// Pobieranie wartości z formularza
	// 	days := r.Form["days"]
	// 	initialTime := r.FormValue("initialTime")
	// 	endTime := r.FormValue("endTime")

	// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 	fmt.Fprintf(w, "Selected days: %s<br>Since: %s<br>Until: %s",
	// 		strings.Join(days, ", "),
	// 		initialTime,
	// 		endTime,
	// 	)

	// })

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
			Pouring     bool `json:"pouring"`
		}{
			WaterLevel1: waterLevel1,
			WaterLevel2: waterLevel2,
			WaterLevel3: waterLevel3,
			ValveSwitch: valveSwitch,
			PumpSwitch:  pumpSwitch,
			Pouring:     pouring,
		}

		json.NewEncoder(w).Encode(status)
	})

	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Form error", http.StatusBadRequest)
			return
		}

		days = r.Form["days"]
		daysToWeekday()
		initialTime := r.FormValue("initialTime")
		endTime := r.FormValue("endTime")

		fmt.Println(isDayToday())

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Selected days: %s\nSince: %s\nUntil: %s",
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

func isDayToday() bool {
	for _, d := range daysWeekday {
		if d == time.Now().Weekday() {
			return true
		}
	}
	return false
}

func daysToWeekday() {

	for _, d := range days {
		switch strings.ToLower(d) {
		case "sunday":
			daysWeekday = append(daysWeekday, time.Sunday)
		case "monday":
			daysWeekday = append(daysWeekday, time.Monday)
		case "tuesday":
			daysWeekday = append(daysWeekday, time.Tuesday)
		case "wednesday":
			daysWeekday = append(daysWeekday, time.Wednesday)
		case "thursday":
			daysWeekday = append(daysWeekday, time.Thursday)
		case "friday":
			daysWeekday = append(daysWeekday, time.Friday)
		case "saturday":
			daysWeekday = append(daysWeekday, time.Saturday)
		default:
			fmt.Println("Nieznany dzień:", d)
		}
	}
}

// -------------------------------wip------------------------------------------
// func stringTimeToTime(timeToChange string) time.Time {
// 	t, _ := time.Parse("15:04", timeToChange)
// 	return t
// }

// func isNowTime() bool {
// 	var now = time.Date(0, time.January, 1, time.Now().Hour(), time.Now().Minute(), 0, 0, time.UTC)
// 	if stringTimeToTime(initialTime).After(now) && stringTimeToTime(endTime).Before(now) {
// 		return true
// 	}
// 	return false
// }

// func counterToStart(initial time.Time) time.Timer {
// 	var timerToStart *time.Timer
// 	var now = time.Date(0, time.January, 1, time.Now().Hour(), time.Now().Minute(), 0, 0, time.UTC)
// 	if isNowTime() {
// 		timerToStart = time.NewTimer(0 * time.Second)
// 	} else {
// 		timerToStart = time.NewTimer(initial.Sub(now))
// 	}
// 	return *timerToStart
// }

// func counterToEnd(end time.Time) time.Timer {
// 	var timerToEnd *time.Timer
// 	var now = time.Date(0, time.January, 1, time.Now().Hour(), time.Now().Minute(), 0, 0, time.UTC)
// 	timerToEnd = time.NewTimer(end.Sub(now))
// 	return *timerToEnd
// }

// func scheduleModeOn() {

// }
// --------------------------------------------------------------------------------
// func timeToStart(){
// 	for _, d := range days{

// 		if <
// 	}
// }

// func isDayMached() {
// 	ticker := time.NewTicker(10*time.Second)
// 	for _, d := range days{
// 		if strings.EqualFold(d, now.Weekday().String)
// 		<- ticker.C
// 	}

// 	ttt := time.NewTimer()
// }

// func isTimeBetween(){
// 	ticker := time.NewTicker(10*time.Second)
// 	for {
// 		if time.Now().Format("15:04")>=initialTime&&time.Now().Format("15:04") <=endTime
// 	}
// }
