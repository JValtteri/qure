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
import Popup from './components/Popup/Popup';


const showLogin = signal( false );
const show = signal({ "selectedEventId": -1, "eventID": -1, "editor": false});
const user = signal({"username": "", "loggedIn": false, "admin": false});
const loadingEvents = signal(false);

const resumeSession = async (
    setServerError: React.Dispatch<React.SetStateAction<string>>,
    setErrorVisible: React.Dispatch<React.SetStateAction<boolean>>
) => {
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

function App() {
    useSignals();
    console.log("App rendered");
    const [errorVisible, setErrorVisible] = useState(false);
    const [serverError, setServerError] = useState("");
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
        resumeSession(setServerError, setErrorVisible);
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
            <Popup show={errorVisible} onHide={() => setErrorVisible(false)} className='error'>
                {serverError}
            </Popup>
        </>
    )
}

export default App
