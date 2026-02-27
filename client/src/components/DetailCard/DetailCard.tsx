import './DetailCard.css';

import { useState, useEffect } from "react";
import { signal, Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import { deleteEvent, type EventResponse, type Timeslot } from "../../api/api";
import { countSlots, posixToDateAndTime } from '../../utils/utils';
import { loadDetails } from '../common/utils';

import Frame from "../common/Frame/Frame";
import ReservationForm from '../ReservationForm/ReservationForm';
import ConfirmDialog from '../common/ConfirmDialog/ConfirmDialog';


const showReservationDialog = signal(false);
const loadingEvents = signal(false);

interface Props {
    show: Signal<{"eventID": string, "editor": boolean, "account": boolean, "inspect": boolean}>;
    user: Signal<{username: string, loggedIn: boolean, role: string}>;
    requestedUpdate: Signal<boolean>;
}

function DetailCard( {show, user, requestedUpdate}: Props ) {
    useSignals();
    const [eventDetails, setEventDetails]         = useState({} as EventResponse)
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);

    const handleClose = () => show.value={"eventID": "none", "editor": false, "account": show.value.account, "inspect": false};

    const loadDetailsHandler = loadDetails(show, loadingEvents, setEventDetails);

    const handleDeleteEvent = () => {
        try {
            deleteEvent(show.value.eventID)
                .then( () => {
                    handleClose();
                });
        } catch (error: any) {
            console.warn(error.message);
        }
    }

    useEffect(() => {
        loadDetailsHandler();
    }, [show.value.eventID, requestedUpdate.value]);

    let totalSlots = 0;
    let totalReservedSlots = 0;
    let timeslots = new Map<number, Timeslot>();
    try {
        timeslots = new Map(Object.entries(eventDetails.Timeslots).map(([k, v]) => [Number(k), v]));
        const slots = countSlots(timeslots);
        totalSlots = slots.totalSlots;
        totalReservedSlots = slots.totalReservedSlots;
    } catch {
        // Ignore errors
    }

    return (
        <Frame
            className={ eventDetails.Draft ? "details yellow" : "details" }
            hidden={show.value.eventID === "none" || show.value.editor || show.value.inspect}
        >
            <div className={"header-container"}>
                <h3>{ `${eventDetails.Name} ${eventDetails.Draft ? "- (Draft)" : ""} ` }</h3>
                <button onClick={ handleClose }>Close</button>
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
            <div id="detail-description">
                {eventDetails.LongDescription}
            </div>
            <hr />
            <div className="detail-footer">
                <button className='selected' onClick={ () => showReservationDialog.value=true } >Reserve</button>
                <div className="footer-text">
                    Slots:
                </div>
                <div className="footer-text">
                    { totalReservedSlots } / { (totalSlots) }
                </div>
            </div>
            <div className={`detail-footer`} hidden={user.value.role != "admin"}>
                <button>Reserve</button>
                <div className="footer-text">
                    Guides:
                </div>
                <div className="footer-text">
                    { eventDetails.Staff } / { eventDetails.StaffSlots }
                </div>
            </div>
            <hr hidden={user.value.role != "admin"} />
            <div className="buttons" hidden={user.value.role != "admin"}>
                <button onClick={ () => show.value={
                    "eventID": show.value.eventID,
                    "editor": true,
                    "account": show.value.account,
                    "inspect": false
                }}>Edit Event</button>
                <button onClick={ () => setShowDeleteDialog(true) } className="red-button" >Delete Event</button>
            </div>
            <br></br>
            <ReservationForm
                showDialog={showReservationDialog}
                eventID={show.value.eventID}
                timeslots={timeslots}
                requestedUpdate={requestedUpdate} user={user}
            />

            <ConfirmDialog
                    hidden={!showDeleteDialog}
                    className='error'
                    confirmBtnName="Confirm Delete Event"
                    confirmBtnClass='red-button'
                    onConfirm={ handleDeleteEvent }
                    onCancel={ ()=>setShowDeleteDialog(false) }
                >
                    <div>
                        <h2 className='delete-dialog'>Deleting Event: <i>"{eventDetails.Name}"</i></h2>
                        {show.value.eventID}
                        <p className='delete-dialog'>Are you sure you want to delete the event?</p>
                        <p className='delete-dialog'><b>This action is not reversible!</b></p>
                    </div>
                </ConfirmDialog>

        </Frame>
    )
}

export default DetailCard;

