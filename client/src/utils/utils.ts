// Function is timezone aware

const SECONDS_IN_DAY = 60*60*24

export function dateAndTimeToPosix(dateValue: string, startTimeValue: string): number {
    const dateTimeString = `${dateValue}T${startTimeValue}`;
    const dateObject = new Date(dateTimeString);
    if (isNaN(dateObject.getTime())) {
        throw new Error(`Invalid date or time input: '${dateTimeString}'`);
    }
    const posixTimestamp = Math.floor(dateObject.getTime() / 1000);
    return posixTimestamp;
}

// Produces a printable european style date and 24 hour time
export function posixToDateAndTime(posix: number): string {
    try {
        const obj = new Date(posix * 1000);
        const str = new Intl.DateTimeFormat("de-DE", {
            dateStyle: "medium",
            timeStyle: "short",
        }).format(obj);
        const time = str.split(", ")[1];
        const date = str.split(", ")[0];
        return `${time} ${date}`;
    } catch(error) {
        return `$Error: ${error}`;
    }
}

// Produces a printable european style 24 hour time
export function posixToTime(posix: number): string {
    try {
        const obj = new Date(posix * 1000);
        const str = new Intl.DateTimeFormat("de-DE", {
        timeStyle: "short",
        }).format(obj);
        return str;
    } catch (error) {
        return `$Error: ${error}`;
    }
}

// Produces a date compatible with date input field value
export function posixToDate(posix: number): string {
    try {
        const d = new Date(posix * 1000);         // POSIX to JS Date (ms)
        return d.toISOString().split('T')[0];     // e.g. "2026-01-09"
    } catch (error) {
        return `$Error: ${error}`;
    }
}

// Returns current POSIX timestamp
export function posixNow(): number {
    let d = new Date();
    return d.getTime()/1000;
}

// Returs true, if given POSIX is in the past
export function isPast(time: number): boolean {
    let now = posixNow();
    return time < now;
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
        totalReservedSlots = totalReservedSlots + data.Reserved;
    }
    return { totalSlots, totalReservedSlots };
}

export function validEmail(email: string): boolean {
    try {
        const [local, domain] = email.split('@');
        if (local.length < 1 || local.length > 64 || domain.length > 63) {
            return false
        }
        const splitDomain = domain.split('.');
        if (splitDomain.length < 2) {
            return false
        }
        for (const elem of splitDomain) {
            if ( elem.length === 0 ) {
                return false;
            }
        }
    } catch (error) {
        console.warn(`Invalid email: '${email}' ; ${error}`)
        return false;
    }
    return true
}
