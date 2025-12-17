import React from 'react'
import { FileText, Trash2 } from "lucide-react";
import { motion } from "framer-motion";
import { useNavigate } from 'react-router-dom';

function DocumentCard({ document, deleteDocument }) {
  const navigate = useNavigate()

  function handleCardClick(){
    navigate(`/document/${document.document_id}`)
  }
  return (
    <motion.div
        onClick={handleCardClick}
        whileHover={{ y: -2 }}
        className="group bg-white rounded-xl border border-gray-200 p-4 shadow-sm hover:shadow-md transition-all cursor-pointer"
    >
        <div className="h-32 bg-gray-50 rounded-lg mb-4 flex items-center justify-center border border-dashed border-gray-200 group-hover:border-blue-300 transition-colors">
            <FileText className="h-8 w-8 text-gray-300 group-hover:text-blue-400 transition-colors" />
        </div>
        <div className="flex items-center justify-between">
            <div>
                <h3 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{document.document_name}</h3>
                <p className="text-xs text-gray-500 mt-1">To do edited at</p>
            </div>
            <button 
                onClick={(e) => {
                    e.stopPropagation()  // Prevent card click from triggering
                    deleteDocument(document.document_id)
                }}
                className="text-gray-400 hover:text-red-600 transition-colors opacity-0 group-hover:opacity-100 p-2 rounded-lg hover:bg-red-50"
            >
                <Trash2 className="h-4 w-4" />
            </button>
        </div>
    </motion.div>
  )
}

export default DocumentCard