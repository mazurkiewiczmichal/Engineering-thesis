const modeSwitch = document.getElementById("modeSwitch");
const scheduleMode = document.getElementById("scheduleMode");
const manualMode = document.getElementById("manualMode");
const waterLevel1 = document.getElementById("waterLevel1");
const waterLevel2 = document.getElementById("waterLevel2");
const waterLevel3 = document.getElementById("waterLevel3");
const soilMoisture = document.getElementById("soilMoisture");
const valveSwitch = document.getElementById("valveSwitch");
const pumpSwitch = document.getElementById("pumpSwitch");
const pouring = document.getElementById("pouring");
const scheduleForm = document.getElementById("scheduleForm");


modeSwitch.onchange = (event) => {
    const isScheduleMode = event.target.checked;
    let url;

    if (isScheduleMode) {
        manualMode.classList.add("hidden");
        scheduleMode.classList.remove("hidden");
        url = "/scheduleMode";
    } else {
        scheduleMode.classList.add("hidden");
        manualMode.classList.remove("hidden");
        url = "/manualMode";
    }

    fetch(url).catch(function(e) {
        console.error(e);
    });
};



    document.addEventListener("DOMContentLoaded", function () {

    });

    pumpSwitch.onchange = (event) => {
        const checked = event.target.checked;
        const url = checked ? "/pumpOn" : "/pumpOff";
        fetch(url).catch((e) => console.error(e));
    };

    valveSwitch.onchange = (event) => {
        const checked = event.target.checked;
        const url = checked ? "/valveOn" : "/valveOff";
        fetch(url).catch((e) => console.error(e));
    };






    // tankLevel1.classList.add("filled");


    fetch('/data')
        .then(res => res.json())
        .then(data => {
            soilMoisture.innerText = "Irrigation: " + data;
        })
        .catch(err => {
            soilMoisture.innerText = "Irigation: failed";
            console.error(err);
        });

    fetch('/status')
        .then(res => res.json())
        .then(data => {
            console.log(data.waterLevel1);
            console.log(data.waterLevel2);
            console.log(data.waterLevel3);
            console.log(data.valveSwitch);
            console.log(data.pumpSwitch);

            if (data.pumpSwitch !== undefined) {
                pumpSwitch.checked = data.pumpSwitch;
            }
            if (data.valveSwitch !== undefined) {
                valveSwitch.checked = data.valveSwitch;
            }


            if (data.waterLevel1) {
                waterLevel1.classList.add("filled");
            } else {
                waterLevel1.classList.remove("filled");
            }
            if (data.waterLevel2) {
                waterLevel2.classList.add("filled");
            } else {
                waterLevel2.classList.remove("filled");
            }
            if (data.waterLevel3) {
                waterLevel3.classList.add("filled");
            } else {
                waterLevel3.classList.remove("filled");
            }
            if (data.pouring) {
                pouring.classList.add("filled");
            } else {
                pouring.classList.remove("filled");
            }

        })
        .catch(err => console.error(err));

    scheduleForm.addEventListener("submit", function (e) {
        e.preventDefault();

        const formData = new FormData(this);

        fetch("/submit", {
            method: "POST",
            body: formData
        })
            .then(res => res.text())
            .then(msg => {
                alert(msg);
            })
            .catch(err => {
                alert("Sending data failed!");
                console.error(err);
            });
    });