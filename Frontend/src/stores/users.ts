import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from 'src/boot/axios'
import type { AxiosError } from 'axios'
import { useAuthStore } from './auth'
import type { User } from '../types/models'

/**
 * Interface for the data required to create a new user.
 */
export interface UserCreateData {
  /** The full name of the user. */
  fullName: string;
  /** The username of the user. */
  username: string;
  /** The email address of the user. */
  email: string;
  /** The password for the new user account. */
  password: string;
  /** The role assigned to the user ('admin' or 'staff'). */
  role: 'admin' | 'staff';
}

/**
 * Interface for the data required to update an existing user.
 */
export interface UserUpdateData {
  /** The updated full name of the user. */
  fullName?: string;
  /** The updated username of the user. */
  username?: string;
  /** The updated email address of the user. */
  email?: string;
  /** The updated role for the user ('admin' or 'staff'). */
  role?: 'admin' | 'staff';
  /** The updated active status of the user. */
  isActive?: boolean;
}

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
      return [] // Return empty array as per JSDoc
    }

    loading.value = true
    error.value = null

    try {
      const response = await api.get<{ users: User[] }>('/api/users')
      
      users.value = response.data.users
      return response.data.users
    } catch (err) {
      const axiosError = err as AxiosError<{ error?: string }>
      console.error('Error fetching users:', axiosError)
      error.value = axiosError.response?.data?.error || 'Failed to fetch users'
      return []
    } finally {
      loading.value = false
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

    try {
      const response = await api.post<{ user: User }>('/api/users', userData)
      
      await fetchUsers() // Refresh the users list
      return response.data.user
    } catch (err) {
      const axiosError = err as AxiosError<{ error?: string }>
      console.error('Error creating user:', axiosError)
      error.value = axiosError.response?.data?.error || 'Failed to create user'
      return null
    } finally {
      loading.value = false
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

    try {
      await api.put(`/api/users/${userId}`, userData)
      
      await fetchUsers()
      return true
    } catch (err) {
      const axiosError = err as AxiosError<{ error?: string }>
      console.error('Error updating user:', axiosError)
      error.value = axiosError.response?.data?.error || 'Failed to update user'
      return false
    } finally {
      loading.value = false
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

    try {
      await api.delete(`/api/users/${userId}`)
      
      await fetchUsers() // Refresh the users list
      return true
    } catch (err) {
      const axiosError = err as AxiosError<{ error?: string }>
      console.error('Error deleting user:', axiosError)
      error.value = axiosError.response?.data?.error || 'Failed to delete user'
      return false
    } finally {
      loading.value = false
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
