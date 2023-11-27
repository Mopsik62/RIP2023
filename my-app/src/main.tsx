import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter} from 'react-router-dom'
import { Routes, Route } from 'react-router-dom';

import 'bootstrap/dist/css/bootstrap.min.css';

import SubstancesPage from './SubstancesPage.tsx'
import SubstancePage from './SubstancePage.tsx'
import Navigation from './components/Navigation'
import Breadcrumbs from './components/Breadcrumbs';

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <BrowserRouter>
            <Navigation />
            <Breadcrumbs />
            <Routes>
                <Route path="/One-pot-front" Component={SubstancesPage} />
                <Route path="/One-pot-front/substance" Component={SubstancePage} />
            </Routes>
        </BrowserRouter>
    </React.StrictMode>,
)
