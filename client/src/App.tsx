import './App.css'

import { lazy, Suspense, useEffect, useState } from 'react';
import { signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import { fetchEvents, type EventResponse } from './api/api';

import Spinner from './components/Spinner/Spinner';
import TitleBar from './components/TitleBar/TitleBar'
import { resumeSession } from './components/common/utils';

const EventList = lazy( () => import('./components/EventList/EventList'));
const DetailCard = lazy(() => import('./components/DetailCard/DetailCard'));
const LoginDialog = lazy(() => import('./components/Login/Login'));
const EventCreation = lazy(() => import('./components/EventCreation/EventCreation'));
const UserForm = lazy( () => import('./components/UserForm/UserForm'));
const Popup = lazy(() => import('./components/Popup/Popup'));


const showLogin = signal( false );
const show = signal({eventID: "none", view: ""});
const user = signal({username: "", loggedIn: false, role: ""});
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
        resumeSession(setServerError, setErrorVisible, user, showLogin);
    }, []);

    return (
        <>
            <div className='view'>
                <TitleBar title='' showLogin={showLogin} user={user} show={show}/>
                <Suspense fallback={<Spinner />}>
                    <EventList show={show} items={events} user={user} update={ updateEventsHandler } />
                </Suspense>
                <Suspense fallback={<Spinner />}>
                    {(show.value.eventID != "none" && show.value.view == "" ) &&
                        <DetailCard show={show} user={user} requestedUpdate={requestedUpdate} />}
                    {["account", "inspect"].includes(show.value.view) &&
                        <UserForm user={user} show={show} />}
                    {show.value.view == "editor" &&
                        <EventCreation show={show} update={ updateEventsHandler } />}
                </Suspense>
            </div>
            <Suspense>
                { showLogin &&
                    <LoginDialog showLogin={showLogin} user={user}/> }
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


function updateEvents(setEvents: React.Dispatch<React.SetStateAction<EventResponse[]>>): () => Promise<void> {
    return async () => {
        if (loadingEvents.value == true) {
            return;
        }
        loadingEvents.value = true;
        try {
            await fetchEvents()
                .then(value => {
                    if (value != null) {
                        setEvents(value);
                    }
                });
        } catch (error: any) {
            console.warn(error.message);
        }
        loadingEvents.value = false;
    };
}
