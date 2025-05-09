<script setup lang="ts">
import { ref, onMounted, computed, watch, defineAsyncComponent } from 'vue';
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

// Add this computed property
const currentSalesTrendData = computed(() => {
  console.log('Computing current sales trend data for period:', timePeriod.value);
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
  console.log('Returning data:', data);
  return data;
});

onMounted(() => {
  // Don't reset chart data on mount, just refresh dashboard data
  // dashboardStore.resetChartData(); // Removing this line
  
  // Refresh general dashboard data (inventories, etc.)
  void dashboardStore.refreshDashboardData();
  
  // Log the total sales value for debugging
  console.log('Total sales value on mount:', dashboardStore.totalSales);
  console.log('Monthly sales trend data on mount:', dashboardStore.salesTrendData);
});

// Add watch for timePeriod changes - no longer needed to fetch data as it's already in the store
watch(
  () => timePeriod.value,
  (newValue) => {
    console.log('Time period changed to:', newValue);
  }
);

// Add watch for totalSales changes
watch(
  () => dashboardStore.totalSales,
  (newValue) => {
    console.log('Total sales changed to:', newValue);
  }
);

// Add watch for salesTrendData changes
watch(
  () => dashboardStore.salesTrendData,
  (newData) => {
    console.log('Sales trend data changed:', newData);
  },
  { deep: true }
);

// Function to force refresh data when page is activated
function onPageActivated() {
  console.log('Dashboard page activated - refreshing data');
  
  try {
    // Force refresh from localStorage
    const monthlySalesData = JSON.parse(localStorage.getItem('monthlySalesTrendData') || '{}');
    const weeklySalesData = JSON.parse(localStorage.getItem('weeklySalesTrendData') || '{}');
    const yearlySalesData = JSON.parse(localStorage.getItem('yearlySalesTrendData') || '{}');
    
    // Check if data is valid before assigning
    if (monthlySalesData && monthlySalesData.labels && monthlySalesData.datasets) {
      dashboardStore.$patch({ salesTrendData: monthlySalesData });
    }
    
    if (weeklySalesData && weeklySalesData.labels && weeklySalesData.datasets) {
      dashboardStore.$patch({ weeklySalesTrendData: weeklySalesData });
    }
    
    if (yearlySalesData && yearlySalesData.labels && yearlySalesData.datasets) {
      dashboardStore.$patch({ yearlySalesTrendData: yearlySalesData });
    }
    
    // Force the computed property to recalculate
    setTimeout(() => {
      const period = timePeriod.value;
      timePeriod.value = period === 'monthly' ? 'weekly' : 'monthly';
      setTimeout(() => {
        timePeriod.value = period;
      }, 5);
    }, 5);
  } catch (error) {
    console.error('Error refreshing dashboard data:', error);
  }
}

// Format total sales correctly
const formattedTotalSales = computed(() => {
  return dashboardStore.formatCurrency(dashboardStore.totalSales);
});

// Create a simplified trend view from the monthly sales data
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

const refreshAll = async () => {
  try {
    console.log('Refreshing all dashboard data...');
    
    // First refresh general dashboard data
    await dashboardStore.refreshDashboardData();
    
    // Explicitly log the current sales data to verify it's being retrieved correctly
    console.log('Current sales data after refresh:');
    console.log('- Monthly:', dashboardStore.salesTrendData);
    console.log('- Weekly:', dashboardStore.weeklySalesTrendData);
    console.log('- Yearly:', dashboardStore.yearlySalesTrendData);
    
    // Force currentSalesTrendData to refresh
    const period = timePeriod.value;
    timePeriod.value = 'monthly';
    setTimeout(() => {
      timePeriod.value = period;
    }, 5);
    
  } catch (error) {
    console.error('Failed to refresh dashboard:', error);
  }
};

// Add a new state for the confirmation dialog
const showClearDataConfirmation = ref(false);

// Add a method to clear all data
const clearAllData = () => {
  console.log('Clearing all sales data...');
  
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
    
    // Force refresh
    setTimeout(() => {
      timePeriod.value = 'weekly';
      setTimeout(() => {
        timePeriod.value = 'monthly';
      }, 10);
    }, 10);
  } else {
    // Show error notification
    $q.notify({
      type: 'negative',
      message: 'Failed to clear sales data',
      position: 'top',
      timeout: 2000
    });
  }
};

</script>

<template>
  <q-page class="q-pa-md" @activated="onPageActivated">
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
              @click="refreshAll" 
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

          <q-card-actions align="right">
            <q-btn flat label="Cancel" color="primary" v-close-popup />
            <q-btn 
              flat 
              label="Clear All Data" 
              color="negative" 
              @click="clearAllData" 
              v-close-popup 
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
