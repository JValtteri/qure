import Frame from "../Frame/Frame"
import './AddCard.css';


interface Props {
    onClick: any;
    hidden?: boolean;
}

function AddCard({onClick, hidden}: Props) {
    return (
        <Frame reactive={true} className="add-card" onClick={onClick} hidden={hidden}>
            ➕
        </Frame>
    )
}

export default AddCard
