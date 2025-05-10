import { api } from 'src/boot/axios';

export const apiService = {
    get: async <T>(url: string, params = {}) => {
        try {
            const response = await api.get<T>(url, { params });
            return response.data;
        } catch (error) {
            console.error(`API GET error for ${url}:`, error);
            throw error;
        }
    },
    post: async <T>(url: string, data = {}) => {
        try {
            const response = await api.post<T>(url, data);
            return response.data;
        } catch (error) {
            console.error(`API POST error for ${url}:`, error);
            throw error;
        }
    },
    put: async <T>(url: string, data = {}) => {
        try {
            const response = await api.put<T>(url, data);
            return response.data;
        } catch (error) {
            console.error(`API PUT error for ${url}:`, error);
            throw error;
        }
    },
    delete: async <T>(url: string) => {
        try {
            const response = await api.delete<T>(url);
            return response.data;
        } catch (error) {
            console.error(`API DELETE error for ${url}:`, error);
            throw error;
        }
    }
}; 
