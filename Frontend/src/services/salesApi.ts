/**
 * API service for sales operations
 */
import { api } from 'src/boot/axios';
import type { CabSalePayload, CabSale, SalesOperationResponse, Sale } from 'src/types/salesTypes';

/**
 * Service for handling sales-related API calls
 */
export const salesService = {
  /**
   * Sell a cab with optional accessories
   * @param cabId - ID of the cab being sold
   * @param payload - Sale details including customer, quantity, and accessories
   * @returns Promise with the sale operation response
   */
  async sellCab(cabId: number, payload: CabSalePayload): Promise<SalesOperationResponse> {
    try {
      const response = await api.post<CabSale>(`/api/cabs/${cabId}/sell`, payload);
      return {
        success: true,
        sale: {
          id: response.data.cabId.toString(), // Using cabId as a placeholder for sale ID
          customerId: response.data.customerId,
          soldBy: 'system', // This would typically come from the authenticated user
          saleDate: response.data.saleDate,
          totalPrice: response.data.totalPrice,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      };
    } catch (error) {
      console.error('Error selling cab:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to process sale'
      };
    }
  },

  /**
   * Get all sales records
   * @returns Promise with array of sales
   */
  async getSales(): Promise<Sale[]> {
    try {
      const response = await api.get<Sale[]>('/api/sales');
      return response.data;
    } catch (error) {
      console.error('Error fetching sales:', error);
      return [];
    }
  },

  /**
   * Get a specific sale by ID
   * @param saleId - ID of the sale to retrieve
   * @returns Promise with the sale details
   */
  async getSaleById(saleId: string): Promise<Sale | null> {
    try {
      const response = await api.get<Sale>(`/api/sales/${saleId}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching sale ${saleId}:`, error);
      return null;
    }
  },

  /**
   * Get sales for a specific customer
   * @param customerId - ID of the customer
   * @returns Promise with array of customer's sales
   */
  async getCustomerSales(customerId: string): Promise<Sale[]> {
    try {
      const response = await api.get<Sale[]>(`/api/customers/${customerId}/sales`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching sales for customer ${customerId}:`, error);
      return [];
    }
  }
};
