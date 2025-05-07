// Available accessory makes/brands
export type AccessoryMake = 'Generic' | 'OEM' | 'Aftermarket' | 'Custom';

// Available accessory colors
export type AccessoryColor = 'Black' | 'White' | 'Silver' | 'Chrome' | 'Custom';

// Available accessory statuses
export type AccessoryStatus = 'In Stock' | 'Low Stock' | 'Out of Stock' | 'Available';

// Input types that allow empty strings for form handling
export type AccessoryMakeInput = AccessoryMake | '';
export type AccessoryColorInput = AccessoryColor | '';

/**
 * Represents an accessory in the inventory system
 */
export interface AccessoryRow {
  /** Unique identifier for the accessory */
  id: number;

  /** Name of the accessory */
  name: string;

  /** Manufacturer/brand of the accessory */
  make: AccessoryMake;

  /** Number of units available in inventory */
  quantity: number;

  /** Price in Philippine Peso (PHP) */
  price: number;

  /** Current inventory status */
  status: AccessoryStatus;

  /** Color of the accessory */
  unit_color: AccessoryColor;

  /** URL or base64 string of the accessory image */
  image: string;
}

/**
 * Type for creating a new accessory (omits the id which is auto-generated)
 * Uses input types that allow empty values for form handling
 */
export interface NewAccessoryInput extends Omit<AccessoryRow, 'id' | 'make' | 'unit_color'> {
  make: AccessoryMakeInput;
  unit_color: AccessoryColorInput;
}

/**
 * Type for updating an existing accessory
 * Makes all fields optional except id
 */
export type UpdateAccessoryInput = Partial<NewAccessoryInput>;

/**
 * Response type for accessory operations
 */
export interface AccessoryOperationResponse {
  success: boolean;
  id?: number;
  data?: AccessoryRow;
  error?: string;
  message?: string;
  statusCode?: number;
  timestamp?: string;
}

/**
 * API response for retrieving multiple accessories
 */
export interface AccessoriesListResponse {
  data: AccessoryRow[];
  count: number;
  page?: number;
  pageSize?: number;
  totalPages?: number;
}

/**
 * API connection error type
 */
export interface ApiError {
  message: string;
  statusCode?: number;
  details?: string;
} 
