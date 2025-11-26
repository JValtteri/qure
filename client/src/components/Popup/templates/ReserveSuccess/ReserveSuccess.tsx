import { posixToDateAndTime } from "../../../../utils/utils";

interface Props {
    size: number;
    time: number;
}

function ReserveSuccess({size, time}: Props) {
    return (
        <>
            <h3>
                Reservation successfull
            </h3>
            Reserved <b>{size}</b> slots for <b>{posixToDateAndTime(time)}</b>.
        </>
    )
}

export default ReserveSuccess;
