import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../Frame/Frame';
import './Login.css';
import { login } from '../../api';
import { useState } from 'react';


interface Props {
    showLogin: Signal<boolean>;
    user: Signal<{"username": string, "loggedIn": boolean}>
}

function LoginDialog({showLogin, user}: Props) {
  useSignals();
  console.log(showLogin.value)

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const submit = async () => {
    let ok = await login(username, password);
    if ( ok === true ) {
      showLogin.value = false;
      user.value = { username: username, loggedIn: true };
    }
  };

  return(
    <Frame className='dialog' hidden={ showLogin.value===false }>
      <label className='email-label' htmlFor="email">Email:</label>
      <input
        className='email'
        type="email"
        id="email"
        onChange={e => setUsername(e.target.value)}
        required
      />
      <label className='password-label' htmlFor="password">Password:</label>
      <input
        className='password'
        type="password"
        id="password"
        onChange={e => setPassword(e.target.value)}
        required
      />
      <div className='buttons'>
        <button onClick={ submit } >Login</button>
        <button onClick={ () => showLogin.value=false }>Cancel</button>
      </div>
    </Frame>
  )
}

export default LoginDialog;
