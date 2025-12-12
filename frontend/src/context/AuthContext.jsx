import { createContext, useContext, useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    // Check if user is logged in (via cookie) on mount
    useEffect(() => {
        checkAuth();
    }, []);

    const checkAuth = async () => {
        try {
            // We'll use the check cookie endpoint or similar to validate session
            // For now, if we have a way to validate, do it here.
            // Based on backend, GET /api/cookie checks the access cookie.
            const res = await fetch('/api/cookie');
            if (res.ok) {
                // We might want to store more user info if available
                // For now just mark as logged in with a placeholder or decoded token if needed
                // The backend returns { value: "jwt_token" }
                // We probably want a /api/me endpoint later to get user details
                setUser({ authenticated: true });
            } else {
                setUser(null);
            }
        } catch (error) {
            console.error("Auth check failed", error);
            setUser(null);
        } finally {
            setLoading(false);
        }
    };

    const login = async (email, password) => {
        const res = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        if (!res.ok) {
            const error = await res.text();
            throw new Error(error || 'Login failed');
        }

        const data = await res.json();
        setUser({ email: data.email });
        return data;
    };

    const register = async (email, password) => {
        const res = await fetch('/api/users', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        if (!res.ok) {
            const error = await res.text();
            throw new Error(error || 'Registration failed');
        }

        return await res.json();
    };

    const logout = async () => {
        // Ideally call a logout endpoint to clear cookies
        // For now just clear local state
        setUser(null);
        // Optional: await fetch('/api/auth/logout', { method: 'POST' }); 
        window.location.href = '/login';
    };

    return (
        <AuthContext.Provider value={{ user, login, register, logout, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
