// src/context/TranslationContext.tsx
// Initial version generated with Devstral-small-2:24b

import { createContext, useContext, useState, type ReactNode, useEffect } from "react";

// Define the structure of our translations
export type Translations = Record<string, Record<string, string>>;

// Define available languages
export type Language = "en" | "fi" | string; // string allows for dynamic language codes

interface TranslationContextType {
    language: Language;
    setLanguage: (lang: Language) => void;
    t: (key: string, params?: Record<string, string>) => string;
    availableLanguages: Language[];
}

const TranslationContext = createContext<TranslationContextType | undefined>(undefined);


interface Props {
    children: ReactNode,
    defaultLanguage?: Language
}

export const TranslationProvider = ( {children, defaultLanguage}: Props ) => {
    const [language, setLanguage] = useState<Language>(defaultLanguage ? defaultLanguage : "en");
    const [translations, setTranslations] = useState<Record<Language, Translations>>({});
    const [availableLanguages, setAvailableLanguages] = useState<Language[]>([]);


    const loadTranslations = async () => {
        try {
            const loadedTranslations: Record<Language, Translations> = {};
            const languages: Language[] = [];

            // Get all translation files from the translations folder
            const translationModules = import.meta.glob('../locales/*.json') as Record<
                string,
                () => Promise<{ default: Translations }>
            >;

            for (const path in translationModules) {
                // Extract language code from filename (e.g., "en.json" -> "en")
                const languageCode = path.split('/').pop()?.split('.')[0] as Language;
                if (languageCode) {
                    languages.push(languageCode);
                    const module = await translationModules[path]();
                    loadedTranslations[languageCode] = module.default;
                }
            }

            setTranslations(loadedTranslations);
            setAvailableLanguages(languages);
        } catch (error) {
            console.error("Error loading translations:", error);
        }
    };

    // Load translations dynamically
    useEffect(() => {
        loadTranslations();
    }, []);

    const t = (key: string, params: Record<string, string> = {}): string => {
        if (!translations[language]) {
            console.warn(`No translations available for language: ${language}`);
            return key;
        }

        const keys = key.split(".");
        let current: any = translations[language];

        for (const k of keys) {
            if (current && current[k]) {
                current = current[k];
            } else {
                console.warn(`Translation key "${key}" not found in language "${language}"`);
                return key;
            }
        }

        // Replace placeholders like {name} with actual values
        let result = current as string;
        Object.entries(params).forEach(([placeholder, value]) => {
            result = result.replace(`{${placeholder}}`, value);
        });

        return result;
    };

    return (
        <TranslationContext.Provider value={{ language, setLanguage, t, availableLanguages }}>
            {children}
        </TranslationContext.Provider>
    );
};

export const useTranslation = (): TranslationContextType => {
    const context = useContext(TranslationContext);
    if (context === undefined) {
        throw new Error("useTranslation must be used within a TranslationProvider");
    }
    return context;
};
