import { api } from '../boot/axios'; // Import the configured Axios instance
import axios from 'axios'; // Import axios for isAxiosError type guard
import type { Customer } from '../types/customerTypes';

// const API_URL = 'http://localhost:8080/api'; // No longer needed if api instance has baseURL

// --- Backend-specific type definitions ---

/**
 * Matches the structure of CustomerResponse from the Go backend.
 */
interface BackendCustomerResponse {
    id: string;
    fullName: string; // Changed from name
    email: string;
    phone: string;
    address?: string; // omitempty in backend
    dateRegistered: string;
    createdAt: string;
    updatedAt: string;
}

/**
 * Matches the structure for a list of customers from the Go backend.
 */
interface BackendCustomerListResponse {
    customers: BackendCustomerResponse[];
}

/**
 * Matches CreateCustomerRequest from the Go backend.
 */
interface BackendCreateCustomerRequest {
    fullName: string; // Changed from name
    email: string;
    phone: string;
    address?: string; // omitempty in backend, maps from frontend's address
}

/**
 * Matches UpdateCustomerRequest from the Go backend.
 * All fields are optional.
 */
interface BackendUpdateCustomerRequest {
    fullName?: string; // Changed from name
    email?: string;
    phone?: string;
    address?: string;
}

// --- Frontend input type (already defined, kept for clarity) ---

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


// --- Mapping Functions ---

/**
 * Maps a BackendCustomerResponse to the frontend Customer type.
 * @param backendCustomer - The customer object from the backend.
 * @returns A frontend Customer object.
 */
function mapBackendToFrontendCustomer(backendCustomer: BackendCustomerResponse): Customer {
    return {
        id: backendCustomer.id,
        fullName: backendCustomer.fullName, // Was backendCustomer.name, now directly uses fullName
        email: backendCustomer.email,
        phone: backendCustomer.phone,
        address: backendCustomer.address || '', // Ensure address is a string
        dateRegistered: backendCustomer.dateRegistered,
        createdAt: backendCustomer.createdAt,
        updatedAt: backendCustomer.updatedAt,
    };
}

/**
 * Maps frontend NewCustomerInput to BackendCreateCustomerRequest.
 * @param frontendInput - The new customer data from the frontend.
 * @returns A BackendCreateCustomerRequest object.
 */
function mapFrontendToBackendCreateRequest(frontendInput: NewCustomerInput): BackendCreateCustomerRequest {
    return {
        fullName: frontendInput.fullName, // Was name: frontendInput.fullName
        email: frontendInput.email,
        phone: frontendInput.phone,
        // Send address if provided, otherwise it will be omitted by backend if empty and omitempty is set (or sent as empty string)
        address: frontendInput.address,
    };
}

/**
 * Maps frontend UpdateCustomerInput to BackendUpdateCustomerRequest.
 * @param frontendInput - The customer update data from the frontend.
 * @returns A BackendUpdateCustomerRequest object.
 */
function mapFrontendToBackendUpdateRequest(frontendInput: UpdateCustomerInput): BackendUpdateCustomerRequest {
    const backendRequest: BackendUpdateCustomerRequest = {};
    if (frontendInput.fullName !== undefined) {
        backendRequest.fullName = frontendInput.fullName; // Was backendRequest.name
    }
    if (frontendInput.email !== undefined) {
        backendRequest.email = frontendInput.email;
    }
    if (frontendInput.phone !== undefined) {
        backendRequest.phone = frontendInput.phone;
    }
    if (frontendInput.address !== undefined) {
        backendRequest.address = frontendInput.address;
    }
    return backendRequest;
}


export const customerApi = {
    /**
     * Fetches all customers from the backend.
     * @returns A promise that resolves to an array of Customer objects.
     */
    getAllCustomers: async (): Promise<Customer[]> => {
        const response = await api.get<BackendCustomerListResponse>('/api/customers'); // Use api and relative path
        // Check if response.data and response.data.customers exist
        console.log('Fetched customers from API', response.data);
        if (response.data && Array.isArray(response.data.customers)) {
            return response.data.customers.map(mapBackendToFrontendCustomer);
        }
        return []; // Return empty array if data is not in expected format
    },

    /**
     * Fetches a specific customer by their ID.
     * @param id - The ID of the customer to fetch.
     * @returns A promise that resolves to a Customer object.
     */
    getCustomerById: async (id: string): Promise<Customer> => {
        const response = await api.get<BackendCustomerResponse>(`/api/customers/${id}`); // Use api and relative path
        return mapBackendToFrontendCustomer(response.data);
    },

    /**
     * Adds a new customer to the backend.
     * @param customerData - The data for the new customer.
     * @returns A promise that resolves to the created Customer object (including backend-generated fields).
     */
    addCustomer: async (customerData: NewCustomerInput): Promise<Customer> => {
        const backendRequestData = mapFrontendToBackendCreateRequest(customerData);
        const response = await api.post<BackendCustomerResponse>('/api/customers', backendRequestData); // Use api and relative path
        return mapBackendToFrontendCustomer(response.data);
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
        const backendRequestData = mapFrontendToBackendUpdateRequest(customerData);
        const response = await api.put<BackendCustomerResponse>(`/api/customers/${id}`, backendRequestData); // Use api and relative path
        return mapBackendToFrontendCustomer(response.data);
    },

    /**
     * Deletes a customer from the backend.
     * @param id - The ID of the customer to delete.
     * @returns A promise that resolves to an object indicating success or failure.
     */
    deleteCustomer: async (id: string): Promise<{ success: boolean; message?: string }> => {
        try {
            // Backend returns 204 No Content on successful deletion
            await api.delete(`/api/customers/${id}`); // Use api and relative path
            return { success: true, message: 'Customer deleted successfully.' };
        } catch (error: unknown) { // Add type annotation for error
            let errorMessage = 'Failed to delete customer.';
            if (axios.isAxiosError(error)) { // Use axios.isAxiosError
                // If backend returns 404, it means customer was not found.
                // For deletion, this can often be treated as success (idempotency).
                if (error.response?.status === 404) {
                    console.warn(`Customer with ID ${id} not found for deletion, or already deleted.`);
                    // Depending on requirements, this might be success: true or success: false with specific message
                    return { success: true, message: 'Customer not found or already deleted.' };
                }
                // Try to get more specific error message from backend response if available
                if (error.response?.data && typeof error.response.data.Error === 'string') {
                    errorMessage = error.response.data.Error;
                } else if (error.response?.data && typeof error.response.data.message === 'string') {
                    errorMessage = error.response.data.message;
                }
            }
            console.error(`Error deleting customer ${id}:`, error);
            return { success: false, message: errorMessage };
        }
    }
}; 
