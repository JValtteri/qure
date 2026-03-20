// src/components/LanguageSelector.tsx
// Initial version generated with Devstral-small-2:24b

import "./LanguageSelector.css";
import { useTranslation } from "../../context/TranslationContext";
import { getLocalStorage } from "../../utils/local_storage";

interface Props {
}

function LanguageSelector({}: Props) {
    const { language, setLanguage } = useTranslation();
    if (getLocalStorage("locale") != "") {
        setLanguage(getLocalStorage("locale"));
    }

    return (
        <div>
            <select
                id="language-selector"
                value={language}
                onChange={e => {
                    setLanguage(e.target.value);
                    localStorage.setItem("locale", e.target.value);
                }}
            >
                <option value="fi">🇫🇮</option>
                <option value="sv">🇸🇪</option>
                <option value="en">🇬🇧</option>
            </select>
        </div>
    );
};

export default LanguageSelector;
