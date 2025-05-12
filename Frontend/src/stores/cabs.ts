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
import { salesService } from 'src/services/salesApi'
import type { CabSalePayload, SalesOperationResponse } from 'src/types/salesTypes'
import { useSearch } from 'src/utils/useSearch'
import { useAuthStore } from './auth'
import { useLogStore } from './logs'
import type { UserSnippet } from 'src/types/logTypes'

const authStore = useAuthStore()
const logStore = useLogStore()

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
      const currentUser = authStore.user;
      if (currentUser) {
        const userSnippet: UserSnippet = {
          id: currentUser.id,
          fullName: currentUser.fullName,
          role: currentUser.role
        };
        await logStore.addLogEntry({
          user: userSnippet,
          actionType: 'Updated',
          details: 'System failed to initialize cabs list',
          status: 'failed',
          isSystemAction: true 
        });
      } else {
        // Log as system action if user is not available
        await logStore.addLogEntry({
          user: { id: 'system', fullName: 'System', role: 'admin'},
          actionType: 'Updated',
          details: 'System failed to initialize cabs list (user not available)',
          status: 'failed',
          isSystemAction: true
        });
      }
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
    // Sort by ID to ensure cabs are displayed with newer items (highest IDs) at the end
    return cabRows.value.slice().sort((a, b) => a.id - b.id)
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
        // Instead of reloading all data, we'll handle the response and update the local state
        // This ensures proper ID handling and positioning of the new cab
        
        // Get the latest data to work with from the server
        // This includes retrieving the newly added cab with its server-assigned ID
        await initializeCabs()

        const currentUser = authStore.user;
        if (currentUser && result.id !== undefined) {
          const userSnippet: UserSnippet = {
            id: currentUser.id,
            fullName: currentUser.fullName,
            role: currentUser.role
          };
          await logStore.addLogEntry({
            user: userSnippet,
            actionType: 'Created',
            details: `Cab ${result.id} created`,
            status: 'successful',
            isSystemAction: false
          });
        } else if (result.id === undefined) {
            console.warn('Cab ID is undefined, cannot log creation.');
        } else {
            console.warn('User not available for logging cab creation.');
        }
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

      if (result.success && cabRows.value) {
        // Instead of reloading all data, just remove the deleted cab from local state
        // This prevents ID reallocation and maintains the integrity of the ID sequence
        const index = cabRows.value.findIndex(c => c.id === id);
        if (index !== -1) {
          cabRows.value.splice(index, 1);
        }
        const currentUser = authStore.user;
        if (currentUser) {
          const userSnippet: UserSnippet = {
            id: currentUser.id,
            fullName: currentUser.fullName,
            role: currentUser.role
          };
          await logStore.addLogEntry({
            user: userSnippet,
            actionType: 'Deleted',
            details: `Cab ${id} deleted`,
            status: 'successful',
            isSystemAction: false
          });
        } else {
          console.warn('User not available for logging cab deletion.');
        }
      }

      return result
    } catch (error) {
      const currentUser = authStore.user;
      if (currentUser) {
        const userSnippet: UserSnippet = {
          id: currentUser.id,
          fullName: currentUser.fullName,
          role: currentUser.role
        };
        await logStore.addLogEntry({
          user: userSnippet,
          actionType: 'Deleted',
          details: `Failed to delete cab ${id}`,
          status: 'failed',
          isSystemAction: false
        });
      } else {
        console.warn('User not available for logging cab deletion failure.');
      }
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

      // Validate make and unit_color before proceeding
      if (updatedCab.make && !makes.includes(updatedCab.make)) {
        return {
          success: false,
          error: `Invalid make value: ${updatedCab.make}`
        }
      }

      if (updatedCab.unit_color && !colors.includes(updatedCab.unit_color)) {
        return {
          success: false,
          error: `Invalid color value: ${updatedCab.unit_color}`
        }
      }

      const result = await cabsService.updateCab(id, updatedCab)

      if (result.success) {
        // Instead of reloading all data, just update the specific cab in the local state
        // This prevents disrupting the ID order and maintains the integrity of the cab list
        const index = cabRows.value.findIndex(c => c.id === id);
        if (index !== -1) {
          // Create a properly typed object by ensuring all properties match the CabsRow type
          const typedUpdatedCab: CabsRow = {
            id: existingCab.id,
            name: updatedCab.name,
            // Cast to proper types (now safe after validation)
            make: updatedCab.make as CabMake,
            unit_color: updatedCab.unit_color as CabColor,
            quantity: updatedCab.quantity,
            price: updatedCab.price,
            status: updatedCab.status,
            image: updatedCab.image
          };
          
          // Update the cab in the array with properly typed data
          cabRows.value[index] = typedUpdatedCab;
          const currentUser = authStore.user;
          if (currentUser) {
            const userSnippet: UserSnippet = {
              id: currentUser.id,
              fullName: currentUser.fullName,
              role: currentUser.role
            };
            await logStore.addLogEntry({
              user: userSnippet,
              actionType: 'Updated',
              details: `Cab ${id} updated`,
              status: 'successful',
              isSystemAction: false
            });
          } else {
            console.warn('User not available for logging cab update.');
          }
        }
      }

      return result
    } catch (error) {
      const currentUser = authStore.user;
      if (currentUser) {
        const userSnippet: UserSnippet = {
          id: currentUser.id,
          fullName: currentUser.fullName,
          role: currentUser.role
        };
        await logStore.addLogEntry({
          user: userSnippet,
          actionType: 'Updated',
          details: `Failed to update cab ${id}`,
          status: 'failed',
          isSystemAction: false
        });
      } else {
        console.warn('User not available for logging cab update failure.');
      }
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

  /**
   * Sell a cab with optional accessories
   * @param cabId - ID of the cab being sold
   * @param payload - Sale details including customer, quantity, and accessories
   * @returns Promise with the sale operation response
   */
  async function sellCab(cabId: number, payload: CabSalePayload): Promise<SalesOperationResponse> {
    try {
      isLoading.value = true;
      
      // Get the existing cab to verify stock
      const existingCab = cabRows.value.find(c => c.id === cabId);
      if (!existingCab) {
        return {
          success: false,
          error: 'Cab not found'
        };
      }
      
      // Verify stock availability
      if (existingCab.quantity < payload.quantity) {
        return {
          success: false,
          error: `Not enough stock. Available: ${existingCab.quantity}, Requested: ${payload.quantity}`
        };
      }
      
      // Process the sale through the API
      const result = await salesService.sellCab(cabId, payload);
      
      if (result.success) {
        // Update local cab inventory after successful sale
        const newQuantity = existingCab.quantity - payload.quantity;
        const newStatus = updateCabStatus(cabId, newQuantity);
        
        // Update the cab in the store
        const updateResult = await updateCab(cabId, {
          quantity: newQuantity,
          status: newStatus
        });
        
        if (!updateResult.success) {
          console.error('Failed to update cab inventory after sale:', updateResult.error);
          // Even if the local update fails, the sale was processed successfully
        }
        const currentUser = authStore.user;
        if (currentUser) {
          const userSnippet: UserSnippet = {
            id: currentUser.id,
            fullName: currentUser.fullName,
            role: currentUser.role
          };
          await logStore.addLogEntry({
            user: userSnippet,
            actionType: 'Updated',
            details: `Cab ${cabId} sold. Quantity: ${payload.quantity}. Customer ID: ${payload.customerId}`,
            status: 'successful',
            isSystemAction: false
          });
        } else {
            console.warn('User not available for logging cab sale.');
        }
      }
      
      return result;
    } catch (error) {
      const currentUser = authStore.user;
      if (currentUser) {
        const userSnippet: UserSnippet = {
          id: currentUser.id,
          fullName: currentUser.fullName,
          role: currentUser.role
        };
        await logStore.addLogEntry({
          user: userSnippet,
          actionType: 'Updated',
          details: `Failed to sell cab ${cabId}. Error: ${error instanceof Error ? error.message : 'Unknown error'}`,
          status: 'failed',
          isSystemAction: false
        });
      } else {
        console.warn('User not available for logging cab sale failure.');
      }
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred during sale'
      };
    } finally {
      isLoading.value = false;
    }
  }

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
    updateCabStatus,
    sellCab
  }
}) 
