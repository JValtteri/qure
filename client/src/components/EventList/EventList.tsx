import type { ReactNode } from 'react';
import { Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";
import Frame from '../Frame/Frame';
import type { Event } from '../../events';
import EventCard from '../EventCard/EventCard';
import './EventList.css';

interface Props {
    items: Event[];
    selectedId: Signal<number>;
}

const makeCard = (event: Event, index: number, selectedId: Signal, selectThis: (index: number) => number  ) => (
        <EventCard
            title={event.name}
            startTime={event.dtStart}
            desc={event.shortDescription}
            time='0'
            slots={event.guestSlots}
            occupied={event.guests}
            key={index}
            onClick={ () => {
                    (selectedId.value = index);
                    console.log(selectedId.value);
            } }
            selected={ selectedId.value == index }
        />
)

function EventList({items, selectedId}: Props) {
    useSignals();
    console.log("List rendered")

    const children: ReactNode[] = (
        items.map( (item: Event, index: number) =>
            makeCard(item, index, selectedId, (index: number) => ( selectedId.value = index ) )
        ));
    return (
        <Frame reactive={false} className='list-body'>
            {items.length === 0 && <p>no item found</p>}
            {children}
        </Frame>
    )
}

export default EventList;
