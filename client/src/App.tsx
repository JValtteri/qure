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
const clientRole = signal( "guest" );

function App() {
  useSignals();
  console.log("App rendered")
  return (
    <>
      <div className='view'>
        <TitleBar title='' role={clientRole} showLogin={showLogin} />
        <EventList items={getEvents()} selectedId={selectedEventId} />
        <DetailCard selectedId={selectedEventId} role={clientRole}>
          <LoremIpsum />
        </DetailCard>
      </div>
      <LoginDialog showLogin={showLogin} />
    </>
  )
}

export default App
