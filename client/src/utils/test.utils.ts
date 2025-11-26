import { describe, expect, test } from '@jest/globals';
import { dateAndTimeToPosix, posixToDateAndTime, posixToTime, validEmail } from '../utils/utils';

describe("dateAndTimeToPosix", () => {
    test("make date", () => {
        let date = "2025-11-10";
        let time = "21:19";
        let out = 1762809540 + timezoneOffset();
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

describe("posixToTime", () => {
    test("parse to time", () => {
        let posix = 1762809540 + timezoneOffset();
        let expected = "21:19";
        expect(posixToTime(posix)).toEqual(expected);
    });
});

describe("validEmail", () => {
    test("minimal valid", () => {
        let email = "a@b.c";
        expect(validEmail(email)).toEqual(true);
    });
    test("invalid domain", () => {
        let email = "a@bc";
        expect(validEmail(email)).toEqual(false);
    });
    test("invalid .. in domain", () => {
        let email = "a@b..c";
        expect(validEmail(email)).toEqual(false);
    });
    test("invalid .. after domain", () => {
        let email = "a@b.c..";
        expect(validEmail(email)).toEqual(false);
    });
    test("missing local", () => {
        let email = "@b.c";
        expect(validEmail(email)).toEqual(false);
    });
    test("missing domain", () => {
        let email = "a@";
        expect(validEmail(email)).toEqual(false);
    });
});

function timezoneOffset(): number {
    const now = new Date();
    const timezoneOffsetMinutes = now.getTimezoneOffset();
    const timezoneOffsetSeconds = timezoneOffsetMinutes*60;
    return timezoneOffsetSeconds;
}
