
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

export const event = {
    name:             "Lorem ipsum",
    shortDescription: "Lorem ipsum dolor sit amet, verear virtute qui an",
    longDescription:  "Lorem ipsum dolor sit amet, verear virtute qui an, eos postea persius deleniti te, ei liber saepe albucius sea. Pri mundi molestie at. At splendide cotidieque eam. Feugiat mediocrem accusamus eu duo.\n\nMea adhuc dissentiunt id, pro te eruditi facilis liberavisse. No duo tation minimum. Id discere lucilius volutpat vim, ei eos volutpat concludaturque. Mea prima quaeque volumus ut, tation possit platonem ne nam. Honestatis mediocritatem ne eum. Ad complectitur signiferumque cum, ut vix putent alienum.",
    dtStart:          "10-10-2025 20:00",
    dtEnd:            "10-10-2025 23:30",
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
    events[1].name = "Dolor sit amet"
    events[2].name = "Verear virtute"
    return events;
}

export function getEvent(): Event {
    return event;
}

