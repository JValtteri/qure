import Frame from "../Frame/Frame";
import './EventCard.css';

interface Props {
    time?: string;
    startTime?: string;
    desc?: string;
    title: string;
    slots: number;
    occupied: number;
    onClick: any;
    selected: boolean;
    className?: string;
}

function EventCard({startTime, desc, title, slots, occupied, onClick, selected, className}: Props) {
    return(
        <Frame reactive={true} className={`${selected ? "card selected" : "card" } ${className}`} onClick={onClick}>
            <div className="primary">{title}</div>
            <div className="secondary right">{`${startTime}`}</div>
            <div className="secondary desc">{desc}</div>
            <div className="secondary right">{`${occupied}/${slots}`}</div>
        </Frame>
    )
}

export default EventCard;
