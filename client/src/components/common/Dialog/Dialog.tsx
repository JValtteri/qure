import Frame from '../Frame/Frame';
import './Dialog.css';
import type { ReactNode } from 'react';


interface Props {
    children?: ReactNode;
    className?: string;
    hidden?: boolean;
}

function Dialog({children, className, hidden}: Props) {

  return(
    <Frame className={className==undefined ? "dialog" : `dialog ${className}`} hidden={ hidden }>
        {children}
    </Frame>
  )
}

export default Dialog;
