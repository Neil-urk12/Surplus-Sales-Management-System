/**
 * Types for the sales functionality
 * These types match the backend models in models.go
 */

// AccessoryForSale represents an accessory included in a cab sale
export interface AccessoryForSale {
  id: number;
  name: string;
  price: number;
  quantity: number;
  unitPrice: number;
}

// CabSalePayload represents the data sent to the backend to record a cab sale
export interface CabSalePayload {
  customerId: string;
  quantity: number;
  accessories: AccessoryForSale[];
}

// CabSale represents a completed cab sale transaction
export interface CabSale {
  cabId: number;
  customerId: string;
  quantity: number;
  accessories: AccessoryForSale[];
  totalPrice: number;
  saleDate: string;
}

// Sale represents a general sale record
export interface Sale {
  id: string;
  customerId: string;
  soldBy: string;
  saleDate: string;
  totalPrice: number;
  createdAt: string;
  updatedAt: string;
}

// SaleItem represents an item in a sale
export interface SaleItem {
  id: string;
  saleId: string;
  itemType: string;
  multiCabId?: string;
  accessoryId?: string;
  materialId?: string;
  quantity: number;
  unitPrice: number;
  subtotal: number;
  createdAt: string;
  updatedAt: string;
}

// SalesOperationResponse represents the response from a sales operation
export interface SalesOperationResponse {
  success: boolean;
  error?: string;
  sale?: Sale;
}

// SellCabResponse represents the response from selling a cab
export interface SellCabResponse {
  success: boolean;
  error?: string;
  cabSale?: CabSale;
}
