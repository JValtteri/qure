import type { ReactNode } from "react";
import { Signal } from "@preact/signals-react";
import { useSignals } from "@preact/signals-react/runtime";

import Frame from "../common/Frame/Frame";

//import "./Inspector.css"


interface Props {
    show: Signal<{"eventID": string, "editor": boolean, "account": boolean, "inspect": boolean}>;
    children?: ReactNode;
    className?: string;
    hidden: boolean;
}


function Inspector({show, children, className, hidden}: Props) {

    useSignals();

    return (
        <Frame hidden={hidden} className={className}>
            {createTable(data)}

            {children}


            // Get Reservations for {show.value.eventID}
            // Loop for in reservations
            // tr th: [ID, seats]
            // --- tr: td: [ID, seats]
        </Frame>
    )
}


function createTable(data): HTMLTableElement {
    data = applySort(data);
    let table = document.createElement('table');
    table = createTitle(table);
    table = populateLines(table, data);
    return table;
}

function applySort(data) {

    return data;
}

function createTitle(table: HTMLTableElement): HTMLTableElement {
    let titleRow = document.createElement('tr');
    for (const title in ["ID", "Time", "Size"]) {
        const titleElm = document.createElement('th');
        titleElm.innerText = title;
        titleRow.appendChild(titleElm);
    }
    return table;
}

function populateLines(table: HTMLTableElement, data): HTMLTableElement {
    for (const line in data) {
        const rowElement = makeLine(line)
        table.appendChild(rowElement)
    }

    return table;
}

function makeLine(line): HTMLTableRowElement {
    const row = document.createElement('tr');
    for (const item in line) {
        const cell = document.createElement('td');
        cell.innerText = item;
        row.appendChild(cell);
    }
    return row;
}


export default Inspector
