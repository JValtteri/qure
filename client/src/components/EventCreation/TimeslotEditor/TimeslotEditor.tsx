import "./TimeslotEditor.css";
import { useEffect, useState } from 'react';
import { type Signal } from "@preact/signals-react";
import { cycleDay, dateAndTimeToPosix } from '../../../utils/utils';
import { useSignals } from "@preact/signals-react/runtime";

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
    const [timeslots, setTimeslots] = useState<TimeslotData[]>([{ time: startTime, groupSize: 0 }]);
    const [initialLoad, setInitialLoad] = useState(true);

    // Update when the timeslot signal is changed externally
    useEffect(() => {
        if (!initialLoad) return;
        if (timeslot.value.size === 0) return;

        const sortedTimeslots = sortTimeslots();

        // Only update if the data is actually different
        const isDifferent = JSON.stringify(sortedTimeslots) !== JSON.stringify(timeslots);
        if (isDifferent) {
            setTimeslots(sortedTimeslots);
        } else {
            console.warn("Really!? New external data was identical")
        }
    }, [timeslot.value]);

    // Update when internal state is changed
    useEffect(() => {
        if (initialLoad) return;
        updateTimeslotSignal();
    }, [timeslots, date, startTime]);


    const sortTimeslots = () => {
        return Array.from(timeslot.value.entries())
            .sort((a, b) => a[0] - b[0])
            .map(([time, data]) => ({
                time: new Date(time * 1000).toISOString().slice(11, 16),
                groupSize: data.Size
            }));
    }

    const updateTimeslotSignal = () => {
        const newTimeslots = new Map<number, { Size: number }>();
        timeslots.forEach((timeslotData, index) => {
            var posixTime = 0;
            try {
                posixTime = convertTime(timeslotData.time);
            } catch (error) {
                console.error(error);
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

    const handleTimeChange = (index: number, newTime: string) => {
        setTimeslots(prev => prev.map((timeslot, i) =>
            i === index ? { ...timeslot, time: newTime } : timeslot
        ));
        updateTimeslotSignal();
    };

    const handleGroupSizeChange = (index: number, newSize: number) => {
        setTimeslots(prev => prev.map((timeslot, i) =>
            i === index ? { ...timeslot, groupSize: newSize } : timeslot
        ));
        updateTimeslotSignal();
    };

    const addTimeslot = () => {
        setTimeslots(prev => [...prev, { time: "", groupSize: 0 }]);
    };

    const renderTimeslotRow = (timeslot: TimeslotData, index: number) => {
        const isFirstRow = index === 0;
        const isLastRow = index === timeslots.length - 1;

        return (
            <div key={index} className='timeslot-editor-row'>
                <label>{index + 1}.</label>
                <input
                    type="time"
                    value={timeslot.time}
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
