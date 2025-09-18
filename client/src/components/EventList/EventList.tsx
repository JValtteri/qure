import type { ReactNode } from 'react';
import { useState } from 'react';
import Frame from '../Frame/Frame';
import type { Event } from '../../events';
import EventCard from '../EventCard/EventCard';
import './EventList.css';

interface Props {
    items: Event[];
}

const makeCard = (event: Event, index: number, selected: number, setSelected: any) => (
        <EventCard
            title={event.name}
            startTime={event.dtStart}
            desc={event.shortDescription}
            time='0'
            slots={event.guestSlots}
            occupied={event.guests}
            key={index}
            onClick={() => setSelected(index)}
            selected={ index == selected }
        />
)

function EventList({items}: Props) {
    const [selected, setSelected] = useState(-1);
    const children: ReactNode[] = (
        items.map( (item: Event, index: number) =>
            makeCard(item, index, selected, setSelected)
        ));
    return (
        <Frame reactive={false} className='list-body'>
            {items.length === 0 && <p>no item found</p>}
            {children}
        </Frame>
    )
}

export default EventList;
