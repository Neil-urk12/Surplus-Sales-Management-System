import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import type { User } from '../types/models'

interface LoginResponse {
  token: string;
  user: User;
  message: string;
}

export const useAuthStore = defineStore('auth', () => {
  // --- State --- Reactive refs
  const token = ref<string | null>(localStorage.getItem('authToken') || null)
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') || 'null'))

  // --- Getters --- Computed properties
  const isAuthenticated = computed(() => !!token.value)

  // --- Actions --- Functions to mutate state

  // Helper to set token in state and localStorage
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('authToken', newToken)
  }

  function setUser(newUser: User) {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function clearAuth() {
    token.value = null
    user.value = null
    localStorage.removeItem('authToken')
    localStorage.removeItem('user')
  }

  async function login(credentials: { email: string; password: string }): Promise<{ success: boolean; message?: string }> {
    clearAuth();
    console.log('Attempting login with:', credentials.email);
    try {
      const response = await axios.post<LoginResponse>('/api/users/login', credentials);

      if (response.status === 200 && response.data.token) {
        const apiToken = response.data.token;
        setToken(apiToken);

        if (response.data.user) {
          setUser(response.data.user);
        }
        // TODO: Alternatively, fetch user details in a separate request using the token

        console.log('Login successful');
        return { success: true };
      } else {
        console.error('Login succeeded but no token received.');
        clearAuth();
        return { success: false, message: 'Authentication failed. Please try again.' };
      }

    } catch (error: unknown) {
      console.error('Login API call failed:', error);
      clearAuth();

      // Check for specific error messages from the backend
      if (axios.isAxiosError(error) && error.response) {
        const status = error.response.status;
        const errorMessage = error.response.data?.error || 'Unknown error occurred';

        // Handle specific status codes
        if (status === 403 && errorMessage.includes('inactive')) {
          return { success: false, message: 'Your account is inactive. Please contact an administrator.' };
        } else if (status === 401) {
          return { success: false, message: errorMessage }; // Return the exact error message from backend
        } else {
          return { success: false, message: errorMessage };
        }
      }
      return { success: false, message: 'An error occurred while connecting to the server. Please try again later.' };
    }
  }

  function logout() {
    clearAuth()
    console.log('Logout action called')
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout,
    setToken,
    setUser,
  }
})
