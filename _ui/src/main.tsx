import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router'
import './index.css'
import Counter from './pages/Counter'
import Login from './pages/Login'
import Sidenav from './components/Sidenav'

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <Sidenav />
    <main className="ml-40">
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/*" element={<Counter />} />
      </Routes>
    </main>

  </BrowserRouter>,
)
