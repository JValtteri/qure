import './App.css'
import EventList from './components/EventList/EventList'
import Frame from './components/Frame/Frame'
import TitleBar from './components/TitleBar/TitleBar'
import { getEvents } from './events'

function App() {

  return (
    <>
      <TitleBar />
      <EventList items={getEvents()}/>
      <Frame hidden></Frame>
    </>
  )
}

export default App
