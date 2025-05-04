import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { 
  MaterialRow, 
  NewMaterialInput, 
  UpdateMaterialInput, 
  MaterialOperationResponse,
  MaterialCategory,
  MaterialSupplier,
  MaterialStatus,
  MaterialCategoryInput,
  MaterialSupplierInput
} from 'src/types/materials'

export type { MaterialRow, NewMaterialInput } from 'src/types/materials'

export const useMaterialsStore = defineStore('materials', () => {
  // State
  const materialRows = ref<MaterialRow[]>([])
  const isLoading = ref(false)

  // Initialize data
  async function initializeMaterials() {
    try {
      isLoading.value = true
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500))
      materialRows.value = [
        {
          name: 'Steel Beams',
          id: 101,
          category: 'Building',
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
      ]
    } catch (error) {
      console.error('Error initializing materials:', error)
    } finally {
      isLoading.value = false
    }
  }

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

  // Use input types that allow empty strings for filters
  const filterCategory = ref<MaterialCategoryInput>('')
  const filterSupplier = ref<MaterialSupplierInput>('')
  const filterStatus = ref<MaterialStatus | ''>('')

  // Available options
  const categories: MaterialCategory[] = ['Lumber', 'Building', 'Electrical', 'Plumbing', 'Hardware']
  const suppliers: MaterialSupplier[] = ['Steel Co.', 'Construction Supplies Inc.', 'Wood Works']
  const statuses: MaterialStatus[] = ['In Stock', 'Low Stock', 'Out of Stock', 'Available']

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

  // Type guard functions
  function isValidMaterialCategory(category: MaterialCategoryInput): category is MaterialCategory {
    return category !== '';
  }

  function isValidMaterialSupplier(supplier: MaterialSupplierInput): supplier is MaterialSupplier {
    return supplier !== '';
  }

  // Actions
  async function addMaterial(material: NewMaterialInput): Promise<MaterialOperationResponse> {
    try {
      isLoading.value = true
      // Validate required fields
      if (!material.category || !material.supplier) {
        throw new Error('Category and supplier are required');
      }

      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const newId = Math.max(...materialRows.value.map(item => item.id)) + 1;
      
      // Create a new material with validated types
      const newMaterial: MaterialRow = {
        id: newId,
        name: material.name,
        category: material.category,
        supplier: material.supplier,
        quantity: material.quantity,
        status: material.status,
        image: material.image
      };

      materialRows.value.push(newMaterial);

      return { success: true, id: newId }
    } catch (error) {
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
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

  async function deleteMaterial(id: number): Promise<MaterialOperationResponse> {
    try {
      isLoading.value = true
      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const index = materialRows.value.findIndex(m => m.id === id);
      if (index !== -1) {
        materialRows.value.splice(index, 1);
        return { success: true };
      }
      throw new Error('Material not found');
    } catch (error) {
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }
    } finally {
      isLoading.value = false
    }
  }

  async function updateMaterial(id: number, material: UpdateMaterialInput): Promise<MaterialOperationResponse> {
    const existingMaterial = materialRows.value.find(m => m.id === id);
    if (!existingMaterial) {
      return {
        success: false,
        error: 'Material not found'
      };
    }
    
    // Create a deep copy of the existing material
    const originalMaterial: MaterialRow = { ...existingMaterial };
    
    try {
      isLoading.value = true;
      // Simulate a brief network delay that would happen in a real API call
      await new Promise(resolve => setTimeout(resolve, 200));

      const index = materialRows.value.findIndex(m => m.id === id);
      
      // Use type guards to validate category and supplier
      const updatedCategory = material.category && isValidMaterialCategory(material.category) ? material.category : originalMaterial.category;
      const updatedSupplier = material.supplier && isValidMaterialSupplier(material.supplier) ? material.supplier : originalMaterial.supplier;

      // Create updated material with all required properties
      const updatedMaterial: MaterialRow = {
        id,
        name: material.name ?? originalMaterial.name,
        category: updatedCategory,
        supplier: updatedSupplier,
        quantity: material.quantity ?? originalMaterial.quantity,
        status: material.status ?? originalMaterial.status,
        image: material.image ?? originalMaterial.image
      };

      try {
        // Attempt to update the material in the store
        materialRows.value[index] = updatedMaterial;
      } catch (updateError) {
        // If the update fails, restore the original material
        console.error('Error updating material in store:', updateError);
        materialRows.value[index] = originalMaterial;
        throw new Error('Failed to update material data');
      }

      return { success: true };
    } catch (error) {
      console.error('Error in updateMaterial:', error);
      // Ensure the original state is restored in case of any error
      const index = materialRows.value.findIndex(m => m.id === id);
      if (index !== -1) {
        materialRows.value[index] = originalMaterial;
      }
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error occurred while updating material'
      };
    } finally {
      isLoading.value = false;
    }
  }

  return {
    // State
    materialRows,
    isLoading,
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
    initializeMaterials,
    addMaterial,
    updateMaterialStatus,
    resetFilters,
    deleteMaterial,
    updateMaterial
  }
})
