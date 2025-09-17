import Frame from '../Frame/Frame';
import './TitleBar.css';

interface Props {
    title?: string
}

function TitleBar({title}: Props) {
    return (
        <Frame>
            {title ? title : "Title"}
            <Frame reactive={true}>
                Login
            </Frame>
        </Frame>
    )
}

export default TitleBar;
