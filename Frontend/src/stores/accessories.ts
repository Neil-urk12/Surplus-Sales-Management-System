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
import { useDashboardStore } from 'src/stores/dashboardStore'
import { v4 as uuidv4 } from 'uuid'

export type { AccessoryRow } from 'src/types/accessories'

export const useAccessoriesStore = defineStore('accessories', () => {
  // State
  const accessoryRows = ref<AccessoryRow[]>([])
  const isLoading = ref(false)
  const apiError = ref<string | null>(null)
  const dashboardStore = useDashboardStore()

  // Initialize data
  async function initializeAccessories() {
    try {
      isLoading.value = true
      apiError.value = null
      
      // Replace with API call using our service
      const accessories = await accessoriesApi.getAllAccessories()
      accessoryRows.value = accessories
    } catch (error: unknown) {
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

  // Helper function to create activity records
  function createActivityRecord(activity: {
    title: string;
    description: string;
    icon: string;
    color: string;
  }) {
    dashboardStore.addActivity({
      id: uuidv4(),
      title: activity.title,
      description: activity.description,
      timestamp: new Date(),
      icon: activity.icon,
      color: activity.color
    });
  }

  // Format price using Intl.NumberFormat
  function formatPrice(price: number): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(price);
  }

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
        
        // Add activity record
        createActivityRecord({
          title: 'Accessory Added',
          description: `Added ${accessory.quantity || 1} units of ${accessory.name}`,
          icon: 'build',
          color: 'info'
        });
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
        // If deletion was successful, find the accessory details from local state
        const index = accessoryRows.value.findIndex(a => a.id === id);
        if (index !== -1) {
          const deletedAccessory = accessoryRows.value[index];
          
          // Remove from local state
          accessoryRows.value.splice(index, 1);
          
          // Add activity record with the accessory details
          if (deletedAccessory) {
            createActivityRecord({
              title: 'Accessory Removed',
              description: `Removed ${deletedAccessory.name} from inventory`,
              icon: 'delete',
              color: 'negative'
            });
          }
        }
        
        return result;
      } else if (result.error) {
        apiError.value = result.error;
      }

      return result
    } catch (error) {
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
      
      // Validate input
      if (!accessoryUpdatePayload) {
        const errMessage = 'Update payload is required';
        apiError.value = errMessage;
        return {
          success: false,
          error: errMessage
        };
      }
      
      // Get the existing accessory to compare later for activity recording
      let originalAccessory = null;
      const existingAccessory = accessoryRows.value.find(a => a.id === id);
      if (existingAccessory) {
        originalAccessory = JSON.parse(JSON.stringify(existingAccessory));
      }

      // Call the API service to update the accessory
      const result = await accessoriesApi.updateAccessory(id, accessoryUpdatePayload)

      if (result.success && result.data) {
        // Update local state with the accessory returned by the API
        const index = accessoryRows.value.findIndex(a => a.id === id);
        if (index !== -1) {
          accessoryRows.value[index] = result.data;
          
          // If we have the original accessory data, create an activity record
          if (originalAccessory) {
            const changes: string[] = [];
            
            // Build changes description
            if (accessoryUpdatePayload.quantity !== undefined && accessoryUpdatePayload.quantity !== originalAccessory.quantity) {
              const difference = accessoryUpdatePayload.quantity - originalAccessory.quantity;
              if (difference > 0) {
                changes.push(`added ${difference} units`);
              } else if (difference < 0) {
                changes.push(`removed ${Math.abs(difference)} units`);
              }
            }
            
            if (accessoryUpdatePayload.name !== undefined && accessoryUpdatePayload.name !== originalAccessory.name) {
              changes.push(`renamed to ${accessoryUpdatePayload.name}`);
            }
            
            if (accessoryUpdatePayload.price !== undefined && accessoryUpdatePayload.price !== originalAccessory.price) {
              changes.push(`price updated to ${formatPrice(accessoryUpdatePayload.price)}`);
            }
            
            // Only add activity if there were actual changes
            if (changes.length > 0) {
              createActivityRecord({
                title: 'Accessory Updated',
                description: `${originalAccessory.name}: ${changes.join(', ')}`,
                icon: 'edit',
                color: 'info'
              });
            }
          }
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
