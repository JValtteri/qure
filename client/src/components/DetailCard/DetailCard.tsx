import './DetailCard.css';

import { useState, useEffect } from "react";
import { Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import Frame from "../common/Frame/Frame";

import { fetchEvent, type EventResponse } from "../../api/api";
import { countSlots, posixToDateAndTime } from '../../utils/utils';


interface Props {
    show: Signal<{ "selectedEventId": number, "eventID": number, "editor": boolean}>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
}

function DetailCard( {show, user}: Props ) {
    useSignals();
    console.log("Detail rendered");
    const [eventDetails, setEventDetails] = useState({} as EventResponse)

    const loadDetailsHandler = loadDetails(show, setEventDetails);
    useEffect(() => {
        loadDetailsHandler();
    }, [show.value.eventID]);

    let totalSlots = 0;
    let totalReservedSlots = 0;

    try {
        let timeslots = new Map(Object.entries(eventDetails.Timeslots).map(([k, v]) => [Number(k), v]));
        let slots = countSlots(timeslots);
        totalSlots = slots.totalSlots;
        totalReservedSlots = slots.totalReservedSlots;
    } catch(error) {
        console.warn(error)
    }

    return (
        <Frame className="details" hidden={show.value.selectedEventId === -1}>
            <div className="header-container">
                <h3>{ eventDetails.Name }</h3>
                <button onClick={ () => show.value={"selectedEventId": -1, "eventID": -1, "editor": false} }>Close</button>
                <div className="detail-time">
                    <div>
                        Start:
                    </div>
                    <div>
                        { posixToDateAndTime(eventDetails.DtStart) }
                    </div>
                    <div>
                        End:
                    </div>
                    <div>
                        { posixToDateAndTime(eventDetails.DtEnd) }
                    </div>
                </div>
            </div>
            <div>
                {eventDetails.LongDescription}
            </div>
            <hr></hr>
            <div className="detail-footer">
                <button>Reserve</button>
                <div className="footer-text">
                    Slots:
                </div>
                <div className="footer-text">
                    { totalReservedSlots } / { (totalSlots) }
                </div>
            </div>
            <div className={`detail-footer`} hidden={!user.value.admin}>
                <button>Reserve</button>
                <div className="footer-text">
                    Guides:
                </div>
                <div className="footer-text">
                    { eventDetails.Staff } / { eventDetails.StaffSlots }
                </div>
            </div>
            <br></br>
        </Frame>
    )
}

export default DetailCard;


function loadDetails(show: Signal<{ selectedEventId: number; eventID: number; editor: boolean; }>, setEventDetails: any) {
    return async () => {
        let details = await fetchEvent(`${show.value.eventID}`);
        setEventDetails(details);
    };
}
