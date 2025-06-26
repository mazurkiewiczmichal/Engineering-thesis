const modeSwitch = document.getElementById("modeSwitch");
const scheduleMode = document.getElementById("scheduleMode");
const manualMode = document.getElementById("manualMode");

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

