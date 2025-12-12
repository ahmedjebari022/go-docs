import axios from 'axios'


const api = axios.create({
    baseURL: 'http:/http://localhost:8081/api',
    timeout:1000,
    withCredentials: true,
    headers:{
        'Content-Type':'application/json'
    }
})



export default api