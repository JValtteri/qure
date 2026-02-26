import './App.css'

import { lazy, Suspense, useEffect, useState } from 'react';
import { signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import { fetchEvents, type EventResponse, authenticate } from './api/api';

import Spinner from './components/Spinner/Spinner';
import TitleBar from './components/TitleBar/TitleBar'

const EventList = lazy( () => import('./components/EventList/EventList'));
const DetailCard = lazy(() => import('./components/DetailCard/DetailCard'));
const LoginDialog = lazy(() => import('./components/Login/Login'));
const EventCreation = lazy(() => import('./components/EventCreation/EventCreation'));
const UserForm = lazy( () => import('./components/UserForm/UserForm'));
const Popup = lazy(() => import('./components/Popup/Popup'));


const showLogin = signal( false );
const show = signal({"eventID": "none", "editor": false, "account": false, "inspect": false});
const user = signal({"username": "", "loggedIn": false, "role": ""});
const requestedUpdate = signal(true);

const loadingEvents = signal(false);


function App() {
    useSignals();

    const [errorVisible, setErrorVisible] = useState(false);
    const [serverError, setServerError] = useState("");
    const [events, setEvents] = useState(new Array<EventResponse>());

    const updateEventsHandler = updateEvents(setEvents);

    useEffect(() => {
        updateEventsHandler();
    }, [requestedUpdate.value]);

    useEffect(() => {
        resumeSession(setServerError, setErrorVisible);
    }, []);

    return (
        <>
            <div className='view'>
                <TitleBar title='' showLogin={showLogin} user={user} showAccount={show}/>
                <Suspense fallback={<Spinner />}>
                    <EventList show={show} items={events} user={user} update={ updateEventsHandler } />
                    {(show.value.eventID != "none" && !show.value.inspect ) && <DetailCard show={show} user={user} requestedUpdate={requestedUpdate} />}
                    {(show.value.account || show.value.inspect) && <UserForm user={user} show={show} />}
                    {show.value.editor && <EventCreation show={show} update={ updateEventsHandler } />}
                </Suspense>
            </div>
            <Suspense>
                { showLogin && <LoginDialog showLogin={showLogin} user={user}/> }
                { errorVisible &&
                    <Popup show={errorVisible} onHide={() => setErrorVisible(false)} className='error'>
                        {serverError}
                    </Popup>
                }
            </Suspense>
        </>
    )
}

export default App


async function resumeSession(
    setServerError: React.Dispatch<React.SetStateAction<string>>,
    setErrorVisible: React.Dispatch<React.SetStateAction<boolean>>
): Promise<void> {
    try {
        const auth = await authenticate();
        if ( auth != null ) {
        showLogin.value = false;
        user.value = { username: auth.User, loggedIn: true, role: auth.Role};
        }
    } catch (error) {
        setServerError(`${error}`);
        setErrorVisible(true);
    }
}

function updateEvents(setEvents: React.Dispatch<React.SetStateAction<EventResponse[]>>): () => Promise<void> {
    return async () => {
        if (loadingEvents.value == true) {
            return;
        }
        loadingEvents.value = true;
        await fetchEvents()
            .then(value => {
                if (value != null) {
                    setEvents(value);
                }
            });
        loadingEvents.value = false;
    };
}
