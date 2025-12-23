/*
 * API module for communicating with the backend server
 */

import EventListResponse from "../components/EventList/EventList";
import { setCookie, ttl } from "../utils/cookie";
import { generalRequest } from "./request";


export type Timeslot = {
    "Size": number, "Reserved": number
};

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
    Timeslots:          Record<number, Timeslot>;
}

export type EventListResponse = Array<EventResponse>;

interface AuthResponse {
    User:           string;
    Authenticated:  boolean;
    IsAdmin:        boolean;
    SessionKey:     string;
    Error:          string;
}

export type ReservationList = Array<ReservationResponse>;

export interface ReservationResponse {
    Id:         string;
    EventID:    number;
    ClientID:   string;
    Size:       number;				// Party size
    Confirmed:  number;				// Reserved size
    Timeslot:   number;
    Expiration: number;
    Error:      string;
    Session:    string;
    Event: {
        ID:         number;
        Name:       string;
        DtStart:    number;
        DtEnd:      number;
    }
}

interface RegistrationResponse {
    SessionKey: string;
    Error:      string;
}

interface SuccessResponse {
    Success:    boolean;
    SessionKey: string;
    Error:      string;
}

export async function fetchEvents(): Promise<EventListResponse> {
    const response = await generalRequest("api/events", "POST");
    const respBody = await response.json() as EventListResponse;
    return respBody;
}

export async function fetchEvent(eventID: string): Promise<EventResponse> {
    const body = {"EventID": eventID};
    const response = await generalRequest("api/event", "POST", body);
    const respBody = await response.json() as EventResponse;
    return respBody;
}

export async function authenticate(): Promise<AuthResponse | null> {
    const response = await generalRequest("/api/session/auth", "POST");
    const respBody = await response.json() as AuthResponse;
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
    const authJson = await response.json() as AuthResponse;
    if (authJson.Authenticated) {
        setCookie("sessionKey", authJson.SessionKey, ttl);
        return authJson;
    }
    return null;
}

export async function logout(): Promise<AuthResponse> {
    const response = await generalRequest("/api/user/logout", "POST");
    const respBody = await response.json() as AuthResponse;
    return respBody;
}

export async function listReservations(): Promise<ReservationList> {
    const response = await generalRequest("/api/user/list", "POST", "");
    const respBody = await response.json() as ReservationList;
    return respBody;
}

export async function editPassword(
    username: string,
    password: string,
    newPassword: string
): Promise<SuccessResponse> {
    const body = {
        User:        username,
        Password:    password,
        NewPassword: newPassword
    };
    const response = await generalRequest("/api/user/change", "POST", body);
    const respBody = await response.json() as SuccessResponse;
    // Session key is renewed when password is changed.
    // This is because all existing sessions are invalidated
    // after password change. This is a security feature.
    if (respBody.Success) {
        setCookie("sessionKey", respBody.SessionKey, ttl);
    }
    return respBody;
}

export async function deleteUser(
    username: string,
    password: string,
): Promise<SuccessResponse> {
    const body = {
        User:        username,
        Password:    password
    };
    const response = await generalRequest("/api/user/delete", "POST", body);
    const respBody = await response.json() as SuccessResponse;
    return respBody;
}

export async function makeReservation (
    email: string,
    size: number,
    eventID: number,
    timeslot: number
): Promise<ReservationResponse> {
    const body = {
        "User": email,
        "Size": size,
        "EventId": eventID,
        "Timeslot": timeslot
    };
    const response = await generalRequest("/api/user/reserve", "POST", body);
    const respBody = await response.json() as ReservationResponse;
    if (respBody.Error == "") {                             // If the user isn't signed in, a session is created on successful
        setCookie("sessionKey", respBody.Session, ttl);     // reservation. The session key is updated here.
    }
    return respBody;
}

export async function loginWithEvent(eventID: string): Promise<AuthResponse> {
    const body = {
        "EventId": eventID
    };
    const response = await generalRequest("/api/res/login", "POST", body);
    const respBody = await response.json() as AuthResponse;
    return respBody;
}

export async function registerUser(email: string, password: string): Promise<RegistrationResponse> {
    const body = {
        "User": email,
        "Password": password
    };
    const response = await generalRequest("/api/user/register", "POST", body);
    const respBody = await response.json() as RegistrationResponse;
    setCookie("sessionKey", respBody.SessionKey, ttl);
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
    timeslots: Map<number, {"Size": number}>
): Promise<EventCreationResponse> {
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
                "Timeslots":        Object.fromEntries(timeslots.entries())
            }
        });
    const response = await generalRequest("/api/admin/create", "POST", body)
    const respBody = await response.json() as EventCreationResponse;
    return respBody;
}
