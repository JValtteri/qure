import type { MouseEventHandler, ReactNode } from 'react';
import './Frame.css';

interface Props {
    children: ReactNode;
    className?: string;
    reactive?: boolean;
    onClick?: MouseEventHandler;
}

function Frame( {children, reactive, className, onClick}: Props) {
    const frameType = reactive ? "reactive" : "frame";
    return (
        <div className={className ? className : frameType} onClick={onClick ?? onClick}>
            {children}
        </div>
    )
}

export default Frame
