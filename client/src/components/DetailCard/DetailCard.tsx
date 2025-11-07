import type { ReactNode } from "react";
import { Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";
import Frame from "../common/Frame/Frame";
import { getEvent } from "../../utils/events";
import './DetailCard.css';


interface Props {
    selectedId: Signal<number>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
    children: ReactNode;
}

function DetailCard( {selectedId, user, children}: Props ) {
    useSignals();
    console.log("Detail rendered");

    const event = getEvent(selectedId.value);
    return (
        <Frame className="details" hidden={selectedId.value === -1}>
            <div className="header-container">
                <h3>{ event.name }</h3>
                <button onClick={ () => selectedId.value = -1 }>Close</button>
                <div className="detail-time">
                    <div>
                        Start:
                    </div>
                    <div>
                        { event.dtStart }
                    </div>
                    <div>
                        End:
                    </div>
                    <div>
                        { event.dtEnd }
                    </div>
                </div>
            </div>
            <div>
                {children}
            </div>
            <hr></hr>
            <div className="detail-footer">
                <button>Reserve</button>
                <div className="footer-text">
                    Slots:
                </div>
                <div className="footer-text">
                    { event.guests } / { event.guestSlots }
                </div>
            </div>
            <div className={`detail-footer`} hidden={!user.value.admin}>
                <button>Reserve</button>
                <div className="footer-text">
                    Guides:
                </div>
                <div className="footer-text">
                    { event.staff } / { event.staffSlots }
                </div>
            </div>
            <br></br>
        </Frame>
    )
}

export default DetailCard;
