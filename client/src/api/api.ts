/*
 * API module for communicating with the backend server
 */

import { setCookie, ttl } from "../utils/cookie";
import { generalRequest } from "./request";


type Timeslot = { [key: number]: {"Size": number} };

export async function fetchEvents(): Promise<Response> {
    const response = await generalRequest("api/events", "GET", "");
    return response;
}

export async function fetchEvent(eventID: string): Promise<Response> {
    const body = JSON.stringify({"EventID": eventID});
    const response = await generalRequest("api/event", "POST", body);
    return response;
}

export async function authenticate() {
    const response = await generalRequest("/api/session/auth", "POST", "");
    return response;
}

export async function login(username: string, password: string): Promise<any> {
    const body = JSON.stringify({
                        user: username,
                        password: password
                    });
    const response = await generalRequest("api/user/login", "POST", body)
    const authJson = await response.json();
    if (authJson.Authenticated) {
        setCookie("sessionKey", authJson.SessionKey, ttl);
        return authJson;
    }
    return null;
}

export async function listReservations() {
    const response = await generalRequest("/api/user/list", "POST", "");
    return response;
}

export async function makeReservation(email: string, size: number, eventID: string, timeslot: number) {
    const body = JSON.stringify({
        "Email": email,
        "Size": size,
        "EventId": eventID,
        "Timeslot": timeslot
    });
    const response = await generalRequest("/api/user/reserve", "POST", body);
    return response;
}

export async function loginWithEvent(eventID: string) {
    const body = JSON.stringify({
        "EventId": eventID
    });
    const response = await generalRequest("/api/res/login", "POST", body);
    return response;
}

export async function registerUser(email: string, password: string) {
    const body = JSON.stringify({
        "User": email,
        "Password": password
    });
    const response = await generalRequest("/api/user/reserve", "POST", body);
    return response;
}

export async function makeEvent(
            name: string, shortDesc: string, longDesc: string, start: number,
            end: number, draft: boolean, staffSlots: number, timeslots: number[],
            groupSize: number
        ): Promise<Response> {
    const timeslotobjs: Timeslot = {};
    for (const slot of timeslots) {
        timeslotobjs[slot] = { "Size": groupSize };
    }
    const body = JSON.stringify({
                        "Name":             name,
                        "ShortDescription": shortDesc,
                        "LongDescription":  longDesc,
                        "Draft":            draft,
                        "DtStart":          start,
                        "DtEnd":            end,
                        "StaffSlots":       staffSlots,
                        "Timeslots":        timeslotobjs
                    })
    const response = await generalRequest("/api/admin/create", "POST", body);
    return response;
}
