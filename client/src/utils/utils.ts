// Function is timezone aware

const SECONDS_IN_DAY = 60*60*24

export function dateAndTimeToPosix(dateValue: string, startTimeValue: string) {
    const dateTimeString = `${dateValue}T${startTimeValue}`;
    const dateObject = new Date(dateTimeString);
    if (isNaN(dateObject.getTime())) {
        throw new Error("Invalid date or time input.");
    }
    const posixTimestamp = Math.floor(dateObject.getTime() / 1000);
    return posixTimestamp;
}

export function posixToDateAndTime(posix: number) {
    const  obj = new Date(posix * 1000);
    const str = new Intl.DateTimeFormat("de-DE", {
        dateStyle: "medium",
        timeStyle: "short",
    }).format(obj);
    let time = str.split(", ")[1];
    let date = str.split(", ")[0];
    return `${time} ${date}`;
}

export function cycleDay(endTT: number) {
    endTT = +SECONDS_IN_DAY;
    return endTT;
}

