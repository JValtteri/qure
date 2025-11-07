import { useState } from "react";
import { useSignals } from "@preact/signals-react/runtime";
import type { Signal } from "@preact/signals-react";
import Frame from "../Frame/Frame";
import "./EventCreation.css";


interface Props {
    show: Signal<boolean>;
}

function EventCreation ({show}: Props) {
    useSignals();
    console.log("EventCreation rendered")

    let [eventName, setEventName] = useState("New Event")
    let [startTime, setStartTime] = useState("")
    let [endTime, setEndTime] = useState("")
    let [slots, setSlots] = useState(1)

    return (
        <Frame className="EventForm" hidden={!show.value}>
            <div className="header">
                <input id="event-name" value={eventName} onChange={e => setEventName(e.target.value)} placeholder="Event Name" required></input>
                <div id="close-box">
                    <button id="close" onClick={ () => show.value=false }>Close</button>
                </div>
            </div>

            <label htmlFor="date">Date</label>
            <input id="date" type="date" value={endTime} onChange={e => setEndTime(e.target.value)} required></input>

            <label htmlFor="start-time">Start Time</label>
            <input id="start-time" type="time" value={startTime} onChange={e => setStartTime(e.target.value)} required></input>

            <label htmlFor="end-time">End Time</label>
            <input id="end-time" type="time" value={endTime} onChange={e => setEndTime(e.target.value)} required></input>

            <label htmlFor="event-description">Event Description</label>
            <textarea id="event-desctiption"></textarea>

            <div className="buttons">
                <button id="publish">Publish</button>
                <button id="save">Save as Draft</button>
            </div>
        </Frame>
    );
}

export default EventCreation;
