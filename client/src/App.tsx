import './App.css'
import { computed, effect, signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import EventList from './components/EventList/EventList'
import LoremIpsum from './components/LoremIpsum/LoremIpsum'
import TitleBar from './components/TitleBar/TitleBar'
import { getEvents } from './events'
import DetailCard from './components/DetailCard/DetailCard';

const selectedEventId = signal( -1 ); // TODO: Add selection system function
// access info selectedEventId.value

function App() {
  useSignals();
  console.log("App rendered")
  return (
    <div className='view'>
      <TitleBar title='' />
      <EventList items={getEvents()} selectedId={selectedEventId} />

      <DetailCard selectedId={selectedEventId}>
        <LoremIpsum />
      </DetailCard>

      <div>
        {selectedEventId.value}
      </div>

    </div>
  )
}

export default App
