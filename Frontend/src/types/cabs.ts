/**
 * Available cab manufacturers
 */
export type CabMake = 'Mazda' | 'Porsche' | 'Toyota' | 'Nissan' | 'Ford';

/**
 * Available cab colors
 */
export type CabColor = 'Black' | 'White' | 'Silver' | 'Red' | 'Blue';

/**
 * Available inventory status values
 */
export type CabStatus = 'In Stock' | 'Low Stock' | 'Out of Stock' | 'Available';

/**
 * Form input types that allow empty values
 */
export type CabMakeInput = CabMake | '';
export type CabColorInput = CabColor | '';

/**
 * Represents a cab in the inventory system
 */
export interface CabsRow {
  /** Unique identifier for the cab */
  id: number;
  
  /** Name of the cab model */
  name: string;
  
  /** Manufacturer of the cab */
  make: CabMake;
  
  /** Number of units available in inventory */
  quantity: number;
  
  /** Price in Philippine Peso (PHP) */
  price: number;
  
  /** Current inventory status */
  status: CabStatus;
  
  /** Color of the cab unit */
  unit_color: CabColor;
  
  /** URL or base64 string of the cab image */
  image: string;
}

/**
 * Type for creating a new cab (omits the id which is auto-generated)
 * Uses input types that allow empty values for form handling
 */
export interface NewCabInput extends Omit<CabsRow, 'id' | 'make' | 'unit_color'> {
  make: CabMakeInput;
  unit_color: CabColorInput;
}

/**
 * Type for updating an existing cab
 */
export type UpdateCabInput = Partial<NewCabInput>;

/**
 * Response type for cab operations
 */
export interface CabOperationResponse {
  success: boolean;
  id?: number;
  error?: string;
}

export interface ImageValidationResult {
  isValid: boolean;
  error?: string;
}

export interface ImageValidationOptions {
  maxBase64Length?: number;  // in bytes
  allowedImageTypes?: string[];
  maxFileSizeMB?: number;
} 