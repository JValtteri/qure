import type { ReactNode } from "react";
import { useSignals } from "@preact/signals-react/runtime";
import Dialog from "../common/Dialog/Dialog";


interface Props {
    children?: ReactNode;
    className?: string;
    show: boolean;
    onHide: any;
}

function Popup({children, className, show, onHide}: Props) {
    useSignals();
    return (
        <Dialog hidden={!show}>
            <pre>
                {children}
            </pre>
            <div className="buttons-center">
                <button id="ok" onClick={ () => onHide() }>Ok</button>
            </div>
        </Dialog>
    )
}

export default Popup
