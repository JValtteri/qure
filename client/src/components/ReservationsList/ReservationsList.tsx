import './ReservationsList.css';

import { type ReactNode } from 'react';
import { type Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import ListCard from '../common/ListCard/ListCard';

import type { ReservationResponse } from "../../api/api";
import { posixToDateAndTime } from '../../utils/utils';


interface Props {
    reservations:   ReservationResponse[];
    selected:       Signal<number>;
    update:         ()=>Promise<void>;
}

function ReservationsList({reservations, selected: selectedSlot, update}: Props) {
    useSignals();
    const children = makeChildren(reservations, selectedSlot, update);
    return (
        <div id="reservation-list">
            {reservations.length === 0 && <p>No reservations found</p>}
            {children}
        </div>
    )
}

export default ReservationsList;


function makeChildren(reservations: ReservationResponse[], selectedSlot: Signal<number>, update: ()=>Promise<void>): ReactNode[] {
    let children: ReactNode[] = [];
    reservations = reservations.sort( (a, b) => a.Timeslot - b.Timeslot );
    children = reservations.map( (item: ReservationResponse, index: number) => {
        console.log()
        return makeListElement(item, index, item.Timeslot, selectedSlot, update);
    })
    return children;
}

function makeListElement(
    reservation:    ReservationResponse,
    index:          number,
    time:           number,
    selectedSlot:   Signal<number>,
    update:         ()=>Promise<void>
) {
    try {
        return makeCard(index, time, reservation.Confirmed, reservation.Size, selectedSlot, update);
    } catch {
        return makeCard(index, time, -1, -1, selectedSlot, update);
    }
};

const makeCard = (index: number, time: number, confirmed: number, size: number, selectedSlot: Signal<number>, update: ()=>Promise<void>) => (
    <ListCard
        title={"< Reservation >"}
        startTime={""}
        desc={posixToDateAndTime(time)}
        slots={confirmed}
        occupied={size}
        key={index}
        onClick={ () => {
            selectedSlot.value = index;
            update();
        } }
        selected={ selectedSlot.value == index }
        className="timeslot-list-card"
    />
)
