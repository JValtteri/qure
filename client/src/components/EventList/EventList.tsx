import './EventList.css';

import type { ReactNode } from 'react';
import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import Frame from '../common/Frame/Frame';
import ListCard from '../common/ListCard/ListCard';
import AddCard from '../AddCard/AddCard';

import type { EventResponse } from '../../api/api';
import { countSlots, posixToDateAndTime } from '../../utils/utils';


interface Props {
    items: EventResponse[];
    show: Signal<{ "selectedEventId": number, "eventID": number, "editor": boolean}>;
    user: Signal<{username: string, loggedIn: boolean, admin: boolean}>;
    update: ()=>Promise<void>
}

function EventList({items, show, user, update}: Props) {
    useSignals();

    items = items.sort( (a, b) => a.DtStart - b.DtStart );
    const children: ReactNode[] = (
        items.map( (item: EventResponse, index: number) => {
            return makeListElement(item, index, show, update);
        })
    );
    return (
        <Frame reactive={false} className='list-body'>
            {items.length === 0 && <p>no item found</p>}
            {children}
            <AddCard onClick={ () => show.value=showEditor() } hidden={!user.value.admin} />
        </Frame>
    )
}

export default EventList;


const makeCard = (event: EventResponse, index: number, slots: number, reserved: number, show: Signal, update: ()=>Promise<void> ) => (
    <ListCard
        title={event.Name}
        startTime={posixToDateAndTime(event.DtStart)}
        desc={event.ShortDescription}
        slots={slots}
        occupied={reserved}
        key={index}
        onClick={ () => {
            show.value = showIndex(index, event.ID)
            update();
        } }
        selected={ show.value.selectedEventId == index }
        className="event-list-card"
    />
)

const showIndex = (index: number, id: number) => ({"selectedEventId": index, "eventID": id, "editor": false});

const showEditor = () => ({"selectedEventId": -1, "eventID": -1, "editor": true});

function makeListElement(
    item: EventResponse,
    index: number,
    show: Signal<{ "selectedEventId": number, "eventID": number, "editor": boolean}>,
    update: ()=>Promise<void>
) {
    const timeslots = new Map(Object.entries(item.Timeslots).map(([k, v]) => [Number(k), v]));
    try {
        const { totalSlots, totalReservedSlots } = countSlots(timeslots);
        return makeCard(item, index, totalSlots, totalReservedSlots, show, update);
    } catch {
        return makeCard(item, index, -1, -1, show, update);
    }
};
