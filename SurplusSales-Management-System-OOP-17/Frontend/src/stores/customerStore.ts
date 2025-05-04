import { defineStore } from 'pinia';
import { ref } from 'vue';
import { uid } from 'quasar';
import type { Customer, Sale, SaleItem } from '../types/models';

// Define an extended Sale type that includes its items
export interface SaleWithItems extends Sale {
  items: SaleItem[];
}

// Define a type for cab purchase details
export interface CabPurchaseDetails {
  cabId: number;
  cabName: string;
  quantity: number;
  unitPrice: number;
  accessories: Array<{
    id: number;
    name: string;
    quantity: number;
    unitPrice: number;
  }>;
}

const initialCustomers: Customer[] = [
  {
    id: 'CUST001',
    fullName: 'Alice Wonderland Store',
    email: 'alice.w.store@example.com',
    phone: '111-222-3333',
    address: '123 Fantasy Lane, Wonderland',
    dateRegistered: new Date().toISOString(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: 'CUST002',
    fullName: 'Bob The Builder Store',
    email: 'bob.b.store@example.com',
    phone: '444-555-6666',
    address: '456 Construction Ave, Builderville',
    dateRegistered: new Date(Date.now() - 86400000 * 5).toISOString(),
    createdAt: new Date(Date.now() - 86400000 * 5).toISOString(),
    updatedAt: new Date(Date.now() - 86400000 * 5).toISOString(),
  },
  {
    id: 'CUST003',
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
  const customers = ref<Customer[]>(initialCustomers);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // --- State for Purchase History ---
  const selectedCustomerHistory = ref<SaleWithItems[]>([]);
  const isLoadingHistory = ref(false);
  const historyError = ref<string | null>(null);
  const customerPurchaseHistory = ref<Record<string, SaleWithItems[]>>({});

  // --- Getters ---
  const getCustomerById = (customerId: string) => {
    return customers.value.find(c => c.id === customerId);
  };

  const validateCustomerId = (customerId: string) => {
    const customer = getCustomerById(customerId);
    return {
      isValid: !!customer,
      customer
    };
  };

  // --- Actions --- 

  // Action to simulate fetching data (can be replaced with API call)
  const fetchCustomers = async () => {
    isLoading.value = true;
    error.value = null;
    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500));
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
      await new Promise(resolve => setTimeout(resolve, 300));
      customers.value.push(newCustomer); 
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
       await new Promise(resolve => setTimeout(resolve, 300));
      const index = customers.value.findIndex(c => c.id === customerId);
      if (index !== -1) {
        const originalCustomer = customers.value[index]; 
        if (originalCustomer) {
          const updatedCustomer: Customer = {
              id: originalCustomer.id, 
              fullName: customerData.fullName,
              email: customerData.email,
              phone: customerData.phone,
              address: customerData.address,
              dateRegistered: originalCustomer.dateRegistered,
              createdAt: originalCustomer.createdAt,
              updatedAt: new Date().toISOString(),
          };
          customers.value[index] = updatedCustomer;
        } else {
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
    console.log(`Fetching history for customer ${customerId}...`);
    try {
      await new Promise(resolve => setTimeout(resolve, 800));
      
      // Get existing history from state or initialize
      let historyData = customerPurchaseHistory.value[customerId] || [];
      
      // If no history exists yet, add dummy data for testing
      if (historyData.length === 0) {
        if (customerId === customers.value[0]?.id) {
          const sale1Id = uid();
          historyData = [
            {
              id: sale1Id,
              customerId: customerId,
              soldBy: 'User1',
              saleDate: new Date(Date.now() - 86400000 * 3).toISOString(),
              totalPrice: 150.75,
              createdAt: new Date(Date.now() - 86400000 * 3).toISOString(),
              updatedAt: new Date(Date.now() - 86400000 * 3).toISOString(),
              multiCabId: '',
              items: [
                { 
                  id: uid(), 
                  saleId: sale1Id, 
                  itemType: 'Accessory', 
                  multiCabId: '', 
                  accessoryId: 'ACC001', 
                  materialId: '', 
                  quantity: 2, 
                  unitPrice: 50.00, 
                  subtotal: 100.00, 
                  createdAt: new Date().toISOString(), 
                  updatedAt: new Date().toISOString(),
                  name: 'Side Mirror'
                },
                { 
                  id: uid(), 
                  saleId: sale1Id, 
                  itemType: 'Material', 
                  multiCabId: '', 
                  accessoryId: '', 
                  materialId: 'MAT005', 
                  quantity: 5, 
                  unitPrice: 10.15, 
                  subtotal: 50.75, 
                  createdAt: new Date().toISOString(), 
                  updatedAt: new Date().toISOString(),
                  name: 'Paint Material'
                },
              ]
            }
          ];
        } else if (customerId === customers.value[1]?.id) {
          const sale2Id = uid();
          historyData = [
            {
              id: sale2Id,
              customerId: customerId,
              soldBy: 'User2',
              saleDate: new Date(Date.now() - 86400000 * 8).toISOString(),
              totalPrice: 25000.00,
              createdAt: new Date(Date.now() - 86400000 * 8).toISOString(),
              updatedAt: new Date(Date.now() - 86400000 * 8).toISOString(),
              multiCabId: '',
              items: [
                { 
                  id: uid(), 
                  saleId: sale2Id, 
                  itemType: 'MultiCab', 
                  multiCabId: 'MC001', 
                  accessoryId: '', 
                  materialId: '', 
                  quantity: 1, 
                  unitPrice: 25000.00, 
                  subtotal: 25000.00, 
                  createdAt: new Date().toISOString(), 
                  updatedAt: new Date().toISOString(),
                  name: 'Suzuki Multicab'
                },
              ]
            }
          ];
        }
        // Store the initial history
        customerPurchaseHistory.value[customerId] = historyData;
      }

      selectedCustomerHistory.value = historyData;
      console.log(`Fetched history for customer ${customerId}:`, historyData);

    } catch (err) {
      historyError.value = `Failed to fetch purchase history for customer ${customerId}.`;
      console.error(err);
    } finally {
      isLoadingHistory.value = false;
    }
  };

  // Action to record a new cab purchase
  const recordCabPurchase = async (customerId: string, purchaseDetails: CabPurchaseDetails) => {
    const validation = validateCustomerId(customerId);
    if (!validation.isValid) {
      throw new Error(`Invalid customer ID: ${customerId}`);
    }

    try {
      // Simulate API delay for saving the purchase
      await new Promise(resolve => setTimeout(resolve, 300));

      const currentDate = new Date().toISOString();
      const saleId = uid();

      // Calculate total price
      const cabTotal = purchaseDetails.quantity * purchaseDetails.unitPrice;
      const accessoriesTotal = purchaseDetails.accessories.reduce(
        (total, acc) => total + (acc.quantity * acc.unitPrice),
        0
      );
      const totalPrice = cabTotal + accessoriesTotal;

      // Create sale record
      const sale: SaleWithItems = {
        id: saleId,
        customerId,
        soldBy: 'Current User', // TODO: Replace with actual user ID
        saleDate: currentDate,
        totalPrice,
        createdAt: currentDate,
        updatedAt: currentDate,
        multiCabId: purchaseDetails.cabId.toString(),
        items: [
          // Add cab as main item
          {
            id: uid(),
            saleId,
            itemType: 'Cab',
            multiCabId: purchaseDetails.cabId.toString(),
            accessoryId: '',
            materialId: '',
            quantity: purchaseDetails.quantity,
            unitPrice: purchaseDetails.unitPrice,
            subtotal: cabTotal,
            createdAt: currentDate,
            updatedAt: currentDate,
            name: purchaseDetails.cabName
          },
          // Add accessories as separate items
          ...purchaseDetails.accessories.map(acc => ({
            id: uid(),
            saleId,
            itemType: 'Accessory',
            multiCabId: '',
            accessoryId: acc.id.toString(),
            materialId: '',
            quantity: acc.quantity,
            unitPrice: acc.unitPrice,
            subtotal: acc.quantity * acc.unitPrice,
            createdAt: currentDate,
            updatedAt: currentDate,
            name: acc.name
          }))
        ]
      };

      // Simulate API call to save the sale
      await new Promise(resolve => setTimeout(resolve, 200));

      // Initialize customer history if it doesn't exist
      if (!customerPurchaseHistory.value[customerId]) {
        customerPurchaseHistory.value[customerId] = [];
      }

      // Add to both the selected history and the stored history
      customerPurchaseHistory.value[customerId] = [sale, ...customerPurchaseHistory.value[customerId]];
      selectedCustomerHistory.value = customerPurchaseHistory.value[customerId];

      console.log(`Recorded purchase for customer ${customerId}:`, sale);
      console.log('Updated purchase history:', customerPurchaseHistory.value[customerId]);

      return {
        success: true,
        sale
      };
    } catch (err) {
      console.error('Error recording purchase:', err);
      throw new Error('Failed to record purchase');
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
    selectedCustomerHistory,
    isLoadingHistory,
    historyError,
    fetchPurchaseHistory,
    recordCabPurchase,
    validateCustomerId,
    getCustomerById
  };
});
