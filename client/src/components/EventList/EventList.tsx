import type { ReactNode } from 'react';
import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../common/Frame/Frame';
import EventCard from './EventCard/EventCard';
import AddCard from '../AddCard/AddCard';
import type { Event } from '../../utils/events';
import './EventList.css';


interface Props {
    items: Event[];
    show: Signal<{ "selectedEventId": number, "editor": boolean}>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
}

const makeCard = (event: Event, index: number, show: Signal ) => (
        <EventCard
            title={event.name}
            startTime={event.dtStart}
            desc={event.shortDescription}
            time='0'
            slots={event.guestSlots}
            occupied={event.guests}
            key={index}
            onClick={ () => (show.value = showIndex(index)) }
            selected={ show.value.selectedEventId == index }
        />
)

const showIndex = (index: number) => {
    return {"selectedEventId": index, "editor": false};
}

const showEditor = () => {
    return {"selectedEventId": -1, "editor": true};
}

function EventList({items, show, user}: Props) {
    useSignals();
    console.log("List rendered")

    const children: ReactNode[] = (
        items.map( (item: Event, index: number) =>
            makeCard(item, index, show )
        ));
    return (
        <Frame reactive={false} className='list-body'>
            {items.length === 0 && <p>no item found</p>}
            {children}
            <AddCard onClick={ () => show.value=showEditor() } hidden={!user.value.admin} />
        </Frame>
    )
}

export default EventList;
