import "./UserInspectCard.css"

import { useEffect, useState } from "react";
import { useSignals } from "@preact/signals-react/runtime";

import { type ClientResponse } from '../../api/api';
import { posixToDateAndTime } from '../../utils/utils';

import Dialog from "../common/Dialog/Dialog";


interface Props {
    client: ClientResponse;
    className?: string;
    hidden: boolean;
    onDelete: (id: string)=>void
    onRoleChange: (role: string)=>void
    onClose: ()=>void;
}

function UserInspectCard({client, className, hidden, onDelete, onRoleChange, onClose}: Props) {
    useSignals();

    const [editMode, setEditMode] = useState(false);
    const [role, setRole] = useState(client.Role);

    useEffect( ()=> {
        if (!hidden) {
            setRole(client.Role)
        }
    }, [client, hidden])

    const handleRoleChange = (role: string) => {
        if (!editMode) return;
        onRoleChange(role);
        setEditMode(false);
        setRole(client.Role);
    }

    return (
        <Dialog hidden={hidden} className={className}>
            <div className="main-text">
                <div>
                    <img src={ './logo.png' } fetchPriority='low' />
                    <h2 className='centered'>{client.Email}</h2>
                    <label className='centered low-profile-label'>{client.Id}</label>

                    <p className='centered'>
                        <select
                            id="roleSelector"
                            className='narrow-input'
                            value={role}
                            onChange={e => {
                                setRole(e.target.value);
                                handleRoleChange(e.target.value);
                            }}
                            disabled={!editMode}
                        >
                            <option value="admin">Admin</option>
                            <option value="guest">Guest</option>
                            <option value="other">...</option>
                        </select>
                        <button id="lock-button" onClick={ () => setEditMode(true) } disabled={editMode}>✏️</button>
                    </p>

                    <p className='centered low-profile-label'>Created</p>
                    <p className='centered'>{posixToDateAndTime(client.CreatedDt)}</p>

                    <label className='centered low-profile-label' hidden={!client.IsTemporary}>Expires</label>
                    <p className='centered' hidden={!client.IsTemporary}>{posixToDateAndTime(client.ExpiresDt)}</p>
                </div>
            </div>
            <div className="buttons-center">
                <button
                    className="centered-button"
                    id="ok"
                    onClick={ () => {
                        onClose()
                        setEditMode(false);
                    }}>
                        Ok
                </button>
                <button
                    className="centered-button red-button"
                    onClick={ () => {
                        onDelete(client.Id)
                        setEditMode(false);
                    } }>
                        Delete
                </button>
            </div>
        </Dialog>
    )
}

export default UserInspectCard
