import { useState } from 'react';

import ConfirmDialog from '../common/ConfirmDialog/ConfirmDialog';


interface Props {
    hidden: boolean;
    onConfirmDelete: (password: string)=>void;
    onCancel: ()=>void;
    userName: string;
}


function ConfirmDeleteDialog({hidden, onConfirmDelete, onCancel, userName}: Props) {
    userName = userName ? userName : ""
    const [password, setPassword] = useState("");

    return(
       <ConfirmDialog
            hidden={hidden}
            className='error'
            confirmBtnName="Confirm Delete Account"
            confirmBtnClass='red-button'
            onConfirm={ () => {
                setPassword("");
                onConfirmDelete(password);
            } }
            onCancel={ onCancel }
        >
            <div>
                <h2 className='delete-dialog'>
                    Deleting Account: <i>"{userName.split('@')[0].toUpperCase()}"</i>
                </h2>
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
    )
}

export default ConfirmDeleteDialog;
