import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
<<<<<<< HEAD
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

export type { CabsRow } from 'src/types/cabs'

export const useCabsStore = defineStore('cabs', () => {
  // State
  const cabRows = ref<CabsRow[]>([])
  const isLoading = ref(false)
  
  // Initialize data
  async function initializeCabs() {
    try {
      isLoading.value = true
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500))
      cabRows.value = [
        {
          name: 'RX‑7',
          id: 1,
          make: 'Mazda',
          quantity: 4,
          price: 7_000_000,
          status: 'In Stock',
          unit_color: 'Black',
          image: 'https://loremflickr.com/600/400/mazda',
        },
        {
          name: '911 GT3',
          id: 2,
          make: 'Porsche',
          quantity: 2,
          price: 10_000_000,
          status: 'In Stock',
          unit_color: 'Black',
          image: 'https://loremflickr.com/600/400/porsche',
        },
        {
          name: '911 GT3',
          id: 3,
          make: 'Porsche',
          quantity: 2,
          price: 10_000_000,
          status: 'Available',
          unit_color: 'Black',
          image: 'https://loremflickr.com/600/400/porsche',
        },
        {
          name: 'Corolla',
          id: 4,
          make: 'Toyota',
          quantity: 2,
          price: 10_000_000,
          status: 'In Stock',
          unit_color: 'Black',
          image: 'https://loremflickr.com/600/400/toyota',
        },
        {
          name: 'Navara',
          id: 5,
          make: 'Nissan',
          quantity: 2,
          price: 10_000_000,
          status: 'In Stock',
          unit_color: 'Black',
          image: 'src/assets/images/Cars/navara.avif',
        },
        {
          name: 'Vios',
          id: 6,
          make: 'Toyota',
          quantity: 2,
          price: 10_000_000,
          status: 'In Stock',
          unit_color: 'Black',
          image: 'https://loremflickr.com/600/400/toyota',
        },
        {
          name: 'Ranger',
          id: 7,
          make: 'Ford',
          quantity: 2,
          price: 10_000_000,
          status: 'In Stock',
          unit_color: 'Black',
          image: 'https://loremflickr.com/600/400/ford',
        },
      ]
    } catch (error) {
      console.error('Error initializing cabs:', error)
    } finally {
      isLoading.value = false
    }
  }
=======

export interface CabsRow {
  name: string
  id: number
  make: string
  quantity: number
  price: number
  status: string
  unit_color: string
  image: string
}

export const useCabsStore = defineStore('cabs', () => {
  // State
  const cabRows = ref<CabsRow[]>([
    {
      name: 'RX‑7',
      id: 1,
      make: 'Mazda',
      quantity: 4,
      price: 7_000_000,
      status: 'In Stock',
      unit_color: 'Black',
      image: 'https://loremflickr.com/600/400/mazda',
    },
    {
      name: '911 GT3',
      id: 2,
      make: 'Porsche',
      quantity: 2,
      price: 10_000_000,
      status: 'In Stock',
      unit_color: 'Black',
      image: 'https://loremflickr.com/600/400/porsche',
    },
    {
      name: '911 GT3',
      id: 3,
      make: 'Porsche',
      quantity: 2,
      price: 10_000_000,
      status: 'Available',
      unit_color: 'Black',
      image: 'https://loremflickr.com/600/400/porsche',
    },
    {
      name: 'Corolla',
      id: 4,
      make: 'Toyota',
      quantity: 2,
      price: 10_000_000,
      status: 'In Stock',
      unit_color: 'Black',
      image: 'https://loremflickr.com/600/400/toyota',
    },
    {
      name: 'Navara',
      id: 5,
      make: 'Nissan',
      quantity: 2,
      price: 10_000_000,
      status: 'In Stock',
      unit_color: 'Black',
      image: 'src/assets/images/Cars/navara.avif',
    },
    {
      name: 'Vios',
      id: 6,
      make: 'Toyota',
      quantity: 2,
      price: 10_000_000,
      status: 'In Stock',
      unit_color: 'Black',
      image: 'https://loremflickr.com/600/400/toyota',
    },
    {
      name: 'Ranger',
      id: 7,
      make: 'Ford',
      quantity: 2,
      price: 10_000_000,
      status: 'In Stock',
      unit_color: 'Black',
      image: 'https://loremflickr.com/600/400/ford',
    },
  ])
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)

  // Search with debounce
  const rawCabSearch = ref('')
  const cabSearch = ref('')
  let debounceTimeout: ReturnType<typeof setTimeout> | null = null

  // Debounce function to update cabSearch after typing stops
  function updateCabSearch(value: string) {
    if (debounceTimeout) {
      clearTimeout(debounceTimeout)
    }

    debounceTimeout = setTimeout(() => {
      cabSearch.value = value
    }, 300) // 300ms debounce delay
  }

  // Watch for changes in rawCabSearch
  watch(rawCabSearch, (newValue) => {
    updateCabSearch(newValue)
  })

<<<<<<< HEAD
  // Use input types that allow empty strings for filters
  const filterMake = ref<CabMakeInput>('')
  const filterColor = ref<CabColorInput>('')
  const filterStatus = ref<CabStatus | ''>('')

  // Available options
  const makes: CabMake[] = ['Mazda', 'Porsche', 'Toyota', 'Nissan', 'Ford']
  const colors: CabColor[] = ['Black', 'White', 'Silver', 'Red', 'Blue']
  const statuses: CabStatus[] = ['In Stock', 'Low Stock', 'Out of Stock', 'Available']
=======
  const filterMake = ref('')
  const filterColor = ref('')
  const filterStatus = ref('')

  // Available options
  const makes = ['Mazda', 'Porsche', 'Toyota', 'Nissan', 'Ford']
  const colors = ['Black', 'White', 'Silver', 'Red', 'Blue']
  const statuses = ['In Stock', 'Low Stock', 'Out of Stock', 'Available']
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)

  // Computed
  const filteredCabRows = computed(() => {
    return cabRows.value.filter(row => {
      const matchesMake = !filterMake.value || row.make === filterMake.value
      const matchesColor = !filterColor.value || row.unit_color === filterColor.value
      const matchesStatus = !filterStatus.value || row.status === filterStatus.value
      const matchesSearch = !cabSearch.value ||
        row.name.toLowerCase().includes(cabSearch.value.toLowerCase()) ||
        row.make.toLowerCase().includes(cabSearch.value.toLowerCase())

      return matchesMake && matchesColor && matchesStatus && matchesSearch
    })
  })

<<<<<<< HEAD
  // Type guard functions
  function isValidCabMake(make: CabMakeInput): make is CabMake {
    return make !== '';
  }

  function isValidCabColor(color: CabColorInput): color is CabColor {
    return color !== '';
  }

  // Actions
  async function addCab(cab: NewCabInput): Promise<CabOperationResponse> {
    try {
      isLoading.value = true
      // Validate required fields
      if (!cab.make || !cab.unit_color) {
        throw new Error('Make and color are required');
      }

      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const newId = Math.max(...cabRows.value.map(item => item.id)) + 1;
      
      // Create a new cab with validated types
      const newCab: CabsRow = {
        id: newId,
        name: cab.name,
        make: cab.make,
        quantity: cab.quantity,
        price: cab.price,
        status: cab.status,
        unit_color: cab.unit_color,
        image: cab.image
      };

      cabRows.value.push(newCab);

      return { success: true, id: newId }
    } catch (error) {
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
=======
  // Actions
  async function addCab(cab: Omit<CabsRow, 'id'>) {
    // Simulate a brief network delay that would happen in a real API call
    await new Promise(resolve => setTimeout(resolve, 200));

    const newId = Math.max(...cabRows.value.map(item => item.id)) + 1
    cabRows.value.push({
      ...cab,
      id: newId
    })

    return { success: true, id: newId }
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
  }

  function updateCabStatus(id: number, quantity: number) {
    const cab = cabRows.value.find(c => c.id === id)
    if (cab) {
      if (quantity === 0) {
        cab.status = 'Out of Stock'
      } else if (quantity <= 2) {
        cab.status = 'Low Stock'
      } else if (quantity <= 5) {
        cab.status = 'In Stock'
      } else {
        cab.status = 'Available'
      }
    }
  }

  function resetFilters() {
    filterMake.value = ''
    filterColor.value = ''
    filterStatus.value = ''
    rawCabSearch.value = ''
    cabSearch.value = ''
  }

<<<<<<< HEAD
  async function deleteCab(id: number): Promise<CabOperationResponse> {
    try {
      isLoading.value = true
      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const index = cabRows.value.findIndex(c => c.id === id);
      if (index !== -1) {
        cabRows.value.splice(index, 1);
        return { success: true };
      }
      throw new Error('Cab not found');
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
    const existingCab = cabRows.value.find(c => c.id === id);
    if (!existingCab) {
      return {
        success: false,
        error: 'Cab not found'
      };
    }
    
    // Create a deep copy of the existing cab
    const originalCab: CabsRow = { ...existingCab };
    
    try {
      isLoading.value = true;
      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const index = cabRows.value.findIndex(c => c.id === id);
      
      // Use type guards to validate make and color
      const updatedMake = cab.make && isValidCabMake(cab.make) ? cab.make : originalCab.make;
      const updatedColor = cab.unit_color && isValidCabColor(cab.unit_color) ? cab.unit_color : originalCab.unit_color;

      // Create updated cab with all required properties
      const updatedCab: CabsRow = {
        id,
        name: cab.name ?? originalCab.name,
        make: updatedMake,
        quantity: cab.quantity ?? originalCab.quantity,
        price: cab.price ?? originalCab.price,
        status: cab.status ?? originalCab.status,
        unit_color: updatedColor,
        image: cab.image ?? originalCab.image
      };

      try {
        // Attempt to update the cab in the store
        cabRows.value[index] = updatedCab;
      } catch (updateError) {
        // If the update fails, restore the original cab
        console.error('Error updating cab in store:', updateError);
        cabRows.value[index] = originalCab;
        throw new Error('Failed to update cab data');
      }

      return { success: true };
    } catch (error) {
      console.error('Error in updateCab:', error);
      // Ensure the original state is restored in case of any error
      const index = cabRows.value.findIndex(c => c.id === id);
      if (index !== -1) {
        cabRows.value[index] = originalCab;
      }
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred while updating cab'
      };
    } finally {
      isLoading.value = false;
    }
=======
  async function deleteCab(id: number) {
    // Simulate a brief network delay that would happen in a real API call
    await new Promise(resolve => setTimeout(resolve, 200));

    const index = cabRows.value.findIndex(c => c.id === id);
    if (index !== -1) {
      cabRows.value.splice(index, 1);
      return { success: true };
    }
    throw new Error('Cab not found');
  }

  async function updateCab(id: number, cab: Omit<CabsRow, 'id'>) {
    // Simulate a brief network delay that would happen in a real API call
    await new Promise(resolve => setTimeout(resolve, 200));

    const index = cabRows.value.findIndex(c => c.id === id);
    if (index !== -1) {
      cabRows.value[index] = {
        ...cab,
        id
      };
      return { success: true };
    }
    throw new Error('Cab not found');
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
  }

  return {
    // State
    cabRows,
<<<<<<< HEAD
    isLoading,
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
    rawCabSearch,
    cabSearch,
    filterMake,
    filterColor,
    filterStatus,
    // Constants
    makes,
    colors,
    statuses,
    // Computed
    filteredCabRows,
    // Actions
<<<<<<< HEAD
    initializeCabs,
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
    addCab,
    updateCabStatus,
    resetFilters,
    deleteCab,
    updateCab
  }
}) 