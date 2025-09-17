import Frame from '../Frame/Frame';
import './TitleBar.css';

interface Props {
    title?: string
}

function TitleBar({title}: Props) {
    return (
        <Frame className='title'>
            {title ? title : "Title"}
            <button>Login</button>
        </Frame>
    )
}

export default TitleBar;
