import { useTranslation } from "../../../../context/TranslationContext";
import { posixToDateAndTime } from "../../../../utils/utils";
import MarkdownRenderer from "../../../MarkdownRenderer/MarkdownRenderer";

interface Props {
    size: number;
    confirmed: number
    time: number;
    code: string;
}

function ReservePartial({size, confirmed, time, code}: Props) {
    const {t} = useTranslation();
    return (
        <>
            <h3 className='centered'>
                {t("notification.queue")}
            </h3>
            <p className='centered'>
                <MarkdownRenderer content={
                    t("notification.reserved partial", {confirmed: String(confirmed), time: posixToDateAndTime(time)})
                } />
            </p>
            <p className='centered'>
                <MarkdownRenderer content={t("notification.queue", {size: String(size-confirmed)})} />
            </p>
            <label className="small-label">{t("notification.your id")}:</label>
            <p className="centered reservation-code">
                {code}
            </p>
        </>
    )
}

export default ReservePartial;
