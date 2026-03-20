import "./EventCreation.css";

import { useEffect, useState } from "react";
import { useSignals } from "@preact/signals-react/runtime";
import { signal, type Signal } from "@preact/signals-react";

import { dateAndTimeToPosix, cycleDay, posixToTime, posixToDate } from "../../utils/utils";
import type { EventResponse } from "../../api/api";
import { makeEvent, editEvent } from "../../api/api";
import { loadDetails } from "../common/utils";
import { useTranslation } from "../../context/TranslationContext";

import Frame from "../common/Frame/Frame";
import Popup from "../Popup/Popup";
import TimeslotEditor from "./TimeslotEditor/TimeslotEditor";


const timeslotSignal = signal<Map<number, {"Size": number}>>(new Map());
const loadingEvents = signal(false);

interface Props {
    show: Signal<{eventID: string, view: string}>;
    update: ()=>Promise<void>
}

function EventCreation ({show, update}: Props) {
    useSignals();
    const {t} = useTranslation();

    // Input state information
    const [eventName, setEventName] = useState("");
    const [startDate, setStartDate] = useState("");
    const [startTime, setStartTime] = useState("");
    const [endTime,   setEndTime]   = useState("");
    const [shortDesc, setShortDesc] = useState("");
    const [longDesc, setLongDesc]   = useState("");

    // State info for editing event
    const [eventId, setEventId]     = useState(show.value.eventID);
    const [eventDetails, setEventDetails] = useState({
        Name: "<undefined>",
        ShortDescription: "<undefined>",
        LongDescription: "<undefined>",
        DtStart: 0,
        DtEnd: 0,
        Timeslots: {

        }
    } as EventResponse); // Makes sure no field is of undefined type

    // Dialog state information
    const [dialogText, setDialogText] = useState("---nothing---");
    const [confiramtionDialogVisible, setConfirmationDialogVisible] = useState(false);

    // Named input elements
    const dateInput  = document.getElementById("date");
    const startInput = document.getElementById("start-time");
    const endInput   = document.getElementById("end-time");

    const loadDetailsHandler = loadDetails(show, loadingEvents, setEventDetails);

    useEffect( () => {
        setEventId(show.value.eventID);
        if (show.value.eventID != "none") {
            loadDetailsHandler();
        } else {
            clearForm();
        }
    }, [show.value.view, show.value.eventID])

    useEffect( () => {
        setEventId(show.value.eventID);
        if (show.value.eventID != "none") {
            populateForm();
        }
    }, [eventDetails])

    const handleSaveEvent = (draft: boolean) => {
        try {
            const startTT = dateAndTimeToPosix(startDate, startTime);
            let endTT = dateAndTimeToPosix(startDate, endTime);
            if (endTT <= startTT) {
                endTT = cycleDay(endTT);
            }
            const timeslots = timeslotSignal.value;
            if (eventId == "none") {
                try {
                    makeEvent(eventName, shortDesc, longDesc, startTT, endTT, draft, 0, timeslots)
                        .then( (value ) => {
                            removeWrongLabelFromInputs(dateInput, startInput, endInput);
                            setDialogText(`Event created.\nEvent ID: ${value.EventID}\n${value.Error}`);
                            clearForm();
                            hideEditor(show);
                            update();
                    });
                } catch (error: any) {
                    setDialogText( `${error.message}\n`);
                    console.warn(error.message);
                }
            } else {
                try {
                    editEvent(eventId, eventName, shortDesc, longDesc, startTT, endTT, draft, 0, timeslots)
                        .then( (value ) => {
                            removeWrongLabelFromInputs(dateInput, startInput, endInput);
                            setDialogText(`Event created.\nEvent ID: ${value.EventID}\n${value.Error}`);
                            clearForm();
                            hideEditor(show);
                            update();
                    });
                } catch (error: any) {
                    setDialogText( `${error.message}\n`);
                    console.warn(error.message);
                }
                setConfirmationDialogVisible(true);
            }

        } catch (error) {
            console.error(error);
            console.error(`Failed to create timestamp from: '${startDate}', '${startTime}', '${endTime}'`);
            labelInputsAsWrong(dateInput, startInput, endInput);
        }
    };

    const clearForm = () => {
        setEventName("");
        setShortDesc("");
        setLongDesc("");
        setStartDate("");
        setStartTime("");
        setEndTime("");
        timeslotSignal.value = new Map()
    };

    const populateForm = () => {
        setEventName(eventDetails.Name);
        setShortDesc(eventDetails.ShortDescription);
        setLongDesc(eventDetails.LongDescription);
        setStartDate(posixToDate(eventDetails.DtStart));
        setStartTime(posixToTime(eventDetails.DtStart));
        setEndTime(posixToTime(eventDetails.DtEnd));
        try {
            timeslotSignal.value = new Map(
                Object.entries(eventDetails.Timeslots)
                    .map( ([k, v]) => [Number(k), {Size: v.Size}] )
            );
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <>
            <Frame className="EventForm" hidden={show.value.view != "editor"}>
                {eventId != "none" && `${t("event.editing")} #${eventId} ${eventDetails.Draft ? `- (${t("event.draft")})` : ""}`}
                <div className="header">
                    <input
                        id="event-name"
                        value={eventName}
                        onChange={e => setEventName(e.target.value)}
                        placeholder={t("event.name")} required>
                    </input>
                    <div id="close-box">
                        <button id="close" onClick={ () => hideEditor(show) }>{t("common.close")}</button>
                    </div>
                </div>

                <label className="form-label" htmlFor="date">{t("event.date")}</label>
                <input id="date" type="date" value={startDate} onChange={e => setStartDate(e.target.value)} required></input>

                <label className="form-label" htmlFor="start-time">{t("event.start")}</label>
                <input id="start-time" type="time" value={startTime} onChange={e => setStartTime(e.target.value)} required></input>

                <label className="form-label" htmlFor="end-time">{t("event.end")}</label>
                <input id="end-time" type="time" value={endTime} onChange={e => setEndTime(e.target.value)} required></input>

                <label className="form-label" htmlFor="short-description">{t("event.short name")}</label>
                <input id="short-desctiption" value={shortDesc} onChange={e => setShortDesc(e.target.value)} required></input>

                <label className="form-label" htmlFor="event-description">{t("event.desc")}</label>
                <textarea id="event-desctiption" value={longDesc} onChange={e => setLongDesc(e.target.value)} required></textarea>

                <label className="form-label" htmlFor="timerslots">{t("event.timeslots")}:</label>
                <div className="timeslots">
                    <TimeslotEditor startTime={startTime} date={startDate} timeslot={timeslotSignal} />
                </div>
                <div className="buttons editor-buttons">
                    <button id="publish" className="selected"  onClick={ () => handleSaveEvent(false) }>{t("event.publish")}</button>
                    <button id="save" className="yellow" onClick={ () => handleSaveEvent(true) }>{t("event.save-draft")}</button>
                </div>
            </Frame>
            <Popup show={confiramtionDialogVisible} onHide={() => setConfirmationDialogVisible(false)}>
                {dialogText}
            </Popup>
        </>
    );
}

export default EventCreation;


const hideEditor = (show: Signal<{eventID: string, view: string}>) => {
    show.value = {"eventID": "none", view: ""};
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
