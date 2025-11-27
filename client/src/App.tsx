import './App.css'

import { useEffect, useState } from 'react';
import { signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import EventList from './components/EventList/EventList'
import TitleBar from './components/TitleBar/TitleBar'
import DetailCard from './components/DetailCard/DetailCard';
import LoginDialog from './components/Login/Login';
import EventCreation from './components/EventCreation/EventCreation';
import Popup from './components/Popup/Popup';

import { fetchEvents, type EventResponse, authenticate } from './api/api';


const showLogin = signal( false );
const show = signal({ "selectedEventId": -1, "eventID": -1, "editor": false});
const user = signal({"username": "", "loggedIn": false, "admin": false});
const requestedUpdate = signal(true);

const loadingEvents = signal(false);


function App() {
    useSignals();

    const [errorVisible, setErrorVisible] = useState(false);
    const [serverError, setServerError] = useState("");
    const [events, setEvents] = useState(new Array<EventResponse>())

    const updateEventsHandler = updateEvents(setEvents);

    useEffect(() => {
        resumeSession(setServerError, setErrorVisible);
        updateEventsHandler();
    }, [show.value, requestedUpdate.value]);

    return (
        <>
            <div className='view'>
                <TitleBar title='' showLogin={showLogin} user={user}/>
                <EventList show={show} items={events} user={user} update={ updateEventsHandler } />
                <DetailCard show={show} user={user} requestedUpdate={requestedUpdate} />
                <EventCreation show={show} update={ updateEventsHandler } />
            </div>
            <LoginDialog showLogin={showLogin} user={user}/>
            <Popup show={errorVisible} onHide={() => setErrorVisible(false)} className='error'>
                {serverError}
            </Popup>
        </>
    )
}

export default App


async function resumeSession(
    setServerError: React.Dispatch<React.SetStateAction<string>>,
    setErrorVisible: React.Dispatch<React.SetStateAction<boolean>>
) {
    try {
        let auth = await authenticate();
        if ( auth != null ) {
        showLogin.value = false;
        user.value = { username: auth.User, loggedIn: true, admin: auth.IsAdmin};
        }
    } catch (error) {
        setServerError(`${error}`);
        setErrorVisible(true);
    }
}

function updateEvents(setEvents: any) {
    return async () => {
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
    };
}
