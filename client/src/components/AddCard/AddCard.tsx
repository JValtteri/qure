import './AddCard.css';


interface Props {
    onClick: any;
    hidden?: boolean;
}

function AddCard({onClick, hidden}: Props) {
    return (
        <button className="add-card" onClick={onClick} hidden={hidden}>
            ➕
        </button>
    )
}

export default AddCard
