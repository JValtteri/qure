import Frame from "../Frame/Frame";
import './EventCard.css'

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
            <div className="primary">{title}</div>
            <div className="secondary right">{`${date}`}</div>
            <div className="secondary desc">{desc}</div>
            <div className="secondary right">{`${occumpied}/${slots}`}</div>
        </Frame>
    )
}

export default EventCard;
