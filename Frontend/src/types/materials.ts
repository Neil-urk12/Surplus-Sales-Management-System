/**
 * Available material categories
 */
export type MaterialCategory = 'Lumber' | 'Building' | 'Electrical' | 'Plumbing' | 'Hardware';

/**
 * Available suppliers
 */
export type MaterialSupplier = 'Steel Co.' | 'Construction Supplies Inc.' | 'Wood Works';

/**
 * Available inventory status values
 */
export type MaterialStatus = 'In Stock' | 'Low Stock' | 'Out of Stock';

/**
 * Form input types that allow empty values
 */
export type MaterialCategoryInput = MaterialCategory | '';
export type MaterialSupplierInput = MaterialSupplier | '';

/**
 * Represents a material in the inventory system
 */
export interface MaterialRow {
  /** Unique identifier for the material */
  id: number;

  /** Name of the material */
  name: string;

  /** Category of the material */
  category: MaterialCategory;

  /** Supplier of the material */
  supplier: MaterialSupplier;

  /** Number of units available in inventory */
  quantity: number;

  /** Current inventory status */
  status: MaterialStatus;

  /** URL or base64 string of the material image */
  image: string;
}

/**
 * Type for creating a new material (omits the id which is auto-generated)
 * Uses input types that allow empty values for form handling
 */
export interface NewMaterialInput extends Omit<MaterialRow, 'id' | 'category' | 'supplier'> {
  category: MaterialCategoryInput;
  supplier: MaterialSupplierInput;
}

/**
 * Type for updating an existing material
 */
export type UpdateMaterialInput = Partial<NewMaterialInput>;

/**
 * Response type for material operations
 */
export interface MaterialOperationResponse {
  success: boolean;
  id?: number;
  error?: string;
} 
