const modeSwitch = document.getElementById("modeSwitch");
const scheduleMode = document.getElementById("scheduleMode");
const manualMode = document.getElementById("manualMode");
const tankLevel1 = document.getElementById("tankLevel1");
const tankLevel2 = document.getElementById("tankLevel2");
const tankLevel3 = document.getElementById("tankLevel3");
const soilMoisture = document.getElementById("soilMoisture");


modeSwitch.onchange = (event) => {
    const isScheduleMode = event.target.checked;
    if (isScheduleMode) {
        manualMode.classList.add("hidden");
        scheduleMode.classList.remove("hidden");
    } else {
        scheduleMode.classList.add("hidden");
        manualMode.classList.remove("hidden");
    }
};

tankLevel1.classList.add("filled");


fetch('/data')
    .then(res => res.json())
    .then(data => {
        soilMoisture.innerText = "Soil moisture: " + data + "%";
    })
    .catch(err => {
        soilMoisture.innerText = "Błąd podczas pobierania danych";
        console.error(err);
    });

