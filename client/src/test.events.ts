import {describe, expect, test} from '@jest/globals';
import type {Event} from './events';
import {getEvent, getEvents} from './events';

const testEvent: Event = {
    name:             "Cool",
    shortDescription: "~",
    longDescription:  "",
    dtStart:          "123",
    dtEnd:            "456",
    staffSlots:       7,
    staff:            6,
    guestSlots:       12,
    guests:           10,
}

describe("Event tests", () => {
    test("getEvent()", () => {
        console.log(`Got: ${getEvent().toString()}`);
        expect(getEvent()).toStrictEqual(testEvent);
    });
    test("getEvents()", () => {
        console.log(`Got: ${getEvents()[0]}`);
        expect(getEvents()[0]).toStrictEqual(testEvent);
    });
});
