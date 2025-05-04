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

export type { AccessoryRow } from 'src/types/accessories'

export const useAccessoriesStore = defineStore('accessories', () => {
  // State
  const accessoryRows = ref<AccessoryRow[]>([])
  const isLoading = ref(false)
  
  // Initialize data
  async function initializeAccessories() {
    try {
      isLoading.value = true
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500))
      accessoryRows.value = [
        {
          name: 'LED Headlights',
          id: 1,
          make: 'OEM',
          quantity: 10,
          price: 15000,
          status: 'In Stock',
          unit_color: 'Chrome',
          image: 'https://loremflickr.com/600/400/headlight',
        },
        {
          name: 'Alloy Wheels',
          id: 2,
          make: 'Aftermarket',
          quantity: 8,
          price: 25000,
          status: 'In Stock',
          unit_color: 'Silver',
          image: 'https://loremflickr.com/600/400/wheel',
        },
        {
          name: 'Seat Covers',
          id: 3,
          make: 'Generic',
          quantity: 15,
          price: 5000,
          status: 'Available',
          unit_color: 'Black',
          image: 'https://loremflickr.com/600/400/seatcover',
        },
      ]
    } catch (error) {
      console.error('Error initializing accessories:', error)
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

  /**
   * Determines the appropriate status based on the quantity
   * @param quantity - The current quantity of the accessory
   * @returns The appropriate AccessoryStatus
   * 
   * Quantity thresholds:
   * - 0: Out of Stock
   * - 1-2: Low Stock
   * - 3-10: In Stock
   * - >10: Available
   */
  function determineAccessoryStatus(quantity: number): AccessoryStatus {
    if (quantity === 0) {
      return 'Out of Stock';
    } else if (quantity <= 2) {
      return 'Low Stock';
    } else if (quantity <= 10) {
      return 'In Stock';
    } else {
      return 'Available';
    }
  }

  // Actions
  async function addAccessory(accessory: NewAccessoryInput): Promise<AccessoryOperationResponse> {
    try {
      isLoading.value = true;
      // Validate required fields
      if (!accessory.make || !accessory.unit_color) {
        throw new Error('Make and color are required');
      }

      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const newId = Math.max(...accessoryRows.value.map(item => item.id)) + 1;
      
      // Create a new accessory with validated types and determined status
      const newAccessory: AccessoryRow = {
        id: newId,
        name: accessory.name,
        make: accessory.make,
        quantity: accessory.quantity,
        price: accessory.price,
        status: determineAccessoryStatus(accessory.quantity), // Set initial status based on quantity
        unit_color: accessory.unit_color,
        image: accessory.image
      };

      accessoryRows.value.push(newAccessory);

      return { success: true, id: newId };
    } catch (error) {
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      };
    } finally {
      isLoading.value = false;
    }
  }

  function updateAccessoryStatus(id: number, quantity: number) {
    const accessory = accessoryRows.value.find(a => a.id === id);
    if (accessory) {
      accessory.status = determineAccessoryStatus(quantity);
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
      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const index = accessoryRows.value.findIndex(a => a.id === id);
      if (index !== -1) {
        accessoryRows.value.splice(index, 1);
        return { success: true };
      }
      throw new Error('Accessory not found');
    } catch (error) {
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
  }

  async function updateAccessory(id: number, accessory: UpdateAccessoryInput): Promise<AccessoryOperationResponse> {
    const existingAccessory = accessoryRows.value.find(a => a.id === id);
    if (!existingAccessory) {
      return {
        success: false,
        error: 'Accessory not found'
      };
    }
    
    // Create a deep copy of the existing accessory
    const originalAccessory: AccessoryRow = { ...existingAccessory };
    
    try {
      isLoading.value = true;
      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const index = accessoryRows.value.findIndex(a => a.id === id);
      if (index === -1) {
        throw new Error('Accessory not found');
      }
      
      // Use type guards to validate make and color
      const updatedMake = accessory.make && isValidAccessoryMake(accessory.make) ? accessory.make : originalAccessory.make;
      const updatedColor = accessory.unit_color && isValidAccessoryColor(accessory.unit_color) ? accessory.unit_color : originalAccessory.unit_color;

      // Create updated accessory with all required properties
      const updatedAccessory: AccessoryRow = {
        id,
        name: accessory.name ?? originalAccessory.name,
        make: updatedMake,
        quantity: accessory.quantity ?? originalAccessory.quantity,
        price: accessory.price ?? originalAccessory.price,
        status: originalAccessory.status, // We'll update this based on quantity
        unit_color: updatedColor,
        image: accessory.image ?? originalAccessory.image
      };

      try {
        // Update the status based on the new quantity
        if (typeof accessory.quantity === 'number') {
          updateAccessoryStatus(id, accessory.quantity);
          const updatedExistingAccessory = accessoryRows.value.find(a => a.id === id);
          if (!updatedExistingAccessory) {
            throw new Error('Failed to update accessory status');
          }
          updatedAccessory.status = updatedExistingAccessory.status;
        }

        // Attempt to update the accessory in the store
        accessoryRows.value[index] = updatedAccessory;
      } catch (updateError) {
        // If the update fails, restore the original accessory
        console.error('Error updating accessory in store:', updateError);
        accessoryRows.value[index] = originalAccessory;
        throw new Error('Failed to update accessory data');
      }

      return { success: true };
    } catch (error) {
      console.error('Error in updateAccessory:', error);
      // Ensure the original state is restored in case of any error
      const index = accessoryRows.value.findIndex(a => a.id === id);
      if (index !== -1) {
        accessoryRows.value[index] = originalAccessory;
      }
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred while updating accessory'
      };
    } finally {
      isLoading.value = false;
    }
  }

  return {
    // State
    accessoryRows,
    isLoading,
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