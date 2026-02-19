import { useSignals } from "@preact/signals-react/runtime";
import Dialog from "../common/Dialog/Dialog";
import "./ReservationCard.css"

import type { ReservationResponse } from '../../api/api';
import { posixToDateAndTime, isPast } from '../../utils/utils';


interface Props {
    reservation: ReservationResponse;
    className?: string;
    email: string
    show: boolean;
    onHide: ()=>void;
    onEdit: (email: string, size: number, eventID: string, timeslot: number)=>void;
    onCancel: (email: string, eventID: string, timeslot: number)=>void;
}

function ReservationCard({reservation, className, email, show, onHide, onEdit, onCancel}: Props) {
    useSignals();
    let inPast = isPast(reservation.Timeslot);
    return (
        <Dialog hidden={!show} className={className}>
            <div className="main-text">
                <div>
                    <img src={ './logo.png' } fetchPriority='low' />
                    <h2 className='centered'>RESERVATION</h2>
                    <pre className='centered'>#{reservation.Id}</pre>
                    <hr className='thin-hr'></hr>
                    <h3 className='centered'>{reservation.Event ? reservation.Event.Name : ""}</h3>
                    <p className='centered'>{posixToDateAndTime(reservation.Timeslot)}</p>
                    <p className='centered low-profile-label'>Group size:</p>
                    <p className='centered big-number'> <b>{reservation.Confirmed}</b></p>
                </div>
            </div>
            <div className="buttons-center">
                <button
                    className="centered-button"
                    id="ok"
                    onClick={ () => onHide() }>
                        Close
                </button>
                <button
                    hidden={inPast}
                    className="centered-button"
                    onClick={ () => onEdit(email, reservation.Size, reservation.EventID, reservation.Timeslot) }>
                        Edit
                </button>
                <button
                    hidden={inPast}
                    className="centered-button red-button"
                    onClick={ () => onCancel(email, reservation.EventID, reservation.Timeslot) }>
                        Delete
                </button>
            </div>
        </Dialog>
    )
}

export default ReservationCard

