import './App.css'
import { useEffect, useState } from 'react';
import { signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import EventList from './components/EventList/EventList'
import TitleBar from './components/TitleBar/TitleBar'
import DetailCard from './components/DetailCard/DetailCard';
import LoginDialog from './components/Login/Login';
import EventCreation from './components/EventCreation/EventCreation';
import { fetchEvents, type EventResponse, authenticate } from './api/api';


const showLogin = signal( false );
const show = signal({ "selectedEventId": -1, "editor": false});
const user = signal({"username": "", "loggedIn": false, "admin": false});
const loadingEvents = signal(false);

const resumeSession = async () => {
    let auth = await authenticate();
    if ( auth != null ) {
      showLogin.value = false;
      user.value = { username: auth.User, loggedIn: true, admin: auth.IsAdmin};
    }
}

function App() {
    useSignals();
    console.log("App rendered");

    const [events, setEvents] = useState(new Array<EventResponse>())

    const updateEvents = async () => {
        if (loadingEvents.value == true) {
            return;
        } else {
            loadingEvents.value = true;
            await fetchEvents()
                .then(value => {
                    if (value != null) {
                        setEvents(value);
                    }
                });
            loadingEvents.value = false;
        }
    }

    useEffect(() => {
        resumeSession();
        updateEvents();
    }, []);

    return (
        <>
            <div className='view'>
                <TitleBar title='' showLogin={showLogin} user={user}/>
                <EventList show={show} items={events} user={user} update={updateEvents} />
                <DetailCard show={show} user={user} />
                <EventCreation show={show} update={updateEvents} />
            </div>
            <LoginDialog showLogin={showLogin} user={user}/>
        </>
    )
}

export default App
