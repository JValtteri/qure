import { useState } from 'react';
import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import Dialog from '../common/Dialog/Dialog';

import { login, registerUser } from '../../api/api';

import './Login.css';


interface Props {
    showLogin: Signal<boolean>;
    user: Signal<{"username": string, "loggedIn": boolean, "admin": boolean}>;
}

function LoginDialog({showLogin, user}: Props) {
  useSignals();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [password2, setPassword2] = useState("");
  const [newAccount, setNewAccount] = useState(false);
  const [registerErr, setRegisterErr] = useState("");

  const emailInput = document.getElementById("email");
  const passInput = document.getElementById("password");
  const passInput2 = document.getElementById("password-confirm");

  const submit = async () => {
    const auth = await login(username, password);
    if ( auth != null ) {
      showLogin.value = false;
      user.value = { username: username, loggedIn: true, admin: auth.IsAdmin};
      removeWrongLabelFromCredentials(emailInput, passInput);
      setPassword("");
    } else {
      labelWrongCredentials(emailInput, passInput);
    }
  };

  const signUp = async () => {
      if (password != password2) {
          passInput2?.classList.add("wrong");
          return;
      }
      passInput2?.classList.remove("wrong");
      registerUser(username, password)
        .then((reg) => {
            setRegisterErr(reg.Error);
            if (reg.Error) {
                return;
            }
            user.value = { username: username, loggedIn: true, admin: false};
            setNewAccount(false);
            showLogin.value=false;
        });
  }

  const handleNewAccount = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewAccount(e.target.checked);
    setRegisterErr("");
  }

  return(
    <Dialog className='login grid' hidden={ showLogin.value===false }>
      <label id='email-label' htmlFor="email">Email:</label>
      <input
        type="email"
        id="email"
        value={username}
        onChange={e => setUsername(e.target.value)}
        required
      />
      <label id='password-label' htmlFor="password">Password:</label>
      <input
        type="password"
        id="password"
        className='password'
        value={password}
        onChange={e => setPassword(e.target.value)}
        required
      />
      <label id='password-confirm-label' htmlFor="password-confirm" hidden={!newAccount}>Confirm password:</label>
      <input
        type="password"
        id="password-confirm"
        value={password2}
        onChange={e => setPassword2(e.target.value)}
        hidden={!newAccount}
      />
      <div className='new-account'>
        <label htmlFor="new-account">New account:</label>
        <input id="new-account" type="checkbox" checked={newAccount} onChange={ handleNewAccount} ></input>
      </div>
      <div className='buttons'>
        <button onClick={submit} hidden={newAccount} className='selected'>Login</button>
        <button onClick={signUp} hidden={!newAccount} className='selected'>Sign Up</button>
        <button onClick={ () => showLogin.value=false }>Cancel</button>
      </div>
      <div id="register-error">{registerErr}</div>
    </Dialog>
  )
}

export default LoginDialog;

function labelWrongCredentials(emailInput: HTMLElement | null, passInput: HTMLElement | null) {
  emailInput?.classList.add("wrong");
  passInput?.classList.add("wrong");
}

function removeWrongLabelFromCredentials(emailInput: HTMLElement | null, passInput: HTMLElement | null) {
  emailInput?.classList.remove("wrong");
  passInput?.classList.remove("wrong");
}
