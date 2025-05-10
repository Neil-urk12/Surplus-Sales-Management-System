import { apiService } from './api';
import type { CabsRow, NewCabInput, CabOperationResponse } from 'src/types/cabs';

const API_PATH = '/api/cabs';

export const cabsService = {
    /**
     * Get all cabs with optional filtering
     */
    getCabs: async (filters?: Record<string, string>) => {
        return await apiService.get<CabsRow[]>(API_PATH, filters);
    },

    /**
     * Get a single cab by ID
     */
    getCabById: async (id: number) => {
        return await apiService.get<CabsRow>(`${API_PATH}/${id}`);
    },

    /**
     * Add a new cab
     */
    addCab: async (cab: NewCabInput): Promise<CabOperationResponse> => {
        try {
            // Let the server generate the ID by sending the cab data without an ID
            const response = await apiService.post<CabsRow>(API_PATH, cab);
            
            // Return the ID assigned by the server
            return {
                success: true,
                id: response.id
            };
        } catch (error) {
            return {
                success: false,
                error: error instanceof Error ? error.message : 'Unknown error occurred'
            };
        }
    },

    /**
     * Update an existing cab
     */
    updateCab: async (id: number, cab: NewCabInput): Promise<CabOperationResponse> => {
        try {
            const response = await apiService.put<CabsRow>(`${API_PATH}/${id}`, cab);
            return {
                success: true,
                id: response.id
            };
        } catch (error) {
            return {
                success: false,
                error: error instanceof Error ? error.message : 'Unknown error occurred'
            };
        }
    },

    /**
     * Delete a cab
     */
    deleteCab: async (id: number): Promise<CabOperationResponse> => {
        try {
            await apiService.delete(`${API_PATH}/${id}`);
            return {
                success: true
            };
        } catch (error) {
            return {
                success: false,
                error: error instanceof Error ? error.message : 'Unknown error occurred'
            };
        }
    }
}; 
