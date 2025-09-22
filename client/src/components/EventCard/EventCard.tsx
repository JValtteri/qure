import { Signal } from "@preact/signals-react";
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
    console.log("Card rendered")
    const baseClasses = selected ? "card selected" : "card";
    const additionalClasses = className==undefined ? "" : className;
    return(
        <Frame reactive={true} className={`${ baseClasses } ${ additionalClasses }`} onClick={onClick}>
            <div className="primary">{title}</div>
            <div className="secondary right">{`${startTime}`}</div>
            <div className="secondary desc">{desc}</div>
            <div className="secondary right">{`${occupied}/${slots}`}</div>
        </Frame>
    )
}

export default EventCard;
