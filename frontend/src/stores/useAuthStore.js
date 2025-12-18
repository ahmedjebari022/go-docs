import { create } from "zustand"
import { persist, createJSONStorage } from 'zustand/middleware'
import { authService } from "../services/authService"

const useAuthStore = create(
    persist(
        (set, get) => ({
           userEmail: '',
           isLoggedIn: false,
           isLoading: true, 

           setUserEmail: (email) => set({ userEmail: email }),
           setIsLoggedIn: (logged) => set({ isLoggedIn: logged }),
           setIsLoading: (loading) => set({ isLoading: loading }),

            initializeAuth: async () => {
                set({ isLoading: true })
                try {
                    const res = await authService.me()
                    if (res.status === 200) {
                        set({ isLoggedIn: true, userEmail: res.data.email })
                        console.log(res.data)
                    } else {
                        set({ isLoggedIn: false, userEmail: "" })
                    }
                    set({ isLoading: false })
                } catch (error) {
                   console.log(error) 
                   set({ isLoading: false, isLoggedIn: false, userEmail: "" })
                }           
            },

            login: async (email, password) => {
                set({ isLoading: true })
                try {
                    const res = await authService.login(email, password) 
                    if (res.status === 200) {
                        set({ isLoggedIn: true, userEmail: email })
                    } else {
                        set({ isLoggedIn: false })
                    }
                    set({ isLoading: false })
                } catch (error) {
                   set({ isLoading: false })
                   console.log(error) 
                }
            },

            logout: async () => {
                set({ isLoading: true })
                try {
                    await authService.logout() 
                    set({ isLoggedIn: false, userEmail: "", isLoading: false })
                } catch (error) {
                   set({ isLoading: false })
                   console.log(error) 
                }
            }
        }),
        {
            name: "auth-storage",
            storage: createJSONStorage(() => localStorage),
            partialize: (state) => ({
                userEmail: state.userEmail,
                isLoggedIn: state.isLoggedIn,
            })
        }
    )
)

export default useAuthStore
