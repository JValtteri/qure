// Function is timezone aware

const SECONDS_IN_DAY = 60*60*24

export function dateAndTimeToPosix(dateValue: string, startTimeValue: string) {
    const dateTimeString = `${dateValue}T${startTimeValue}`;
    const dateObject = new Date(dateTimeString);
    if (isNaN(dateObject.getTime())) {
        throw new Error(`Invalid date or time input: '${dateTimeString}'`);
    }
    const posixTimestamp = Math.floor(dateObject.getTime() / 1000);
    return posixTimestamp;
}

export function posixToDateAndTime(posix: number) {
    try {
        const  obj = new Date(posix * 1000);
        const str = new Intl.DateTimeFormat("de-DE", {
            dateStyle: "medium",
            timeStyle: "short",
        }).format(obj);
        let time = str.split(", ")[1];
        let date = str.split(", ")[0];
        return `${time} ${date}`;
    } catch(error) {
        return `$Error: ${error}`;
    }
}

export function cycleDay(endTT: number) {
    endTT = endTT + SECONDS_IN_DAY;
    return endTT;
}

export function countSlots(timeslots: Map<number, { Size: number; Reserved: number; }>) {
    let totalSlots = 0;
    let totalReservedSlots = 0;
    for (const [_, data] of timeslots) {
        totalSlots = totalSlots + data.Size;
        totalReservedSlots = totalReservedSlots = data.Reserved;
    }
    return { totalSlots, totalReservedSlots };
}
