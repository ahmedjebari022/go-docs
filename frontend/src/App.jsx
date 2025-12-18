import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard'
import useAuthStore from './stores/useAuthStore';
import { useEffect } from 'react';
import Document from './pages/Document';

function App() {
    const {initializeAuth, isLoggedIn, isLoading} = useAuthStore();
    useEffect(() => {
      initializeAuth()
    },[])



    return (
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />}/>
          <Route path="/register" element={<Register/>}/>
          <Route path="/dashboard" element={isLoggedIn ? <Dashboard/> : <Login/>} />
          <Route path="/document/:documentId" element={isLoggedIn? <Document /> : <Login/>} />
        </Routes>
      </BrowserRouter>
    )
}

export default App
