import { posixToDateAndTime } from "../../../../utils/utils";

interface Props {
    size: number;
    confirmed: number
    time: number;
}

function ReservePartial({size, confirmed, time}: Props) {
    return (
        <>
            <h3>
                You are in queue
            </h3>
            <p>
                Reserved <b>{confirmed}</b> slots for <b>{posixToDateAndTime(time)}</b>.
            </p>
            <p>
                If more slots open up, the remaining <b>{size-confirmed}</b> will be added.
            </p>
        </>
    )
}

export default ReservePartial;
