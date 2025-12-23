import { posixToDateAndTime } from "../../../../utils/utils";

interface Props {
    size: number;
    confirmed: number
    time: number;
    code: string;
}

function ReservePartial({size, confirmed, time, code}: Props) {
    return (
        <>
            <h3 className='centered'>
                You are in queue
            </h3>
            <p className='centered'>
                Reserved <b>{confirmed}</b> place(s) for <b>{posixToDateAndTime(time)}</b>.
            </p>
            <p className='centered'>
                and <b>{size-confirmed}</b> place(s) is in queue.
            </p>
            <label className="small-label">Your reservation ID:</label>
            <p className="centered reservation-code">
                {code}
            </p>
        </>
    )
}

export default ReservePartial;
