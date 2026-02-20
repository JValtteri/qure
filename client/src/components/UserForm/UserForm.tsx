import './UserForm.css';

import { useEffect, useState } from 'react';
import { signal, Signal } from "@preact/signals-react";
import { useSignals } from '@preact/signals-react/runtime';

import Frame from '../common/Frame/Frame';
import Popup from '../Popup/Popup';
import ReservationsList from '../ReservationsList/ReservationsList';
import ReservationCard from '../ReservationCard/ReservationCard';

import { deleteUser, editPassword, listReservations, amendReservation, cancelReservation } from '../../api/api';
import type { ReservationResponse } from '../../api/api';
import ConfirmDialog from '../common/ConfirmDialog/ConfirmDialog';
import { Reservation } from '../../utils/reservations';


const selectedReservation = signal("none");

interface Props {
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
    show: Signal<{"eventID": string, "editor": boolean, "account": boolean}>;

}

function UserForm({user, show}: Props) {
    useSignals();

    const [mode, setMode] = useState(0);
    const [password, setPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [newPassword2, setNewPassword2] = useState("");
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [showPopup, setShowPopup] = useState(false);
    const [popupMessage, setPopupMessage] = useState("none");

    const newPasswordField = document.getElementById("new-password");
    const newPasswordField2 = document.getElementById("new-password2");
    const currentPasswordField = document.getElementById("current-password");

    const [reservations, setReservations] = useState(new Array<ReservationResponse>());


    const updateReservationsHandler = updateReservations(setReservations);

    useEffect(() => {
        updateReservationsHandler();
    }, [show.value]);


    const removeHighlights = () => {
        currentPasswordField?.classList.remove("wrong");
        newPasswordField?.classList.remove("wrong");
        newPasswordField2?.classList.remove("wrong");
    }

    const handleClose = () => {
        setPassword("");
        setNewPassword("");
        setNewPassword2("");
        removeHighlights();
        show.value = {"eventID": "none", "editor": false, "account": false}
    }

    const handleDeleteSelf = async () => {
        let resp = await deleteUser(user.value.username, password);
        if (resp.Success) {
            setPopupMessage("Success");
            user.value = { username: "", loggedIn: false, admin: false};
            setNewPassword("");
            setNewPassword2("");
            removeHighlights();
        } else {
            setPopupMessage(`Error: ${resp.Error}`);
        }
        setPassword("");
        setShowDeleteDialog(false);
        setShowPopup(true);
    }

    const handleCloseConfirm = () => {
        setShowPopup(false);
        if (user.value.username == "") {
            handleClose();
        }
    }

    const handleCloseCard = () => {
        selectedReservation.value = "none";
        if (user.value.username == "") {
            handleClose();
        }
    }

    const amendReservationHandler = async (reservationID: string, email: string, size: number, eventID: string, timeslot: number) => {
        const reservation = await amendReservation(reservationID, email, size, eventID, timeslot);
            if (reservation.Error != "" ) {
                setPopupMessage(reservation.Error);
            } else {
                setPopupMessage("Success");
            }
            selectedReservation.value = "none";
            setShowPopup(true);
    }

    const cancelReservationHandler = async (reservationID: string, email: string, eventID: string, timeslot: number) => {
        const reservation = await cancelReservation(reservationID, email, eventID, timeslot);
            if (reservation.Error != "" ) {
                setPopupMessage(reservation.Error);
            } else {
                setPopupMessage("Success");
            }
            selectedReservation.value = "none";
            setShowPopup(true);
    }

    const handlePasswordChange = async () => {
        if (user.value.username == undefined) {
            return;
        }
        if (newPassword != newPassword2) {
            newPasswordField?.classList.add("wrong");
            newPasswordField2?.classList.add("wrong");
            setPassword("");
            return;
        }
        removeHighlights();
        let resp = await editPassword(user.value.username, password, newPassword);
        if (resp.Success) {
            setPassword("");
            setNewPassword("");
            setNewPassword2("");
            setPopupMessage("Success");
        } else {
            currentPasswordField?.classList.add("wrong");
            newPasswordField?.classList.add("wrong");
            newPasswordField2?.classList.add("wrong");
            setPopupMessage(`Error: ${resp.Error}`);
        }
        setShowPopup(true);
    }

    return (
        <Frame className="details" hidden={!show.value.account}>
            <div className="header-container">
                <h3>{ user.value.username }</h3>
                <button onClick={handleClose} >Close</button>
            </div>

            <div id='tabs' className='grid account-tab'>
                <button onClick={()=> setMode(0)} className={mode==0 ? 'selected' : ''}>
                    <input type='checkbox' checked={mode==0} readOnly></input> Reservations
                </button>
                <button onClick={()=> setMode(1)} className={mode==1 ? 'selected' : ''}>
                    <input type='checkbox' checked={mode==1} readOnly></input> Edit Account
                </button>
            </div>

            <div hidden={mode != 0}>
                <h3>Reservations</h3>
                <ReservationsList reservations={reservations} selected={selectedReservation} update={updateReservationsHandler} />
            </div>

            <div id='account-editor' className='grid account-tab' hidden={mode != 1}>
                <label id='password-label' htmlFor="password">Current password:</label>
                <input
                    type="password"
                    id="current-password"
                    className='password'
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                    required
                />
                <label id='new-pass-label' htmlFor="password">New password:</label>
                <input
                    type="password"
                    id="new-password"
                    className='password'
                    value={newPassword}
                    onChange={e => setNewPassword(e.target.value)}
                />
                <label id='password-confirm-label' htmlFor="password-confirm">Confirm password:</label>
                <input
                    type="password"
                    id="new-password2"
                    value={newPassword2}
                    onChange={e => setNewPassword2(e.target.value)}
                />
                <button id={"delete-account"} className="red-button" onClick={ ()=>setShowDeleteDialog(true) }>Delete Account</button>
                <button id={"apply-button"} className='selected' onClick={ handlePasswordChange }>Apply</button>
            </div>

            <ConfirmDialog
                hidden={!showDeleteDialog}
                className='error'
                confirmBtnName="Confirm Delete Account"
                confirmBtnClass='red-button'
                onConfirm={ handleDeleteSelf }
                onCancel={ ()=>setShowDeleteDialog(false) }
            >
                <div>
                    <h2 className='delete-dialog'>Deleting Account: <i>"{user.value.username.split('@')[0].toUpperCase()}"</i></h2>
                    <p className='delete-dialog'>Are you sure you want to delete your account?</p>
                    <p className='delete-dialog'><b>This action is not reversible!</b></p>
                    <input
                        type="password"
                        value={password}
                        placeholder='Password'
                        onChange={e => setPassword(e.target.value)}
                    />
                </div>
            </ConfirmDialog>

            <Popup children={ popupMessage } show={ showPopup } onHide={ handleCloseConfirm } />
            <ReservationCard
                reservation={renderReservationCard(reservations)}
                email={user.value.username}
                show={selectedReservation.value != "none"}
                onHide={ handleCloseCard }
                onEdit={ amendReservationHandler }
                onCancel={ cancelReservationHandler }
            />

        </Frame>
    )
}

export default UserForm;


function updateReservations(setReservations: React.Dispatch<React.SetStateAction<Array<ReservationResponse>>>): () => Promise<void> {
    return async () => {
        await listReservations()
            .then(value => {
                if (value != null) {
                    setReservations(value);
                }
            });
    };
}

function renderReservationCard(reservations: Array<ReservationResponse>): ReservationResponse {
    let reservation = {} as ReservationResponse;
    for (const res of reservations) {
        if (res.EventID == selectedReservation.value) {
            reservation = res;
            break;
        }
    }
    return reservation;
}
