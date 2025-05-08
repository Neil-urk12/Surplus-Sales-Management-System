import axios from 'axios';
import type { Customer } from '../types/customerTypes';

const API_URL = 'http://localhost:8080/api'; // Adjust if your backend URL is different

/**
 * Input type for creating a new customer.
 * The backend will typically set ID, CreatedAt, UpdatedAt, and DateRegistered.
 */
export type NewCustomerInput = Omit<Customer, 'id' | 'createdAt' | 'updatedAt' | 'dateRegistered'>;

/**
 * Input type for updating an existing customer.
 * Allows partial updates of specified fields.
 */
export type UpdateCustomerInput = Partial<NewCustomerInput>;

/**
 * Optional: If the backend wraps list responses in a data object.
 */
// export interface CustomersListResponse {
//   data: Customer[];
//   total?: number;
//   // Add other pagination/metadata if applicable
// }

export const customerApi = {
    /**
     * Fetches all customers from the backend.
     * @returns A promise that resolves to an array of Customer objects.
     */
    getAllCustomers: async (): Promise<Customer[]> => {
        // If backend wraps response: const response = await axios.get<CustomersListResponse>(`${API_URL}/customers`);
        // return response.data.data || [];
        const response = await axios.get<Customer[]>(`${API_URL}/customers`);
        return response.data || [];
    },

    /**
     * Fetches a specific customer by their ID.
     * @param id - The ID of the customer to fetch.
     * @returns A promise that resolves to a Customer object.
     */
    getCustomerById: async (id: string): Promise<Customer> => {
        const response = await axios.get<Customer>(`${API_URL}/customers/${id}`);
        return response.data;
    },

    /**
     * Adds a new customer to the backend.
     * @param customerData - The data for the new customer.
     * @returns A promise that resolves to the created Customer object (including backend-generated fields).
     */
    addCustomer: async (customerData: NewCustomerInput): Promise<Customer> => {
        const response = await axios.post<Customer>(`${API_URL}/customers`, customerData);
        return response.data;
    },

    /**
     * Updates an existing customer on the backend.
     * @param id - The ID of the customer to update.
     * @param customerData - The data to update for the customer.
     * @returns A promise that resolves to the updated Customer object.
     */
    updateCustomer: async (
        id: string,
        customerData: UpdateCustomerInput
    ): Promise<Customer> => {
        const response = await axios.put<Customer>(`${API_URL}/customers/${id}`, customerData);
        return response.data;
    },

    /**
     * Deletes a customer from the backend.
     * @param id - The ID of the customer to delete.
     * @returns A promise that resolves to an object indicating success or failure.
     */
    deleteCustomer: async (id: string): Promise<{ success: boolean; message?: string }> => {
        try {
            await axios.delete(`${API_URL}/customers/${id}`);
            return { success: true, message: 'Customer deleted successfully.' };
        } catch (error) {
            let errorMessage = 'Failed to delete customer.';
            if (axios.isAxiosError(error)) {
                if (error.response?.status === 404) {
                    console.log('Customer not found (404 response), considered deleted.');
                    return { success: true, message: 'Customer already deleted or not found.' };
                }
                if (error.response?.data?.message) {
                    errorMessage = error.response.data.message;
                }
            }
            // For other errors, return failure
            return { success: false, message: errorMessage };
        }
    }
}; 
