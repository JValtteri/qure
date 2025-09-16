import type { MouseEventHandler, ReactNode } from 'react';
import './Frame.css';

interface Props {
    children?: ReactNode;
    className?: string;
    reactive?: boolean;
    onClick?: MouseEventHandler;
    hidden?: boolean;
}

function Frame( {children, reactive, className, onClick, hidden}: Props) {
    const frameType = reactive ? "reactive" : "frame";
    return (
        <div className={className ? className : frameType} onClick={onClick ?? onClick} hidden={hidden}>
            {children}
        </div>
    )
}

export default Frame
