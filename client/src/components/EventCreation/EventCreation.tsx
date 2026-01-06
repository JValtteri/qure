import "./EventCreation.css";

import { useState } from "react";
import { useSignals } from "@preact/signals-react/runtime";
import { signal, type Signal } from "@preact/signals-react";

import Frame from "../common/Frame/Frame";
import Popup from "../Popup/Popup";
import TimeslotEditor from "./TimeslotEditor/TimeslotEditor";

import { dateAndTimeToPosix, cycleDay } from "../../utils/utils";
import { makeEvent } from "../../api/api";


const timeslotSignal = signal<Map<number, {"Size": number}>>(new Map());

interface Props {
    show: Signal<{"eventID": number, "editor": boolean}>;
    update: ()=>Promise<void>
}

function EventCreation ({show, update}: Props) {
    useSignals();

    // Input state information
    const [eventName, setEventName] = useState("New Event");
    const [startDate, setStartDate] = useState("");
    const [startTime, setStartTime] = useState("");
    const [endTime,   setEndTime]   = useState("");
    const [shortDesc, setShortDesc] = useState("");
    const [longDesc, setLongDesc]   = useState("");

    // Dialog state information
    const [dialogText, setDialogText] = useState("---nothing---");
    const [confiramtionDialogVisible, setConfirmationDialogVisible] = useState(false);

    // Named input elements
    const dateInput  = document.getElementById("date");
    const startInput = document.getElementById("start-time");
    const endInput   = document.getElementById("end-time");

    const handleSaveEvent = (draft: boolean) => {
        try {
            const startTT = dateAndTimeToPosix(startDate, startTime);
            let endTT = dateAndTimeToPosix(startDate, endTime);
            if (endTT <= startTT) {
                endTT = cycleDay(endTT);
            }
            const timeslots = timeslotSignal.value;
            makeEvent(eventName, shortDesc, longDesc, startTT, endTT, draft, 0, timeslots)
                .then( (value ) => {
                    removeWrongLabelFromInputs(dateInput, startInput, endInput);
                    setConfirmationDialogVisible(true);
                    setDialogText( `Event created.\nEvent ID: ${value.EventID}\n${value.Error}`);
                    clearForm();
                    hideEditor(show);
                    update();
                });
        } catch (error) {
            console.error(error);
            console.error(`Failed to create timestamp from: '${startDate}', '${startTime}', '${endTime}'`);
            labelInputsAsWrong(dateInput, startInput, endInput);
        }
    };

    const clearForm = () => {
        setEventName("New Event");
        setShortDesc("");
        setLongDesc("");
        setStartDate("");
        setStartTime("");
        setEndTime("");
        timeslotSignal.value = new Map()
    }

    return (
        <>
            <Frame className="EventForm" hidden={!show.value.editor}>
                <div className="header">
                    <input id="event-name" value={eventName} onChange={e => setEventName(e.target.value)} placeholder="Event Name" required></input>
                    <div id="close-box">
                        <button id="close" onClick={ () => hideEditor(show) }>Close</button>
                    </div>
                </div>

                <label className="form-label" htmlFor="date">Date</label>
                <input id="date" type="date" value={startDate} onChange={e => setStartDate(e.target.value)} required></input>

                <label className="form-label" htmlFor="start-time">Start Time</label>
                <input id="start-time" type="time" value={startTime} onChange={e => setStartTime(e.target.value)} required></input>

                <label className="form-label" htmlFor="end-time">End Time</label>
                <input id="end-time" type="time" value={endTime} onChange={e => setEndTime(e.target.value)} required></input>

                <label className="form-label" htmlFor="short-description">Short Description</label>
                <input id="short-desctiption" value={shortDesc} onChange={e => setShortDesc(e.target.value)} required></input>

                <label className="form-label" htmlFor="event-description">Event Description</label>
                <textarea id="event-desctiption" onChange={e => setLongDesc(e.target.value)} required></textarea>

                <label className="form-label" htmlFor="timerslots">Timeslots:</label>
                <div className="timeslots">
                    <TimeslotEditor startTime={startTime} date={startDate} timeslot={timeslotSignal} />
                </div>
                <div className="buttons editor-buttons">
                    <button id="publish" onClick={ () => handleSaveEvent(false) }>Publish</button>
                    <button id="save" onClick={ () => handleSaveEvent(true) }>Save as Draft</button>
                </div>
            </Frame>
            <Popup show={confiramtionDialogVisible} onHide={() => setConfirmationDialogVisible(false)}>
                {dialogText}
            </Popup>
        </>
    );
}

export default EventCreation;


const hideEditor = (show: Signal<{"eventID": number, "editor": boolean}>) => {
    show.value = {"eventID": -1, "editor": false};
}

function labelInputsAsWrong(dateInput: HTMLElement | null, startInput: HTMLElement | null, endInput: HTMLElement | null) {
    dateInput?.classList.add("wrong");
    startInput?.classList.add("wrong");
    endInput?.classList.add("wrong");
}

function removeWrongLabelFromInputs(dateInput: HTMLElement | null, startInput: HTMLElement | null, endInput: HTMLElement | null) {
    dateInput?.classList.remove("wrong");
    startInput?.classList.remove("wrong");
    endInput?.classList.remove("wrong");
}

