import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../common/Frame/Frame';
import { logout } from '../../api/api';
import './TitleBar.css';
import { clearCookie } from '../../utils/cookie';

interface Props {
    title?: string
    icon?: string
    showLogin: Signal<boolean>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>
}

function TitleBar({title, icon, showLogin, user}: Props) {
    useSignals();
    console.log("Title rendered")

    const handleLogout = () => {
        logout();
        clearCookie("sessionKey");
        user.value = { username: "", loggedIn: false, admin: false};
    };

    const handleLogin = () => showLogin.value=true

    return (
        <Frame className='title'>
            <img src={ icon ? icon : './logo.png' } />
            <div />
            <span>
                {title ? title : "< Title >"}
            </span>
            <div>
                {user.value.username} {user.value.admin && "(admin)"}
            </div>
            <button hidden={user.value.loggedIn === false} onClick={ handleLogout }>Logout</button>
            <button hidden={user.value.loggedIn === true}  onClick={ handleLogin }>Login</button>
        </Frame>
    )
}

export default TitleBar;
