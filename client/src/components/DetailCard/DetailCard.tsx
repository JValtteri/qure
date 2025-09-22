import type { ReactNode } from "react";
import { Signal } from "@preact/signals-react";
import Frame from "../Frame/Frame";
import './DetailCard.css';
import { getEvent } from "../../events";

interface Props {
    selectedId: Signal<number>;
    role: Signal<string>;
    children: ReactNode
}

function DetailCard( {selectedId, role, children}: Props ) {
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
            <div className={`detail-footer ${ role.value === "guest" && "hidden"}`}>
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
