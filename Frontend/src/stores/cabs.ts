import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type {
  CabsRow,
  NewCabInput,
  UpdateCabInput,
  CabOperationResponse,
  CabMake,
  CabColor,
  CabStatus,
  CabMakeInput,
  CabColorInput
} from 'src/types/cabs'
import { cabsService } from 'src/services/cabsService'
import { useSearch } from 'src/utils/useSearch'

export type { CabsRow } from 'src/types/cabs'

export const useCabsStore = defineStore('cabs', () => {
  // State
  const cabRows = ref<CabsRow[]>([])
  const isLoading = ref(false)

  // Use input types that allow empty strings for filters
  const filterMake = ref<CabMakeInput>('')
  const filterColor = ref<CabColorInput>('')
  const filterStatus = ref<CabStatus | ''>('')

  // Available options
  const makes: CabMake[] = ['Mazda', 'Porsche', 'Toyota', 'Nissan', 'Ford']
  const colors: CabColor[] = ['Black', 'White', 'Silver', 'Red', 'Blue']
  const statuses: CabStatus[] = ['In Stock', 'Low Stock', 'Out of Stock']

  // Setup search with the composable
  const search = useSearch({
    debounceTime: 300,
    onSearch: () => {
      // Only trigger API reload if we're doing server-side filtering
      // If using client-side filtering exclusively, you could comment this out
      void initializeCabs()
    }
  })

  // Initialize data
  async function initializeCabs() {
    try {
      isLoading.value = true
      const filters: Record<string, string> = {}

      if (filterMake.value) filters.make = filterMake.value
      if (filterColor.value) filters.unit_color = filterColor.value
      if (filterStatus.value) filters.status = filterStatus.value

      // Only include search if it has a non-empty value
      if (search.searchValue.value.trim()) {
        filters.search = search.searchValue.value.trim()
      }

      const cabs = await cabsService.getCabs(filters)
      cabRows.value = cabs
      return { success: true }
    } catch (error) {
      console.error('Error initializing cabs:', error)
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
  }

  // Computed
  const filteredCabRows = computed(() => {
    return cabRows.value
  })

  // Actions
  async function addCab(cab: NewCabInput): Promise<CabOperationResponse> {
    try {
      isLoading.value = true
      // Validate required fields
      if (!cab.make || !cab.unit_color) {
        throw new Error('Make and color are required');
      }

      const result = await cabsService.addCab(cab)

      if (result.success) {
        // Reload data from API if successful
        await initializeCabs()
      }

      return result
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
  }

  function updateCabStatus(id: number, quantity: number) {
    // This function will be used locally after API operations to update status
    let newStatus: CabStatus;
    if (quantity === 0) {
      newStatus = 'Out of Stock';
    } else if (quantity <= 7) {
      newStatus = 'Low Stock';
    } else { // quantity > 7
      newStatus = 'In Stock';
    }
    return newStatus;
  }

  async function resetFilters() {
    filterMake.value = ''
    filterColor.value = ''
    filterStatus.value = ''
    search.clearSearch()
    // Reload data without filters
    await initializeCabs()
  }

  async function deleteCab(id: number): Promise<CabOperationResponse> {
    try {
      isLoading.value = true
      const result = await cabsService.deleteCab(id)

      if (result.success) {
        // Reload data from API if successful
        await initializeCabs()
      }

      return result
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
  }

  async function updateCab(id: number, cab: UpdateCabInput): Promise<CabOperationResponse> {
    try {
      isLoading.value = true

      // Get the existing cab to merge with updates
      const existingCab = cabRows.value.find(c => c.id === id)
      if (!existingCab) {
        return {
          success: false,
          error: 'Cab not found'
        }
      }

      // Merge existing cab with updates
      const updatedCab: NewCabInput = {
        name: cab.name ?? existingCab.name,
        make: cab.make ?? existingCab.make,
        quantity: cab.quantity ?? existingCab.quantity,
        price: cab.price ?? existingCab.price,
        unit_color: cab.unit_color ?? existingCab.unit_color,
        status: cab.status ?? existingCab.status,
        image: cab.image ?? existingCab.image
      }

      const result = await cabsService.updateCab(id, updatedCab)

      if (result.success) {
        // Reload data from API if successful
        await initializeCabs()
      }

      return result
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
  }

  // Watch for filter changes and reload data
  watch([filterMake, filterColor, filterStatus], () => {
    void initializeCabs()
  })

  return {
    // State
    cabRows,
    isLoading,
    search,
    filterMake,
    filterColor,
    filterStatus,

    // Computed
    filteredCabRows,

    // Constants
    makes,
    colors,
    statuses,

    // Actions
    initializeCabs,
    addCab,
    updateCab,
    deleteCab,
    resetFilters,
    updateCabStatus
  }
}) 
