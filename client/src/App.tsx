import './App.css'
import EventList from './components/EventList/EventList'
import Frame from './components/Frame/Frame'
import { getEvents } from './events'

function App() {

  return (
    <>
      <EventList items={getEvents()}/>
      <Frame hidden></Frame>
    </>
  )
}

export default App
