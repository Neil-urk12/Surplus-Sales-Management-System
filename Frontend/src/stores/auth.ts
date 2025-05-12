/**
 * Authentication Store
 * 
 * This store manages the authentication state of the application, including:
 * - User authentication status
 * - Auth token storage and management
 * - User profile information
 * - Login and logout operations
 * 
 * It uses localStorage for persistent storage of authentication data
 * between page refreshes or application restarts.
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '../types/models'
import { loginUser } from '../services/authApi'

export const useAuthStore = defineStore('auth', () => {
  // --- State --- Reactive refs
  /**
   * JWT authentication token stored in state and localStorage
   */
  const token = ref<string | null>(localStorage.getItem('authToken') || null)
  
  /**
   * Currently authenticated user data
   */
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') || 'null'))

  // --- Getters --- Computed properties
  /**
   * Whether the user is currently authenticated
   */
  const isAuthenticated = computed(() => !!token.value)
  
  /**
   * The ID of the currently authenticated user, or empty string if not authenticated
   */
  const getUserId = computed(() => user.value?.id || '')

  // --- Actions --- Functions to mutate state

  /**
   * Sets the authentication token in both state and localStorage
   * @param {string} newToken - The JWT token to store
   */
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('authToken', newToken)
  }

  /**
   * Sets the user data in both state and localStorage
   * @param {User} newUser - The user object to store
   */
  function setUser(newUser: User) {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  /**
   * Clears all authentication data from state and localStorage
   */
  function clearAuth() {
    token.value = null
    user.value = null
    localStorage.removeItem('authToken')
    localStorage.removeItem('user')
  }

  /**
   * Authenticates a user with the backend API
   * 
   * @param {Object} credentials - User credentials
   * @param {string} credentials.username - Username for login
   * @param {string} credentials.password - Password for login
   * @returns {Promise<{success: boolean, message?: string}>} Result object with success status and optional error message
   * 
   * @example
   * const { success, message } = await login({ username: 'admin', password: 'password' });
   * if (success) {
   *   // Login successful
   * } else {
   *   // Show error message
   * }
   */
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

  /**
   * Logs out the current user by clearing authentication data
   */
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
