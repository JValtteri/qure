import Frame from '../Frame/Frame';
import './Dialog.css';
import type { ReactNode } from 'react';


interface Props {
    children?: ReactNode;
    hidden?: boolean;
}

function Dialog({children, hidden}: Props) {

  return(
    <Frame className='dialog' hidden={ hidden }>
        {children}
    </Frame>
  )
}

export default Dialog;
