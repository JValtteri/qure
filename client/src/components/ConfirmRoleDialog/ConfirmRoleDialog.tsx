import { useState } from 'react';

import { useTranslation } from '../../context/TranslationContext';

import ConfirmDialog from '../common/ConfirmDialog/ConfirmDialog';


interface Props {
    hidden: boolean;
    onConfirm: (password: string)=>void;
    onCancel: ()=>void;
    userName: string;
    role: string;
}


function ConfirmRoleDialog({hidden, onConfirm, onCancel, userName, role}: Props) {
    const { t } = useTranslation();
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
                    {t("warning.changeing role", {name: userName, role: role})}
                </h3>
                <p className='dialog-text'>{t("warning.are you sure")}</p>
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

export default ConfirmRoleDialog;
