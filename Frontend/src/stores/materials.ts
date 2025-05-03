import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

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
  ])

  const materialSearch = ref('')
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
  function addMaterial(material: Omit<MaterialRow, 'id'>) {
    const newId = Math.max(...materialRows.value.map(item => item.id)) + 1
    materialRows.value.push({
      ...material,
      id: newId
    })
  }

  function updateMaterialStatus(id: number, quantity: number) {
    const material = materialRows.value.find(m => m.id === id)
    if (material) {
      if (quantity === 0) {
        material.status = 'Out of Stock'
      } else if (quantity < 10) {
        material.status = 'Low Stock'
      } else {
        material.status = 'In Stock'
      }
    }
  }

  function resetFilters() {
    filterCategory.value = ''
    filterSupplier.value = ''
    filterStatus.value = ''
    materialSearch.value = ''
  }

  return {
    // State
    materialRows,
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
    resetFilters
  }
}) 