import { defineBoot } from '#q-app/wrappers';
import axios, { type AxiosInstance } from 'axios';

declare module 'vue' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
    $api: AxiosInstance;
  }
}

// Create axios instance with base URL
const api = axios.create({ baseURL: 'http://localhost:8080' });

// Add a request interceptor to include the auth token in all requests
api.interceptors.request.use(
  (config) => {
    // Get the token from localStorage
    const token = localStorage.getItem('authToken');

    // If token exists, add it to the Authorization header
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`, config.data);
    return config;
  },
  (error) => {
    console.error('[API Request Error]', error);
    return Promise.reject(error instanceof Error ? error : new Error(String(error)));
  }
);

// Add interceptors for debugging
api.interceptors.response.use(
  (response) => {
    console.log(`[API Response] ${response.config.method?.toUpperCase()} ${response.config.url}`, response.data);
    return response;
  },
  (error) => {
    if (error.response) {
      console.error(`[API Response Error] ${error.config?.method?.toUpperCase()} ${error.config?.url}`, error.response.data);
    } else {
      console.error('[API Response Error]', error);
    }
    return Promise.reject(error instanceof Error ? error : new Error(String(error)));
  }
);

export default defineBoot(({ app }) => {
  // for use inside Vue files through this.$axios and this.$api
  app.config.globalProperties.$axios = axios;
  app.config.globalProperties.$api = api;
});

export { api };
