import Frame from '../Frame/Frame';
import './TitleBar.css';

interface Props {
    title?: string
    icon?: string
}

function TitleBar({title, icon}: Props) {
    return (
        <Frame className='title'>
            <img src={ icon ? icon : './logo.png' } />
            <span>
                {title ? title : "< Title >"}
            </span>
            <button>Login</button>
        </Frame>
    )
}

export default TitleBar;
