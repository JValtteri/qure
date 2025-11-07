/**
 * @jest-environment jsdom
 */

import {describe, expect, test} from '@jest/globals';
import {setCookie, getCookie, ttl} from './cookie';


beforeEach(() => {
    document.cookie = "";
});

afterEach(() => {
    document.cookie = "";
});

describe("set and getCookie()", () => {
    let firstCookieName = "test-cookie";
    let secondCookieName = "test-other-cookie";
    let thirdCookieName = "test-no-cookie";
    let value1 = "test123";
    let value2 = "test456";
    test("set and get first cookie", () => {
        setCookie(firstCookieName, value1, ttl);
        let got = getCookie(firstCookieName);
        expect(got).toEqual(value1);
    });
    test("set and get second cookie", () => {
        setCookie(secondCookieName, value2, ttl);
        let got = getCookie(secondCookieName);
        expect(got).toEqual(value2);
    });
    test("set and get no cookie", () => {
        let got = getCookie(thirdCookieName);
        expect(got).toEqual("");
    });
});
