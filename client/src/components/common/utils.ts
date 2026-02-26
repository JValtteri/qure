import { Signal } from "@preact/signals-react";
import { fetchEvent, type EventResponse } from "../../api/api";


export function loadDetails(
    show: Signal<{ eventID: string; editor: boolean; }>,
    loadingEvents: Signal<boolean>,
    setEventDetails: React.Dispatch<React.SetStateAction<EventResponse>>
) {
    return async () => {
        // If no event is selected, don't make a request
        if (show.value.eventID === "none") {
            return;
        }
        if (loadingEvents.value == true) {
            return;
        }
        loadingEvents.value = true;
        const details = await fetchEvent(`${show.value.eventID}`);
        setEventDetails(details);
        loadingEvents.value = false;
    };
}
