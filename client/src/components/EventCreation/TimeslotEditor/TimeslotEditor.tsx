import "./TimeslotEditor.css"

import { useEffect, useState, type ReactNode } from 'react';
import { type Signal } from "@preact/signals-react";

import { cycleDay, dateAndTimeToPosix } from '../../../utils/utils';
import { useSignals } from "@preact/signals-react/runtime";


interface Props {
    startTime: string;
    date: string;
    timeslot: Signal<Map<number, {"Size": number}>>;
}

function TimeslotEditor({startTime, date, timeslot}: Props) {
    useSignals();
    // Timeslot data
    const [times, setTimes] = useState<{ [label: number]: string }>({[0]: startTime});
    const [groupSizes, setGroupSizes] = useState<{ [label: number]: number }>({[0]: 0});
    // Aux data
    const [defaultGroupSize, setDefaultGroupSize] = useState(0);
    const [size, setSize] = useState(1);

    useEffect( () => {
        updateTimeslots();
    }, [times, groupSizes, date])

    const renderTimeslots = (size: number, inputs: {[label: number]: string;}, handleInputChange: (index: number, element: HTMLInputElement)=>void) => {
        inputs[0] = startTime;
        let out = new Array<ReactNode>;
        out.push(
            <div key={0} className='timeslot-editor-row'>
                <label>1.</label>
                <input type="time" value={startTime} disabled />
                <label className="timeslot-label" htmlFor="group-size">Group Size:</label>
                <input className="group-size" type="number" value={groupSizes[0]} min={1} onChange={ (event) => handleGroupSizeChange(0, event.target) } required></input>
                <button id="newline-btn" onClick={handleAddInput} hidden={size != 1}>+</button>
            </div>
        )
        for (let index = 1; index < size; ++index) {
            out.push(
                <div key={index} className='timeslot-editor-row'>
                    <label>{index+1}.</label>
                    <input type="time" value={inputs[index]} onChange={ (event) => handleInputChange(index, event.target) } />
                    <label className="timeslot-label" htmlFor="group-size">Group Size:</label>
                    <input className="group-size" type="number" value={groupSizes[index]} min={1} onChange={ (event) => handleGroupSizeChange(index, event.target) } required></input>
                    <button id="newline-btn" onClick={handleAddInput} hidden={size != index+1}><b>+</b></button>
                </div>
            )
        }
        return out;
    }
    const handleAddInput = () => {
        setSize(size+1);
        setTimes( { ...times, [size]: ""} );
        setGroupSizes( { ...groupSizes, [size]: defaultGroupSize} )
    };
    const handleInputChange = (index: number, element: HTMLInputElement) => {
        setTimes((prev) => ({ ...prev, [index]: element.value }));
    };
    const handleGroupSizeChange = (index: number, element: HTMLInputElement) => {
        if (size == index+1) {
            setDefaultGroupSize(Number(element.value));
        }
        setGroupSizes((prev) => ({ ...prev, [index]: Number(element.value) }));
    };
    const convertTime = (timeStr: string): number => {
        if (timeStr === "" || date === "") {
            console.warn(`Unable to create time from: '${timeStr}', '${date}'`)
            return 0;
        }
        let start = dateAndTimeToPosix(date, startTime);
        let thisTime = dateAndTimeToPosix(date, timeStr);
        if (thisTime < start) {
            thisTime = cycleDay(thisTime);
        }
        return thisTime;
    };
    const updateTimeslots = () => {
        let newTimeslots: Map<number, {"Size": number}> = new Map();
        for (let i = 0; i < size; ++i) {
            let thisTime = convertTime(times[i]);
            newTimeslots.set(thisTime, {"Size": groupSizes[i]});
        }
        timeslot.value = newTimeslots;
    };
    return (
        <>
            {renderTimeslots(size, times, handleInputChange)}
        </>
    );
};

export default TimeslotEditor;
