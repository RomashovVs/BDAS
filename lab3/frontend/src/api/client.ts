import axios from 'axios';

const apiClient = axios.create({
    baseURL: '/api', // Будет проксироваться через Nginx
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    }
});

export const api = {
    getHealth() {
        return apiClient.get('/health');
    }
};