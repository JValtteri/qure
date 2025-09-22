import type { ReactNode } from "react";
import { Signal } from "@preact/signals-react";
import Frame from "../Frame/Frame";

interface Props {
    selectedId: Signal<number>;
    children: ReactNode
}

function DetailCard( {selectedId, children}: Props ) {
    console.log("Detail rendered")
    return (
        <Frame hidden={selectedId.value === -1}>
            {children}
        </Frame>
    )
}

export default DetailCard;
