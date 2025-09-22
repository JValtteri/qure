import { Signal } from '@preact/signals-react';
import Frame from '../Frame/Frame';
import './TitleBar.css';

interface Props {
    title?: string
    icon?: string
    role: Signal<string>
}

function TitleBar({title, icon, role}: Props) {
    return (
        <Frame className='title'>
            <img src={ icon ? icon : './logo.png' } />
            <span>
                {title ? title : "< Title >"}
            </span>
            <button onClick={ () => role.value = "test" }>Login</button>
        </Frame>
    )
}

export default TitleBar;
