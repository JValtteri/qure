import './UserForm.css';

import { useEffect, useState, lazy, Suspense } from 'react';
import { signal, Signal } from "@preact/signals-react";
import { useSignals } from '@preact/signals-react/runtime';

import { deleteUser, editPassword, listReservations, amendReservation, cancelReservation } from '../../api/api';
import type { ReservationResponse } from '../../api/api';
import { useTranslation } from '../../context/TranslationContext';

import Frame from '../common/Frame/Frame';
import Popup from '../Popup/Popup';
import ReservationsList from '../ReservationsList/ReservationsList';
import ReservationCard from '../ReservationCard/ReservationCard';
import ConfirmDeleteDialog from '../ConfirmDeleteDialog/ConfirmDeleteDialog';

import Spinner from '../Spinner/Spinner';
const Inspector = lazy( () => import('../Inspector/Inspector'));
const UserListView = lazy( () => import('../UserListView/UserListView'));



const selectedReservation = signal("none");
const loadingEvents = signal(false);

interface Props {
    user: Signal<{username: string, loggedIn: boolean, role: string}>;
    show: Signal<{eventID: string, view: string}>;
}

function UserForm({user, show}: Props) {
    useSignals();
    const {t} = useTranslation();

    const [mode, setMode] = useState(0);
    const [adminMode, setAdminMode] = useState(1);
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
    }, [show.value, showPopup]);

    useEffect(()=> {
        if (mode == 2 && adminMode == 1) {
            handleEnableInspector();
        } else {
            handleTurnOffInspector();
        }
    }, [mode, adminMode])

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
        setAdminMode(1);
        show.value = {eventID: "none", view: ""}
    }

    const handleDeleteSelf = async (password: string) => {
        let resp = null;
        try {
            resp = await deleteUser(user.value.username, password);
            if (resp.Success) {
                setPopupMessage("Success");
                user.value = { username: "", loggedIn: false, role: ""};
                show.value = {eventID: "none", view: ""}
                setNewPassword("");
                setNewPassword2("");
                removeHighlights();
            } else {
                setPopupMessage(`Error: ${resp.Error}`);
            }
            setPassword("");
            setShowDeleteDialog(false);
        } catch (error: any) {
            setPopupMessage(`Error: ${error.message}`);
            console.warn(error.message);
        }
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

    const amendReservationHandler = async (
        reservationID: string,
        email: string,
        size: number,
        eventID: string,
    timeslot: number) => {
        try {
            const reservation = await amendReservation(reservationID, email, size, eventID, timeslot);
            if (reservation.Error != "" ) {
                setPopupMessage(reservation.Error);
            } else {
                setPopupMessage("Success");
            }
            selectedReservation.value = "none";
        } catch (error: any) {
            setPopupMessage(error.message);
            console.warn(error.message);
        }
        setShowPopup(true);
    }

    const cancelReservationHandler = async (reservationID: string, email: string, eventID: string, timeslot: number) => {
        try {
            const reservation = await cancelReservation(reservationID, email, eventID, timeslot);
            if (reservation.Error != "" ) {
                setPopupMessage(reservation.Error);
            } else {
                setPopupMessage("Success");
            }
            selectedReservation.value = "none";
        } catch (error: any) {
            setPopupMessage(error.message);
            console.warn(error.message);
        }
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
        let resp = null;
        try {
            resp = await editPassword(user.value.username, password, newPassword);
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
        } catch (error: any) {
            setPopupMessage(`Error: ${error.message}`);
            console.warn(error.message);
        }
        setShowPopup(true);
    }

    const handleEnableInspector = () => {
        show.value = {eventID: "none", view: "inspect"};
        setAdminMode(1);
    };

    function handleTurnOffInspector() {
        show.value = { eventID: "none", view: "account" };
    }

    const handleUserList = () => {
        show.value = {eventID: "none", view: "account"};
        setAdminMode(2);
    };

    const isVisible = (show: Signal) => {
        let visible = ["account", "inspect"].includes(show.value.view);
        return visible;
    }

    return (
        <Frame className="details, inner-frame" hidden={!isVisible(show)}>
            <div className="header-container">
                <h3>{ user.value.username }</h3>
                <button onClick={handleClose} >{t("common.close")}</button>
            </div>

            {/* Buttons */}

            <div id='top-tabs' className='tabs'>
                <button
                    onClick={()=> setMode(0) }
                    className={mode==0 ? 'selected' : ''}>
                    <input type='checkbox' checked={mode==0} readOnly></input>
                    {t("tools.reservations")}
                </button>
                <button
                    onClick={()=> setMode(1) }
                    className={mode==1 ? 'selected' : ''}>
                    <input type='checkbox' checked={mode==1} readOnly></input>
                    {t("tools.edit account")}
                </button>
                <button
                    onClick={()=> setMode(2) }
                    className={mode==2 ? 'selected' : ''}
                    hidden={user.value.role != "admin" && user.value.role != "staff"}>
                    <input type='checkbox' checked={mode==2} readOnly></input>
                    {t("tools.admin tools")}
                </button>
            </div>

            {/* Views */}

            <Frame hidden={mode == 2}>

                <div hidden={mode != 0}>
                    <h3>{t("tools.reservations")}</h3>
                    <ReservationsList
                        reservations={reservations}
                        selected={selectedReservation}
                        update={updateReservationsHandler} />
                </div>

                <div hidden={mode != 1} id='account-editor' className='grid account-tab'>
                    <label id='password-label' htmlFor="password">{t("user.current-password")}:</label>
                    <input
                        type="password"
                        id="current-password"
                        className='password'
                        value={password}
                        onChange={e => setPassword(e.target.value)}
                        required
                    />
                    <label id='new-pass-label' htmlFor="password">{t("user.new-password")}:</label>
                    <input
                        type="password"
                        id="new-password"
                        className='password'
                        value={newPassword}
                        onChange={e => setNewPassword(e.target.value)}
                    />
                    <label id='password-confirm-label' htmlFor="password-confirm">{t("user.confirm-password")}:</label>
                    <input
                        type="password"
                        id="new-password2"
                        value={newPassword2}
                        onChange={e => setNewPassword2(e.target.value)}
                    />
                    <button
                        id={"delete-account"}
                        className="red-button"
                        onClick={ ()=>setShowDeleteDialog(true) }
                        >{t("user.delete account")}</button>
                    <button
                        id={"apply-button"}
                        className='selected'
                        onClick={ handlePasswordChange }
                        >{t("common.apply")}</button>
                </div>
            </Frame>

            <Suspense fallback={<Spinner />}>
                {["admin", "staff"].includes(user.value.role) &&
                <Frame className='collapsed-frame' hidden={mode != 2}>
                    <div hidden={mode != 2} className='tabs'>
                        <button
                            id={"search-reservations"}
                            onClick={ handleEnableInspector }
                            className={adminMode==1 ? 'selected' : ''}
                            >{t("tools.reservations")}</button>
                        <button
                            id={"users-reservations"}
                            hidden={ user.value.role != "admin" }
                            onClick={ handleUserList }
                            className={adminMode==2 ? 'selected' : ''}
                            >{t("tools.all users")}</button>
                    </div>

                    <div hidden={adminMode != 1}>
                        <Inspector show={show} hidden={show.value.view != "inspect"} />
                    </div>

                    <div hidden={adminMode != 2}>
                        <UserListView active={adminMode === 2} setShowPopup={setShowPopup} setPopupMessage={setPopupMessage} />
                    </div>
                </Frame>
                }
            </Suspense>

            <ConfirmDeleteDialog
                hidden={!showDeleteDialog}
                userName={user.value.username}
                onConfirmDelete={handleDeleteSelf}
                onCancel={ ()=>setShowDeleteDialog(false) }
            />

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


function updateReservations(
    setReservations: React.Dispatch<React.SetStateAction<Array<ReservationResponse>>>
): () => Promise<void> {
    return async () => {
        if (loadingEvents.value == true) {
            return;
        }
        loadingEvents.value = true;
        try {
            await listReservations()
            .then(value => {
                if (value != null) {
                    setReservations(value);
                }
            });
        } catch (error: any) {
            console.warn(error.message);
        }
        loadingEvents.value = false;
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
