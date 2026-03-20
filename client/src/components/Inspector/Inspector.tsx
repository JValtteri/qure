import "./Inspector.css"

import { useState, useEffect } from "react";
import { Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import { listEventReservations, type ReservationResponse } from '../../api/api';
import { useTranslation } from "../../context/TranslationContext";

import Frame from "../common/Frame/Frame";
import GenericTable from "../common/GenericTable/GenericTable";


interface Props {
    show: Signal<{eventID: string, view: string}>;
    className?: string;
    hidden: boolean;
}

function Inspector({show, className, hidden}: Props) {
    const [reservations, setReservations] = useState(new Array<ReservationResponse>());
    const updateReservationsHandler = updateReservations(show.value.eventID, setReservations);
    useSignals();
    const {t} = useTranslation();


    useEffect(() => {
        if (show.value.view == "inspect") {
            updateReservationsHandler();
        }
    }, [show.value.eventID]);

    return (
        <Frame hidden={hidden} className={className}>
            <div className="table-container">
                <h2>{ reservations.length > 0 ? reservations.at(0)?.Event.Name : ""}</h2>
                <GenericTable
                    data={filterDuplicate(reservations)}
                    rowKey={'Id'}
                    columns={["Id", "Timeslot", "Size"]}
                    onRowClick={ ()=>{} }
                    filterable={true}
                    sortable={true}
                    defaultSortColumn={'Timeslot'}
                    interpretBigNumbersAs='time'
                />
                <p className="noMatches noTopPad">
                    {show.value.eventID == "none" && t("tools.select-event")}
                </p>
            </div>
        </Frame>
    )
}


function updateReservations(id: string, setReservations: React.Dispatch<React.SetStateAction<Array<ReservationResponse>>>): () => Promise<void> {
    return async () => {
        try {
            await listEventReservations(id)
                .then(value => {
                    if (value != null) {
                        setReservations(value);
                    } else {
                        setReservations([])
                    }
                });
        } catch (error: any) {
            console.warn(error.message);
        }
    };
}

function filterDuplicate(data: ReservationResponse[]) {
    const seenIds = new Set<number | string>();
    const uniqueData = data.filter((row) => {   // Filter duplicate rows
        if (seenIds.has(row.Id)) {
            return false;
        }
        seenIds.add(row.Id);
        return true;
    });
    return uniqueData;
}

/*
// Snippet of the old code. This is how the highlighting was done:
function populateLines(data: ReservationResponse[]): ReactNode {
    const firstTimeslot = data.at(0)?.Timeslot;
    const uniqueData = filterDuplicate(data);
    return uniqueData.map((row) => (
        <tr key={row.Id} className={row.Timeslot == firstTimeslot ? "first-in-queue" : ""}>
            <td>#{row.Id}</td>
            <td>{posixToTime(row.Timeslot)}</td>
            <td>{row.Size}</td>
        </tr>
    ));
}
*/

export default Inspector
