import React, { useState } from 'react'
import { X, UserPlus } from 'lucide-react'

function ShareModal({onSubmit, onClose}) {
    const [formData, setFormData] = useState({
        userId: "",
        permission: "viewer"
    })
    
    function handleInputChange(e){
        const {name, value} = e.target
        setFormData({
            ...formData,
            [name]: value
        })
    }

    function handleSubmit(e){
        e.preventDefault()
        onSubmit(formData.userId, formData.permission)
    }

    return (
        // Backdrop overlay
        <div 
            className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
            onClick={onClose}
        >
            {/* Modal card */}
            <div 
                className="bg-white rounded-xl shadow-2xl max-w-md w-full p-6"
                onClick={(e) => e.stopPropagation()}
            >
                {/* Header */}
                <div className="flex items-center justify-between mb-6">
                    <div className="flex items-center gap-2">
                        <UserPlus className="h-5 w-5 text-blue-600" />
                        <h2 className="text-xl font-bold text-gray-900">Share Document</h2>
                    </div>
                    <button 
                        onClick={onClose}
                        className="text-gray-400 hover:text-gray-600 transition-colors"
                    >
                        <X className="h-5 w-5" />
                    </button>
                </div>

                {/* Form */}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            User Email
                        </label>
                        <input 
                            name="userId" 
                            value={formData.userId} 
                            type="text" 
                            onChange={handleInputChange}
                            className="w-full px-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                            placeholder="user@example.com"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Permission Level
                        </label>
                        <select 
                            name="permission" 
                            value={formData.permission}
                            onChange={handleInputChange}
                            className="w-full px-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all bg-white"
                        >
                            <option value="viewer">Viewer - Can only view</option>
                            <option value="editor">Editor - Can edit</option>
                        </select>
                    </div>

                    {/* Buttons */}
                    <div className="flex gap-3 pt-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="flex-1 px-4 py-2 border border-gray-200 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
                        >
                            Share
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default ShareModal