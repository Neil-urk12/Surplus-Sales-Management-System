import { api } from 'src/boot/axios'
import type { AxiosError } from 'axios'
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

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

/**
 * Fetches the list of users from the API.
 * @returns {Promise<ApiResponse<User[]>>} A promise that resolves with the API response containing users array.
 */
export async function fetchUsers(): Promise<ApiResponse<User[]>> {
  try {
    const response = await api.get<{ users: User[] }>('/api/users');
    return {
      success: true,
      data: response.data.users
    };
  } catch (err) {
    const axiosError = err as AxiosError<{ error?: string }>;
    console.error('Error fetching users:', axiosError);
    return {
      success: false,
      error: axiosError.response?.data?.error || 'Failed to fetch users'
    };
  }
}

/**
 * Creates a new user via an API request.
 * @param {UserCreateData} userData - The data for the new user.
 * @returns {Promise<ApiResponse<User>>} A promise that resolves with the API response containing the created user.
 */
export async function createUser(userData: UserCreateData): Promise<ApiResponse<User>> {
  try {
    const response = await api.post<{ user: User }>('/api/users', userData);
    return {
      success: true,
      data: response.data.user
    };
  } catch (err) {
    const axiosError = err as AxiosError<{ error?: string }>;
    console.error('Error creating user:', axiosError);
    return {
      success: false,
      error: axiosError.response?.data?.error || 'Failed to create user'
    };
  }
}

/**
 * Updates an existing user via an API request.
 * @param {string} userId - The ID of the user to update.
 * @param {UserUpdateData} userData - The data to update the user with.
 * @returns {Promise<ApiResponse<void>>} A promise that resolves with the API response.
 */
export async function updateUser(userId: string, userData: UserUpdateData): Promise<ApiResponse<void>> {
  try {
    await api.put(`/api/users/${userId}`, userData);
    return { success: true };
  } catch (err) {
    const axiosError = err as AxiosError<{ error?: string }>;
    console.error('Error updating user:', axiosError);
    return {
      success: false,
      error: axiosError.response?.data?.error || 'Failed to update user'
    };
  }
}

/**
 * Deletes a user via an API request.
 * @param {string} userId - The ID of the user to delete.
 * @returns {Promise<ApiResponse<void>>} A promise that resolves with the API response.
 */
export async function deleteUser(userId: string): Promise<ApiResponse<void>> {
  try {
    await api.delete(`/api/users/${userId}`);
    return { success: true };
  } catch (err) {
    const axiosError = err as AxiosError<{ error?: string }>;
    console.error('Error deleting user:', axiosError);
    return {
      success: false,
      error: axiosError.response?.data?.error || 'Failed to delete user'
    };
  }
}
