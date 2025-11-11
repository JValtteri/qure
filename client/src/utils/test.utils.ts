import { describe, expect, test } from '@jest/globals';
import { dateAndTimeToPosix, posixToDateAndTime } from '../utils/utils';

describe("dateAndTimeToPosix", () => {
    test("make date", () => {
        let date = "2025-11-10";
        let time = "21:19";
        let out = 1762809540 + timezoneOffset();
                //1762795140
        expect(dateAndTimeToPosix(date, time)).toEqual(out);
    });
});

describe("posixToDateAndTime", () => {
    test("parse to date", () => {
        let posix = 1762809540 + timezoneOffset();
        let expected = "21:19 10.11.2025";
        expect(posixToDateAndTime(posix)).toEqual(expected);
    });
});

function timezoneOffset(): number {
    const now = new Date();
    const timezoneOffsetMinutes = now.getTimezoneOffset();
    const timezoneOffsetSeconds = timezoneOffsetMinutes*60;
    return timezoneOffsetSeconds;
}
