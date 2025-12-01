import { Suspense } from 'react';
import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import Frame from '../common/Frame/Frame';
import Spinner from '../Spinner/Spinner';

import { logout } from '../../api/api';
import { clearCookie } from '../../utils/cookie';

import './TitleBar.css';


interface Props {
    title?: string
    icon?: string
    showLogin: Signal<boolean>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>
}

function TitleBar({title, icon, showLogin, user}: Props) {
    useSignals();

    const handleLogout = () => {
        logout();
        clearCookie("sessionKey");
        user.value = { username: "", loggedIn: false, admin: false};
    };

    const handleLogin = () => showLogin.value=true;

    return (
        <Frame className='title'>
            <Suspense fallback={<Spinner />}>
                <img src={ icon ? icon : './logo.png' } />
            </Suspense>
            <div />
            <span id='title'>
                {title ? title : "< Title >"}
            </span>
            <div>
                <div id='user'>
                    {user.value.username.split('@')[0]}
                    {user.value.admin && "(admin)"}
                </div>
            </div>

            <div hidden={user.value.loggedIn} />
            <button id='menu-button' hidden={!user.value.loggedIn}>≡</button>

            <button hidden={user.value.loggedIn === false} onClick={ handleLogout }>Logout</button>
            <button hidden={user.value.loggedIn === true}  onClick={ handleLogin }>Login</button>
        </Frame>
    )
}

export default TitleBar;
