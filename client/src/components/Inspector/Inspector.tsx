import { useState, useEffect, type ReactNode } from "react";
import { Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import { listEventReservations, type ReservationResponse } from '../../api/api';

import Frame from "../common/Frame/Frame";

import "./Inspector.css"


interface Props {
    show: Signal<{"eventID": string, "editor": boolean, "account": boolean, "inspect": boolean}>;
    className?: string;
    hidden: boolean;
}

function Inspector({show, className, hidden}: Props) {
    const [reservations, setReservations] = useState(new Array<ReservationResponse>());
    const updateReservationsHandler = updateReservations(show.value.eventID, setReservations);
    useSignals();


    useEffect(() => {
        if (show.value.inspect) {
            updateReservationsHandler();
        }
    }, [show.value.eventID]);

    return (
        <Frame hidden={hidden} className={className}>
            <h2>{ reservations.length > 0 ? reservations.at(0)?.Event.Name : ""}</h2>
            {createTable(reservations)}
        </Frame>
    )
}


function updateReservations(id: string, setReservations: React.Dispatch<React.SetStateAction<Array<ReservationResponse>>>): () => Promise<void> {
    return async () => {
        await listEventReservations(id)
            .then(value => {
                if (value != null) {
                    setReservations(value);
                } else {
                    setReservations([])
                }
            });
    };
}

function createTable(data: ReservationResponse[]): ReactNode {
    data = applySort(data);
    const thead = createTitle();
    const tbody = populateLines(data);
    return (
        <table>
            <thead>
                <tr>{thead}</tr>
            </thead>
            <tbody>
                {tbody}
            </tbody>
        </table>
    );
}

function applySort(data: ReservationResponse[]): ReservationResponse[] {
    return data.sort((a: any, b: any) => b.Size - a.Size);
}

function createTitle(): ReactNode[] {
    const titles = ["Reservation", "Time", "Size"]
    return titles.map(title => (
        <th key={title}>{title}</th>
    ));
}

function populateLines(data: ReservationResponse[]): ReactNode {
    const seenIds = new Set<number | string>();
    const uniqueData = data.filter((row) => {   // Filter duplicate rows
        if (seenIds.has(row.Id)) {
            return false;
        }
        seenIds.add(row.Id);
        return true;
    });
    return uniqueData.map((row) => (
        <tr key={row.Id}>
            <td>#{row.Id}</td>
            <td>{row.Timeslot}</td>
            <td>{row.Size}</td>
        </tr>
    ));
}

export default Inspector
