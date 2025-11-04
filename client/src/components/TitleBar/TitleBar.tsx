import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../Frame/Frame';
import './TitleBar.css';

interface Props {
    title?: string
    icon?: string
    role: Signal<string>
    showLogin: Signal<boolean>;
}

function TitleBar({title, icon, role, showLogin}: Props) {
    useSignals();
    return (
        <Frame className='title'>
            <img src={ icon ? icon : './logo.png' } />
            <span>
                {title ? title : "< Title >"}
            </span>
            <button onClick={ () => showLogin.value=true }>Login</button>
        </Frame>
    )
}

export default TitleBar;
