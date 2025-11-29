
export class Reservation {
  #clientId: string;
  #email: string;
  id: number;
  eventID: number;
  slotID: number;
  guests: number;
  discountedGuests: number;

  constructor(eventID: number, email: string, guests: number, discountedGuests: number) {
    this.id = -1;
    this.eventID = eventID;
    this.slotID = -1;
    this.#clientId = "";
    this.#email = email;
    this.guests = guests;
    this.discountedGuests = discountedGuests;
    if ( guests < 1 && discountedGuests < 1 ) {
        throw new Error("Reservation must have at least one guest");
    }
    // TODO: Chack for empty reservation
  }

  getEmail(authorisation: string) {
    if (authorisation) {    // TODO: Create a check for authorisation
        return this.#email;
    }
    return "";
  }

  resetEmail(email: string) {
    this.#email = email;
  }

  cancel(clientId: string) {
    if (clientId === this.#clientId ) {
        this.resetEmail("");
        this.id = -1;
        this.eventID = -1;
        this.#clientId = "";
        this.guests = 0;
        this.discountedGuests = 0;
        return true;
    }
    return false;
  }
}

export function makeReservation(eventID: number,
                                email: string,
                                guests: number,
                                discountedGuests: number
                            ): Reservation {
    return new Reservation(eventID, email, guests, discountedGuests);
}
