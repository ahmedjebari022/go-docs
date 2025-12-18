import React, { useEffect, useState } from 'react'
import { useQueryClient, useQuery, useMutation, } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import documentService from '../services/documentService'
import { ChevronDown, Plus, User, Bold, Italic, Underline, Type, Palette } from 'lucide-react'
import useAuthStore from '../stores/useAuthStore'
import ShareModal from '../components/ShareModal'
import permissionService from '../services/permissionService'
import {TextEditor} from '../components/TextEditor'

function Document() {
  const { documentId } = useParams()
  const queryClient = useQueryClient() 
  const [documentContent, setDocumentContent] = useState({})
  const {userEmail, logout} = useAuthStore()
  const [showModal, setShowModal] = useState(false)
  const { data: document, isLoading, isError }= useQuery({
    queryKey:['document', documentId],
    queryFn: () => documentService.get(documentId).then(res => res.data) 
  })
  async function addCollaborator(userId,permission){
     try {
      const res = permissionService.add(userId,permission, documentId)
      console.log(res)
      if (res.status === 201){
        setShowModal(false)
      }
     } catch (error) {
      console.log(error) 
     }
  }
  async function handleSaveDocument(){
    try {
      console.log(documentContent)
      const res = await documentService.update(documentId,JSON.stringify(documentContent))
      console.log(res)
    } catch (error) {
      console.log(error) 
    }
  }

  function handleChangeDocument(value){
    setDocumentContent(value)
  }

  if (isLoading){
    return <div>Loading document...</div>
  }
  if (isError){
    return <div>Loading loading document</div>
  }
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navbar */}
      <nav className="bg-white border-b border-gray-200 fixed top-0 left-0 right-0 z-50">
        <div className="flex items-center justify-between px-4 py-3">
          {/* Left side - Dropdown and Document name */}
          <div className="flex items-center gap-4">
            {/* Dropdown */}
            <button onClick={handleSaveDocument} className="flex items-center gap-2 px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors">
              <span className="text-sm font-medium text-gray-700">Save</span>
            </button>

            {/* Document Name */}
            <h1 className="text-lg font-semibold text-gray-900">
              {document.meta_data.document_name}
            </h1>
          </div>

          {/* Middle - Text formatting tools */}
          <div className="flex items-center gap-2">
            {/* Font selector */}
            <button className="flex items-center gap-2 px-3 py-2 border border-gray-200 hover:bg-gray-50 rounded-lg transition-colors">
              <Type className="h-4 w-4 text-gray-600" />
              <span className="text-sm text-gray-700">Arial</span>
              <ChevronDown className="h-4 w-4 text-gray-500" />
            </button>

            {/* Font size */}
            <button className="flex items-center gap-2 px-3 py-2 border border-gray-200 hover:bg-gray-50 rounded-lg transition-colors">
              <span className="text-sm text-gray-700">14</span>
              <ChevronDown className="h-4 w-4 text-gray-500" />
            </button>

            {/* Divider */}
            <div className="h-6 w-px bg-gray-200"></div>

            {/* Text styling buttons */}
            <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
              <Bold className="h-4 w-4 text-gray-600" />
            </button>
            <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
              <Italic className="h-4 w-4 text-gray-600" />
            </button>
            <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
              <Underline className="h-4 w-4 text-gray-600" />
            </button>

            {/* Divider */}
            <div className="h-6 w-px bg-gray-200"></div>

            {/* Color picker */}
            <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
              <Palette className="h-4 w-4 text-gray-600" />
            </button>
          </div>

          {/* Right side - User and Share */}
          <div className="flex items-center gap-3">
            {/* Share button */}
            <button className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors">
              <Plus className="h-4 w-4" />
              <span onClick={()=>setShowModal(true)}className="text-sm font-medium">Share</span>
            </button>

            {/* User avatar */}
            <div className="w-9 h-9 bg-blue-100 rounded-full flex items-center justify-center">
            {userEmail[0].toUpperCase() || "U"}
            </div>
          </div>
        </div>
      </nav>

      {/* Main content area */}
      <main className="pt-20 px-4">
        <div className="max-w-4xl mx-auto bg-white shadow-sm min-h-screen p-12">
            <TextEditor
              userEmail={userEmail}
              documentId={documentId}
              handleUpdate={handleChangeDocument}              

              />
        </div>
      </main>
      {showModal && (
        <ShareModal 
          onClose={()=>setShowModal(false)} 
          onSubmit={addCollaborator}
        />
      )}
    </div>
  )
}
export default Document