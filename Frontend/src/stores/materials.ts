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
import { AxiosError } from 'axios';
import { useSearch } from 'src/utils/useSearch';
import { operationNotifications } from 'src/utils/notifications';

export type { MaterialRow, NewMaterialInput } from 'src/types/materials'

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
    try {
      isLoading.value = true
      const response = await api.get<MaterialRow[]>('/api/materials', {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })
      materialRows.value = response.data
    } catch (error) {
      console.error('Error initializing materials:', error)
      materialRows.value = []
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Setup search with the composable
   */
  const search = useSearch({
    onSearch: (value) => {
      materialSearch.value = value;
    }
  });

  /**
   * Reactive state for the raw material search input value.
   * @type {string}
   */
  const rawMaterialSearch = ref('')
  /**
   * Reactive state for the debounced material search value.
   * This value is updated after a short delay from rawMaterialSearch.
   * @type {string}
   */
  const materialSearch = ref('')

  /**
   * Updates the search value and triggers the search.
   * @param {string} value - The new search value.
   */
  function updateSearch(value: string) {
    rawMaterialSearch.value = value;
    materialSearch.value = value;
  }

  /**
   * Watches for changes in rawMaterialSearch and updates materialSearch.
   */
  watch(rawMaterialSearch, (newValue) => {
    materialSearch.value = newValue;
  })

  /**
   * Reactive state for filtering materials by category.
   * @type {MaterialCategoryInput}
   */
  const filterCategory = ref<MaterialCategoryInput>('')
  /**
   * Reactive state for filtering materials by supplier.
   * @type {MaterialSupplierInput}
   */
  const filterSupplier = ref<MaterialSupplierInput>('')
  /**
   * Reactive state for filtering materials by status.
   * @type {MaterialStatus | ''}
   */
  const filterStatus = ref<MaterialStatus | ''>('')

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
    const filtered = materialRows.value.filter(row => {
      const matchesCategory = !filterCategory.value || row.category === filterCategory.value
      const matchesSupplier = !filterSupplier.value || row.supplier === filterSupplier.value
      const matchesStatus = !filterStatus.value || row.status === filterStatus.value
      const matchesSearch = !materialSearch.value ||
        row.name.toLowerCase().includes(materialSearch.value.toLowerCase()) ||
        row.category.toLowerCase().includes(materialSearch.value.toLowerCase()) ||
        row.supplier.toLowerCase().includes(materialSearch.value.toLowerCase())

      return matchesCategory && matchesSupplier && matchesStatus && matchesSearch
    })
    
    // Sort by ID to ensure items are in order with new items (highest IDs) at the end
    return filtered.sort((a, b) => a.id - b.id)
  })

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

      // Find the highest ID to ensure new materials are always at the end
      // and to avoid reusing deleted IDs
      const maxId = materialRows.value.length > 0
        ? Math.max(...materialRows.value.map(m => m.id))
        : 0;
      
      // If the response doesn't include an ID, assign a new one
      if (!response.data.id) {
        response.data.id = maxId + 1;
      }

      // Add to the end of the array
      materialRows.value.push(response.data)
      operationNotifications.add.success(`material: ${material.name}`);
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
    filterCategory.value = ''
    filterSupplier.value = ''
    filterStatus.value = ''
    rawMaterialSearch.value = ''
    materialSearch.value = ''
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

    const existingMaterial = materialRows.value.find(m => m.id === id);
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

      // Remove the material from the local state without affecting other IDs
      const index = materialRows.value.findIndex(m => m.id === id);
      if (index !== -1) {
        materialRows.value.splice(index, 1);
      }
      
      operationNotifications.delete.success('material');
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

    const existingMaterial = materialRows.value.find(m => m.id === id);
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
      operationNotifications.update.success(`material: ${materialUpdate.name}`);
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

  return {
    materialRows,
    isLoading,
    rawMaterialSearch,
    materialSearch,
    search,
    filterCategory,
    filterSupplier,
    filterStatus,
    categories,
    suppliers,
    statuses,
    filteredMaterialRows,
    initializeMaterials,
    addMaterial,
    updateMaterialStatus,
    resetFilters,
    deleteMaterial,
    updateMaterial,
    updateSearch
  }
})