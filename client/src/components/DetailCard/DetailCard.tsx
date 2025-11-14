import { useState, type ReactNode, useEffect } from "react";
import { Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";
import Frame from "../common/Frame/Frame";
import { fetchEvent, type EventResponse } from "../../api/api";
import './DetailCard.css';


interface Props {
    show: Signal<{ "selectedEventId": number, "eventID": number, "editor": boolean}>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
}

function DetailCard( {show, user}: Props ) {
    useSignals();
    console.log("Detail rendered");
    const [eventDetails, setEventDetails] = useState({} as EventResponse)

    const loadDetails = async () => {
        let details = await fetchEvent(`${show.value.eventID}`);
        setEventDetails(details);
    }

    useEffect(() => {
        loadDetails();
    }, [show.value.eventID]);

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
                        { eventDetails.DtStart }
                    </div>
                    <div>
                        End:
                    </div>
                    <div>
                        { eventDetails.DtEnd }
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
                    { 999 } / { (999) }
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
