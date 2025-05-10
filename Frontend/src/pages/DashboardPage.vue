<script setup lang="ts">
import { ref, onMounted, computed, defineAsyncComponent, watchEffect } from 'vue';
import { useQuasar } from 'quasar';
import { useDashboardStore } from '../stores/dashboardStore';

const MetricsCard = defineAsyncComponent(() => import('../components/dashboard/MetricsCard.vue'));
const ChartSection = defineAsyncComponent(() => import('../components/dashboard/ChartSection.vue'));
const RecentActivitiesSection = defineAsyncComponent(() => import('../components/dashboard/RecentActivitiesSection.vue'));

const $q = useQuasar();
const dashboardStore = useDashboardStore();

// Local state for UI controls
const timePeriod = ref<'weekly' | 'monthly' | 'yearly'>('monthly');
const chartType = ref<'line' | 'bar' | 'area'>('line');
// Add last refresh timestamp to prevent excessive refreshes
const lastRefreshTimestamp = ref(Date.now());
// Add last clear timestamp for rate limiting
const lastClearTimestamp = ref(0);
// Add confirmation text for data clearing
const confirmationText = ref('');
const requiredConfirmation = 'DELETE';
// Add a flag to control debug logging in development mode only
const isDev = process.env.NODE_ENV === 'development';

// Define proper interfaces for chart data
interface ChartDataset {
  label: string;
  data: number[];
  borderColor?: string;
  backgroundColor?: string[];
  tension?: number;
  fill?: boolean;
}

interface ChartData {
  labels: string[];
  datasets: ChartDataset[];
}

// Helper function to validate ChartData structure
function isValidChartData(data: unknown): boolean {
  if (!data || typeof data !== 'object') return false;
  
  const chartData = data as Partial<ChartData>;
  
  return !!chartData && 
         Array.isArray(chartData.labels) &&
         Array.isArray(chartData.datasets) &&
         chartData.datasets.length > 0 &&
         chartData.datasets.every((dataset: Partial<ChartDataset>) => 
           typeof dataset.label === 'string' && 
           Array.isArray(dataset.data)
         );
}

// Add this computed property as readonly
const currentSalesTrendData = computed(() => {
  if (isDev) console.log('Computing current sales trend data for period:', timePeriod.value);
  let data;
  switch (timePeriod.value) {
    case 'weekly':
      data = dashboardStore.weeklySalesTrendData;
      break;
    case 'yearly':
      data = dashboardStore.yearlySalesTrendData;
      break;
    case 'monthly':
    default:
      data = dashboardStore.salesTrendData;
      break;
  }
  if (isDev) console.log('Returning data:', data);
  return data;
});

onMounted(() => {
  // Refresh general dashboard data (inventories, etc.)
  void dashboardStore.refreshDashboardData();
  
  if (isDev) {
    console.log('Total sales value on mount:', dashboardStore.totalSales);
    console.log('Monthly sales trend data on mount:', dashboardStore.salesTrendData);
  }
});

// Replace multiple watches with a single watchEffect
watchEffect(() => {
  if (isDev) {
    console.log('Time period:', timePeriod.value);
    console.log('Total sales:', dashboardStore.totalSales);
    console.log('Sales trend data:', dashboardStore.salesTrendData);
  }
});

// Improved page activation handler with throttling
function activated() {
  const now = Date.now();
  // Only refresh if at least 2 seconds have passed since last refresh
  if (now - lastRefreshTimestamp.value < 2000) {
    if (isDev) console.log('Skipping refresh, too soon since last refresh');
    return;
  }
  
  if (isDev) console.log('Dashboard page activated - refreshing data');
  lastRefreshTimestamp.value = now;
  
  try {
    // Force refresh dashboard data first
    void dashboardStore.refreshDashboardData();
    
    // Force refresh from localStorage with validation
    const monthlySalesData = JSON.parse(localStorage.getItem('monthlySalesTrendData') || '{}');
    const weeklySalesData = JSON.parse(localStorage.getItem('weeklySalesTrendData') || '{}');
    const yearlySalesData = JSON.parse(localStorage.getItem('yearlySalesTrendData') || '{}');
    
    // Check if data is valid before assigning
    if (isValidChartData(monthlySalesData)) {
      dashboardStore.$patch({ salesTrendData: monthlySalesData });
    } else if (isDev) {
      console.warn('Invalid monthly sales data in localStorage, reinitializing with defaults');
      const defaultData = dashboardStore.resetChartData();
      dashboardStore.$patch({ salesTrendData: defaultData });
    }
    
    if (isValidChartData(weeklySalesData)) {
      dashboardStore.$patch({ weeklySalesTrendData: weeklySalesData });
    } else if (isDev) {
      console.warn('Invalid weekly sales data in localStorage, reinitializing with defaults');
      const defaultData = dashboardStore.resetChartData();
      dashboardStore.$patch({ weeklySalesTrendData: defaultData });
    }
    
    if (isValidChartData(yearlySalesData)) {
      dashboardStore.$patch({ yearlySalesTrendData: yearlySalesData });
    } else if (isDev) {
      console.warn('Invalid yearly sales data in localStorage, reinitializing with defaults');
      const defaultData = dashboardStore.resetChartData();
      dashboardStore.$patch({ yearlySalesTrendData: defaultData });
    }

    // Use a more reliable approach to force chart re-render
    // First change the period to something different
    const currentPeriod = timePeriod.value;
    timePeriod.value = currentPeriod === 'monthly' ? 'weekly' : 'monthly';
    
    // Wait for the next animation frame to ensure the DOM has updated
    requestAnimationFrame(() => {
      // Then change it back to original value
      timePeriod.value = currentPeriod;
      
      // And force chart type to refresh as well
      const currentType = chartType.value;
      chartType.value = currentType === 'line' ? 'bar' : 'line';
      
      // Give it a moment to process
      setTimeout(() => {
        chartType.value = currentType;
      }, 50);
    });
  } catch (error) {
    console.error('Error refreshing dashboard data:', error);
    $q.notify({
      type: 'negative',
      message: 'Error refreshing dashboard data',
      position: 'top',
      timeout: 2000
    });
  }
}

// Format total sales correctly - made readonly
const formattedTotalSales = computed(() => {
  return dashboardStore.formatCurrency(dashboardStore.totalSales);
});

// Create a simplified trend view from the monthly sales data - made readonly
const calculatedSalesTrend = computed(() => {
  const monthlyData = dashboardStore.salesTrendData;
  if (!monthlyData || !monthlyData.datasets || !monthlyData.datasets[2] || !monthlyData.datasets[2].data) {
    return [0, 0, 0, 0, 0, 0, 0];
  }
  
  // Get the total items sold from the last 7 months (or fewer if not enough data)
  const totalSalesData = monthlyData.datasets[2].data;
  const currentMonth = new Date().getMonth();
  
  // Create an array of the last 7 months' data (for the sparkline)
  const result: number[] = [];
  for (let i = 0; i < 7; i++) {
    // Get data from the correct month, wrapping around to the previous year if needed
    const monthIndex = (currentMonth - i + 12) % 12;
    result.unshift(totalSalesData[monthIndex] || 0);
  }
  
  return result;
});

// Renamed from refreshAll to refreshDashboardData for clarity
const refreshDashboardData = async () => {
  try {
    if (isDev) console.log('Refreshing all dashboard data...');
    lastRefreshTimestamp.value = Date.now();
    
    // First refresh general dashboard data
    await dashboardStore.refreshDashboardData();
    
    // Explicitly log the current sales data to verify it's being retrieved correctly
    if (isDev) {
      console.log('Current sales data after refresh:');
      console.log('- Monthly:', dashboardStore.salesTrendData);
      console.log('- Weekly:', dashboardStore.weeklySalesTrendData);
      console.log('- Yearly:', dashboardStore.yearlySalesTrendData);
    }
    
    // Force currentSalesTrendData to refresh using a more reliable approach
    const currentPeriod = timePeriod.value;
    timePeriod.value = 'monthly';
    requestAnimationFrame(() => {
      timePeriod.value = currentPeriod;
    });
    
  } catch (error) {
    console.error('Failed to refresh dashboard:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to refresh dashboard data',
      position: 'top',
      timeout: 2000
    });
  }
};

// Computed property to check if the confirmation text is valid
const isConfirmationValid = computed(() => {
  return confirmationText.value === requiredConfirmation;
});

// Add a new state for the confirmation dialog
const showClearDataConfirmation = ref(false);

// Add a method to clear all data with improved error handling and rate limiting
const clearAllData = () => {
  if (isDev) console.log('Clearing all sales data...');
  
  // Add rate limiting - only allow one clear operation every 60 seconds
  const now = Date.now();
  const timeSinceLastClear = now - lastClearTimestamp.value;
  const minTimeBetweenClears = 60 * 1000; // 60 seconds
  
  if (timeSinceLastClear < minTimeBetweenClears) {
    const waitTimeSeconds = Math.ceil((minTimeBetweenClears - timeSinceLastClear) / 1000);
    $q.notify({
      type: 'warning',
      message: `Please wait ${waitTimeSeconds} seconds before attempting to clear data again`,
      position: 'top',
      timeout: 3000
    });
    return;
  }
  
  // Verify text confirmation
  if (confirmationText.value !== requiredConfirmation) {
    $q.notify({
      type: 'negative',
      message: `Please type "${requiredConfirmation}" to confirm data deletion`,
      position: 'top',
      timeout: 3000
    });
    return;
  }
  
  try {
    // Update last clear timestamp
    lastClearTimestamp.value = now;
    
    // Use the store's clearAllSalesData method
    const result = dashboardStore.clearAllSalesData();
    
    if (result) {
      // Show success notification
      $q.notify({
        type: 'positive',
        message: 'All sales data has been reset to zero',
        position: 'top',
        timeout: 2000
      });
      
      // Reset local state
      timePeriod.value = 'monthly';
      chartType.value = 'line';
      lastRefreshTimestamp.value = Date.now();
      confirmationText.value = ''; // Clear confirmation text
      
      // Force refresh using requestAnimationFrame for better timing
      requestAnimationFrame(() => {
        timePeriod.value = 'weekly';
        requestAnimationFrame(() => {
          timePeriod.value = 'monthly';
        });
      });
    } else {
      // Show error notification
      $q.notify({
        type: 'negative',
        message: 'Failed to clear sales data',
        position: 'top',
        timeout: 2000
      });
    }
  } catch (error) {
    console.error('Error clearing sales data:', error);
    $q.notify({
      type: 'negative',
      message: 'Error clearing sales data: ' + (error instanceof Error ? error.message : String(error)),
      position: 'top',
      timeout: 3000
    });
  }
};

</script>

<template>
  <q-page class="q-pa-md" @activated="activated">
    <div class="q-px-lg">
      <div class="row items-center q-mb-md">
        <div class="col">
          <h1 class="text-h4 q-mb-none">Dashboard</h1>
          <p class="text-subtitle1 q-mt-xs">System Overview</p>
        </div>
        <div class="col-auto">
          <div class="row q-gutter-sm">
            <q-btn 
              :class="[
                $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white'
              ]" 
              @click="refreshDashboardData" 
              :loading="dashboardStore.isLoading"
            >
              <q-icon :name="'refresh'" :color="$q.dark.isActive ? 'black' : 'white'" />
              <span :class="$q.dark.isActive ? 'text-black' : 'text-white'">Refresh</span>
            </q-btn>
            <q-btn 
              outline
              color="negative"
              @click="showClearDataConfirmation = true"
            >
              <q-icon name="restart_alt" class="q-mr-xs" />
              Clear Data
            </q-btn>
          </div>
        </div>
      </div>
      
      <!-- Confirmation Dialog for Clearing Data -->
      <q-dialog v-model="showClearDataConfirmation" persistent>
        <q-card style="min-width: 350px">
          <q-card-section class="row items-center">
            <q-avatar icon="warning" color="negative" text-color="white" />
            <span class="q-ml-sm text-h6">Clear All Sales Data?</span>
          </q-card-section>

          <q-card-section>
            This will permanently remove all historical sales data and reset all charts to zero.
            This action cannot be undone.
          </q-card-section>
          
          <q-card-section>
            <p class="text-caption q-mb-sm">To confirm, please type "{{ requiredConfirmation }}" below:</p>
            <q-input 
              v-model="confirmationText" 
              label="Confirmation" 
              outlined 
              dense
              autofocus
              :color="isConfirmationValid ? 'positive' : 'negative'"
            />
          </q-card-section>

          <q-card-actions align="right">
            <q-btn flat label="Cancel" color="primary" v-close-popup />
            <q-btn 
              flat 
              label="Clear All Data" 
              color="negative" 
              @click="clearAllData" 
              :disabled="!isConfirmationValid"
              v-close-popup="isConfirmationValid" 
            />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Key Metrics -->
      <div class="row q-col-gutter-md q-mb-lg">
        <!-- Total Inventory Value -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Total Inventory Value"
            icon="attach_money"
            :value="dashboardStore.formatCurrency(dashboardStore.totalInventoryValue)"
            :trend-data="dashboardStore.inventoryTrendData"
            :is-loading="dashboardStore.isLoading"
          />
        </div>

        <!-- Total Sales -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Total Sales"
            icon="shopping_cart"
            :value="formattedTotalSales"
            :trend-data="calculatedSalesTrend"
            :is-loading="dashboardStore.isLoading"
          />
        </div>

        <!-- Inventory Items -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Total Inventory Items"
            icon="inventory_2"
            :value="dashboardStore.totalInventoryItems"
            :trend-percentage="7.2"
            additional-info="Total across all categories"
            :trend-data="dashboardStore.inventoryTrendData"
            :is-loading="dashboardStore.isLoading"
          />
        </div>

        <!-- Low Stock Items -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Low Stock Items"
            icon="warning"
            :value="dashboardStore.lowStockItems"
            :trend-percentage="-4"
            :additional-info="`${dashboardStore.criticalItems} critical`"
            :trend-data="dashboardStore.lowStockTrendData"
            :is-loading="dashboardStore.isLoading"
          />
        </div>
      </div>

      <!-- Charts Section - Using the new component -->
      <chart-section
        :sales-trend-data="currentSalesTrendData"
        :inventory-data="dashboardStore.inventoryDistribution"
        v-model:time-period="timePeriod"
        v-model:chart-type="chartType"
        :is-dark="$q.dark.isActive"
        :is-loading="dashboardStore.isLoading"
      />

      <!-- Activity Feed Section -->
      <div class="row q-col-gutter-md q-mb-xl">
        <div class="col-12">
          <recent-activities-section :activities="dashboardStore.recentActivities" :is-loading="dashboardStore.isLoading" />
        </div>
      </div>
    </div>
  </q-page>
</template>

<style scoped>
.text-h4 {
  font-weight: 600;
}

.metrics-card {
  height: 100%;
}

.q-page {
  padding: 16px !important;
}

.q-px-lg {
  padding-left: 24px !important;
  padding-right: 24px !important;
}

/* Responsive adjustments for metrics cards */
@media (max-width: 599px) {
  .metrics-card {
    min-height: 140px;
  }
  .q-page {
    padding: 12px !important;
  }
  .q-px-lg {
    padding-left: 16px !important;
    padding-right: 16px !important;
  }
}
</style>
