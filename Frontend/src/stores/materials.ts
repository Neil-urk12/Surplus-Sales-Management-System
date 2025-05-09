import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type {
  MaterialRow,
  NewMaterialInput,
  UpdateMaterialInput,
  MaterialOperationResponse,
  MaterialCategory,
  MaterialSupplier,
  MaterialStatus,
  MaterialCategoryInput,
  MaterialSupplierInput
} from 'src/types/materials'
import { api } from 'boot/axios'
import { useAuthStore } from 'src/stores/auth'
import { useDashboardStore } from 'src/stores/dashboardStore'
import { AxiosError } from 'axios';
import type { QTableProps } from 'quasar'
import { v4 as uuidv4 } from 'uuid'

export type { MaterialRow, NewMaterialInput } from 'src/types/materials'

/** Default number of rows to display per page in the materials table */
const DEFAULT_ROWS_PER_PAGE = 10

// Activity constants
const ACTIVITY = {
  INVENTORY_ADDED: 'Inventory Added',
  INVENTORY_UPDATED: 'Inventory Updated',
  INVENTORY_REMOVED: 'Inventory Removed'
}

// Icon constants
const ICON = {
  INVENTORY: 'inventory_2',
  EDIT: 'edit',
  DELETE: 'delete'
}

// Color constants
const COLOR = {
  INFO: 'info',
  NEGATIVE: 'negative'
}

/**
 * Sanitizes a string to prevent XSS vulnerabilities
 * @param {string} str - String to sanitize
 * @returns {string} Sanitized string
 */
function sanitizeString(str: string | number): string {
  if (typeof str !== 'string') {
    str = String(str);
  }
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

/**
 * Pinia store for managing material data.
 * Provides state, computed properties, and actions for materials,
 * including fetching, adding, updating, deleting, and filtering.
 */
export const useMaterialsStore = defineStore('materials', () => {
  /**
   * Reactive state holding the array of material rows.
   * @type {MaterialRow[]}
   */
  const materialRows = ref<MaterialRow[]>([])
  /**
   * Reactive state indicating if an asynchronous operation is in progress.
   * @type {boolean}
   */
  const isLoading = ref(false)
  /**
   * Instance of the authentication store.
   */
  const authStore = useAuthStore()
  const dashboardStore = useDashboardStore()

  /**
   * Initializes the material data by fetching it from the API.
   * Requires an authentication token.
   * Handles loading state and error logging.
   */
  async function initializeMaterials() {
    if (!authStore.token) {
      console.error('No auth token found for initializing materials.')
      return
    }
    await onRequest({ pagination: pagination.value })
  }

  /**
   * Reactive state for the raw material search input value.
   * @type {string}
   */
  const rawMaterialSearch = ref('')

  // Watch for changes in rawMaterialSearch and trigger search with debounce
  let searchDebounceTimeout: ReturnType<typeof setTimeout> | null = null;
  watch(rawMaterialSearch, (newValue) => {
    console.log('Search input changed:', newValue);
    if (searchDebounceTimeout) {
      clearTimeout(searchDebounceTimeout);
    }
    searchDebounceTimeout = setTimeout(() => {
      console.log('Debounce timer expired, triggering search with query:', newValue);
      void onRequest({
        pagination: {
          ...pagination.value,
          page: 1 // Reset to first page on new search
        }
      });
    }, 300);
  });

  /**
   * Reactive state for filtering materials by category.
   * @type {MaterialCategoryInput | 'All'}
   */
  const filterCategory = ref<MaterialCategoryInput | 'All'>('All')
  /**
   * Reactive state for filtering materials by supplier.
   * @type {MaterialSupplierInput | 'All'}
   */
  const filterSupplier = ref<MaterialSupplierInput | 'All'>('All')
  /**
   * Reactive state for filtering materials by status.
   * @type {MaterialStatus | 'All'}
   */
  const filterStatus = ref<MaterialStatus | 'All'>('All')

  /**
   * Array of available material categories.
   * @type {MaterialCategory[]}
   */
  const categories: MaterialCategory[] = ['Lumber', 'Building', 'Electrical', 'Plumbing', 'Hardware']
  /**
   * Array of available material suppliers.
   * @type {MaterialSupplier[]}
   */
  const suppliers: MaterialSupplier[] = ['Steel Co.', 'Construction Supplies Inc.', 'Wood Works']
  /**
   * Array of available material statuses.
   * @type {MaterialStatus[]}
   */
  const statuses: MaterialStatus[] = ['In Stock', 'Low Stock', 'Out of Stock']

  /**
   * Computed property that returns the material rows filtered based on
   * category, supplier, status, and search term.
   * @returns {MaterialRow[]} The filtered array of material rows.
   */
  const filteredMaterialRows = computed(() => {
    return materialRows.value.filter(row => {
      const matchesCategory = !filterCategory.value || row.category === filterCategory.value
      const matchesSupplier = !filterSupplier.value || row.supplier === filterSupplier.value
      const matchesStatus = !filterStatus.value || row.status === filterStatus.value
      const matchesSearch = !rawMaterialSearch.value ||
        row.name.toLowerCase().includes(rawMaterialSearch.value.toLowerCase()) ||
        row.category.toLowerCase().includes(rawMaterialSearch.value.toLowerCase()) ||
        row.supplier.toLowerCase().includes(rawMaterialSearch.value.toLowerCase())

      return matchesCategory && matchesSupplier && matchesStatus && matchesSearch
    })
  })

  /**
   * Records a material activity in the dashboard.
   * @param {string} title - Activity title
   * @param {string} description - Activity description
   * @param {string} icon - Icon for the activity
   * @param {string} color - Color for the activity
   */
  function trackMaterialActivity(title: string, description: string, icon: string, color: string) {
    dashboardStore.addActivity({
      id: uuidv4(),
      title: title,
      description: description,
      timestamp: new Date(),
      icon: icon,
      color: color
    });
  }

  /**
   * Adds a new material by making an API call.
   * Requires an authentication token.
   * Adds the newly created material to the local store upon success.
   * Handles loading state and error handling.
   * @param {NewMaterialInput} material - The material data to add.
   * @returns {Promise<MaterialOperationResponse>} A promise resolving to the operation response.
   */
  async function addMaterial(material: NewMaterialInput): Promise<MaterialOperationResponse> {
    if (!authStore.token) {
      console.error('No auth token found for adding material.')
      return { success: false, error: 'Authentication required.' }
    }
    if (!material.name || !material.category || !material.supplier || material.quantity == null || !material.status) {
      return { success: false, error: 'Missing required material fields.' };
    }

    try {
      isLoading.value = true
      const response = await api.post<MaterialRow>('/api/materials', material, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })

      // Add the material returned from the server (with server-generated ID)
      materialRows.value.push(response.data)
      
      // Add activity record with sanitized values
      trackMaterialActivity(
        ACTIVITY.INVENTORY_ADDED,
        `Added ${sanitizeString(material.quantity)} units of ${sanitizeString(material.name)}`,
        ICON.INVENTORY,
        COLOR.INFO
      );
      
      return { success: true, id: response.data.id }
    } catch (error: unknown) {
      console.error('Error adding material:', error)
      let errorMessage = 'Unknown error occurred';
      if (error instanceof AxiosError) {
        errorMessage = error.response?.data?.error || error.message;
      } else if (error instanceof Error) {
        errorMessage = error.message;
      }
      return {
        success: false,
        error: errorMessage
      }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Updates the status of a material based on its quantity.
   * Finds the material by ID and updates its status in the local store.
   * @param {number} id - The ID of the material to update.
   * @param {number} quantity - The new quantity of the material.
   */
  function updateMaterialStatus(id: number, quantity: number) {
    const material = materialRows.value.find(m => m.id === id)
    if (material) {
      if (quantity === 0) {
        material.status = 'Out of Stock'
      } else if (quantity <= 10) {
        material.status = 'Low Stock'
      } else {
        material.status = 'In Stock'
      }
    }
  }

  /**
   * Resets all material filters and search terms to their default empty states.
   */
  function resetFilters() {
    filterCategory.value = 'All'
    filterSupplier.value = 'All'
    filterStatus.value = 'All'
    rawMaterialSearch.value = ''
  }

  /**
   * Deletes a material by making an API call.
   * Requires an authentication token.
   * Removes the material from the local store upon success.
   * Handles loading state and error handling.
   * @param {number} id - The ID of the material to delete.
   * @returns {Promise<MaterialOperationResponse>} A promise resolving to the operation response.
   */
  async function deleteMaterial(id: number): Promise<MaterialOperationResponse> {
    if (!authStore.token) {
      console.error('No auth token found for deleting material.')
      return { success: false, error: 'Authentication required.' }
    }

    const existingMaterial: MaterialRow | undefined = materialRows.value.find(m => m.id === id);
    if (!existingMaterial) {
      console.warn(`Material with ID ${id} not found locally for deletion.`);
      return { success: false, error: 'Material not found locally.' };
    }

    try {
      isLoading.value = true;
      await api.delete(`/api/materials/${id}`, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      });

      const index = materialRows.value.indexOf(existingMaterial);
      if (index !== -1) {
        materialRows.value.splice(index, 1);
      }
      
      // Add activity record with sanitized values
      trackMaterialActivity(
        ACTIVITY.INVENTORY_REMOVED,
        `Removed ${sanitizeString(existingMaterial.name)} from inventory`,
        ICON.DELETE,
        COLOR.NEGATIVE
      );

      return { success: true };
    } catch (error: unknown) {
      console.error('Error deleting material:', error);
      let errorMessage = 'Unknown error occurred while deleting material';
      if (error instanceof AxiosError) {
        errorMessage = error.response?.data?.error || error.message;
      } else if (error instanceof Error) {
        errorMessage = error.message;
      }

      return {
        success: false,
        error: errorMessage
      };
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Tracks changes between the original material and the update.
   * Creates a descriptive message about what changed.
   * @param {MaterialRow} existingMaterial - The original material data
   * @param {UpdateMaterialInput} materialUpdate - The updated material data
   * @returns {string[]} Array of change descriptions
   */
  function trackMaterialChanges(existingMaterial: MaterialRow, materialUpdate: UpdateMaterialInput): string[] {
    const changes: string[] = [];
    
    if (materialUpdate.quantity !== undefined && existingMaterial.quantity !== undefined) {
      const oldQuantity = Number(existingMaterial.quantity);
      const newQuantity = Number(materialUpdate.quantity);
      
      if (newQuantity !== oldQuantity) {
        const difference = newQuantity - oldQuantity;
        if (difference > 0) {
          changes.push(`added ${sanitizeString(difference)} units`);
        } else if (difference < 0) {
          changes.push(`removed ${sanitizeString(-difference)} units`);
        }
      }
    }
    
    if (materialUpdate.name !== undefined && materialUpdate.name !== existingMaterial.name) {
      changes.push(`renamed to ${sanitizeString(materialUpdate.name)}`);
    }
    
    return changes;
  }

  /**
   * Updates an existing material by making an API call.
   * Requires an authentication token.
   * Updates the material in the local store upon success.
   * Implements a rollback mechanism for the local state if the API call fails.
   * Handles loading state and error handling.
   * @param {number} id - The ID of the material to update.
   * @param {UpdateMaterialInput} materialUpdate - The updated material data.
   * @returns {Promise<MaterialOperationResponse>} A promise resolving to the operation response.
   */
  async function updateMaterial(id: number, materialUpdate: UpdateMaterialInput): Promise<MaterialOperationResponse> {
    if (!authStore.token) {
      console.error('No auth token found for updating material.')
      return { success: false, error: 'Authentication required.' }
    }

    const existingMaterial: MaterialRow | undefined = materialRows.value.find(m => m.id === id);
    if (!existingMaterial) {
      console.warn(`Material with ID ${id} not found locally for update.`);
      return { success: false, error: 'Material not found locally.' };
    }

    const index = materialRows.value.indexOf(existingMaterial);
    const originalMaterial = JSON.parse(JSON.stringify(existingMaterial));

    try {
      isLoading.value = true;
      const response = await api.put<MaterialRow>(`/api/materials/${id}`, materialUpdate, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      });

      materialRows.value[index] = response.data;
      
      // Track changes and add activity record if there were changes
      const changes = trackMaterialChanges(existingMaterial, materialUpdate);
      
      if (changes.length > 0) {
        trackMaterialActivity(
          ACTIVITY.INVENTORY_UPDATED,
          `${sanitizeString(existingMaterial.name)}: ${changes.join(', ')}`,
          ICON.EDIT,
          COLOR.INFO
        );
      }
      
      return { success: true };
    } catch (error: unknown) {
      console.error('Error updating material:', error);
      materialRows.value[index] = originalMaterial;

      let errorMessage = 'Unknown error occurred while updating material';
      if (error instanceof AxiosError) {
        errorMessage = error.response?.data?.error || error.message;
      } else if (error instanceof Error) {
        errorMessage = error.message;
      }

      return {
        success: false,
        error: errorMessage
      };
    } finally {
      isLoading.value = false;
    }
  }

  const pagination = ref({
    sortBy: 'name',
    descending: false,
    page: 1,
    rowsPerPage: DEFAULT_ROWS_PER_PAGE,
    rowsNumber: 0
  });

  async function onRequest(props: { pagination: QTableProps['pagination'] }) {
    if (!props.pagination) {
      console.log('No pagination provided, skipping request');
      return;
    }
    if (!authStore.token) {
      console.error('No auth token found for fetching materials.');
      return;
    }

    const { page = 1, rowsPerPage = DEFAULT_ROWS_PER_PAGE } = props.pagination;
    const params = {
      page,
      limit: rowsPerPage,
      search: rawMaterialSearch.value,
      category: filterCategory.value === 'All' ? '' : filterCategory.value,
      supplier: filterSupplier.value === 'All' ? '' : filterSupplier.value,
      status: filterStatus.value === 'All' ? '' : filterStatus.value
    };

    try {
      isLoading.value = true;
      console.log('Making API request with params:', params);

      const response = await api.get('/api/materials/paginated', {
        params,
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      });

      console.log('API response:', {
        materials: response.data.materials,
        total: response.data.total,
        page: response.data.page,
        limit: response.data.limit,
        totalPages: response.data.totalPages
      });

      materialRows.value = response.data.materials;
      pagination.value = {
        ...pagination.value,
        page,
        rowsPerPage,
        rowsNumber: response.data.total
      };
    } catch (error) {
      console.error('Error fetching paginated materials:', error);
      if (error instanceof AxiosError) {
        console.error('API error details:', {
          status: error.response?.status,
          data: error.response?.data,
          config: {
            url: error.config?.url,
            params: error.config?.params,
            headers: error.config?.headers
          }
        });
      }
      materialRows.value = [];
      pagination.value = {
        ...pagination.value,
        rowsNumber: 0
      };
    } finally {
      isLoading.value = false;
    }
  }

  return {
    // State
    materialRows,
    isLoading,
    rawMaterialSearch,
    filterCategory,
    filterSupplier,
    filterStatus,
    pagination,

    // Constants
    categories,
    suppliers,
    statuses,

    // Computed
    filteredMaterialRows,

    // Actions
    initializeMaterials,
    addMaterial,
    updateMaterial,
    deleteMaterial,
    resetFilters,
    updateMaterialStatus,
    onRequest
  }
})
