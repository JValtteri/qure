import { posixToDateAndTime } from "../../../../utils/utils";

interface Props {
    size: number;
    time: number;
    code: string;
}

function ReserveSuccess({size, time, code}: Props) {
    return (
        <>
            <h3 className='centered'>
                Reservation successfull
            </h3>
            <p className='centered'>
                Reserved <b>{size}</b> place(s) for <b>{posixToDateAndTime(time)}</b>.
            </p>
            <label className="small-label">Your reservation ID:</label>
            <p className="centered reservation-code">
                {code}
            </p>
        </>
    )
}

export default ReserveSuccess;
