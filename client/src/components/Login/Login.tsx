import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Dialog from '../Dialog/Dialog';
import './Login.css';
import { login } from '../../api';
import { useState } from 'react';


interface Props {
    showLogin: Signal<boolean>;
    user: Signal<{"username": string, "loggedIn": boolean}>;
}

function LoginDialog({showLogin, user}: Props) {
  useSignals();
  console.log(showLogin.value)

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const emailInput = document.getElementById("email");
  const passInput = document.getElementById("password");

  const submit = async () => {
    let ok = await login(username, password);
    if ( ok === true ) {
      showLogin.value = false;
      user.value = { username: username, loggedIn: true };
      emailInput?.classList.remove("wrong");
      passInput?.classList.remove("wrong");
      setPassword("");
    } else {
      emailInput?.classList.add("wrong");
      passInput?.classList.add("wrong");
    }
  };

  return(
    <Dialog hidden={ showLogin.value===false }>
      <label className='email-label' htmlFor="email">Email:</label>
      <input
        type="email"
        id="email"
        value={username}
        onChange={e => setUsername(e.target.value)}
        required
      />
      <label className='password-label' htmlFor="password">Password:</label>
      <input
        type="password"
        id="password"
        value={password}
        onChange={e => setPassword(e.target.value)}
        required
      />
      <div className='buttons'>
        <button onClick={ submit } >Login</button>
        <button onClick={ () => showLogin.value=false }>Cancel</button>
      </div>
    </Dialog>
  )
}

export default LoginDialog;
