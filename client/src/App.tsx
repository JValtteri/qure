import './App.css'
import { signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import EventList from './components/EventList/EventList'
import LoremIpsum from './components/LoremIpsum/LoremIpsum'
import TitleBar from './components/TitleBar/TitleBar'
import DetailCard from './components/DetailCard/DetailCard';
import LoginDialog from './components/Login/Login';
import { getEvents } from './events'

const selectedEventId = signal( -1 );
const showLogin = signal( false );
const user = signal({"username": "", "loggedIn": false, "admin": false});

function App() {
  useSignals();
  console.log("App rendered")
  return (
    <>
      <div className='view'>
        <TitleBar title='' showLogin={showLogin} user={user}/>
        <EventList items={getEvents()} selectedId={selectedEventId} user={user} />
        <DetailCard selectedId={selectedEventId} user={user}>
          <LoremIpsum />
        </DetailCard>
      </div>
      <LoginDialog showLogin={showLogin} user={user}/>
    </>
  )
}

export default App
