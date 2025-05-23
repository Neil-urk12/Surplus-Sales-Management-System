export interface User {
  id: string;
  fullName: string;
  username: string;
  email: string;
  role: 'admin' | 'staff';
  createdAt: string;
  updatedAt: string;
  isActive: boolean;
}

export interface Customer {
  id: string;
  fullName: string;
  email: string;
  phone: string;
  address: string;
  dateRegistered: string;
  createdAt: string;
  updatedAt: string;
}

export interface Sale {
  id: string;
  customerId: string;
  soldBy: string;
  saleDate: string;
  totalPrice: number;
  createdAt: string;
  updatedAt: string;
  multiCabId: string;
}

export interface SaleItem {
  id: string;
  saleId: string;
  itemType: string;
  multiCabId: string;
  accessoryId: string;
  materialId: string;
  quantity: number;
  unitPrice: number;
  subtotal: number;
  createdAt: string;
  updatedAt: string;
}

// Extended SaleItem interface that includes name
export interface ExtendedSaleItem extends SaleItem {
  name?: string;
}

export interface StockTransaction {
  id: string;
  userId: string;
  timestamp: string;
  type: string;
  quantityChange: number;
  remarks: string;
  accessoryId: string;
  materialId: string;
  createdAt: string;
  updatedAt: string;
}

export interface Accessory {
  id: string;
  name: string;
  quantity: number;
  createdAt: string;
  updatedAt: string;
}

export interface MultiCabAccessory {
  id: string;
  multiCabId: string;
  accessoryId: string;
  quantityAdded: number;
  dateApplied: string;
  createdAt: string;
  updatedAt: string;
}

export interface Material {
  id: string;
  name: string;
  quantity: number;
  createdAt: string;
  updatedAt: string;
}

export interface MultiCabMaterial {
  id: string;
  multiCabId: string;
  materialId: string;
  quantityUsed: number;
  dateApplied: string;
  createdAt: string;
  updatedAt: string;
}

export interface MultiCab {
  id: string;
  make: string;
  model: string;
  year: number;
  color: string;
  condition: string;
  price: number;
  status: string;
  dateAdded: string;
  serialNumber: string;
  createdAt: string;
  updatedAt: string;
}
