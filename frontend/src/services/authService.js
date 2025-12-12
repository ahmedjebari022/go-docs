import api from "./api"

export const authService = {

    register : async (email, password) => {
        return await api.post("/users", {
            email,
            password
        })
    },

    login : async (email, password) => {
        return await api.post("/auth/login", {
            email,
            password 
        })
    },

    logout : async() => {
        return await api.post("/auth/logout")
    },

    me : async() => {
        return await api.post("users/me")
    }

}

