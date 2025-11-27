import { useState } from 'react';
import { signal, Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import Dialog from '../common/Dialog/Dialog';
import TimeslotList from '../TimeslotList/TimeslotList';
import Popup from '../Popup/Popup';
import ReserveFailed from '../Popup/templates/ReserveFailed/ReserveFailed';
import ReservePartial from '../Popup/templates/ReservePartial/ReservePartial';

import { makeReservation, type Timeslot } from '../../api/api';
import { validEmail } from '../../utils/utils';

import './ReservationForm.css';
import ReserveSuccess from '../Popup/templates/ReserveSuccess/ReserveSuccess';


const selectedSlot = signal(-1);

interface Props {
    showDialog: Signal<boolean>;
    eventID:    number;
    timeslots:  Map<number, Timeslot>;
}

function ReservationForm({showDialog, eventID, timeslots}: Props) {
    useSignals();
    const [email, setEmail] = useState("");
    const [groupSize, setGroupSize] = useState(0);
    const [reservationConfirmationVisible, setReservationConfirmationVisible] = useState(false);
    const [reservationConfiramtion, setReservationConfiramtion] = useState(<></>);

    const reserveHandler = async () => {
        if ( selectedSlot.value === -1 ) {
            setReservationConfiramtion(<>Please select a group <b>timeslot</b> and try again.</>);
            setReservationConfirmationVisible(true);
            return;
        } else if ( groupSize < 1 ) {
            setReservationConfiramtion(<>Please select a <b>size</b> for the group and try again.</>);
            setReservationConfirmationVisible(true);
            return;
        } else if ( !validEmail(email) ) {
            setReservationConfiramtion(<>Please enter a valid <b>email</b> and try again.</>);
            setReservationConfirmationVisible(true);
            return;
        }
        let reservation = await makeReservation(email, groupSize, eventID, selectedSlot.value);
        if (reservation.Error != "" ) {
            setReservationConfiramtion(ReserveFailed({
                error: reservation.Error
            }));
        } else if ( reservation.Confirmed < reservation.Size ) {
            setReservationConfiramtion(ReservePartial({
                confirmed: reservation.Confirmed,
                size: reservation.Size,
                time: reservation.Timeslot
            }));
            showDialog.value=false;
        } else {
            setReservationConfiramtion(ReserveSuccess({
                size: reservation.Size,
                time: reservation.Timeslot
            }));
            showDialog.value=false;
        }
        setReservationConfirmationVisible(true);
    };

    return(
        <>
            <Dialog className='reservation' hidden={ showDialog.value===false }>
                <div className="header-container">
                    <h3>Reservation</h3>
                    <button onClick={ () => showDialog.value=false }>Cancel</button>
                </div>
                <div></div>

                <label>Select Group</label>
                <div>
                    <TimeslotList timeslots={timeslots} selectedSlot={selectedSlot} />
                </div>

                <label className="form-label" htmlFor="reserve-email">Email</label>
                <input id="reserve-email" type="email" value={email} placeholder='example@email.com' onChange={e => setEmail(e.target.value)} required ></input>

                <label className="form-label" htmlFor="group-size">Group Size</label>
                <input id="group-size" type="number" value={groupSize} min={1} placeholder='example@email.com' onChange={e => setGroupSize(Number(e.target.value))} required ></input>
                <hr></hr>
                <div className='buttons 2'>
                    <button className='selected' onClick={ ()=>reserveHandler() }>Reserve</button>
                </div>
            </Dialog>
            <Popup show={reservationConfirmationVisible} onHide={() => setReservationConfirmationVisible(false)}>
                {reservationConfiramtion}
            </Popup>
        </>
    );
}

export default ReservationForm;
