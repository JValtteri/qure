import { useTranslation } from "../../context/TranslationContext"
import PrivacyPolicy from "../PrivacyPolicy/PrivacyPolicy"

interface Props {
    hidden?: boolean
    onChange: (e: boolean)=>void
    id?: string
    className?: string
}

function PolicyAccept({hidden, onChange, id, className}: Props) {
    const {t} = useTranslation();
    const handleChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
        onChange(e.target.checked)
    }

    return (
        <label id={id} className={`gdpr-label ${className}`} htmlFor="gdpr-accept" hidden={hidden}>
            {t("login.policy")} <PrivacyPolicy /> <input type='checkbox' onChange={handleChanged}></input>
        </label>
    )
}

export default PolicyAccept;
