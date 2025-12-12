import React from 'react'
import { Search, FileText, Plus, Bell, Settings, LogOut, Grid, List as ListIcon, Filter } from "lucide-react";
import { motion } from "framer-motion";

function DocumentCard({document}) {
  return (
    
    <motion.div
        whileHover={{ y: -2 }}
        className="group bg-white rounded-xl border border-gray-200 p-4 shadow-sm hover:shadow-md transition-all cursor-pointer"
    >
        <div className="h-32 bg-gray-50 rounded-lg mb-4 flex items-center justify-center border border-dashed border-gray-200 group-hover:border-blue-300 transition-colors">
            <FileText className="h-8 w-8 text-gray-300 group-hover:text-blue-400 transition-colors" />
        </div>
        <div>
            <h3 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">test</h3>
            <p className="text-xs text-gray-500 mt-1">To do edited at</p>
        </div>
    </motion.div>
  )
}

export default DocumentCard