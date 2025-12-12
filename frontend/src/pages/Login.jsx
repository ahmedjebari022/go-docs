import React, { useState, useEffect } from 'react'
import { authService } from '../services/authService';
import { Link, useNavigate } from 'react-router-dom';
import { FileText, Mail, Lock, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";
import useAuthStore from '../stores/useAuthStore';

function Login() {
    const navigate = useNavigate();
    const { login, isLoggedIn } = useAuthStore();
    const [formData, setFormData] = useState({
        email: "",
        password: "",
    });
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (isLoggedIn) {
            navigate('/dashboard');
        }
    }, [isLoggedIn, navigate]);

    async function handleSubmit(e) {
        e.preventDefault();
        setError("");
        setLoading(true);
        
        try {
            await login(formData.email, formData.password);
        } catch (err) {
            setError("Invalid email or password");
        } finally {
            setLoading(false);
        }
    }

    function handleFormChange(e) {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex items-center justify-center p-4">
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
                className="w-full max-w-md"
            >
                {/* Logo */}
                <div className="flex justify-center mb-8">
                    <div className="flex items-center gap-2">
                        <div className="bg-blue-600 p-2 rounded-lg">
                            <FileText className="h-6 w-6 text-white" />
                        </div>
                        <span className="text-2xl font-bold bg-gradient-to-r from-blue-700 to-indigo-700 bg-clip-text text-transparent">
                            GoDocs
                        </span>
                    </div>
                </div>

                {/* Card */}
                <div className="bg-white rounded-2xl shadow-xl border border-gray-100 p-8">
                    <div className="mb-8">
                        <h1 className="text-2xl font-bold text-gray-900 mb-2">Welcome back</h1>
                        <p className="text-gray-500">Sign in to continue to your documents</p>
                    </div>

                    {error && (
                        <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-600 rounded-lg text-sm">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Email
                            </label>
                            <div className="relative">
                                <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                                <input
                                    type="email"
                                    name="email"
                                    value={formData.email}
                                    onChange={handleFormChange}
                                    className="w-full pl-10 pr-4 py-2.5 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                                    placeholder="you@example.com"
                                    required
                                />
                            </div>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Password
                            </label>
                            <div className="relative">
                                <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                                <input
                                    type="password"
                                    name="password"
                                    value={formData.password}
                                    onChange={handleFormChange}
                                    className="w-full pl-10 pr-4 py-2.5 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                                    placeholder="••••••••"
                                    required
                                />
                            </div>
                        </div>

                        <motion.button
                            whileHover={{ scale: 1.01 }}
                            whileTap={{ scale: 0.99 }}
                            type="submit"
                            disabled={loading}
                            className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2.5 rounded-lg transition-colors flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {loading ? "Signing in..." : "Sign in"}
                            {!loading && <ArrowRight className="h-4 w-4" />}
                        </motion.button>
                    </form>

                    <div className="mt-6 text-center">
                        <p className="text-sm text-gray-600">
                            Don't have an account?{" "}
                            <Link to="/register" className="text-blue-600 hover:text-blue-700 font-medium">
                                Sign up
                            </Link>
                        </p>
                    </div>
                </div>

                <p className="text-center text-xs text-gray-500 mt-8">
                    By signing in, you agree to our Terms of Service and Privacy Policy
                </p>
            </motion.div>
        </div>
    );
}

export default Login;