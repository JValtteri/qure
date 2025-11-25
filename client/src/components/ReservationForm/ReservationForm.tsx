import { signal, Signal } from '@preact/signals-react';
import { useSignals } from "@preact/signals-react/runtime";

import Dialog from '../common/Dialog/Dialog';
import TimeslotList from '../TimeslotList/TimeslotList';

import './ReservationForm.css';

import type { Timeslot } from '../../api/api';
import { useState } from 'react';


const selectedSlot = signal(-1);

interface Props {
    showDialog: Signal<boolean>;
    timeslots: Map<number, Timeslot>;
}

function ReservationForm({showDialog, timeslots}: Props) {
    useSignals();
    const [email, setEmail] = useState("");
    const [groupSize, setGroupSize] = useState(0);

    return(
        <Dialog className='reservation' hidden={ showDialog.value===false }>
            <div className="header-container">
                <h3>Reservation</h3>
                <button onClick={ () => showDialog.value=false }>Cancel</button>
            </div>
            <div></div>

            <label>Select Group</label>
            <div>
                <TimeslotList timeslots={timeslots} selectedSlot={selectedSlot} />
            </div>

            <label className="form-label" htmlFor="reserve-email">Email</label>
            <input id="reserve-email" type="email" value={email} placeholder='example@email.com' onChange={e => setEmail(e.target.value)} required ></input>

            <label className="form-label" htmlFor="group-size">Group Size</label>
            <input id="group-size" type="number" value={groupSize} min={1} placeholder='example@email.com' onChange={e => setGroupSize(Number(e.target.value))} required ></input>
            <hr></hr>
            <div className='buttons 2'>
                <button className='selected'>Reserve</button>
            </div>
        </Dialog>
    );
}

export default ReservationForm;
