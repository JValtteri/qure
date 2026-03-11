import "./TimeslotEditor.css";

import { useEffect, useState } from 'react';
import type { Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import { cycleDay, dateAndTimeToPosix, posixToTime } from '../../../utils/utils';

interface Props {
    startTime: string;
    date: string;
    timeslot: Signal<Map<number, { Size: number }>>;
}

interface TimeslotData {
    time: string;
    groupSize: number;
}

function TimeslotEditor({ startTime, date, timeslot }: Props) {
    useSignals();

    const [initialLoad, setInitialLoad] = useState(true);
    const [timeslots, setTimeslots] = useState<TimeslotData[]>([{ time: startTime, groupSize: 0 }]);
    const [lastGroupSize, setLastGroupSize] = useState(0);

    // Update when the timeslot signal is changed externally
    useEffect(() => {
        if (!initialLoad) return;
        if (timeslot.value.size === 0) return;
        const sortedTimeslots = sortTimeslots();
        setTimeslots(sortedTimeslots);
    }, [timeslot.value]);

    // Update when internal state is changed
    useEffect(() => {
        if (initialLoad) return;
        updateTimeslotSignal();
    }, [timeslots, date, startTime]);


    const sortTimeslots = () => {
        return Array.from(timeslot.value.entries())
            //.sort((a, b) => a[0] - b[0])              // Why the sort???
            .map(([time, data]) => ({
                time: posixToTime(time),
                groupSize: data.Size
            }));
    }

    const updateTimeslotSignal = () => {
        const newTimeslots = new Map<number, { Size: number }>();
        timeslots.forEach((timeslotData, _) => {
            var posixTime = 0;
            try {
                posixTime = convertTime(timeslotData.time);
            } catch (error) {
                console.error(`Time parsing failed: ${error}`);
            }
            newTimeslots.set(posixTime, { Size: timeslotData.groupSize });
        });
        timeslot.value = newTimeslots;
        setInitialLoad(false);
    }

    const convertTime = (timeStr: string): number => {
        if (timeStr === "" || date === "") {
            return 0;
        }
        const start = dateAndTimeToPosix(date, startTime);
        let thisTime = dateAndTimeToPosix(date, timeStr);
        if (thisTime < start) {
            thisTime = cycleDay(thisTime);
        }
        return thisTime;
    };

    const handleTimeChange = (index: number, newTime: string) => {      // Is this really the best way?
        setTimeslots(prev => prev.map((timeslot, i) =>
            i === index ? { ...timeslot, time: newTime } : timeslot
        ));
        updateTimeslotSignal();
    };

    const handleGroupSizeChange = (index: number, newSize: number) => {      // Is this really the best way?
        setTimeslots(prev => prev.map((timeslot, i) =>
            i === index ? { ...timeslot, groupSize: newSize } : timeslot
        ));
        updateTimeslotSignal();
        setLastGroupSize(newSize);
    };

    const addTimeslot = () => {
        setTimeslots(prev => [...prev, { time: "", groupSize: lastGroupSize }]);
    };

    const removeTimeslot = (index: number) => {
        setTimeslots(prev => prev.filter((_, i) => i !== index));
        updateTimeslotSignal();
    };

    const renderTimeslotRow = (timeslot: TimeslotData, index: number) => {
        const isFirstRow = index === 0;
        const isLastRow = index === timeslots.length - 1;

        return (
            <div key={index} className='timeslot-editor-row'>
                <label>{index + 1}.</label>
                <input
                    type="time"
                    value={ isFirstRow ? startTime : timeslot.time}
                    onChange={(e) => handleTimeChange(index, e.target.value)}
                    disabled={isFirstRow}
                />
                <label className="timeslot-label" htmlFor="group-size">Group Size:</label>
                <input
                    className="group-size"
                    type="number"
                    value={timeslot.groupSize}
                    min={1}
                    onChange={(e) => handleGroupSizeChange(index, Number(e.target.value))}
                    required
                />
                {!isFirstRow && (
                    <button id="newline-btn" onClick={ () => removeTimeslot(index) }>
                        <b>-</b>
                    </button>
                )
                }
                {isLastRow && (
                    <button id="newline-btn" onClick={addTimeslot}>
                        <b>+</b>
                    </button>
                )}
            </div>
        );
    };

    return (
        <>
            {timeslots.map((timeslot, index) => (
                renderTimeslotRow(timeslot, index)
            ))}
        </>
    );
}

export default TimeslotEditor;
