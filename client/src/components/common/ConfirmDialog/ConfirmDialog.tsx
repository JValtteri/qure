import './ConfirmDialog.css'

import type { ReactNode } from "react";
import Dialog from "../Dialog/Dialog";


interface Props {
    hidden:             boolean,
    className?:         string,
    children:           ReactNode,
    onCancel:           ()=>void,
    onConfirm:          ()=>void,
    confirmBtnName?:    string,
    confirmBtnClass?:   string
}

function ConfirmDialog({hidden, className, children, onCancel, onConfirm, confirmBtnName, confirmBtnClass}: Props) {
    return (
        <Dialog hidden={hidden} className={className}>
            {children}
            <div className='grid confirm-buttons'>
                <button className={confirmBtnClass} onClick={ onConfirm }>{confirmBtnName ? confirmBtnName : "Confirm" }</button>
                <button onClick={ onCancel }>Cancel</button>
            </div>
        </Dialog>
    )
}

export default ConfirmDialog;
