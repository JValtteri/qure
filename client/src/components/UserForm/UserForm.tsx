import { useSignals } from '@preact/signals-react/runtime';
import Frame from '../common/Frame/Frame';
import './UserForm.css';

import { Signal } from "@preact/signals-react";
import { useState } from 'react';
import Dialog from '../common/Dialog/Dialog';


interface Props {
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
    show: Signal<{ "selectedEventId": number, "eventID": number, "editor": boolean, "account": boolean}>;

}

function UserForm({user, show}: Props) {
    useSignals();

    const [mode, setMode] = useState(0);
    const [password, setPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [newPassword2, setNewPassword2] = useState("");
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);

    const handleClose = () => {
        show.value = {"selectedEventId": -1, "eventID": -1, "editor": false, "account": false}
    }

    const handleDeleteSelf = () => {
        console.warn("Not implemented: Delete self");
    }

    return (
        <Frame className="details" hidden={!show.value.account}>
            <div className="header-container">
                <h3>{ user.value.username }</h3>
                <button onClick={handleClose} >Close</button>
            </div>

            <div id='tabs' className='grid'>
                <button onClick={()=> setMode(0)} className={mode==0 ? 'selected' : ''}>
                    <input type='checkbox' checked={mode==0} readOnly></input> Reservations
                </button>
                <button onClick={()=> setMode(1)} className={mode==1 ? 'selected' : ''}>
                    <input type='checkbox' checked={mode==1} readOnly></input> Edit Account
                </button>
            </div>

            <div hidden={mode != 0}>
                <label>Reservations</label>
                <p>aa</p>
                <p>aa</p>
                <p>aa</p>
            </div>

            <div id='account-editor' className='grid' hidden={mode != 1}>
                <label id='password-label' htmlFor="password">Current password:</label>
                <input
                    type="password"
                    id="password"
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
                    id="password-confirm"
                    value={newPassword2}
                    onChange={e => setNewPassword2(e.target.value)}
                />
                <button id={"delete-account"} className="delete-account" onClick={ ()=>setShowDeleteDialog(true) }>Delete Account</button>
                <button id={"apply-button"} className='selected'>Apply</button>
            </div>

            <Dialog hidden={!showDeleteDialog} className='error'>
                <div>
                    <h2 className='delete-dialog'>Deleting Account: <i>"{user.value.username.split('@')[0].toUpperCase()}"</i></h2>
                    <p className='delete-dialog'>Are you sure you want to delete your account?</p>
                    <p className='delete-dialog'><b>This action is not reversible!</b></p>
                </div>
                <div className='grid delete-buttons'>
                    <button className="delete-account" onClick={ handleDeleteSelf }>Confirm Delete Account</button>
                    <button onClick={ ()=>setShowDeleteDialog(false) }>Cancel</button>
                </div>
            </Dialog>

        </Frame>
    )
}

export default UserForm;
