import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Customer } from '../types/customerTypes.ts';
import type { Sale, ExtendedSaleItem } from '../types/models';
import { customerApi } from '../services/customerApi';
import type { NewCustomerInput, UpdateCustomerInput } from '../services/customerApi.ts';

// Define an extended Sale type that includes its items
export interface SaleWithItems extends Sale {
  items: ExtendedSaleItem[];
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

export const useCustomerStore = defineStore('customer', () => {
  // --- State ---
  const customers = ref<Customer[]>([]);
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
      const fetchedCustomers = await customerApi.getAllCustomers();
      customers.value = fetchedCustomers;
      console.log('Fetched customers from API', customers.value);
    } catch (err) {
      error.value = 'Failed to fetch customers.';
      console.error('Fetch customers error:', err);
      customers.value = []; // Clear customers on error or set to empty if API returns nothing suitable
    } finally {
      isLoading.value = false;
    }
  };

  // Action to add a new customer
  const addCustomer = async (customerData: NewCustomerInput) => {
    isLoading.value = true;
    error.value = null;
    try {
      const createdCustomer = await customerApi.addCustomer(customerData);
      customers.value.push(createdCustomer);
      console.log('Added customer via API');
    } catch (err) {
      error.value = 'Failed to add customer.';
      console.error('Add customer error:', err);
    } finally {
      isLoading.value = false;
    }
  };

  // Action to update an existing customer
  const updateCustomer = async (customerId: string, customerData: UpdateCustomerInput) => {
    isLoading.value = true;
    error.value = null;
    try {
      const updatedCustomerData = await customerApi.updateCustomer(customerId, customerData);
      const index = customers.value.findIndex(c => c.id === customerId);
      if (index !== -1) {
        customers.value[index] = updatedCustomerData;
        console.log('Updated customer via API');
      } else {
        error.value = `Customer with ID ${customerId} not found in local store.`;
        console.error(error.value);
      }
    } catch (err) {
      error.value = 'Failed to update customer.';
      console.error('Update customer error:', err);
    } finally {
      isLoading.value = false;
    }
  };

  // Action to delete a customer
  const deleteCustomer = async (customerId: string) => {
    isLoading.value = true;
    error.value = null;
    try {
      const result = await customerApi.deleteCustomer(customerId);
      if (result.success) {
        customers.value = customers.value.filter(c => c.id !== customerId);
        console.log('Deleted customer via API');
      } else {
        error.value = result.message || 'Failed to delete customer from API.';
        console.error(error.value);
      }
    } catch (err) {
      error.value = 'Failed to delete customer.';
      console.error('Delete customer error:', err);
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
      await new Promise(resolve => setTimeout(resolve, 800)); // Simulate API delay

      // TODO: Replace with actual API call to fetch purchase history for the customerId
      // Example: const historyData = await customerApi.getPurchaseHistory(customerId);

      // For now, clear existing history or set to empty. 
      // When API is integrated, this will be populated by API response.
      selectedCustomerHistory.value = [];
      customerPurchaseHistory.value[customerId] = [];

      console.log(`Purchase history for customer ${customerId} would be fetched here.`);

    } catch (err) {
      historyError.value = `Failed to fetch purchase history for customer ${customerId}.`;
      console.error(err);
      selectedCustomerHistory.value = []; // Clear history on error
      if (customerPurchaseHistory.value[customerId]) {
        customerPurchaseHistory.value[customerId] = [];
      }
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
      await new Promise(resolve => setTimeout(resolve, 300));

      const currentDate = new Date().toISOString();
      // const saleId = uid(); // Backend should generate Sale ID

      // Calculate total price
      const cabTotal = purchaseDetails.quantity * purchaseDetails.unitPrice;
      const accessoriesTotal = purchaseDetails.accessories.reduce(
        (total, acc) => total + (acc.quantity * acc.unitPrice),
        0
      );
      const totalPrice = cabTotal + accessoriesTotal;

      // Create sale record - Note: ID will be assigned by the backend
      const saleToCreate = {
        customerId,
        soldBy: 'Current User', // TODO: Replace with actual user ID
        saleDate: currentDate,
        totalPrice,
        multiCabId: purchaseDetails.cabId.toString(),
        items: [
          // Add cab as main item - Note: Item ID will be assigned by the backend
          {
            // id: uid(), // Backend should generate Item ID
            itemType: 'Cab',
            multiCabId: purchaseDetails.cabId.toString(),
            accessoryId: '',
            materialId: '',
            quantity: purchaseDetails.quantity,
            unitPrice: purchaseDetails.unitPrice,
            subtotal: cabTotal,
            name: purchaseDetails.cabName
          },
          // Add accessories as separate items - Note: Item ID will be assigned by the backend
          ...purchaseDetails.accessories.map(acc => ({
            // id: uid(), // Backend should generate Item ID
            itemType: 'Accessory',
            multiCabId: '',
            accessoryId: acc.id.toString(),
            materialId: '',
            quantity: acc.quantity,
            unitPrice: acc.unitPrice,
            subtotal: acc.quantity * acc.unitPrice,
            name: acc.name
          }))
        ]
      };

      // Simulate API call to save the sale
      // TODO: Replace with actual API call: const savedSale = await customerApi.createSale(saleToCreate);
      console.log('Simulating API call to save sale:', saleToCreate);
      await new Promise(resolve => setTimeout(resolve, 200));

      // Mocking a backend response with IDs for now
      const mockSavedSale: SaleWithItems = {
        ...saleToCreate,
        id: `sale_${Date.now()}`,
        createdAt: currentDate,
        updatedAt: currentDate,
        items: saleToCreate.items.map((item, index) => ({
          ...item,
          id: `item_${Date.now()}_${index}`,
          saleId: `sale_${Date.now()}`, // This would be the actual saleId from backend
          createdAt: currentDate,
          updatedAt: currentDate,
        })) as ExtendedSaleItem[] // Cast needed because original items don't have id, saleId, etc.
      };

      // Initialize customer history if it doesn't exist
      if (!customerPurchaseHistory.value[customerId]) {
        customerPurchaseHistory.value[customerId] = [];
      }

      // Add to both the selected history and the stored history
      customerPurchaseHistory.value[customerId] = [mockSavedSale, ...customerPurchaseHistory.value[customerId]];
      selectedCustomerHistory.value = customerPurchaseHistory.value[customerId];

      console.log(`Recorded purchase for customer ${customerId}:`, mockSavedSale);
      console.log('Updated purchase history:', customerPurchaseHistory.value[customerId]);

      return {
        success: true,
        sale: mockSavedSale
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
