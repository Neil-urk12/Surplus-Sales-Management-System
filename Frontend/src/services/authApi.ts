import { api } from 'src/boot/axios'
import type { AxiosError } from 'axios'
import type { User } from '../types/models'

interface LoginResponse {
  token: string;
  user: User;
  message: string;
}

interface ErrorResponse {
  error?: string;
  statusCode?: number;
}

interface LoginCredentials {
  username: string;
  password: string;
}

interface LoginResult {
  success: boolean;
  message?: string;
  token?: string;
  user?: User;
}

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
