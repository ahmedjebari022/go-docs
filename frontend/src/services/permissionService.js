import api from "./api"



const permissionService = {
    add: async (user_id, role, documentId) => {
        return api.post(`/permissions/${documentId}`,{
            user_id,
            role,
        })
    },
    delete: async (userId) => {
        return api.delete("/permission/:documentid",{
            userId
        })
    }

}
export default permissionService