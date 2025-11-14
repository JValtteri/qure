import { useState, type ReactNode } from 'react';
import "./TimeslotEditor.css"


interface Props {
    firstTime: string;
}

function TimeslotEditor({firstTime}: Props) {
    const [inputs, setInputs] = useState<{ [label: number]: string }>({[0]: firstTime});
    const [groupSizes, setGroupSizes] = useState<{ [label: number]: number }>({[0]: 0});
    const [defaultGroupSize, setDefaultGroupSize] = useState(0);
    const [size, setSize] = useState(1);

    const renderTimeslots = (size: number, inputs: {[label: number]: string;}, handleInputChange: (index: number, element: HTMLInputElement)=>void) => {
        inputs[0] = firstTime;
        let out = new Array<ReactNode>;
        out.push(
            <div key={0} className='timeslot-editor-row'>
                <label>1.</label>
                <input type="time" value={firstTime} disabled />
                <label className="timeslot-label" htmlFor="group-size">Group Size:</label>
                <input className="group-size" type="number" value={groupSizes[0]} min={1} onChange={ (event) => handleGroupSizeChange(0, event.target) } required></input>
                <button id="newline-btn" onClick={handleAddInput} hidden={size != 1}>+</button>
            </div>
        )
        for (let index = 1; index < size; ++index) {
            out.push(
                <div key={index} className='timeslot-editor-row'>
                    <label>{index+1}.</label>
                    <input type="time" value={inputs[index]} onChange={ (event) => handleInputChange(index, event.target) } />
                    <label className="timeslot-label" htmlFor="group-size">Group Size:</label>
                    <input className="group-size" type="number" value={groupSizes[index]} min={1} onChange={ (event) => handleGroupSizeChange(index, event.target) } required></input>
                    <button id="newline-btn" onClick={handleAddInput} hidden={size != index+1}><b>+</b></button>
                </div>
            )
        }
        return out;
    }
    const handleAddInput = () => {
        setSize(size+1);
        setInputs( { ...inputs, [size]: ""} );
        setGroupSizes( { ...groupSizes, [size]: defaultGroupSize} )
    };
    const handleInputChange = (index: number, element: HTMLInputElement) => {
        setInputs((prev) => ({ ...prev, [index]: element.value }));
    };
    const handleGroupSizeChange = (index: number, element: HTMLInputElement) => {
        if (size == index+1) {
            setDefaultGroupSize(Number(element.value));
        }
        setGroupSizes((prev) => ({ ...prev, [index]: Number(element.value) }));
    };
    return (
        <>
            {renderTimeslots(size, inputs, handleInputChange)}
        </>
    );
};

export default TimeslotEditor;
