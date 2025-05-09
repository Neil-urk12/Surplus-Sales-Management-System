/**
 * API service for sales operations
 */
import { api } from 'src/boot/axios';
import type { CabSalePayload, CabSale, SalesOperationResponse, Sale, SaleItem, SellCabResponse } from 'src/types/salesTypes';

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
  async sellCab(cabId: number, payload: CabSalePayload): Promise<SellCabResponse> {
    try {
      const response = await api.post<CabSale>(`/api/cabs/${cabId}/sell`, payload);
      return {
        success: true,
        cabSale: response.data
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
  },

  /**
   * Get all items for a specific sale
   * @param saleId - ID of the sale
   * @returns Promise with array of sale items
   */
  async getSaleItems(saleId: string): Promise<SaleItem[]> {
    try {
      const response = await api.get<SaleItem[]>(`/api/sales/${saleId}/items`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching items for sale ${saleId}:`, error);
      return [];
    }
  },

  /**
   * Create a new sale
   * @param saleData - Data for the new sale
   * @returns Promise with the sale operation response
   */
  async createSale(saleData: Sale): Promise<SalesOperationResponse> {
    try {
      const response = await api.post<Sale>('/api/sales', saleData);
      return {
        success: true,
        sale: response.data
      };
    } catch (error) {
      console.error('Error creating sale:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to create sale'
      };
    }
  },

  /**
   * Update an existing sale
   * @param saleId - ID of the sale to update
   * @param saleData - Updated sale data
   * @returns Promise with the sale operation response
   */
  async updateSale(saleId: string, saleData: Partial<Sale>): Promise<SalesOperationResponse> {
    try {
      const response = await api.put<Sale>(`/api/sales/${saleId}`, saleData);
      return {
        success: true,
        sale: response.data
      };
    } catch (error) {
      console.error(`Error updating sale ${saleId}:`, error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to update sale'
      };
    }
  },

  /**
   * Delete a sale
   * @param saleId - ID of the sale to delete
   * @returns Promise with the sale operation response
   */
  async deleteSale(saleId: string): Promise<SalesOperationResponse> {
    try {
      await api.delete(`/api/sales/${saleId}`);
      return { success: true };
    } catch (error) {
      console.error(`Error deleting sale ${saleId}:`, error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to delete sale'
      };
    }
  }
};
