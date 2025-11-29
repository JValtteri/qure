/**
 * @jest-environment jsdom
 */

import {describe, expect, test} from '@jest/globals';
import {setCookie, getCookie, ttl, clearCookie} from './cookie';


beforeEach(() => {
    document.cookie = "";
});

afterEach(() => {
    document.cookie = "";
});

describe("set and getCookie()", () => {
    const firstCookieName = "test-cookie";
    const secondCookieName = "test-other-cookie";
    const thirdCookieName = "test-no-cookie";
    const value1 = "test123";
    const value2 = "test456";
    test("set and get first cookie", () => {
        setCookie(firstCookieName, value1, ttl);
        const got = getCookie(firstCookieName);
        expect(got).toEqual(value1);
    });
    test("set and get second cookie", () => {
        setCookie(secondCookieName, value2, ttl);
        const got = getCookie(secondCookieName);
        expect(got).toEqual(value2);
    });
    test("set and get no cookie", () => {
        const got = getCookie(thirdCookieName);
        expect(got).toEqual("");
    });
});

describe("set and clearCookie()", () => {
    const cookieName = "test-cookie";
    const value = "test123";
    test("set and clear first cookie", () => {
        setCookie(cookieName, value, ttl);
        clearCookie(cookieName);
        const got = getCookie(cookieName);
        expect(got).toEqual("");
    });
});

