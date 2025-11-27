import './TimeslotList.css';

import { type ReactNode } from 'react';
import { type Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import ListCard from '../common/ListCard/ListCard';

import type { Timeslot } from "../../api/api";
import { posixToDateAndTime, posixToTime } from '../../utils/utils';

interface Props {
    timeslots:     Map<Number, Timeslot>;
    selectedSlot:  Signal<number>;
    requestUpdate: Signal<boolean>;
}

function TimeslotList({timeslots, selectedSlot, requestUpdate}: Props) {
    useSignals();
    const children = makeChildren(timeslots, selectedSlot, requestUpdate);
    return (
        <>
            {timeslots.size === 0 && <p>No timeslots found</p>}
            {children}
        </>
    )
}

export default TimeslotList;


const makeCard = (time: number, slots: number, reserved: number, selectedSlot: Signal<number>, requestUpdate: Signal<boolean>) => (
    <ListCard
        title={posixToTime(time)}
        startTime={""}
        desc={posixToDateAndTime(time)}
        slots={slots}
        occupied={reserved}
        key={time}
        onClick={ () => {
            selectedSlot.value = time;
            requestUpdate.value = !requestUpdate.value;
        } }
        selected={ selectedSlot.value == time }
        className="timeslot-list-card"
    />
)

function makeListElement(
    timeslot:       Timeslot,
    index:          number,
    selectedSlot:   Signal<number>,
    requestUpdate:  Signal<boolean>
) {
    try {
        return makeCard(index, timeslot.Size, timeslot.Reserved, selectedSlot, requestUpdate);
    } catch {
        return makeCard(index, -1, -1, selectedSlot, requestUpdate);
    }
};

function makeChildren(timeslots: Map<Number, Timeslot>, selectedSlot: Signal<number>, requestUpdate: Signal<boolean>): ReactNode[] {
    let children: ReactNode[] = [];
    for (const [time, timeslot] of timeslots) {
        children.push(makeListElement(timeslot, time.valueOf(), selectedSlot, requestUpdate));
    }
    return children;
}
