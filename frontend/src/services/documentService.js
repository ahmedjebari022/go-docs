import api from "./api"


const documentService = {
    create: async (name) => {
        return await api.post("/documents",{
            name
        })
    },

    getAll: async () => {
        return await api.get("/documents")
    },

    delete: async (documentId) => {
        return await api.delete(`/documents/${documentId}`)
    },

    get: async (documentId) => {
        return await api.get(`/documents/${documentId}`)
    },

    update: async (documentId, documentContent) => {
        return await api.put(`/documents/${documentId}`,
            documentContent
        )
    }

    


}


export default documentService