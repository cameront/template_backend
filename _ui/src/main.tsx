import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router'
import './index.css'
import App from './App'

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <Routes>
      <Route index path='*' element={<App />} />
    </Routes>
  </BrowserRouter>,
)
