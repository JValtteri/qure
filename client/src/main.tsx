import './index.css'

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import App from './App.tsx'
import { TranslationProvider } from './context/TranslationContext.tsx'


createRoot(document.getElementById('root')!).render(

  <StrictMode>
        <TranslationProvider defaultLanguage='fi'>
            <App />
        </TranslationProvider>
    </StrictMode>
)
