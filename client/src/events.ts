
export interface Event{
    name:             string;
    shortDescription: string;
    longDescription:  string;
    dtStart:          string;
    dtEnd:            string;
    staffSlots:       number;
    staff:            number;
    guestSlots:       number;
    guests:           number;
}

const event = {
        name:             "Cool",
        shortDescription: "~",
        longDescription:  "",
        dtStart:          "123",
        dtEnd:            "456",
        staffSlots:       7,
        staff:            6,
        guestSlots:       12,
        guests:           10,
}

export function getEvents(): Event[] {
    const events: Event[] = [];
    events.push(event);
    events.push(structuredClone(event));
    events.push(structuredClone(event));
    events[1].name = "Beans"
    events[2].name = "Foo"
    return events;
}

export function getEvent(): Event {
    return event;
}

