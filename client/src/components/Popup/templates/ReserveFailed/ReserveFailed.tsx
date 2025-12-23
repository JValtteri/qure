
interface Props {
    error: string;
}

function ReserveFailed({error}: Props) {
    return (
        <>
            <h3 className='centered'>
                Reservation Failed
            </h3>
        <p className='centered'>{error}</p>
        </>
    );
}

export default ReserveFailed;
