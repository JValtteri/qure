import { Signal } from "@preact/signals-react";
import { authenticate, fetchEvent, type EventResponse } from "../../api/api";


export function loadDetails(
    show: Signal<{ eventID: string; view: string}>,
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

        try {
            const details = await fetchEvent(`${show.value.eventID}`);
            setEventDetails(details);
        } catch (error: any) {
            console.warn(error.message);
        }
        loadingEvents.value = false;
    };
}

export async function resumeSession(
    setServerError: React.Dispatch<React.SetStateAction<string>> | undefined,
    setErrorVisible: React.Dispatch<React.SetStateAction<boolean>> | undefined,
    user: Signal<{ username: string, loggedIn: boolean, role: string}>,
    showLogin: Signal<boolean> | undefined
): Promise<void> {
    try {
        const auth = await authenticate();
        if ( auth != null ) {
        if (showLogin) {
            showLogin.value = false;
        }
        user.value = { username: auth.User, loggedIn: true, role: auth.Role};
        }
    } catch (error) {
        setServerError && setServerError(`${error}`);
        setErrorVisible && setErrorVisible(true);
    }
}
