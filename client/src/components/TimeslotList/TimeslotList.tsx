import './TimeslotList.css';

import { type ReactNode } from 'react';
import { useSignals } from "@preact/signals-react/runtime";

import ListCard from '../common/ListCard/ListCard';

import type { Timeslot } from "../../api/api";
import { posixToTime } from '../../utils/utils';

interface Props {
    timeslots: Map<Number, Timeslot>;
}

function TimeslotList({timeslots}: Props) {
    useSignals();
    console.log("Timeslots rendered");
    //show: Signal<{ "selectedEventId": number, "eventID": number, "editor": boolean}>;
    const children = makeChildren(timeslots);
    return (
        <>
            {timeslots.size === 0 && <p>No timeslots found</p>}
            {children}
        </>
    )
}

export default TimeslotList;


const makeCard = (time: number, slots: number, reserved: number) => (
    <ListCard
        title={posixToTime(time)}
        startTime={""}
        desc={""}
        slots={slots}
        occupied={reserved}
        key={time}
        onClick={ () => {
            //show.value = showIndex(index, event.ID)
            //update(); ////////////////////////////////////////
        } }
        selected={false}//{ show.value.selectedEventId == index }
        className="timeslot-list-card"
    />
)

function makeListElement(
    timeslot: Timeslot,
    index: number,
) {
    try {
        return makeCard(index, timeslot.Size, timeslot.Reserved);
    } catch {
        return makeCard(index, -1, -1);
    }
};

function makeChildren(timeslots: Map<Number, Timeslot>): ReactNode[] {
    console.log("Rendered Children")
    let children: ReactNode[] = [];
    for (const [time, timeslot] of timeslots) {
        children.push(makeListElement(timeslot, time.valueOf()));
    }
    return children;
}
