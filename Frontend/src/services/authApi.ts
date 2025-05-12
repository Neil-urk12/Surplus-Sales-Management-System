/**
 * Authentication API service
 * 
 * This module provides functions for interacting with the authentication-related
 * endpoints of the backend API. It handles API requests, error processing,
 * and data formatting for authentication operations.
 */
import { api } from 'src/boot/axios'
import type { AxiosError } from 'axios'
import type { User } from '../types/models'

/**
 * Response structure from the login API endpoint
 */
interface LoginResponse {
  token: string;
  user: User;
  message: string;
}

/**
 * Error response structure from API endpoints
 */
interface ErrorResponse {
  error?: string;
  statusCode?: number;
}

/**
 * Structure for login credentials sent to the API
 */
interface LoginCredentials {
  username: string;
  password: string;
}

/**
 * Standardized return structure for the login function
 */
interface LoginResult {
  success: boolean;
  message?: string;
  token?: string;
  user?: User;
}

/**
 * Attempts to authenticate a user with the backend API
 * 
 * @param {LoginCredentials} credentials - User login credentials (username and password)
 * @returns {Promise<LoginResult>} Result object containing success status, auth token, user data and/or error message
 * 
 * @example
 * const result = await loginUser({ username: 'john', password: 'secret' });
 * if (result.success) {
 *   // Authentication successful, token and user available in result
 * } else {
 *   // Authentication failed, message contains reason
 * }
 */
export async function loginUser(credentials: LoginCredentials): Promise<LoginResult> {
  console.log('Attempting login with:', credentials.username);
  try {
    const response = await api.post<LoginResponse>('/api/users/login', credentials);

    if (response.status === 200 && response.data.token) {
      console.log('Login successful');
      return { 
        success: true, 
        token: response.data.token,
        user: response.data.user
      };
    } else {
      console.error('Login succeeded but no token received.');
      return { success: false, message: 'Authentication failed. Please try again.' };
    }

  } catch (error: unknown) {
    console.error('Login API call failed:', error);

    // Check for specific error messages from the backend
    if (error && typeof error === 'object' && 'isAxiosError' in error && (error as AxiosError).response) {
      const axiosError = error as AxiosError<ErrorResponse>;
      const status = axiosError.response?.status;
      const errorMessage = axiosError.response?.data?.error || 'Unknown error occurred';

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
