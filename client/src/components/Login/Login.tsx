import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../Frame/Frame';
import './Login.css';


interface Props {
    showLogin: Signal<boolean>;
}

function LoginDialog({showLogin}: Props) {
  useSignals();
  console.log(showLogin.value)

  return(
    <Frame className='dialog' hidden={ showLogin.value===false }>
      <label className='email-label' htmlFor="email">Email:</label>
      <input
        className='email'
        type="email"
        id="email"
        required
      />
      <label className='password-label' htmlFor="password">Password:</label>
      <input
        className='password'
        type="password"
        id="password"
        required
      />
      <div className='buttons'>
        <button>Login</button>
        <button onClick={ () => showLogin.value=false }>Cancel</button>
      </div>
    </Frame>
  )
}

export default LoginDialog;
