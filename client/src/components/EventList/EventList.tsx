import type { ReactNode } from 'react';
import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../common/Frame/Frame';
import EventCard from './EventCard/EventCard';
import AddCard from '../AddCard/AddCard';
import './EventList.css';
import { type EventResponse } from '../../api/api';
import { posixToDateAndTime } from '../../utils/utils';


interface Props {
    items: EventResponse[];
    show: Signal<{ "selectedEventId": number, "editor": boolean}>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
    update: ()=>Promise<void>
}

const makeCard = (event: EventResponse, index: number, show: Signal, update: ()=>Promise<void> ) => (
        <EventCard
            title={event.Name}
            startTime={posixToDateAndTime(event.DtStart)}
            desc={event.ShortDescription}
            time='0'
            slots={999} ///////////
            occupied={999} //////////////
            key={index}
            onClick={ () => {
                show.value = showIndex(index)
                update();
            } }
            selected={ show.value.selectedEventId == index }
        />
)

const showIndex = (index: number) => {
    return {"selectedEventId": index, "editor": false};
}

const showEditor = () => {
    return {"selectedEventId": -1, "editor": true};
}

function EventList({items, show, user, update}: Props) {
    useSignals();
    console.log("List rendered")

    const children: ReactNode[] = (
        items.map( (item: EventResponse, index: number) =>
            makeCard(item, index, show, update)
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
