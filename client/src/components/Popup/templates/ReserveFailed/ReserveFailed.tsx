import { useTranslation } from "../../../../context/TranslationContext";

interface Props {
    error: string;
}

function ReserveFailed({error}: Props) {
    const {t} = useTranslation();
    return (
        <>
            <h3 className='centered'>
                {t("error.reservation failed")}
            </h3>
            <p className='centered'>{t(error)}</p>
        </>
    );
}

export default ReserveFailed;
