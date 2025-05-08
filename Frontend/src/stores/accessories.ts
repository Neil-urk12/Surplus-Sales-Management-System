import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type {
  AccessoryRow,
  NewAccessoryInput,
  UpdateAccessoryInput,
  AccessoryOperationResponse,
  AccessoryMake,
  AccessoryColor,
  AccessoryStatus,
  AccessoryMakeInput,
  AccessoryColorInput
} from 'src/types/accessories'
import { accessoriesApi } from 'src/services/accessoriesApi'
import axios from 'axios'
import { useSearch } from 'src/utils/useSearch'

export type { AccessoryRow } from 'src/types/accessories'

export const useAccessoriesStore = defineStore('accessories', () => {
  // State
  const accessoryRows = ref<AccessoryRow[]>([])
  const isLoading = ref(false)
  const apiError = ref<string | null>(null)

  // Initialize data
  async function initializeAccessories() {
    try {
      isLoading.value = true
      apiError.value = null

      // Replace with API call using our service
      const accessories = await accessoriesApi.getAllAccessories()
      accessoryRows.value = accessories
    } catch (error: unknown) {
      console.error('Error initializing accessories:', error);
      if (axios.isAxiosError(error)) { // Axios specific error handling
        apiError.value = `Failed to load accessories: ${error.response?.status} - ${error.response?.data?.message || 'Unknown error'}`;
      } else if (error instanceof Error) {
        apiError.value = `Failed to load accessories: ${error.message}`;
      } else {
        apiError.value = 'Failed to load accessories: Unknown error';
      }
    } finally {
      isLoading.value = false
    }
  }

  const search = useSearch({
    onSearch: (value) => {
      accessorySearch.value = value;
    }
  });

  const rawAccessorySearch = ref('')
  const accessorySearch = ref('')
  let debounceTimeout: ReturnType<typeof setTimeout> | null = null

  // Debounce function to update accessorySearch after typing stops
  function updateAccessorySearch(value: string) {
    if (debounceTimeout) {
      clearTimeout(debounceTimeout)
    }

    debounceTimeout = setTimeout(() => {
      accessorySearch.value = value
    }, 300) // 300ms debounce delay
  }

  // Watch for changes in rawAccessorySearch
  watch(rawAccessorySearch, (newValue) => {
    updateAccessorySearch(newValue)
  })

  // Use input types that allow empty strings for filters
  const filterMake = ref<AccessoryMakeInput>('')
  const filterColor = ref<AccessoryColorInput>('')
  const filterStatus = ref<AccessoryStatus | ''>('')

  // Available options
  const makes: AccessoryMake[] = ['Generic', 'OEM', 'Aftermarket', 'Custom']
  const colors: AccessoryColor[] = ['Black', 'White', 'Silver', 'Chrome', 'Custom']
  const statuses: AccessoryStatus[] = ['In Stock', 'Low Stock', 'Out of Stock']

  // Computed
  const filteredAccessoryRows = computed(() => {
    return accessoryRows.value.filter(row => {
      const matchesMake = !filterMake.value || row.make === filterMake.value
      const matchesColor = !filterColor.value || row.unit_color === filterColor.value
      const matchesStatus = !filterStatus.value || row.status === filterStatus.value
      const matchesSearch = !accessorySearch.value ||
        row.name.toLowerCase().includes(accessorySearch.value.toLowerCase()) ||
        row.make.toLowerCase().includes(accessorySearch.value.toLowerCase())

      return matchesMake && matchesColor && matchesStatus && matchesSearch
    })
  })

  // Actions
  async function addAccessory(accessory: NewAccessoryInput): Promise<AccessoryOperationResponse> {
    try {
      isLoading.value = true
      apiError.value = null

      // Validate required fields (client-side)
      if (!accessory.make || !accessory.unit_color) {
        const errMessage = 'Make and color are required';
        apiError.value = errMessage;
        return {
          success: false,
          error: errMessage
        };
      }

      // Call the API service to add the accessory
      const result = await accessoriesApi.addAccessory(accessory)

      if (result.success && result.data) {
        // Add the accessory returned by the API (includes ID, status, timestamps etc.)
        accessoryRows.value.push(result.data);
      } else if (result.error) {
        apiError.value = result.error;
      }

      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred while adding accessory'
      apiError.value = errorMessage
      return {
        success: false,
        error: errorMessage
      }
    } finally {
      isLoading.value = false
    }
  }

  function resetFilters() {
    filterMake.value = ''
    filterColor.value = ''
    filterStatus.value = ''
    rawAccessorySearch.value = ''
    accessorySearch.value = ''
    if (search && typeof search.clearSearch === 'function') {
      search.clearSearch();
    }
  }

  async function deleteAccessory(id: number): Promise<AccessoryOperationResponse> {
    try {
      isLoading.value = true
      apiError.value = null

      // Call the API service to delete the accessory
      const result = await accessoriesApi.deleteAccessory(id)

      if (result.success) {
        // Update local state
        const index = accessoryRows.value.findIndex(a => a.id === id);
        if (index !== -1) {
          accessoryRows.value.splice(index, 1);
        }
        return result;
      } else if (result.error) {
        apiError.value = result.error;
      }

      return result
    } catch (error) {
      console.error('Error in store.deleteAccessory:', error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred while deleting accessory'
      apiError.value = errorMessage
      return {
        success: false,
        error: errorMessage
      }
    } finally {
      isLoading.value = false
    }
  }

  async function updateAccessory(id: number, accessoryUpdatePayload: UpdateAccessoryInput): Promise<AccessoryOperationResponse> {
    try {
      isLoading.value = true
      apiError.value = null

      // Call the API service to update the accessory
      const result = await accessoriesApi.updateAccessory(id, accessoryUpdatePayload)

      if (result.success && result.data) {
        // Update local state with the accessory returned by the API
        const index = accessoryRows.value.findIndex(a => a.id === id);
        if (index !== -1) {
          accessoryRows.value[index] = result.data;
        }
      } else if (result.error) {
        apiError.value = result.error;
      }

      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred while updating accessory'
      apiError.value = errorMessage
      return {
        success: false,
        error: errorMessage
      }
    } finally {
      isLoading.value = false
    }
  }

  return {
    // State
    accessoryRows,
    isLoading,
    apiError,
    rawAccessorySearch,
    accessorySearch,
    search,
    filterMake,
    filterColor,
    filterStatus,
    // Constants
    makes,
    colors,
    statuses,
    // Computed
    filteredAccessoryRows,
    // Actions
    initializeAccessories,
    addAccessory,
    resetFilters,
    deleteAccessory,
    updateAccessory,
    updateAccessorySearch
  }
}) 
