import { describe, expect, test } from '@jest/globals';
import { dateAndTimeToPosix, posixToDateAndTime, posixToTime, validEmail } from '../utils/utils';

describe("dateAndTimeToPosix", () => {
    test("make date", () => {
        const date = "2025-11-10";
        const time = "21:19";
        const out = 1762809540 + timezoneOffset();
        expect(dateAndTimeToPosix(date, time)).toEqual(out);
    });
});

describe("posixToDateAndTime", () => {
    test("parse to date", () => {
        const posix = 1762809540 + timezoneOffset();
        const expected = "21:19 10.11.2025";
        expect(posixToDateAndTime(posix)).toEqual(expected);
    });
});

describe("posixToTime", () => {
    test("parse to time", () => {
        const posix = 1762809540 + timezoneOffset();
        const expected = "21:19";
        expect(posixToTime(posix)).toEqual(expected);
    });
});

describe("validEmail", () => {
    test("minimal valid", () => {
        const email = "a@b.c";
        expect(validEmail(email)).toEqual(true);
    });
    test("invalid domain", () => {
        const email = "a@bc";
        expect(validEmail(email)).toEqual(false);
    });
    test("invalid .. in domain", () => {
        const email = "a@b..c";
        expect(validEmail(email)).toEqual(false);
    });
    test("invalid .. after domain", () => {
        const email = "a@b.c..";
        expect(validEmail(email)).toEqual(false);
    });
    test("missing local", () => {
        const email = "@b.c";
        expect(validEmail(email)).toEqual(false);
    });
    test("missing domain", () => {
        const email = "a@";
        expect(validEmail(email)).toEqual(false);
    });
});

function timezoneOffset(): number {
    const now = new Date();
    const timezoneOffsetMinutes = now.getTimezoneOffset();
    const timezoneOffsetSeconds = timezoneOffsetMinutes*60;
    return timezoneOffsetSeconds;
}
