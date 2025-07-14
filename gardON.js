const modeSwitch = document.getElementById("modeSwitch");
const scheduleMode = document.getElementById("scheduleMode");
const manualMode = document.getElementById("manualMode");
const tankLevel1 = document.getElementById("tankLevel1");
const tankLevel2 = document.getElementById("tankLevel2");
const tankLevel3 = document.getElementById("tankLevel3");

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

