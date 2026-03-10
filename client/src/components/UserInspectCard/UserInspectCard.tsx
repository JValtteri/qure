import { useSignals } from "@preact/signals-react/runtime";
import Dialog from "../common/Dialog/Dialog";
import "./UserInspectCard.css"

import type { ClientResponse } from '../../api/api';
import { posixToDateAndTime } from '../../utils/utils';


interface Props {
    client: ClientResponse;
    className?: string;
    hidden: boolean;
    onDelete: (id: string)=>void
    onClose: ()=>void;
}

function UserInspectCard({client, className, hidden, onDelete, onClose}: Props) {
    useSignals();

    return (
        <Dialog hidden={hidden} className={className}>
            <div className="main-text">
                <div>
                    <img src={ './logo.png' } fetchPriority='low' />
                    <h2 className='centered'>{client.Email}</h2>
                    <label className='centered low-profile-label'>{client.Id}</label>

                    <p className='centered'>
                        <select className='narrow-input' value={client.Role} disabled>
                            <option value="admin">Admin</option>
                            <option value="guest">Guest</option>
                            <option value="other">...</option>
                        </select>
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
                    }}>
                        Ok
                </button>
                <button
                    className="centered-button red-button"
                    onClick={ () => onDelete(client.Id) }>
                        Delete
                </button>
            </div>
        </Dialog>
    )
}

export default UserInspectCard
