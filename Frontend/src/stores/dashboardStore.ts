import { defineStore } from 'pinia';
import { ref, computed, onMounted, watch } from 'vue';
import { useCabsStore } from './cabs';
import { useMaterialsStore } from './materials';
import { useAccessoriesStore } from './accessories';

// Define chart dataset types
interface ChartDataset {
  label: string;
  data: number[];
  borderColor: string;
  tension: number;
  fill: boolean;
}

interface ChartData {
  labels: string[];
  datasets: ChartDataset[];
}

// Helper functions to initialize empty sales data
function initializeMonthlyData(): ChartData {
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  
  return {
    labels: months,
    datasets: [
      {
        label: 'Cabs Sold',
        data: Array(12).fill(0),
        borderColor: '#7B66FF',
        tension: 0.4,
        fill: false
      },
      {
        label: 'Accessories Sold',
        data: Array(12).fill(0),
        borderColor: '#FFB534',
        tension: 0.4,
        fill: false
      },
      {
        label: 'Total Items Sold',
        data: Array(12).fill(0),
        borderColor: '#FF6B6B',
        tension: 0.4,
        fill: false
      }
    ]
  };
}

function initializeWeeklyData(): ChartData {
  return {
    labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
    datasets: [
      {
        label: 'Cabs Sold',
        data: Array(4).fill(0),
        borderColor: '#7B66FF',
        tension: 0.4,
        fill: false
      },
      {
        label: 'Accessories Sold',
        data: Array(4).fill(0),
        borderColor: '#FFB534',
        tension: 0.4,
        fill: false
      },
      {
        label: 'Total Items Sold',
        data: Array(4).fill(0),
        borderColor: '#FF6B6B',
        tension: 0.4,
        fill: false
      }
    ]
  };
}

function initializeYearlyData(): ChartData {
  const currentYear = new Date().getFullYear();
  const years = [currentYear - 4, currentYear - 3, currentYear - 2, currentYear - 1, currentYear].map(String);
  
  return {
    labels: years,
    datasets: [
      {
        label: 'Cabs Sold',
        data: Array(5).fill(0),
        borderColor: '#7B66FF',
        tension: 0.4,
        fill: false
      },
      {
        label: 'Accessories Sold',
        data: Array(5).fill(0),
        borderColor: '#FFB534',
        tension: 0.4,
        fill: false
      },
      {
        label: 'Total Items Sold',
        data: Array(5).fill(0),
        borderColor: '#FF6B6B',
        tension: 0.4,
        fill: false
      }
    ]
  };
}

// Function to load sales data from localStorage or initialize with mock data
function loadSalesData(key: string, initFn: () => ChartData): ChartData {
  // Try to load from localStorage first
  const savedData = localStorage.getItem(key);
  
  if (savedData) {
    try {
      return JSON.parse(savedData) as ChartData;
    } catch (error) {
      console.error(`Error parsing ${key} from localStorage:`, error);
    }
  }
  
  // Return mock data if no valid data in localStorage
  const mockData = initFn();
  localStorage.setItem(key, JSON.stringify(mockData));
  return mockData;
}

export const useDashboardStore = defineStore('dashboard', () => {
  // Connect to inventory stores
  const cabsStore = useCabsStore();
  const materialsStore = useMaterialsStore();
  const accessoriesStore = useAccessoriesStore();

  // Dashboard state
  const isLoading = ref(false);
  const totalInventoryValue = ref(0);
  // Initialize totalSales from localStorage or default to zero
  const localStorageTotalSales = localStorage.getItem('totalSales');
  console.log('Loading totalSales from localStorage:', localStorageTotalSales);

  // Make sure we're parsing it as a valid number or use a default value of 0
  const parsedValue = localStorageTotalSales ? Number(localStorageTotalSales) : 0;
  const totalSales = ref(isNaN(parsedValue) ? 0 : parsedValue);

  console.log('Initialized totalSales with value:', totalSales.value);
  const inventoryTrendData = ref([3200, 3300, 3400, 3500, 3600, 3700, 3800]);
  
  // Initialize sales trend data from localStorage
  const salesTrendData = ref<ChartData>(
    loadSalesData('monthlySalesTrendData', initializeMonthlyData)
  );
  
  const weeklySalesTrendData = ref<ChartData>(
    loadSalesData('weeklySalesTrendData', initializeWeeklyData)
  );
  
  const yearlySalesTrendData = ref<ChartData>(
    loadSalesData('yearlySalesTrendData', initializeYearlyData)
  );
  
  // Watch for changes in sales trend data and save to localStorage
  watch(salesTrendData, (newValue) => {
    localStorage.setItem('monthlySalesTrendData', JSON.stringify(newValue));
  }, { deep: true });
  
  watch(weeklySalesTrendData, (newValue) => {
    localStorage.setItem('weeklySalesTrendData', JSON.stringify(newValue));
  }, { deep: true });
  
  watch(yearlySalesTrendData, (newValue) => {
    localStorage.setItem('yearlySalesTrendData', JSON.stringify(newValue));
  }, { deep: true });
  
  // Watch for changes in totalSales and save to localStorage
  watch(totalSales, (newValue) => {
    localStorage.setItem('totalSales', newValue.toString());
  });
  
  // Computed properties that automatically update when inventory changes
  const cabsCount = computed(() => {
    return cabsStore.cabRows.reduce((total, cab) => total + cab.quantity, 0);
  });
  
  const materialsCount = computed(() => {
    return materialsStore.materialRows.reduce((total, material) => total + material.quantity, 0);
  });
  
  const accessoriesCount = computed(() => {
    return accessoriesStore.accessoryRows.reduce((total, accessory) => total + accessory.quantity, 0);
  });
  
  // Total inventory items across all categories
  const totalInventoryItems = computed(() => {
    return cabsCount.value + materialsCount.value + accessoriesCount.value;
  });

  // Count of low stock items across all categories
  const lowStockItems = computed(() => {
    const lowStockCabs = cabsStore.cabRows.filter(cab => cab.status === 'Low Stock').length;
    const lowStockMaterials = materialsStore.materialRows.filter(material => material.status === 'Low Stock').length;
    const lowStockAccessories = accessoriesStore.accessoryRows.filter(accessory => accessory.status === 'Low Stock').length;
    
    return lowStockCabs + lowStockMaterials + lowStockAccessories;
  });

  // Count of critical items (out of stock)
  const criticalItems = computed(() => {
    const outOfStockCabs = cabsStore.cabRows.filter(cab => cab.status === 'Out of Stock').length;
    const outOfStockMaterials = materialsStore.materialRows.filter(material => material.status === 'Out of Stock').length;
    const outOfStockAccessories = accessoriesStore.accessoryRows.filter(accessory => accessory.status === 'Out of Stock').length;
    
    return outOfStockCabs + outOfStockMaterials + outOfStockAccessories;
  });

  // Dynamic trend data for low stock items
  const lowStockTrendData = computed(() => {
    // Generate historic-like data based on current value
    // Last value is always the current lowStockItems count
    const baseValue = lowStockItems.value;
    const variance = Math.max(5, Math.floor(baseValue * 0.2)); // 20% variance or at least 5
    
    return [
      baseValue + Math.floor(Math.random() * variance) - Math.floor(variance / 2),
      baseValue + Math.floor(Math.random() * variance) - Math.floor(variance / 2),
      baseValue + Math.floor(Math.random() * variance) - Math.floor(variance / 2),
      baseValue + Math.floor(Math.random() * variance) - Math.floor(variance / 2),
      baseValue + Math.floor(Math.random() * variance) - Math.floor(variance / 2),
      baseValue + Math.floor(Math.random() * variance) - Math.floor(variance / 2),
      lowStockItems.value // Current value
    ];
  });

  // Inventory distribution for chart
  const inventoryDistribution = computed(() => ({
    labels: ['Cabs', 'Materials', 'Accessories'],
    datasets: [{
      label: 'Inventory Distribution',
      data: [cabsCount.value, materialsCount.value, accessoriesCount.value],
      backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56']
    }]
  }));
  
  // Reset chart data to mock data values
  function resetChartData() {
    // Create fresh mock data
    salesTrendData.value = initializeMonthlyData();
    weeklySalesTrendData.value = initializeWeeklyData();
    yearlySalesTrendData.value = initializeYearlyData();
    
    // Update localStorage
    localStorage.setItem('monthlySalesTrendData', JSON.stringify(salesTrendData.value));
    localStorage.setItem('weeklySalesTrendData', JSON.stringify(weeklySalesTrendData.value));
    localStorage.setItem('yearlySalesTrendData', JSON.stringify(yearlySalesTrendData.value));
    
    // Return the monthly data for immediate use
    return salesTrendData.value;
  }
  
  // Activities tracking
  interface Activity {
    id: string;
    title: string;
    description: string;
    timestamp: Date;
    icon: string;
    color: string;
  }

  const recentActivities = ref<Activity[]>([
    {
      id: '1',
      title: 'New Car Sold',
      description: 'Multicab X1 sold to John Smith',
      timestamp: new Date(),
      icon: 'directions_car',
      color: 'positive'
    },
    {
      id: '2',
      title: 'Low Stock Alert',
      description: 'Multicab X2 stock is below minimum threshold',
      timestamp: new Date(Date.now() - 1800000), // 30 minutes ago
      icon: 'warning',
      color: 'warning'
    },
    {
      id: '3',
      title: 'Inventory Updated',
      description: 'Added 5 units of Multicab X3',
      timestamp: new Date(Date.now() - 3600000), // 1 hour ago
      icon: 'inventory',
      color: 'info'
    },
    {
      id: '4',
      title: 'Payment Received',
      description: 'Payment received for Invoice #1234',
      timestamp: new Date(Date.now() - 7200000), // 2 hours ago
      icon: 'payments',
      color: 'positive'
    },
    {
      id: '5',
      title: 'Customer Inquiry',
      description: 'New inquiry for Multicab X5 model',
      timestamp: new Date(Date.now() - 10800000), // 3 hours ago
      icon: 'contact_support',
      color: 'primary'
    }
  ]);

  // Methods
  async function refreshDashboardData() {
    isLoading.value = true;
    try {
      // Refresh data from all inventory stores
      await Promise.all([
        cabsStore.initializeCabs(),
        materialsStore.initializeMaterials(),
        accessoriesStore.initializeAccessories()
      ]);
      
      // Update other dashboard metrics
      totalInventoryValue.value = calculateTotalInventoryValue();
      
      return true;
    } catch (error) {
      console.error('Error refreshing dashboard data:', error);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  function calculateTotalInventoryValue(): number {
    const cabsValue = cabsStore.cabRows.reduce((total, cab) => total + (cab.price * cab.quantity), 0);
    // Materials don't have a price property, so we'll exclude them from the calculation
    const accessoriesValue = accessoriesStore.accessoryRows.reduce((total, accessory) => total + (accessory.price * accessory.quantity), 0);
    
    // Since we don't have price for materials, we'll exclude them from the monetary calculation
    // You may want to add a default price per material or modify the MaterialRow type to include price
    return cabsValue + accessoriesValue;
  }

  /**
   * Updates the sales trend data for a specific time period
   * @param {object} saleInfo - Information about the sale
   * @param {number} saleInfo.quantity - Quantity of items sold
   * @param {string} saleInfo.type - Type of item sold ('cab', 'accessory', etc.)
   */
  function updateSalesTrendData({
    quantity = 1,
    type
  }: {
    quantity?: number,
    type: 'cab' | 'accessory' | 'material' | 'other'
  }) {
    try {
      // Force quantity to be a number and at least 1
      const validQuantity = Math.max(1, Number(quantity) || 1);
      
      // Log the incoming parameters
      console.log(`Updating sales trend data: Type: ${type}, Quantity: ${validQuantity} (original: ${quantity})`);
      
      // Get current date information
      const now = new Date();
      const currentMonth = now.getMonth();
      
      // More robust method for calculating week within the month
      // This uses ISO week but then normalizes to 1-4 for our data structure
      const dayOfMonth = now.getDate();
      const totalDaysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate();
      // Calculate week as a position within the month (1-4)
      const currentWeek = Math.min(4, Math.ceil(dayOfMonth / (totalDaysInMonth / 4)));
      
      const currentYear = now.getFullYear();
      
      console.log(`Current date info: Month: ${currentMonth}, Week: ${currentWeek}, Year: ${currentYear}`);
      
      // Ensure salesTrendData is properly loaded from localStorage before making changes
      const monthlyData = loadSalesData('monthlySalesTrendData', initializeMonthlyData);
      const weeklyData = loadSalesData('weeklySalesTrendData', initializeWeeklyData);
      const yearlyData = loadSalesData('yearlySalesTrendData', initializeYearlyData);
      
      console.log('Initial data loaded from localStorage:');
      console.log('Monthly data before update:', JSON.stringify(monthlyData));
      
      // Update the appropriate dataset based on sale type for monthly data
      if (type === 'cab') {
        // Update car sales for current month (track quantity instead of amount)
        if (monthlyData.datasets[0] && monthlyData.datasets[0].data) {
          const oldValue = monthlyData.datasets[0].data[currentMonth] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          monthlyData.datasets[0].data[currentMonth] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} cabs to month ${currentMonth}: ${oldValue} → ${monthlyData.datasets[0].data[currentMonth]}`);
        } else {
          console.error('Monthly cab dataset is invalid:', monthlyData.datasets[0]);
        }
      } else if (type === 'accessory') {
        // Update accessory sales for current month (track quantity instead of amount)
        if (monthlyData.datasets[1] && monthlyData.datasets[1].data) {
          const oldValue = monthlyData.datasets[1].data[currentMonth] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          monthlyData.datasets[1].data[currentMonth] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} accessories to month ${currentMonth}: ${oldValue} → ${monthlyData.datasets[1].data[currentMonth]}`);
        } else {
          console.error('Monthly accessory dataset is invalid:', monthlyData.datasets[1]);
        }
      }
      
      // Update total quantity for current month
      if (monthlyData.datasets[2] && monthlyData.datasets[2].data) {
        const oldValue = monthlyData.datasets[2].data[currentMonth] || 0;
        // Use assignment instead of += to avoid "object possibly undefined" error
        monthlyData.datasets[2].data[currentMonth] = oldValue + validQuantity;
        console.log(`Adding ${validQuantity} total items to month ${currentMonth}: ${oldValue} → ${monthlyData.datasets[2].data[currentMonth]}`);
      } else {
        console.error('Monthly total dataset is invalid:', monthlyData.datasets[2]);
      }
      
      // Update weekly data
      if (currentWeek >= 1 && currentWeek <= 4) {
        const weekIndex = currentWeek - 1;
        
        if (type === 'cab' && weeklyData.datasets[0] && weeklyData.datasets[0].data) {
          const oldValue = weeklyData.datasets[0].data[weekIndex] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          weeklyData.datasets[0].data[weekIndex] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} cabs to week ${currentWeek}: ${oldValue} → ${weeklyData.datasets[0].data[weekIndex]}`);
        } else if (type === 'accessory' && weeklyData.datasets[1] && weeklyData.datasets[1].data) {
          const oldValue = weeklyData.datasets[1].data[weekIndex] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          weeklyData.datasets[1].data[weekIndex] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} accessories to week ${currentWeek}: ${oldValue} → ${weeklyData.datasets[1].data[weekIndex]}`);
        }
        
        // Update total for the week
        if (weeklyData.datasets[2] && weeklyData.datasets[2].data) {
          const oldValue = weeklyData.datasets[2].data[weekIndex] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          weeklyData.datasets[2].data[weekIndex] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} total items to week ${currentWeek}: ${oldValue} → ${weeklyData.datasets[2].data[weekIndex]}`);
        }
      }
      
      // Update yearly data
      const yearLabels = yearlyData.labels || [];
      const currentYearIndex = yearLabels.findIndex((year: string) => year === currentYear.toString());
      
      if (currentYearIndex !== -1) {
        if (type === 'cab' && yearlyData.datasets[0] && yearlyData.datasets[0].data) {
          const oldValue = yearlyData.datasets[0].data[currentYearIndex] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          yearlyData.datasets[0].data[currentYearIndex] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} cabs to year ${currentYear}: ${oldValue} → ${yearlyData.datasets[0].data[currentYearIndex]}`);
        } else if (type === 'accessory' && yearlyData.datasets[1] && yearlyData.datasets[1].data) {
          const oldValue = yearlyData.datasets[1].data[currentYearIndex] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          yearlyData.datasets[1].data[currentYearIndex] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} accessories to year ${currentYear}: ${oldValue} → ${yearlyData.datasets[1].data[currentYearIndex]}`);
        }
        
        // Update total for the year
        if (yearlyData.datasets[2] && yearlyData.datasets[2].data) {
          const oldValue = yearlyData.datasets[2].data[currentYearIndex] || 0;
          // Use assignment instead of += to avoid "object possibly undefined" error
          yearlyData.datasets[2].data[currentYearIndex] = oldValue + validQuantity;
          console.log(`Adding ${validQuantity} total items to year ${currentYear}: ${oldValue} → ${yearlyData.datasets[2].data[currentYearIndex]}`);
        }
      }
      
      // Force the reactive reference to update by assigning new objects
      salesTrendData.value = { ...monthlyData };
      weeklySalesTrendData.value = { ...weeklyData };
      yearlySalesTrendData.value = { ...yearlyData };
      
      // Immediately persist to localStorage
      localStorage.setItem('monthlySalesTrendData', JSON.stringify(monthlyData));
      localStorage.setItem('weeklySalesTrendData', JSON.stringify(weeklyData));
      localStorage.setItem('yearlySalesTrendData', JSON.stringify(yearlyData));
      
      console.log('Monthly data after update:', JSON.stringify(monthlyData));
      console.log('Sales trend data has been updated and saved to localStorage');
      
      return true;
    } catch (error) {
      console.error('Error updating sales trend data:', error);
      return false;
    }
  }

  /**
   * Record a sale in the dashboard
   * @param {object} saleInfo - Information about the sale
   * @param {string} saleInfo.itemName - Name of the item sold
   * @param {number} saleInfo.amount - Sale amount in currency
   * @param {number} saleInfo.quantity - Quantity of items sold
   * @param {string} saleInfo.type - Type of item sold ('cab', 'accessory', etc.)
   * @param {string} [saleInfo.customerName] - Optional customer name
   */
  function recordSale({ itemName, amount, quantity = 1, type, customerName = 'a customer' }: { 
    itemName: string, 
    amount: number,
    quantity?: number,
    type: 'cab' | 'accessory' | 'material' | 'other',
    customerName?: string
  }) {
    try {
      // Ensure amount is a number and not NaN
      const validAmount = Number(amount) || 0;
      // Ensure quantity is a positive number
      const validQuantity = Math.max(1, Number(quantity) || 1);
      
      // Log incoming data
      console.log('recordSale called with:', { itemName, amount, quantity, type, customerName });
      console.log('Validated values:', { validAmount, validQuantity });
      
      // Update total sales amount
      totalSales.value += validAmount;
      
      // Log for debugging
      console.log(`Recording sale: ${itemName}, Amount: ${validAmount}, Quantity: ${validQuantity}, Total Sales now: ${totalSales.value}`);
      
      // Save to localStorage
      localStorage.setItem('totalSales', totalSales.value.toString());
      
      // Update sales trend data with quantity only
      const updateResult = updateSalesTrendData({ quantity: validQuantity, type });
      console.log('Update sales trend data result:', updateResult);
      
      // Add an activity record
      addActivity({
        id: Date.now().toString(),
        title: `New ${type === 'cab' ? 'Car' : type} Sold`,
        description: `${validQuantity} ${itemName} sold to ${customerName}`,
        timestamp: new Date(),
        icon: type === 'cab' ? 'directions_car' : 'shopping_cart',
        color: 'positive'
      });
      
      return true;
    } catch (error) {
      console.error('Error recording sale:', error);
      return false;
    }
  }

  function addActivity(activity: Activity) {
    recentActivities.value = [activity, ...recentActivities.value.slice(0, 4)];
  }

  function formatCurrency(value: number): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      maximumFractionDigits: 0
    }).format(value);
  }

  /**
   * Completely clears all sales data and resets to empty charts.
   * This should only be used when you want to start fresh with no historical data.
   * @param {boolean} confirmed - Whether the operation is confirmed by the user
   * @returns {boolean} - Whether the operation was successful
   */
  function clearAllSalesData(confirmed = false) {
    try {
      // Check if the operation is confirmed
      if (!confirmed) {
        console.warn('Clear all sales data operation was not confirmed. No data was deleted.');
        // Return false to indicate that the operation was not executed
        return false;
      }
      
      // Reset total sales to 0
      totalSales.value = 0;
      localStorage.removeItem('totalSales');
      
      // Reset all chart data to empty (zeros)
      salesTrendData.value = initializeMonthlyData();
      weeklySalesTrendData.value = initializeWeeklyData();
      yearlySalesTrendData.value = initializeYearlyData();
      
      // Clear all sales data from localStorage
      localStorage.removeItem('monthlySalesTrendData');
      localStorage.removeItem('weeklySalesTrendData');
      localStorage.removeItem('yearlySalesTrendData');
      
      // Save the empty data to localStorage
      localStorage.setItem('monthlySalesTrendData', JSON.stringify(salesTrendData.value));
      localStorage.setItem('weeklySalesTrendData', JSON.stringify(weeklySalesTrendData.value));
      localStorage.setItem('yearlySalesTrendData', JSON.stringify(yearlySalesTrendData.value));
      
      console.log('All sales data has been cleared and reset to empty values');
      
      // Add an activity to show the reset
      addActivity({
        id: Date.now().toString(),
        title: 'Sales Data Reset',
        description: 'All historical sales data has been cleared',
        timestamp: new Date(),
        icon: 'restart_alt',
        color: 'info'
      });
      
      return true;
    } catch (error) {
      console.error('Error clearing sales data:', error);
      return false;
    }
  }

  // Initialize dashboard data
  onMounted(() => {
    void refreshDashboardData();
  });

  return {
    // State
    isLoading,
    totalInventoryValue,
    totalSales,
    lowStockItems,
    recentActivities,
    inventoryTrendData,
    salesTrendData,
    weeklySalesTrendData,
    yearlySalesTrendData,
    
    // Computed
    cabsCount,
    materialsCount,
    accessoriesCount,
    totalInventoryItems,
    inventoryDistribution,
    lowStockTrendData,
    criticalItems,
    
    // Methods
    refreshDashboardData,
    calculateTotalInventoryValue,
    recordSale,
    addActivity,
    formatCurrency,
    updateSalesTrendData,
    resetChartData,
    clearAllSalesData
  };
}); 