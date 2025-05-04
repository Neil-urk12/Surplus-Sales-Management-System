import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export interface MaterialRow {
  name: string
  id: number
  category: string
  supplier: string
  quantity: number
  status: string
  image: string
}

export const useMaterialsStore = defineStore('materials', () => {
  // State
  const materialRows = ref<MaterialRow[]>([
    {
      name: 'Steel Beams',
      id: 101,
      category: 'Lumber',
      supplier: 'Steel Co.',
      quantity: 25,
      status: 'In Stock',
      image: 'https://loremflickr.com/600/400/steel',
    },
    {
      name: 'Concrete Mix',
      id: 102,
      category: 'Building',
      supplier: 'Construction Supplies Inc.',
      quantity: 50,
      status: 'Low Stock',
      image: 'https://loremflickr.com/600/400/concrete',
    },
    {
      name: 'Lumber',
      id: 103,
      category: 'Lumber',
      supplier: 'Wood Works',
      quantity: 100,
      status: 'Available',
      image: 'https://loremflickr.com/600/400/lumber',
    },
    {
      name: 'Concrete Mix',
      id: 102,
      category: 'Building',
      supplier: 'Construction Supplies Inc.',
      quantity: 50,
      status: 'Low Stock',
      image: 'https://loremflickr.com/600/400/concrete',
    },
    {
      name: 'Lumber',
      id: 103,
      category: 'Lumber',
      supplier: 'Wood Works',
      quantity: 100,
      status: 'Available',
      image: 'https://loremflickr.com/600/400/lumber',
    },
    {
      name: 'Concrete Mix',
      id: 102,
      category: 'Building',
      supplier: 'Construction Supplies Inc.',
      quantity: 50,
      status: 'Low Stock',
      image: 'https://loremflickr.com/600/400/concrete',
    },
    {
      name: 'Lumber',
      id: 103,
      category: 'Lumber',
      supplier: 'Wood Works',
      quantity: 100,
      status: 'Available',
      image: 'https://loremflickr.com/600/400/lumber',
    },
    {
      name: 'Concrete Mix',
      id: 102,
      category: 'Building',
      supplier: 'Construction Supplies Inc.',
      quantity: 50,
      status: 'Low Stock',
      image: 'https://loremflickr.com/600/400/concrete',
    },
    {
      name: 'Lumber',
      id: 103,
      category: 'Lumber',
      supplier: 'Wood Works',
      quantity: 100,
      status: 'Available',
      image: 'https://loremflickr.com/600/400/lumber',
    },
  ])

  // Search with debounce
  const rawMaterialSearch = ref('')
  const materialSearch = ref('')
  let debounceTimeout: ReturnType<typeof setTimeout> | null = null

  // Debounce function to update materialSearch after typing stops
  function updateMaterialSearch(value: string) {
    if (debounceTimeout) {
      clearTimeout(debounceTimeout)
    }

    debounceTimeout = setTimeout(() => {
      materialSearch.value = value
    }, 300) // 300ms debounce delay
  }

  // Watch for changes in rawMaterialSearch
  watch(rawMaterialSearch, (newValue) => {
    updateMaterialSearch(newValue)
  })

  const filterCategory = ref('')
  const filterSupplier = ref('')
  const filterStatus = ref('')

  // Available options
  const categories = ['Lumber', 'Building', 'Electrical', 'Plumbing', 'Hardware']
  const suppliers = ['Steel Co.', 'Construction Supplies Inc.', 'Wood Works']
  const statuses = ['In Stock', 'Low Stock', 'Out of Stock', 'Available']

  // Computed
  const filteredMaterialRows = computed(() => {
    return materialRows.value.filter(row => {
      const matchesCategory = !filterCategory.value || row.category === filterCategory.value
      const matchesSupplier = !filterSupplier.value || row.supplier === filterSupplier.value
      const matchesStatus = !filterStatus.value || row.status === filterStatus.value
      const matchesSearch = !materialSearch.value ||
        row.name.toLowerCase().includes(materialSearch.value.toLowerCase()) ||
        row.category.toLowerCase().includes(materialSearch.value.toLowerCase()) ||
        row.supplier.toLowerCase().includes(materialSearch.value.toLowerCase())

      return matchesCategory && matchesSupplier && matchesStatus && matchesSearch
    })
  })

  // Actions
  async function addMaterial(material: Omit<MaterialRow, 'id'>) {
    // Simulate a brief network delay that would happen in a real API call
    await new Promise(resolve => setTimeout(resolve, 200));

    const newId = Math.max(...materialRows.value.map(item => item.id)) + 1
    materialRows.value.push({
      ...material,
      id: newId
    })

    return { success: true, id: newId }
  }

  function updateMaterialStatus(id: number, quantity: number) {
    const material = materialRows.value.find(m => m.id === id)
    if (material) {
      if (quantity === 0) {
        material.status = 'Out of Stock'
      } else if (quantity <= 10) {
        material.status = 'Low Stock'
      } else if (quantity <= 50) {
        material.status = 'In Stock'
      } else {
        material.status = 'Available'
      }
    }
  }

  function resetFilters() {
    filterCategory.value = ''
    filterSupplier.value = ''
    filterStatus.value = ''
    rawMaterialSearch.value = ''
    materialSearch.value = ''
  }

  async function deleteMaterial(id: number) {
    // Simulate a brief network delay that would happen in a real API call
    await new Promise(resolve => setTimeout(resolve, 200));

    const index = materialRows.value.findIndex(m => m.id === id);
    if (index !== -1) {
      materialRows.value.splice(index, 1);
      return { success: true };
    }
    throw new Error('Material not found');
  }

  async function updateMaterial(id: number, material: Omit<MaterialRow, 'id'>) {
    // Simulate a brief network delay that would happen in a real API call
    await new Promise(resolve => setTimeout(resolve, 200));

    const index = materialRows.value.findIndex(m => m.id === id);
    if (index !== -1) {
      materialRows.value[index] = {
        ...material,
        id
      };
      return { success: true };
    }
    throw new Error('Material not found');
  }

  return {
    // State
    materialRows,
    rawMaterialSearch,
    materialSearch,
    filterCategory,
    filterSupplier,
    filterStatus,
    // Constants
    categories,
    suppliers,
    statuses,
    // Computed
    filteredMaterialRows,
    // Actions
    addMaterial,
    updateMaterialStatus,
    resetFilters,
    deleteMaterial,
    updateMaterial
  }
})
