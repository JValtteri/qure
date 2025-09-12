import Frame from "../Frame/Frame";

interface Props {
    time?: string;
    date?: string;
    desc?: string;
    title: string;
    slots: number;
    occumpied: number;
    onClick: any;
    className: string;
}

function EventCard({time, date, desc, title, slots, occumpied, onClick, className}: Props) {
    return(
        <Frame reactive={true} className={className} onClick={onClick}>
            <div>{title}</div>
            <div>{desc}</div>
            <div>{`${date} ${time}`}</div>
            <div>{`${occumpied}/${slots}`}</div>
        </Frame>
    )
}

export default EventCard;
