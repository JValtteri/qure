import './ReservationForm.css';

import { useEffect, useState } from 'react';
import { signal, Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import { makeReservation, type Timeslot } from '../../api/api';
import { validEmail } from '../../utils/utils';

import Dialog from '../common/Dialog/Dialog';
import TimeslotList from '../TimeslotList/TimeslotList';
import Popup from '../Popup/Popup';
import ReserveFailed from '../Popup/templates/ReserveFailed/ReserveFailed';
import ReservePartial from '../Popup/templates/ReservePartial/ReservePartial';
import ReserveSuccess from '../Popup/templates/ReserveSuccess/ReserveSuccess';
import PolicyAccept from '../PolicyAccept/PolicyAccept';
import PrivacyPolicy from '../PrivacyPolicy/PrivacyPolicy';
import { resumeSession } from '../common/utils';
import { useTranslation } from '../../context/TranslationContext';
import MarkdownRenderer from '../MarkdownRenderer/MarkdownRenderer';


const selectedSlot = signal(-1);

interface Props {
    showDialog:       Signal<boolean>;
    eventID:          string;
    timeslots:        Map<number, Timeslot>;
    requestedUpdate:  Signal<boolean>;
    user:             Signal<{username: string, loggedIn: boolean, role: string}>;
}

function ReservationForm({showDialog, eventID, timeslots, requestedUpdate, user}: Props) {
    useSignals();
    const { t } = useTranslation();
    const [email, setEmail] = useState( user.value.loggedIn ? user.value.username : "");
    const [groupSize, setGroupSize] = useState(0);
    const [reservationConfirmationVisible, setReservationConfirmationVisible] = useState(false);
    const [reservationConfiramtion, setReservationConfiramtion] = useState<React.ReactNode>(null);
    const [policyAccepted, setPolicyAccepted] = useState(false);

    useEffect( () => {
        setEmail( user.value.loggedIn ? user.value.username : "");
    }, [user.value.username]);

    const requestUpdate = ()  => requestedUpdate.value = !requestedUpdate.value; // Request update of slot information

    const reserveHandler = async () => {
        requestUpdate();
        if ( policyAccepted === false && user.value.loggedIn === false ) {
            setReservationConfiramtion(<>
                {t("error.concent.accept")} <PrivacyPolicy /> {t("error.concent.again")}.
            </>);
            setReservationConfirmationVisible(true);
            return;
        } else if ( selectedSlot.value === -1 ) {
            setReservationConfiramtion(<>
                <MarkdownRenderer content={t("error.missing group")} />
            </>);
            setReservationConfirmationVisible(true);
            return;
        } else if ( groupSize < 1 ) {
            setReservationConfiramtion(<>
                <MarkdownRenderer content={t("error.missing size")} />
            </>);
            setReservationConfirmationVisible(true);
            return;
        } else if ( !validEmail(email) ) {
            setReservationConfiramtion(<>
                <MarkdownRenderer content={t("error.missing email")} />
            </>);
            setReservationConfirmationVisible(true);
            return;
        }

        let reservation = null;
        try {
            reservation = await makeReservation(email, groupSize, eventID, selectedSlot.value);
            if (reservation.Error != "" ) {
                setReservationConfiramtion(<ReserveFailed error={reservation.Error} />);
            } else if ( reservation.Confirmed < reservation.Size ) {
                setReservationConfiramtion(<ReservePartial
                    confirmed={reservation.Confirmed}
                    size={reservation.Size}
                    time={reservation.Timeslot}
                    code={reservation.Id}
                />);
                showDialog.value=false;
            } else {
                setReservationConfiramtion(<ReserveSuccess
                    size={reservation.Size}
                    time={reservation.Timeslot}
                    code={reservation.Id}
                />);
                showDialog.value=false;
            }
            if (!user.value.loggedIn) {
                resumeSession(undefined, undefined, user, undefined);
            }
        } catch (error: any) {
            setReservationConfiramtion(<ReserveFailed error={error.message} />);
            console.warn(error.message);
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
                    <h3>{t("event.reservation")}</h3>
                    <button onClick={ handleClose }>{t("common.cancel")}</button>
                </div>
                <div></div>

                <label>{t("event.select timeslot")}</label>
                <TimeslotList timeslots={timeslots} selectedSlot={selectedSlot} requestUpdate={requestedUpdate} />

                <label className="form-label" htmlFor="reserve-email">{t("user.email")}</label>
                <input
                    id="reserve-email"
                    type="email"
                    value={email}
                    placeholder='example@email.com'
                    onChange={e => setEmail(e.target.value)}
                    required
                    disabled={user.value.loggedIn}
                />
                <label className="form-label" htmlFor="group-size">{t("event.reservation size")}</label>
                <input
                    id="group-size"
                    type="number"
                    value={groupSize}
                    min={1}
                    onChange={e => setGroupSize(Number(e.target.value))}
                    required
                />
                <PolicyAccept className="form-label" hidden={user.value.loggedIn} onChange={setPolicyAccepted} />

                <hr></hr>
                <div className='buttons 2'>
                    <button className='selected' onClick={()=>{
                        reserveHandler();
                        requestUpdate();
                        }}
                    >{t("event.reserve")}</button>
                </div>
            </Dialog>
            <Popup show={reservationConfirmationVisible} onHide={() => setReservationConfirmationVisible(false)}>
                {reservationConfiramtion}
            </Popup>
        </>
    );
}

export default ReservationForm;
