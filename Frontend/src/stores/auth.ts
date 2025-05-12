import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '../types/models'
import { loginUser } from '../services/authApi'

export const useAuthStore = defineStore('auth', () => {
  // --- State --- Reactive refs
  const token = ref<string | null>(localStorage.getItem('authToken') || null)
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') || 'null'))

  // --- Getters --- Computed properties
  const isAuthenticated = computed(() => !!token.value)
  const getUserId = computed(() => user.value?.id || '')

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

  async function login(credentials: { username: string; password: string }): Promise<{ success: boolean; message?: string }> {
    clearAuth();
    
    const result = await loginUser(credentials);
    
    if (result.success && result.token && result.user) {
      setToken(result.token);
      setUser(result.user);
      return { success: true };
    } else {
      clearAuth();
      return { 
        success: false, 
        message: result.message || 'Authentication failed. Please try again.' 
      };
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
    getUserId,
    login,
    logout,
    setToken,
    setUser,
  }
})
