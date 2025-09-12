import './App.css'
import EventList from './components/EvenList/EventList'
import { getEvents } from './events'

function App() {

  return (
    <>
      <EventList items={getEvents()}/>
    </>
  )
}

export default App
