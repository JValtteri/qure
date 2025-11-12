/*
 * API module for communicating with the backend server
 */

import { setCookie, ttl } from "../utils/cookie";
import { generalRequest } from "./request";


type Timeslot = { [key: number]: {"Size": number} };

export async function fetchEvents(): Promise<Response> {
    return await generalRequest("api/events", "GET", "");
}

export async function fetchEvent(eventID: string): Promise<Response> {
    const body = JSON.stringify({"EventID": eventID});
    return await generalRequest("api/event", "POST", body);
}

export async function authenticate() {
    return await generalRequest("/api/session/auth", "POST", "");
}

export async function login(username: string, password: string): Promise<any> {
    const body = JSON.stringify({
                        user: username,
                        password: password
                    });
    const authJson = await generalRequest("api/user/login", "POST", body)
    if (authJson.Authenticated) {
        setCookie("sessionKey", authJson.SessionKey, ttl);
        return authJson;
    }
    return null;
}

export async function listReservations() {
    return await generalRequest("/api/user/list", "POST", "");
}

export async function makeReservation(email: string, size: number, eventID: string, timeslot: number) {
    const body = JSON.stringify({
        "Email": email,
        "Size": size,
        "EventId": eventID,
        "Timeslot": timeslot
    });
    return await generalRequest("/api/user/reserve", "POST", body);
}

export async function loginWithEvent(eventID: string) {
    const body = JSON.stringify({
        "EventId": eventID
    });
    return await generalRequest("/api/res/login", "POST", body);
}

export async function registerUser(email: string, password: string) {
    const body = JSON.stringify({
        "User": email,
        "Password": password
    });
    return await generalRequest("/api/user/reserve", "POST", body);
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
    return await generalRequest("/api/admin/create", "POST", body);
}
