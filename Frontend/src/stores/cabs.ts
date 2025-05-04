import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

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

  const filterMake = ref('')
  const filterColor = ref('')
  const filterStatus = ref('')

  // Available options
  const makes = ['Mazda', 'Porsche', 'Toyota', 'Nissan', 'Ford']
  const colors = ['Black', 'White', 'Silver', 'Red', 'Blue']
  const statuses = ['In Stock', 'Low Stock', 'Out of Stock', 'Available']

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
  }

  return {
    // State
    cabRows,
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
    addCab,
    updateCabStatus,
    resetFilters,
    deleteCab,
    updateCab
  }
}) 