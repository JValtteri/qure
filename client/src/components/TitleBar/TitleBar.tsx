import Frame from '../Frame/Frame';
import './TitleBar.css';

interface Props {
    title?: string
}

function TitleBar({title}: Props) {
    return (
        <Frame className='title'>
            <img src='./logo.png' />
            <span>
                {title ? title : "< Title >"}
            </span>
            <button>Login</button>
        </Frame>
    )
}

export default TitleBar;
