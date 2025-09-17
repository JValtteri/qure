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
    const additionalClasses = className==undefined ? "" : className;
    return (
        <div className={`${frameType} ${additionalClasses}`} onClick={onClick ?? onClick} hidden={hidden}>
            {children}
        </div>
    )
}

export default Frame
