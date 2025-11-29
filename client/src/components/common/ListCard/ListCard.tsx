import './ListCard.css';

import Frame from "../Frame/Frame";


interface Props {
    startTime?: string;
    desc?: string;
    title: string;
    slots: number;
    occupied: number;
    onClick: ()=>void;
    selected: boolean;
    className?: string;
}

function ListCard({startTime, desc, title, slots, occupied, onClick, selected, className}: Props) {
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

export default ListCard;
