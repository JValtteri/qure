import { Signal } from "@preact/signals-react";
import { fetchEvent, type EventResponse } from "../../api/api";


export function loadDetails(
    show: Signal<{ eventID: number; editor: boolean; }>,
    setEventDetails: React.Dispatch<React.SetStateAction<EventResponse>>
) {
    return async () => {
        // If no event is selected, don't make a request
        if (show.value.eventID === -1) {
            return;
        }
        const details = await fetchEvent(`${show.value.eventID}`);
        setEventDetails(details);
    };
}
