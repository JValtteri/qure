import { useTranslation } from "../../../../context/TranslationContext";
import { posixToDateAndTime } from "../../../../utils/utils";
import MarkdownRenderer from "../../../MarkdownRenderer/MarkdownRenderer";

interface Props {
    size: number;
    time: number;
    code: string;
}

function ReserveSuccess({size, time, code}: Props) {
    const {t} = useTranslation();
    return (
        <>
            <h3 className='centered'>
                {t("notification.reservation successfull")}
            </h3>
            <p className='centered'>
                <MarkdownRenderer content={t("notification.reserved places", {size: String(size), time: posixToDateAndTime(time)})} />
            </p>
            <label className="small-label">Your reservation ID:</label>
            <p className="centered reservation-code">
                {code}
            </p>
        </>
    )
}

export default ReserveSuccess;
