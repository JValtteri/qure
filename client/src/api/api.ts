/*
 * API module for communicating with the backend server
 */

import EventListResponse from "../components/EventList/EventList";
import { setCookie, ttl } from "../utils/cookie";
import { generalRequest } from "./request";


type Timeslot = { [key: number]: {"Size": number} };

interface EventCreationResponse {
    EventID:  number;
    Error:    string;
}

export interface EventResponse {
    ID:                 number
    Name:               string;
    ShortDescription:   string;
    LongDescription:    string;
    Draft:              boolean;
    DtStart:            number;
    DtEnd:              number;
    StaffSlots:         number;
    Staff:              number;
    Timeslots:          any;
}

export type EventListResponse = Array<EventResponse>;

interface AuthResponse {
    User:           string;
    Authenticated:  boolean;
    IsAdmin:        boolean;
    SessionKey:     string;
    Error:          string;
}

type ReservationList = Array<ReservationResponse>

interface ReservationResponse {
    Id:         string;
    EventID:    number;
    ClientID:   string;
    Size:       number;				// Party size
    Confirmed:  number;				// Reserved size
    Timeslot:   number;
    Expiration: number;
    Error:      string;
}

interface RegistrationResponse {
    SessionKey: string;
    Error:      string;
}


export async function fetchEvents(): Promise<EventListResponse> {
    let response = await generalRequest("api/events", "GET");
    let respBody = await response.json() as EventListResponse;
    return respBody;
}

export async function fetchEvent(eventID: string): Promise<EventResponse> {
    const body = {"EventID": eventID};
    let response = await generalRequest("api/event", "POST", body);
    let respBody = await response.json() as EventResponse;
    return respBody;
}

export async function authenticate(): Promise<AuthResponse | null> {
    let response = await generalRequest("/api/session/auth", "POST");
    let respBody = await response.json() as AuthResponse;
    if (respBody.Authenticated) {
        return respBody;
    }
    return null;
}

export async function login(username: string, password: string): Promise<AuthResponse | null> {
    const body = {
        user: username,
        password: password
    };
    const response = await generalRequest("api/user/login", "POST", body)
    let authJson = await response.json() as AuthResponse;
    if (authJson.Authenticated) {
        setCookie("sessionKey", authJson.SessionKey, ttl);
        return authJson;
    }
    return null;
}

export async function logout(): Promise<AuthResponse> {
    let response = await generalRequest("/api/user/logout", "POST");
    let respBody = await response.json() as AuthResponse;
    return respBody;
}

export async function listReservations(): Promise<ReservationList> {
    let response = await generalRequest("/api/user/list", "POST", "");
    let respBody = await response.json() as ReservationList;
    return respBody;
}

export async function makeReservation(
    email: string,
    size: number,
    eventID: string,
    timeslot: number
): Promise<ReservationResponse> {
    const body = {
        "Email": email,
        "Size": size,
        "EventId": eventID,
        "Timeslot": timeslot
    };
    let response = await generalRequest("/api/user/reserve", "POST", body);
    let respBody = await response.json() as ReservationResponse;
    return respBody;
}

export async function loginWithEvent(eventID: string): Promise<AuthResponse> {
    const body = {
        "EventId": eventID
    };
    let response = await generalRequest("/api/res/login", "POST", body);
    let respBody = await response.json() as AuthResponse;
    return respBody;
}

export async function registerUser(email: string, password: string): Promise<RegistrationResponse> {
    const body = {
        "User": email,
        "Password": password
    };
    let response = await generalRequest("/api/user/reserve", "POST", body);
    let respBody = await response.json() as RegistrationResponse;
    return respBody;
}

export async function makeEvent(
    name: string,
    shortDesc: string,
    longDesc: string,
    start: number,
    end: number,
    draft: boolean,
    staffSlots: number,
    timeslots: number[],
    groupSize: number
): Promise<EventCreationResponse> {
    const timeslotobjs: Timeslot = {};
    for (const slot of timeslots) {
        timeslotobjs[slot] = { "Size": groupSize };
    }
    const body = (
        {
            "Event": {
                "Name":             name,
                "ShortDescription": shortDesc,
                "LongDescription":  longDesc,
                "Draft":            draft,
                "DtStart":          start,
                "DtEnd":            end,
                "StaffSlots":       staffSlots,
                "Timeslots":        timeslotobjs
            }
        })
    let response = await generalRequest("/api/admin/create", "POST", body)
    let respBody = await response.json() as EventCreationResponse;
    return respBody;
}
