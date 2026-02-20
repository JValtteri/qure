import { useEffect, useState } from 'react';
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
    onEdit: (reservationID: string, email: string, size: number, eventID: string, timeslot: number)=>void;
    onCancel: (reservationID: string, email: string, eventID: string, timeslot: number)=>void;
}

function ReservationCard({reservation, className, email, show, onHide, onEdit, onCancel}: Props) {
    useSignals();

    const [size, setSize] = useState(reservation.Size);
    const [editing, setEditing] = useState(false);


    useEffect(() => {
        setSize(reservation.Size);
    }, [show]);

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
                    <p className='centered big-number'>
                        <b hidden={editing}>{reservation.Confirmed}</b>
                        <input
                            hidden={!editing}
                            className='small-input'
                            type="number"
                            value={size}
                            min={1}
                            step={1}
                            onChange={e => setSize(Number(e.target.value))}
                        />
                    </p>
                </div>
            </div>
            <div className="buttons-center">
                <button
                    className="centered-button"
                    id="ok"
                    onClick={ () => {
                        setEditing(false)
                        onHide()
                        }
                    }>
                        { editing ? "Back" : "Ok" }
                </button>
                <button
                    hidden={inPast || editing}
                    className="centered-button"
                    onClick={ () => setEditing(true) }>
                        Edit
                </button>
                <button
                    hidden={!editing}
                    className="centered-button"
                    onClick={ () => {
                        setEditing(false)
                        onEdit(reservation.Id, email, reservation.Size, reservation.EventID, reservation.Timeslot)
                    }
                }>
                        Update
                </button>
                <button
                    hidden={!editing}
                    className="centered-button red-button"
                    onClick={ () => onCancel(reservation.Id, email, reservation.EventID, reservation.Timeslot) }>
                        Delete
                </button>
            </div>
        </Dialog>
    )
}

export default ReservationCard

