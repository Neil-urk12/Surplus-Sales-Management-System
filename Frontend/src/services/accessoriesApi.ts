import axios from 'axios';
import type { AccessoryRow, NewAccessoryInput, AccessoryOperationResponse } from 'src/types/accessories';

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
        // TODO: Replace with actual API call
        console.log('API call: Get all accessories');
        // Simulating API call
        return new Promise((resolve) => {
            setTimeout(() => {
                // Mock data
                resolve([
                    {
                        id: 1,
                        name: 'Premium Steering Wheel',
                        make: 'OEM',
                        quantity: 25,
                        price: 4500,
                        status: 'In Stock',
                        unit_color: 'Black',
                        image: ''
                    },
                    {
                        id: 2,
                        name: 'Sport Seats',
                        make: 'Aftermarket',
                        quantity: 10,
                        price: 12000,
                        status: 'Low Stock',
                        unit_color: 'Black',
                        image: ''
                    },
                    {
                        id: 3,
                        name: 'LED Headlights',
                        make: 'Custom',
                        quantity: 0,
                        price: 8500,
                        status: 'Out of Stock',
                        unit_color: 'White',
                        image: ''
                    }
                ]);
            }, 800);
        });

        // When ready to connect to backend:
        // const response = await axios.get<AccessoryRow[]>(`${API_URL}/accessories`);
        // return response.data;
    },

    /**
     * Get a specific accessory by ID
     * @param id Accessory ID
     * @returns Promise with accessory data
     */
    getAccessoryById: async (id: number): Promise<AccessoryRow> => {
        // TODO: Replace with actual API call
        console.log(`API call: Get accessory by ID: ${id}`);

        // Simulating API call
        return new Promise((resolve, reject) => {
            setTimeout(() => {
                // Mock data
                const mockAccessory = {
                    id,
                    name: 'Premium Steering Wheel',
                    make: 'OEM',
                    quantity: 25,
                    price: 4500,
                    status: 'In Stock',
                    unit_color: 'Black',
                    image: ''
                };
                resolve(mockAccessory);
            }, 500);
        });

        // When ready to connect to backend:
        // const response = await axios.get<AccessoryRow>(`${API_URL}/accessories/${id}`);
        // return response.data;
    },

    /**
     * Add a new accessory
     * @param accessory Accessory data to add
     * @returns Promise with operation response
     */
    addAccessory: async (accessory: NewAccessoryInput): Promise<AccessoryOperationResponse> => {
        // TODO: Replace with actual API call
        console.log('API call: Add accessory', accessory);

        // Simulating API call
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve({
                    success: true,
                    id: Math.floor(Math.random() * 1000) + 10 // Generate a random ID
                });
            }, 1000);
        });

        // When ready to connect to backend:
        // const response = await axios.post<AccessoryOperationResponse>(
        //   `${API_URL}/accessories`,
        //   accessory
        // );
        // return response.data;
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
        // TODO: Replace with actual API call
        console.log(`API call: Update accessory ID ${id}`, accessory);

        // Simulating API call
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve({
                    success: true,
                    id
                });
            }, 1000);
        });

        // When ready to connect to backend:
        // const response = await axios.put<AccessoryOperationResponse>(
        //   `${API_URL}/accessories/${id}`,
        //   accessory
        // );
        // return response.data;
    },

    /**
     * Delete an accessory
     * @param id Accessory ID to delete
     * @returns Promise with operation response
     */
    deleteAccessory: async (id: number): Promise<AccessoryOperationResponse> => {
        // TODO: Replace with actual API call
        console.log(`API call: Delete accessory ID ${id}`);

        // Simulating API call
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve({
                    success: true
                });
            }, 1000);
        });

        // When ready to connect to backend:
        // const response = await axios.delete<AccessoryOperationResponse>(
        //   `${API_URL}/accessories/${id}`
        // );
        // return response.data;
    }
}; 
