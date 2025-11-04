import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../Frame/Frame';
import './TitleBar.css';

interface Props {
    title?: string
    icon?: string
    role: Signal<string>
    showLogin: Signal<boolean>;
    user: Signal<{username: string, loggedIn: boolean}>
}

function TitleBar({title, icon, role, showLogin, user}: Props) {
    useSignals();
    console.log("Title rendered")

    const logout = () => {
        user.value = { username: "", loggedIn: false };
    };

    const login = () => showLogin.value=true

    return (
        <Frame className='title'>
            <img src={ icon ? icon : './logo.png' } />
            <div />
            <span>
                {title ? title : "< Title >"}
            </span>
            <div>
                {user.value.username}
            </div>
            <button hidden={user.value.loggedIn === false} onClick={ logout }>Logout</button>
            <button hidden={user.value.loggedIn === true}  onClick={ login }>Login</button>
        </Frame>
    )
}

export default TitleBar;
