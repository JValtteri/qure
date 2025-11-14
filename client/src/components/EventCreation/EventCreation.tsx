import { useState } from "react";
import { useSignals } from "@preact/signals-react/runtime";
import { type Signal } from "@preact/signals-react";

import Frame from "../common/Frame/Frame";
import Popup from "../Popup/Popup";
import TimeslotEditor from "./TimeslotEditor/TimeslotEditor";

import { dateAndTimeToPosix, cycleDay } from "../../utils/utils";
import { makeEvent } from "../../api/api";

import "./EventCreation.css";


interface Props {
    show: Signal<{ "selectedEventId": number, "editor": boolean}>;
    update: ()=>Promise<void>
}

const hideEditor = () => {
    return {"selectedEventId": -1, "editor": false};
}

function EventCreation ({show, update}: Props) {
    useSignals();

    let [eventName, setEventName] = useState("New Event");
    let [startDate, setStartDate] = useState("");
    let [startTime, setStartTime] = useState("");
    let [endTime, setEndTime] = useState("");
    let [shortDesc, setShortDesc] = useState("");
    let [longDesc, setLongDesc] = useState("");

    let [dialogText, setDialogText] = useState("---nothing---");
    let [dialogVisible, setDialogVisible] = useState(false);

    const dateInput = document.getElementById("date");
    const startInput = document.getElementById("start-time");
    const endInput = document.getElementById("end-time");

    console.log("EventCreation rendered");

    const genericHandleSaveEvent = (draft: boolean) => {
        try {
            let startTT = dateAndTimeToPosix(startDate, startTime);
            let endTT = dateAndTimeToPosix(startDate, endTime);
            if (endTT <= startTT) {
                endTT = cycleDay(endTT);
            }
            let timeslots = [startTT]
            makeEvent(eventName, shortDesc, longDesc, startTT, endTT, draft, 0, timeslots, 999)
                .then( (value ) => {
                    dateInput?.classList.remove("wrong");
                    startInput?.classList.remove("wrong");
                    endInput?.classList.remove("wrong");
                    setDialogVisible(true);
                    setDialogText( `Event created.\nEvent ID: ${value.EventID}\n${value.Error}`);
                    update();
                });
        } catch (error) {
            console.error(error);
            console.error(`Failed to create timestamp from: '${startDate}', '${startTime}', '${endTime}'`);
            dateInput?.classList.add("wrong");
            startInput?.classList.add("wrong");
            endInput?.classList.add("wrong");
        }
    };

    return (
        <Frame className="EventForm" hidden={!show.value.editor}>
            <div className="header">
                <input id="event-name" value={eventName} onChange={e => setEventName(e.target.value)} placeholder="Event Name" required></input>
                <div id="close-box">
                    <button id="close" onClick={ () => show.value = hideEditor() }>Close</button>
                </div>
            </div>

            <label className="form-label" htmlFor="date">Date</label>
            <input id="date" type="date" value={startDate} onChange={e => setStartDate(e.target.value)} required></input>

            <label className="form-label" htmlFor="start-time">Start Time</label>
            <input id="start-time" type="time" value={startTime} onChange={e => setStartTime(e.target.value)} required></input>

            <label className="form-label" htmlFor="end-time">End Time</label>
            <input id="end-time" type="time" value={endTime} onChange={e => setEndTime(e.target.value)} required></input>

            <label className="form-label" htmlFor="short-description">Short Description</label>
            <input id="short-desctiption" onChange={e => setShortDesc(e.target.value)} required></input>

            <label className="form-label" htmlFor="event-description">Event Description</label>
            <textarea id="event-desctiption" onChange={e => setLongDesc(e.target.value)} required></textarea>

            <label className="form-label" htmlFor="timerslots">Timeslots:</label>
            <div className="timeslots">
                <TimeslotEditor firstTime={startTime} />
            </div>
            <div className="buttons editor-buttons">
                <button id="publish" onClick={ () => genericHandleSaveEvent(false) }>Publish</button>
                <button id="save" onClick={ () => genericHandleSaveEvent(true) }>Save as Draft</button>
            </div>
            <Popup show={dialogVisible} onHide={() => setDialogVisible(false)}>
                {dialogText}
            </Popup>
        </Frame>
    );
}

export default EventCreation;
