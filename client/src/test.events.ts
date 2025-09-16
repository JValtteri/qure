import {describe, expect, test} from '@jest/globals';
import type {Event} from './events';
import {getEvent, getEvents, event} from './events';

const testEvent: Event = event;

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
