import { useState } from 'react';

import ConfirmDialog from '../common/ConfirmDialog/ConfirmDialog';


interface Props {
    hidden: boolean;
    onConfirm: (password: string)=>void;
    onCancel: ()=>void;
    userName: string;
    role: string;
}


function ConfirmRoleDialog({hidden, onConfirm, onCancel, userName, role}: Props) {
    userName = userName ? userName : ""
    const [password, setPassword] = useState("");

    return(
       <ConfirmDialog
            hidden={hidden}
            confirmBtnClass='selected'
            onConfirm={ () => {
                setPassword("");
                onConfirm(password);
            } }
            onCancel={ () => {
                setPassword("");
                onCancel() ;
            }}
        >
            <div>
                <h3 className='dialog-text'>
                    Changeing: '{userName}' to '{role}'
                </h3>
                <p className='dialog-text'>Are you sure you want to do this?</p>
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

export default ConfirmRoleDialog;
