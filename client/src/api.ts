/*
 * API module for communicating with the backend server
 */

import { setCookie, ttl } from "./cookie";


export async function fetchEvents(): Promise<Response> {
    try {
        const response = await fetch("api/events");
        return response;
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
        return new Response();
    }
}

export async function login(username: string, password: string): Promise<any> {
    const response = await requestLogin(username, password);
    const authJson = await response.json();
    if (authJson.Authenticated) {
        setCookie("sessionKey", authJson.SessionKey, ttl);
        return authJson
    }
    return null;
}

async function requestLogin(username: string, password: string): Promise<Response> {
    try {
        const response = await fetch("api/user/login", {
            method: "POST",
            body: JSON.stringify({
                user: username,
                password: password
            }),
            headers: {
                "Content-type": "application/json; charset=UTF-8"
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`)
        }
        return response;
    } catch (error) {
        console.error(error);
        return new Response();
    }
}
