import { useEffect, useState } from 'react';
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
import PolicyAccept from '../PolicyAccept/PolicyAccept';
import PrivacyPolicy from '../PrivacyPolicy/PrivacyPolicy';


const selectedSlot = signal(-1);

interface Props {
    showDialog:       Signal<boolean>;
    eventID:          string;
    timeslots:        Map<number, Timeslot>;
    requestedUpdate:  Signal<boolean>;
    user:             Signal<{username: string, loggedIn: boolean, admin: boolean}>;
}

function ReservationForm({showDialog, eventID, timeslots, requestedUpdate, user}: Props) {
    useSignals();
    const [email, setEmail] = useState( user.value.loggedIn ? user.value.username : "");
    const [groupSize, setGroupSize] = useState(0);
    const [reservationConfirmationVisible, setReservationConfirmationVisible] = useState(false);
    const [reservationConfiramtion, setReservationConfiramtion] = useState(<></>);
    const [policyAccepted, setPolicyAccepted] = useState(false);

    useEffect( () => {
        setEmail( user.value.loggedIn ? user.value.username : "");
    }, [user.value.username]);

    const requestUpdate = ()  => requestedUpdate.value = !requestedUpdate.value; // Request update of slot information

    const reserveHandler = async () => {
        requestUpdate();
        if ( policyAccepted === false && user.value.loggedIn === false ) {
            setReservationConfiramtion(<>Please accept the <PrivacyPolicy /> and try again.</>);
            setReservationConfirmationVisible(true);
            return;
        } else if ( selectedSlot.value === -1 ) {
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
        const reservation = await makeReservation(email, groupSize, eventID, selectedSlot.value);
        if (reservation.Error != "" ) {
            setReservationConfiramtion(ReserveFailed({
                error: reservation.Error
            }));
        } else if ( reservation.Confirmed < reservation.Size ) {
            setReservationConfiramtion(ReservePartial({
                confirmed: reservation.Confirmed,
                size: reservation.Size,
                time: reservation.Timeslot,
                code: reservation.Id
            }));
            showDialog.value=false;
        } else {
            setReservationConfiramtion(ReserveSuccess({
                size: reservation.Size,
                time: reservation.Timeslot,
                code: reservation.Id
            }));
            showDialog.value=false;
        }
        setReservationConfirmationVisible(true);
    };

    const handleClose = () => {
        requestUpdate();
        showDialog.value=false;
    }

    return(
        <>
            <Dialog className='reservation' hidden={ showDialog.value===false }>
                <div className="header-container">
                    <h3>Reservation</h3>
                    <button onClick={ handleClose }>Cancel</button>
                </div>
                <div></div>

                <label>Select Timeslot</label>
                <TimeslotList timeslots={timeslots} selectedSlot={selectedSlot} requestUpdate={requestedUpdate} />

                <label className="form-label" htmlFor="reserve-email">Email</label>
                <input
                    id="reserve-email"
                    type="email"
                    value={email}
                    placeholder='example@email.com'
                    onChange={e => setEmail(e.target.value)}
                    required
                    disabled={user.value.loggedIn}
                />
                <label className="form-label" htmlFor="group-size">Reservation Size</label>
                <input
                    id="group-size"
                    type="number"
                    value={groupSize}
                    min={1}
                    placeholder='example@email.com'
                    onChange={e => setGroupSize(Number(e.target.value))}
                    required
                />
                <PolicyAccept className="form-label" hidden={user.value.loggedIn} onChange={setPolicyAccepted} />

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
