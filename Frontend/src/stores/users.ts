import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'
import type { User } from '../types/models'
import { 
  fetchUsers as apiFetchUsers, 
  createUser as apiCreateUser, 
  updateUser as apiUpdateUser, 
  deleteUser as apiDeleteUser,
  type UserCreateData,
  type UserUpdateData
} from '../services/usersApi'

/**
 * Pinia store for managing user data and related actions.
 * Handles fetching, creating, updating, and deleting users.
 */
export const useUsersStore = defineStore('users', () => {
  /**
   * Ref holding the array of user objects.
   * @type {import('vue').Ref<User[]>}
   */
  const users = ref<User[]>([])
  /**
   * Ref indicating if an API request is in progress.
   * @type {import('vue').Ref<boolean>}
   */
  const loading = ref(false)
  /**
   * Ref holding any error message from API requests.
   * @type {import('vue').Ref<string | null>}
   */
  const error = ref<string | null>(null)

  /**
   * Fetches the list of users from the API.
   * Requires authentication. Populates the `users` state.
   * @async
   * @returns {Promise<User[]>} A promise that resolves with the array of users, or an empty array if an error occurs or authentication fails.
   */
  async function fetchUsers(): Promise<User[]> {
    const authStore = useAuthStore()
    if (!authStore.token) {
      error.value = 'Authentication required'
      return []
    }

    loading.value = true
    error.value = null

    const response = await apiFetchUsers()
    
    loading.value = false
    
    if (response.success && response.data) {
      users.value = response.data
      return response.data
    } else {
      error.value = response.error || 'Failed to fetch users'
      return []
    }
  }

  /**
   * Creates a new user via an API request.
   * Requires authentication. Refreshes the user list on success.
   * @async
   * @param {UserCreateData} userData - The data for the new user.
   * @returns {Promise<User | null>} A promise that resolves with the created user object, or null if an error occurs or authentication fails.
   */
  async function createUser(userData: UserCreateData): Promise<User | null> {
    const authStore = useAuthStore()
    if (!authStore.token) {
      error.value = 'Authentication required'
      return null
    }

    loading.value = true
    error.value = null

    const response = await apiCreateUser(userData)
    
    loading.value = false
    
    if (response.success && response.data) {
      await fetchUsers() // Refresh the users list
      return response.data
    } else {
      error.value = response.error || 'Failed to create user'
      return null
    }
  }

  /**
   * Updates an existing user via an API request.
   * Requires authentication. Refreshes the user list on success.
   * @async
   * @param {string} userId - The ID of the user to update.
   * @param {UserUpdateData} userData - The data to update the user with.
   * @returns {Promise<boolean>} A promise that resolves with true if the update was successful, false otherwise.
   */
  async function updateUser(userId: string, userData: UserUpdateData): Promise<boolean> {
    const authStore = useAuthStore()
    if (!authStore.token) {
      error.value = 'Authentication required'
      return false
    }

    loading.value = true
    error.value = null
    
    const response = await apiUpdateUser(userId, userData)
    
    loading.value = false
    
    if (response.success) {
      await fetchUsers()
      return true
    } else {
      error.value = response.error || 'Failed to update user'
      return false
    }
  }

  /**
   * Deletes a user via an API request.
   * Requires authentication. Refreshes the user list on success.
   * @async
   * @param {string} userId - The ID of the user to delete.
   * @returns {Promise<boolean>} A promise that resolves with true if the deletion was successful, false otherwise.
   */
  async function deleteUser(userId: string): Promise<boolean> {
    const authStore = useAuthStore()
    if (!authStore.token) {
      error.value = 'Authentication required'
      return false
    }

    loading.value = true
    error.value = null
    
    const response = await apiDeleteUser(userId)
    
    loading.value = false
    
    if (response.success) {
      await fetchUsers() // Refresh the users list
      return true
    } else {
      error.value = response.error || 'Failed to delete user'
      return false
    }
  }

  return {
    users,
    loading,
    error,
    fetchUsers,
    createUser,
    updateUser,
    deleteUser
  }
})
