/*
 * Cookie(s)
 */

export const ttl = 30;     // cookie max life in days
const milliSecondInDay = 1000*60*60*24;

export function setCookie(name: string, value: string, ttl: number) {
    const d = new Date();
    d.setTime(d.getTime() + (ttl*milliSecondInDay));
    const expires = "expires="+ d.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";SameSite=Lax" + ";path=/";
}

export function clearCookie(name: string) {
    const value = "";
    const expires = "expires=Thu, 01 Jan 1970 00:00:00 GMT";
    document.cookie = name + "=" + value + ";" + expires + ";SameSite=Lax" + ";path=/";
}

export function getCookie(name: string): string {
    name = name + "=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookies = decodedCookie.split(';');
    for(let i = 0; i < cookies.length; i++) {
        let c = cookies[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}
