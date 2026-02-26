import { Suspense } from 'react';
import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import Frame from '../common/Frame/Frame';
import Spinner from '../Spinner/Spinner';

import { logout } from '../../api/api';
import { clearCookie } from '../../utils/cookie';

import './TitleBar.css';


interface Props {
    title?: string;
    icon?: string;
    showLogin: Signal<boolean>;
    user: Signal<{username: string, loggedIn: boolean, role: string}>;
    showAccount: Signal<{"eventID": string, "editor": boolean, "account": boolean, "inspect": boolean}>;
}


function TitleBar({title, icon, showLogin, user, showAccount}: Props) {
    useSignals();

    const handleLogout = () => {
        logout();
        clearCookie("sessionKey");
        user.value = { username: "", loggedIn: false, role: ""};
        showAccount.value = {eventID: "none", editor: false, account: false, inspect: false};
    };

    const handleLogin = () => showLogin.value=true;

    const handleHamburgerMenu = () => {
        showAccount.value = {
            "eventID": "none",
            "editor": false,
            "account": !(showAccount.value.account || showAccount.value.inspect),
            "inspect": false
        }
    }

    return (
        <Frame className='title'>
            <Suspense fallback={<Spinner />}>
                <img src={ icon ? icon : './logo.png' } fetchPriority='low' />
            </Suspense>
            <div />
            <span id='title'>
                {title ? title : "< Title >"}
            </span>
            <div>
                <div id='user'>
                    {user.value.username.split('@')[0]}
                    {user.value.role === "admin" && "(admin)"}
                </div>
            </div>

            <div hidden={user.value.loggedIn} />
            <button
                id='menu-button'
                onClick={ handleHamburgerMenu }
                hidden={!user.value.loggedIn}
                className={showAccount.value.account ? 'selected' : ''
            }>≡</button>

            <button hidden={user.value.loggedIn === false} onClick={ handleLogout }>Logout</button>
            <button hidden={user.value.loggedIn === true}  onClick={ handleLogin }>Login</button>
        </Frame>
    )
}

export default TitleBar;
