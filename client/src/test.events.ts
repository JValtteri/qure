import {describe, expect, test} from '@jest/globals';
import type {Event} from './events';
import {getEvent, getEvents, event} from './events';

const testEvent: Event = event;

describe("getEvent()", () => {
    test("Event matches", () => {
        console.log(`Got: ${getEvent(0).toString()}`);
        expect(getEvent(0)).toStrictEqual(testEvent);
    });
});
describe("getEvents()", () => {
    test("First event matches", () => {
        console.log(`Got: ${getEvents()[0]}`);
        expect(getEvents()[0]).toStrictEqual(testEvent);
    });
    test("Third event's name matches", () => {
        console.log(`Got: ${getEvents()[0]}`);
        expect(getEvents()[2].name).toStrictEqual("Verear virtute");
    });
});