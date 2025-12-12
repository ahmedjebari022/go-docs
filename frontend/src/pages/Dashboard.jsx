import { useAuth } from "../context/AuthContext";
import { Link } from "react-router-dom";
import { Search, FileText, Plus, Bell, Settings, LogOut, Grid, List as ListIcon, Filter } from "lucide-react";
import { useState } from "react";
import { motion } from "framer-motion";

export default function Dashboard() {
    const { user, logout } = useAuth();
    const [viewMode, setViewMode] = useState("grid"); // 'grid' or 'list'

    return (
        <div className="min-h-screen bg-gray-50 flex flex-col">
            {/* Navbar */}
            <nav className="fixed top-0 left-0 right-0 h-16 bg-white border-b border-gray-200 z-50 flex items-center justify-between px-4 sm:px-6 lg:px-8">
                {/* Left: Logo */}
                <div className="flex items-center gap-2">
                    <div className="bg-blue-600 p-1.5 rounded-lg">
                        <FileText className="h-5 w-5 text-white" />
                    </div>
                    <span className="text-xl font-bold bg-gradient-to-r from-blue-700 to-indigo-700 bg-clip-text text-transparent">
                        GoDocs
                    </span>
                </div>

                {/* Middle: Search Bar */}
                <div className="flex-1 max-w-2xl mx-8 hidden md:block">
                    <div className="relative group">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <Search className="h-4 w-4 text-gray-400 group-focus-within:text-blue-500 transition-colors" />
                        </div>
                        <input
                            type="text"
                            className="block w-full pl-10 pr-3 py-2 border border-gray-200 rounded-lg leading-5 bg-gray-100 placeholder-gray-500 focus:outline-none focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 sm:text-sm"
                            placeholder="Search documents..."
                        />
                    </div>
                </div>

                {/* Right: User Actions */}
                <div className="flex items-center gap-4">
                    <button className="text-gray-500 hover:text-gray-700 transition-colors relative">
                        <Bell className="h-5 w-5" />
                        <span className="absolute top-0 right-0 block h-2 w-2 rounded-full bg-red-500 ring-2 ring-white" />
                    </button>

                    <div className="h-8 w-px bg-gray-200 mx-1"></div>

                    <div className="flex items-center gap-3">
                        <div className="text-right hidden sm:block">
                            <p className="text-sm font-medium text-gray-900">{user?.email || "User"}</p>
                            <p className="text-xs text-gray-500">Free Plan</p>
                        </div>
                        <button
                            onClick={logout}
                            className="h-9 w-9 rounded-full bg-gradient-to-tr from-blue-100 to-indigo-100 border border-white shadow-sm flex items-center justify-center text-blue-700 font-semibold hover:ring-2 hover:ring-blue-500 transition-all"
                        >
                            {user?.email?.[0].toUpperCase() || "U"}
                        </button>
                    </div>
                </div>
            </nav>

            {/* Main Content */}
            <main className="flex-1 pt-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
                {/* Header Section */}
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900">Documents</h1>
                        <p className="text-sm text-gray-500 mt-1">Manage and collaborate on your files</p>
                    </div>

                    <motion.button
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.98 }}
                        className="inline-flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg shadow-sm transition-colors"
                    >
                        <Plus className="h-4 w-4 mr-2" />
                        New Document
                    </motion.button>
                </div>

                {/* Filters & Controls */}
                <div className="flex items-center justify-between mb-6 bg-white p-2 rounded-xl border border-gray-200 shadow-sm">
                    <div className="flex items-center gap-2">
                        <button className="px-3 py-1.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors">All</button>
                        <button className="px-3 py-1.5 text-sm font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded-md transition-colors">Shared</button>
                        <button className="px-3 py-1.5 text-sm font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 rounded-md transition-colors">Private</button>
                    </div>

                    <div className="flex items-center gap-2 border-l border-gray-200 pl-2">
                        <button
                            onClick={() => setViewMode("grid")}
                            className={`p-1.5 rounded-md transition-colors ${viewMode === "grid" ? "bg-blue-50 text-blue-600" : "text-gray-400 hover:text-gray-600"}`}
                        >
                            <Grid className="h-4 w-4" />
                        </button>
                        <button
                            onClick={() => setViewMode("list")}
                            className={`p-1.5 rounded-md transition-colors ${viewMode === "list" ? "bg-blue-50 text-blue-600" : "text-gray-400 hover:text-gray-600"}`}
                        >
                            <ListIcon className="h-4 w-4" />
                        </button>
                    </div>
                </div>

                {/* Documents Grid - Placeholder for now */}
                <div className={`grid gap-6 ${viewMode === "grid" ? "grid-cols-1 sm:grid-cols-2 lg:grid-cols-3" : "grid-cols-1"}`}>
                    {/* Example Card 1 */}
                    <motion.div
                        whileHover={{ y: -2 }}
                        className="group bg-white rounded-xl border border-gray-200 p-4 shadow-sm hover:shadow-md transition-all cursor-pointer"
                    >
                        <div className="h-32 bg-gray-50 rounded-lg mb-4 flex items-center justify-center border border-dashed border-gray-200 group-hover:border-blue-300 transition-colors">
                            <FileText className="h-8 w-8 text-gray-300 group-hover:text-blue-400 transition-colors" />
                        </div>
                        <div>
                            <h3 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">Project Proposal</h3>
                            <p className="text-xs text-gray-500 mt-1">Edited 2 hours ago</p>
                        </div>
                    </motion.div>

                    {/* Example Card 2 */}
                    <motion.div
                        whileHover={{ y: -2 }}
                        className="group bg-white rounded-xl border border-gray-200 p-4 shadow-sm hover:shadow-md transition-all cursor-pointer"
                    >
                        <div className="h-32 bg-gray-50 rounded-lg mb-4 flex items-center justify-center border border-dashed border-gray-200 group-hover:border-blue-300 transition-colors">
                            <FileText className="h-8 w-8 text-gray-300 group-hover:text-blue-400 transition-colors" />
                        </div>
                        <div>
                            <h3 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">Meeting Notes</h3>
                            <p className="text-xs text-gray-500 mt-1">Edited yesterday</p>
                        </div>
                    </motion.div>

                    {/* Example Card 3 */}
                    <motion.div
                        whileHover={{ y: -2 }}
                        className="group bg-white rounded-xl border border-gray-200 p-4 shadow-sm hover:shadow-md transition-all cursor-pointer"
                    >
                        <div className="h-32 bg-gray-50 rounded-lg mb-4 flex items-center justify-center border border-dashed border-gray-200 group-hover:border-blue-300 transition-colors">
                            <FileText className="h-8 w-8 text-gray-300 group-hover:text-blue-400 transition-colors" />
                        </div>
                        <div>
                            <h3 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">Design Specs</h3>
                            <p className="text-xs text-gray-500 mt-1">Edited 5 days ago</p>
                        </div>
                    </motion.div>
                </div>
            </main>
        </div>
    );
}
