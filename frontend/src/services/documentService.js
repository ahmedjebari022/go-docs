import api from "./api"


const documentService = {
    create: async (name) => {
        return await api.post("/documents",name)
    },

    getAll: async () => {
        return await api.get("/documents")
    },

    delete: async (documentId) => {
        return await api.delete(`/documents/${documentId}`)
    }


}


export default documentService