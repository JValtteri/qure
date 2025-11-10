// Function is timezone aware
export function dateAndTimeToPosix(dateValue: string, startTimeValue: string) {
    const dateTimeString = `${dateValue}T${startTimeValue}`;
    try {
        const dateObject = new Date(dateTimeString);
        if (isNaN(dateObject.getTime())) {
            throw new Error("Invalid date or time input.");
        }
        const posixTimestamp = Math.floor(dateObject.getTime() / 1000);
        return posixTimestamp;
    } catch (error) {
        console.error(error);
        return null;
    }
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
