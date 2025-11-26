
interface Props {
    error: string;
}

function ReserveFailed({error}: Props) {
    return (
        <>
            <h3>
                Reservation Failed
            </h3>
        {error}
        </>
    );
}

export default ReserveFailed;
