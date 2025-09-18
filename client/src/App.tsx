import './App.css'
import EventList from './components/EventList/EventList'
import Frame from './components/Frame/Frame'
import LoremIpsum from './components/LoremIpsum/LoremIpsum'
import TitleBar from './components/TitleBar/TitleBar'
import { getEvents } from './events'

function App() {

  return (
    <div className='view'>
      <TitleBar title='' />
      <EventList items={getEvents()}/>
      <Frame hidden>
        <LoremIpsum />
      </Frame>
    </div>
  )
}

export default App
