import { defineBoot } from '#q-app/wrappers';
import axios, { type AxiosInstance } from 'axios';

declare module 'vue' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
    $api: AxiosInstance;
  }
}

const baseURL = import.meta.env.VITE_API_BASE_URL

if (!baseURL) {
  throw new Error("Missing VITE_API_BASE_URL in .env file")
}

// Create axios instance with base URL
export const api = axios.create({ baseURL });

// Add a request interceptor to include the auth token in all requests
api.interceptors.request.use(
  (config) => {
    // Get the token from localStorage
    const token = localStorage.getItem('authToken');

    // If token exists, add it to the Authorization header
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(new Error(error));
  });

export default defineBoot(({ app }) => {
  // for use inside Vue files through this.$axios and this.$api
  app.config.globalProperties.$axios = axios;
  app.config.globalProperties.$api = api;
});