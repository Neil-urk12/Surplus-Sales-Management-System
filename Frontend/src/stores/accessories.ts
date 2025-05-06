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
    } catch (error) {
      console.error('Error initializing accessories:', error)
      apiError.value = error instanceof Error ? error.message : 'Failed to load accessories'
    } finally {
      isLoading.value = false
    }
  }

  // Search with debounce
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
  const statuses: AccessoryStatus[] = ['In Stock', 'Low Stock', 'Out of Stock', 'Available']

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

  // Type guard functions
  function isValidAccessoryMake(make: AccessoryMakeInput): make is AccessoryMake {
    return make !== '';
  }

  function isValidAccessoryColor(color: AccessoryColorInput): color is AccessoryColor {
    return color !== '';
  }

  // Actions
  async function addAccessory(accessory: NewAccessoryInput): Promise<AccessoryOperationResponse> {
    try {
      isLoading.value = true
      apiError.value = null

      // Validate required fields
      if (!accessory.make || !accessory.unit_color) {
        throw new Error('Make and color are required');
      }

      // Call the API service to add the accessory
      const result = await accessoriesApi.addAccessory(accessory)

      if (result.success && result.id) {
        // Create a new accessory with validated types
        const newAccessory: AccessoryRow = {
          id: result.id,
          name: accessory.name,
          make: accessory.make as AccessoryMake,
          quantity: accessory.quantity,
          price: accessory.price,
          status: accessory.quantity > 0 ? 'In Stock' : 'Out of Stock',
          unit_color: accessory.unit_color as AccessoryColor,
          image: accessory.image
        };

        // Add to local state
        accessoryRows.value.push(newAccessory);
      }

      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred'
      apiError.value = errorMessage
      return {
        success: false,
        error: errorMessage
      }
    } finally {
      isLoading.value = false
    }
  }

  function updateAccessoryStatus(id: number, quantity: number) {
    const accessory = accessoryRows.value.find(a => a.id === id)
    if (accessory) {
      if (quantity === 0) {
        accessory.status = 'Out of Stock'
      } else if (quantity <= 2) {
        accessory.status = 'Low Stock'
      } else if (quantity <= 5) {
        accessory.status = 'In Stock'
      } else {
        accessory.status = 'Available'
      }
    }
  }

  function resetFilters() {
    filterMake.value = ''
    filterColor.value = ''
    filterStatus.value = ''
    rawAccessorySearch.value = ''
    accessorySearch.value = ''
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
      }

      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred'
      apiError.value = errorMessage
      return {
        success: false,
        error: errorMessage
      }
    } finally {
      isLoading.value = false
    }
  }

  async function updateAccessory(id: number, accessory: UpdateAccessoryInput): Promise<AccessoryOperationResponse> {
    try {
      isLoading.value = true
      apiError.value = null

      // Call the API service to update the accessory
      const result = await accessoriesApi.updateAccessory(id, accessory)

      if (result.success) {
        // Update local state
        const index = accessoryRows.value.findIndex(a => a.id === id);
        if (index !== -1) {
          const existingAccessory = accessoryRows.value[index];

          // Use type guards to validate make and color
          const updatedMake = accessory.make && isValidAccessoryMake(accessory.make)
            ? accessory.make
            : existingAccessory.make;

          const updatedColor = accessory.unit_color && isValidAccessoryColor(accessory.unit_color)
            ? accessory.unit_color
            : existingAccessory.unit_color;

          // Update properties
          accessoryRows.value[index] = {
            ...existingAccessory,
            name: accessory.name ?? existingAccessory.name,
            make: updatedMake,
            quantity: accessory.quantity ?? existingAccessory.quantity,
            price: accessory.price ?? existingAccessory.price,
            unit_color: updatedColor,
            image: accessory.image ?? existingAccessory.image
          };

          // Update status based on quantity
          if (typeof accessory.quantity === 'number') {
            updateAccessoryStatus(id, accessory.quantity);
          }
        }
      }

      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred'
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
    updateAccessoryStatus,
    resetFilters,
    deleteAccessory,
    updateAccessory
  }
}) 
