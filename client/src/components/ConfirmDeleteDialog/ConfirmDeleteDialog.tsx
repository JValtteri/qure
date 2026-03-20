import { useState } from 'react';

import { useTranslation } from '../../context/TranslationContext';

import ConfirmDialog from '../common/ConfirmDialog/ConfirmDialog';


interface Props {
    hidden: boolean;
    onConfirmDelete: (password: string)=>void;
    onCancel: ()=>void;
    userName: string;
}


function ConfirmDeleteDialog({hidden, onConfirmDelete, onCancel, userName}: Props) {
    userName = userName ? userName : ""
    const {t} = useTranslation();
    const [password, setPassword] = useState("");

    return(
       <ConfirmDialog
            hidden={hidden}
            className='error'
            confirmBtnName={t("user.confirm delete account")}
            confirmBtnClass='red-button'
            onConfirm={ () => {
                setPassword("");
                onConfirmDelete(password);
            } }
            onCancel={ () => {
                setPassword("");
                onCancel() ;
            }}
        >
            <div>
                <h2 className='dialog-text'>
                    {t("warning.confirm-delete-account", {name: userName.split('@')[0].toUpperCase()})}
                </h2>
                <p className='dialog-text'>{t("warning.confirm-delete-account", {name: userName})}</p>
                <p className='dialog-text'><b>{t("warning.no-takebacks")}</b></p>
                <input
                    type="password"
                    value={password}
                    placeholder={t("common.password")}
                    onChange={e => setPassword(e.target.value)}
                />
            </div>
        </ConfirmDialog>
    )
}

export default ConfirmDeleteDialog;
