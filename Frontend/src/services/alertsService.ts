import { accessoriesApi } from './accessoriesApi';
import { api } from 'boot/axios';
import type { MaterialRow } from 'src/types/materials';
import { useAuthStore } from 'src/stores/auth';
import { useCabsStore } from 'src/stores/cabs';

/**
 * Category types for inventory items
 */
export type AlertCategory = 'Cabs' | 'Accessories' | 'Materials';

/**
 * Status types for inventory alerts
 */
export type AlertStatus = 'Low Stock' | 'Out of Stock';

/**
 * Interface for inventory alert data
 */
export interface InventoryAlert {
  category: AlertCategory;
  status: AlertStatus;
  count: number;
}

/**
 * Service for handling system alerts related to inventory
 */
export const alertsService = {
  /**
   * Fetches all inventory items with low stock or out of stock status
   * @returns Promise with array of inventory alerts
   */
  getInventoryAlerts: async (): Promise<InventoryAlert[]> => {
    try {
      const authStore = useAuthStore();
      const cabsStore = useCabsStore();

      // Initialize cabs if needed
      if (cabsStore.cabRows.length === 0) {
        await cabsStore.initializeCabs();
      }

      // Fetch data from all three inventory types
      const [accessories, materials] = await Promise.all([
        accessoriesApi.getAllAccessories(),
        api.get<MaterialRow[]>('/api/materials', {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        }).then(response => response.data)
      ]);

      // Use the cabs from the store
      const cabs = cabsStore.cabRows;

      // Process alerts for each category
      const alerts: InventoryAlert[] = [];

      // Process cabs alerts
      const lowStockCabs = cabs.filter(cab => cab.status === 'Low Stock');
      const outOfStockCabs = cabs.filter(cab => cab.status === 'Out of Stock');

      if (lowStockCabs.length > 0) {
        alerts.push({
          category: 'Cabs',
          status: 'Low Stock',
          count: lowStockCabs.length
        });
      }

      if (outOfStockCabs.length > 0) {
        alerts.push({
          category: 'Cabs',
          status: 'Out of Stock',
          count: outOfStockCabs.length
        });
      }

      // Process accessories alerts
      const lowStockAccessories = accessories.filter(accessory => accessory.status === 'Low Stock');
      const outOfStockAccessories = accessories.filter(accessory => accessory.status === 'Out of Stock');

      if (lowStockAccessories.length > 0) {
        alerts.push({
          category: 'Accessories',
          status: 'Low Stock',
          count: lowStockAccessories.length
        });
      }

      if (outOfStockAccessories.length > 0) {
        alerts.push({
          category: 'Accessories',
          status: 'Out of Stock',
          count: outOfStockAccessories.length
        });
      }

      // Process materials alerts
      const lowStockMaterials = materials.filter(material => material.status === 'Low Stock');
      const outOfStockMaterials = materials.filter(material => material.status === 'Out of Stock');

      if (lowStockMaterials.length > 0) {
        alerts.push({
          category: 'Materials',
          status: 'Low Stock',
          count: lowStockMaterials.length
        });
      }

      if (outOfStockMaterials.length > 0) {
        alerts.push({
          category: 'Materials',
          status: 'Out of Stock',
          count: outOfStockMaterials.length
        });
      }

      return alerts;
    } catch (error) {
      console.error('Error fetching inventory alerts:', error);
      return [];
    }
  }
};
