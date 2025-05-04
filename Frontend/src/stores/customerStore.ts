import { defineStore } from 'pinia';
import { ref } from 'vue';
import { uid } from 'quasar';
import type { Customer, Sale, SaleItem } from '../types/models';

// Define an extended Sale type that includes its items
export interface SaleWithItems extends Sale {
  items: SaleItem[];
}

// Define the initial dummy data directly in the store for now
// In a real app, this would likely be fetched via an action
const initialCustomers: Customer[] = [
  {
    id: uid(),
    fullName: 'Alice Wonderland Store',
    email: 'alice.w.store@example.com',
    phone: '111-222-3333',
    address: '123 Fantasy Lane, Wonderland',
    dateRegistered: new Date().toISOString(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: uid(),
    fullName: 'Bob The Builder Store',
    email: 'bob.b.store@example.com',
    phone: '444-555-6666',
    address: '456 Construction Ave, Builderville',
    dateRegistered: new Date(Date.now() - 86400000 * 5).toISOString(),
    createdAt: new Date(Date.now() - 86400000 * 5).toISOString(),
    updatedAt: new Date(Date.now() - 86400000 * 5).toISOString(),
  },
  {
    id: uid(),
    fullName: 'Charlie Chaplin Store',
    email: 'charlie.c.store@example.com',
    phone: '777-888-9999',
    address: '789 Silent Film St, Hollywood',
    dateRegistered: new Date(Date.now() - 86400000 * 10).toISOString(),
    createdAt: new Date(Date.now() - 86400000 * 10).toISOString(),
    updatedAt: new Date(Date.now() - 86400000 * 10).toISOString(),
  },
];

export const useCustomerStore = defineStore('customer', () => {
  // --- State --- 
  const customers = ref<Customer[]>(initialCustomers); // Start with dummy data
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // --- State for Purchase History ---
  const selectedCustomerHistory = ref<SaleWithItems[]>([]);
  const isLoadingHistory = ref(false);
  const historyError = ref<string | null>(null);

  // --- Actions --- 

  // Action to simulate fetching data (can be replaced with API call)
  const fetchCustomers = async () => {
    isLoading.value = true;
    error.value = null;
    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500));
      // In real app: 
      // const response = await api.get('/customers');
      // customers.value = response.data;
      // For now, just ensure the initial data is set (it already is)
      console.log('Fetched customers (simulated)');
    } catch (err) {
      error.value = 'Failed to fetch customers.';
      console.error(err);
    } finally {
      isLoading.value = false;
    }
  };

  // Action to add a new customer
  const addCustomer = async (customerData: Omit<Customer, 'id' | 'createdAt' | 'updatedAt' | 'dateRegistered'>) => {
    isLoading.value = true;
    error.value = null;
    try {
      const currentDate = new Date().toISOString();
      const newCustomer: Customer = {
        id: uid(), // Generate ID locally for now
        ...customerData,
        dateRegistered: currentDate,
        createdAt: currentDate,
        updatedAt: currentDate,
      };
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 300));
      // In real app: 
      // const response = await api.post('/customers', newCustomer);
      // customers.value.push(response.data); // Add the customer returned by the API (with backend ID)
      customers.value.push(newCustomer); // Add locally for now
    } catch (err) {
      error.value = 'Failed to add customer.';
      console.error(err);
    } finally {
      isLoading.value = false;
    }
  };

  // Action to update an existing customer
  const updateCustomer = async (customerId: string, customerData: Omit<Customer, 'id' | 'createdAt' | 'updatedAt' | 'dateRegistered'>) => {
    isLoading.value = true;
    error.value = null;
    try {
       // Simulate API delay
       await new Promise(resolve => setTimeout(resolve, 300));
      // In real app: 
      // const response = await api.put(`/customers/${customerId}`, { ...customerData });
      // const updatedCustomerFromApi = response.data;
      // customers.value = customers.value.map(c => c.id === customerId ? updatedCustomerFromApi : c);

      // Update locally for now
      const index = customers.value.findIndex(c => c.id === customerId);
      if (index !== -1) {
        // Access original customer safely within the check
        const originalCustomer = customers.value[index]; 
        
        // Explicitly check if originalCustomer is defined before using it
        if (originalCustomer) {
          // Construct the new object now that originalCustomer is confirmed defined
          const updatedCustomer: Customer = {
              id: originalCustomer.id, // Required
              fullName: customerData.fullName, // Updated
              email: customerData.email, // Updated
              phone: customerData.phone, // Updated
              address: customerData.address, // Updated
              dateRegistered: originalCustomer.dateRegistered, // Required, kept from original
              createdAt: originalCustomer.createdAt, // Required, kept from original
              updatedAt: new Date().toISOString(), // Updated
          };
          customers.value[index] = updatedCustomer;
        } else {
          // Handle unexpected case where customer at index is undefined (log error)
          console.error(`Customer with id ${customerId} found at index ${index}, but the value was undefined.`);
          error.value = 'An internal error occurred while updating the customer.';
        }
      }
    } catch (err) {
      error.value = 'Failed to update customer.';
      console.error(err);
    } finally {
      isLoading.value = false;
    }
  };

  // Action to delete a customer
  const deleteCustomer = (customerId: string) => {
    isLoading.value = true;
    error.value = null;
    try {
        // await api.delete(`/customers/${customerId}`);
        customers.value = customers.value.filter(c => c.id !== customerId);
    } catch (err) {
      error.value = 'Failed to delete customer.';
      console.error(err);
    } finally {
      isLoading.value = false;
    }
  };

  // --- Action for Purchase History ---
  const fetchPurchaseHistory = async (customerId: string) => {
    isLoadingHistory.value = true;
    historyError.value = null;
    selectedCustomerHistory.value = []; // Clear previous history
    console.log(`Fetching history for customer ${customerId}...`);
    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 800));

      // --- Dummy History Data (Replace with API call) ---
      // In a real app: const response = await api.get(`/customers/${customerId}/sales?include=items`);
      // const historyData = response.data;

      let historyData: SaleWithItems[] = [];
      // Example: Return some dummy sales only for specific dummy customers
      if (customerId === customers.value[0]?.id) { // Alice Wonderland Store
        const sale1Id = uid();
        historyData = [
          {
            id: sale1Id,
            customerId: customerId,
            soldBy: 'User1', // Example user ID or name
            saleDate: new Date(Date.now() - 86400000 * 3).toISOString(), // 3 days ago
            totalPrice: 150.75,
            createdAt: new Date(Date.now() - 86400000 * 3).toISOString(),
            updatedAt: new Date(Date.now() - 86400000 * 3).toISOString(),
            items: [
              { id: uid(), saleId: sale1Id, itemType: 'Accessory', multiCabId: '', accessoryId: 'ACC001', materialId: '', quantity: 2, unitPrice: 50.00, subtotal: 100.00, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
              { id: uid(), saleId: sale1Id, itemType: 'Material', multiCabId: '', accessoryId: '', materialId: 'MAT005', quantity: 5, unitPrice: 10.15, subtotal: 50.75, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
            ]
          }
        ];
      } else if (customerId === customers.value[1]?.id) { // Bob The Builder Store
         const sale2Id = uid();
         historyData = [
           {
             id: sale2Id,
             customerId: customerId,
             soldBy: 'User2',
             saleDate: new Date(Date.now() - 86400000 * 8).toISOString(), // 8 days ago
             totalPrice: 25000.00,
             createdAt: new Date(Date.now() - 86400000 * 8).toISOString(),
             updatedAt: new Date(Date.now() - 86400000 * 8).toISOString(),
             items: [
                { id: uid(), saleId: sale2Id, itemType: 'MultiCab', multiCabId: 'MC001', accessoryId: '', materialId: '', quantity: 1, unitPrice: 25000.00, subtotal: 25000.00, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
             ]
           }
         ];
      }
      // --- End Dummy Data ---

      selectedCustomerHistory.value = historyData;
      console.log(`Fetched history for customer ${customerId}:`, historyData);

    } catch (err) {
      historyError.value = `Failed to fetch purchase history for customer ${customerId}.`;
      console.error(err);
    } finally {
      isLoadingHistory.value = false;
    }
  };

  return {
    customers,
    isLoading,
    error,
    fetchCustomers,
    addCustomer,
    updateCustomer,
    deleteCustomer,
    // Expose history state and actions
    selectedCustomerHistory,
    isLoadingHistory,
    historyError,
    fetchPurchaseHistory,
  };
});
