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
  const user = ref<User | null>(null)

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
  }

  function clearAuth() {
    token.value = null
    user.value = null
    localStorage.removeItem('authToken')
  }

  async function login(credentials: { email: string; password: string }): Promise<boolean> {
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
        return true;
      } else {
        console.error('Login succeeded but no token received.');
        clearAuth();
        return false;
      }

    } catch (error) {
      console.error('Login API call failed:', error);
      clearAuth();
      return false;
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
