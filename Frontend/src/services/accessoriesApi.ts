import axios from 'axios';
import type { AccessoryRow, NewAccessoryInput, AccessoryOperationResponse, AccessoriesListResponse } from 'src/types/accessories';

// Base API URL - replace with actual backend URL when connecting
const API_URL = 'http://localhost:8080/api';

/**
 * Service for handling accessory-related API calls
 */
export const accessoriesApi = {
    /**
     * Get all accessories
     * @returns Promise with array of accessories
     */
    getAllAccessories: async (): Promise<AccessoryRow[]> => {
        const response = await axios.get<AccessoriesListResponse>(`${API_URL}/accessories`);
        return response.data.data || [];
    },

    /**
     * Get a specific accessory by ID
     * @param id Accessory ID
     * @returns Promise with accessory data
     */
    getAccessoryById: async (id: number): Promise<AccessoryRow> => {
        const response = await axios.get<AccessoryRow>(`${API_URL}/accessories/${id}`);
        return response.data;
    },

    /**
     * Add a new accessory
     * @param accessory Accessory data to add
     * @returns Promise with operation response
     */
    addAccessory: async (accessory: NewAccessoryInput): Promise<AccessoryOperationResponse> => {
        console.log('Accessory being sent to API:', JSON.stringify(accessory)); // ADD THIS
        const response = await axios.post<AccessoryOperationResponse>(
            `${API_URL}/accessories`,
            accessory
        );
        console.log('API response data for addAccessory:', JSON.stringify(response.data)); // ADD THIS
        return response.data;
    },

    /**
     * Update an existing accessory
     * @param id Accessory ID
     * @param accessory Updated accessory data
     * @returns Promise with operation response
     */
    updateAccessory: async (
        id: number,
        accessory: Partial<NewAccessoryInput>
    ): Promise<AccessoryOperationResponse> => {
        const response = await axios.put<AccessoryOperationResponse>(
            `${API_URL}/accessories/${id}`,
            accessory
        );
        return response.data;
    },

    /**
     * Delete an accessory
     * @param id Accessory ID to delete
     * @returns Promise with operation response
     */
    deleteAccessory: async (id: number): Promise<AccessoryOperationResponse> => {
        try {
            const response = await axios.delete<AccessoryOperationResponse>(
                `${API_URL}/accessories/${id}`
            );
            if (response.status === 204) {
                return { success: true };
            }
            return response.data;
        } catch (error) {
            // If we get a 404 after trying to delete, the item was already deleted
            // Consider this a success case
            if (axios.isAxiosError(error) && error.response?.status === 404) {
                console.log('Accessory already deleted (404 response)');
                return { success: true };
            }
            throw error; // Re-throw other errors to be handled by the caller
        }
    }
}; 
