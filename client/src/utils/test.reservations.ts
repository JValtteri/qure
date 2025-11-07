import { describe, expect, test } from '@jest/globals';
import { Reservation, makeReservation } from './reservations';

const exampleEmail = "example@example";
const reservation: Reservation = new Reservation(123, exampleEmail, 1, 1);

describe("makeReservation(123)", () => {
    test("make new", () => {
        const newReservation = makeReservation(123, exampleEmail, 1, 1)
        expect(newReservation).toStrictEqual(reservation);
    });
    test("make empty", () => {
        expect( () => makeReservation(123, exampleEmail, 0, 0)).toThrow();
    });
});

describe("Reservation.getEmail()", () => {
    test("authorized", () => {
        expect(reservation.getEmail("authorized"))
                .toEqual(exampleEmail);
    });
    test("not authorized", () => {
        expect(reservation.getEmail(""))
                .toStrictEqual("");
    });
});

describe("Reservation.resetEmail()", () => {
    const newReservation = makeReservation(123, exampleEmail, 1, 1);

    test("resetEmail('foo@bar')", () => {
        newReservation.resetEmail("foo@bar")
        expect(newReservation.getEmail("authorized"))
                .toEqual("foo@bar");
    });
    test("resetEmail('')", () => {
        newReservation.resetEmail("")
        expect(newReservation.getEmail("authorized"))
                .toStrictEqual("");
    });
});

describe("Reservation.cancel()", () => {
    test("correct client", () => {
        const newReservation = makeReservation(123, exampleEmail, 1, 1);
        expect(newReservation.cancel("")).toEqual(true);
        expect(newReservation.eventID).toEqual(-1);
    });
    test("invalid client", () => {
        const newReservation = makeReservation(123, exampleEmail, 1, 1);
        expect(newReservation.cancel("foo")).toEqual(false);
    });
});
