<script setup lang="ts">
import { ref, onMounted, computed, watch, defineAsyncComponent } from 'vue';
import { useQuasar } from 'quasar';
const MetricsCard = defineAsyncComponent(() => import('../components/dashboard/MetricsCard.vue'));
const ChartSection = defineAsyncComponent(() => import('../components/dashboard/ChartSection.vue'));
const RecentActivitiesSection = defineAsyncComponent(() => import('../components/dashboard/RecentActivitiesSection.vue'));

const $q = useQuasar();

// State
const isLoading = ref(false);
const totalInventoryValue = ref(2748560);
const totalSales = ref(845675);
const inventoryItems = ref(3845);
const lowStockItems = ref(24);
const timePeriod = ref<'weekly' | 'monthly' | 'yearly'>('monthly');
const chartType = ref<'line' | 'bar' | 'area'>('line');

// Mock data for demonstration
const salesTrendData = ref({
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
  datasets: [
    {
      label: 'Cars',
      data: [14000, 13500, 14500, 15000, 14500, 15500, 16000, 15500, 16500, 17000, 17500, 18000],
      borderColor: '#7B66FF',
      tension: 0.4,
      fill: false
    },
    {
      label: 'Accessories',
      data: [3000, 2800, 3000, 3000, 2900, 3000, 3100, 3000, 3100, 3200, 3200, 3300],
      borderColor: '#FFB534',
      tension: 0.4,
      fill: false
    },
    {
      label: 'Total',
      data: [21000, 20000, 21500, 22000, 21000, 22500, 23000, 22500, 24000, 25000, 26000, 27000],
      borderColor: '#FF6B6B',
      tension: 0.4,
      fill: false
    }
  ]
});

const inventoryData = ref({
  labels: ['Cabs', 'Materials', 'Accessories'],
  datasets: [{
    label: 'Inventory Distribution',
    data: [300, 150, 100],
    backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56']
  }]
});

// Add this interface before the recentActivities definition
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

// Add these after the salesTrendData definition
const weeklySalesTrendData = ref({
  labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
  datasets: [
    {
      label: 'Cars',
      data: [3500, 3800, 3600, 4000],
      borderColor: '#7B66FF',
      tension: 0.4,
      fill: false
    },
    {
      label: 'Accessories',
      data: [800, 750, 900, 850],
      borderColor: '#FFB534',
      tension: 0.4,
      fill: false
    },
    {
      label: 'Total',
      data: [4300, 4550, 4500, 4850],
      borderColor: '#FF6B6B',
      tension: 0.4,
      fill: false
    }
  ]
});

const yearlySalesTrendData = ref({
  labels: ['2020', '2021', '2022', '2023'],
  datasets: [
    {
      label: 'Cars',
      data: [160000, 180000, 200000, 220000],
      borderColor: '#7B66FF',
      tension: 0.4,
      fill: false
    },
    {
      label: 'Accessories',
      data: [35000, 38000, 40000, 42000],
      borderColor: '#FFB534',
      tension: 0.4,
      fill: false
    },
    {
      label: 'Total',
      data: [195000, 218000, 240000, 262000],
      borderColor: '#FF6B6B',
      tension: 0.4,
      fill: false
    }
  ]
});

// Add this computed property
const currentSalesTrendData = computed(() => {
  console.log('Computing current sales trend data for period:', timePeriod.value);
  let data;
  switch (timePeriod.value) {
    case 'weekly':
      data = weeklySalesTrendData.value;
      break;
    case 'yearly':
      data = yearlySalesTrendData.value;
      break;
    case 'monthly':
    default:
      data = salesTrendData.value;
      break;
  }
  console.log('Returning data:', data);
  return data;
});

// Add this type definition at the top of the script section
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

// Methods
function formatCurrency(value: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0
  }).format(value);
}

// Update the function signature with the type
function addNewActivity(activity: Activity) {
  // Add new activity at the beginning and keep only 5 items
  recentActivities.value = [activity, ...recentActivities.value.slice(0, 4)];
}

async function refreshData() {
  isLoading.value = true;
  try {
    // TODO: Implement actual data fetching
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Mock data update
    totalInventoryValue.value = 2748560;
    totalSales.value = 845675;
    inventoryItems.value = 3845;
    lowStockItems.value = 24;

    // Example of adding a new activity
    const newActivity = {
      id: Date.now().toString(),
      title: 'Data Refreshed',
      description: 'Dashboard data has been updated',
      timestamp: new Date(),
      icon: 'refresh',
      color: 'info'
    };
    addNewActivity(newActivity);
  } catch (error) {
    console.error('Error refreshing dashboard data:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to refresh dashboard data'
    });
  } finally {
    isLoading.value = false;
  }
}

const fetchTimeData = () => {
  isLoading.value = true;
  try {
    console.log('Fetching data for period:', timePeriod.value);
    // Simulated data - replace with actual API calls in production
    switch (timePeriod.value) {
      case 'weekly': {
        const weeklyData: ChartData = {
          labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
          datasets: [
            {
              label: 'Cars',
              data: generateRandomData(4, 3000, 4500),
              borderColor: '#7B66FF',
              tension: 0.4,
              fill: false
            },
            {
              label: 'Accessories',
              data: generateRandomData(4, 700, 1000),
              borderColor: '#FFB534',
              tension: 0.4,
              fill: false
            }
          ]
        };
        
        if (weeklyData.datasets[0] && weeklyData.datasets[1]) {
          const carsData = [...weeklyData.datasets[0].data];
          const accessoriesData = [...weeklyData.datasets[1].data];
          const totalData = carsData.map((val, idx) => val + (accessoriesData[idx] || 0));
          
          weeklyData.datasets.push({
            label: 'Total',
            data: totalData,
            borderColor: '#FF6B6B',
            tension: 0.4,
            fill: false
          });
        }
        
        weeklySalesTrendData.value = weeklyData;
        break;
      }
      
      case 'monthly': {
        const monthlyData: ChartData = {
          labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
          datasets: [
            {
              label: 'Cars',
              data: generateRandomData(12, 14000, 18000),
              borderColor: '#7B66FF',
              tension: 0.4,
              fill: false
            },
            {
              label: 'Accessories',
              data: generateRandomData(12, 2800, 3300),
              borderColor: '#FFB534',
              tension: 0.4,
              fill: false
            }
          ]
        };
        
        if (monthlyData.datasets[0] && monthlyData.datasets[1]) {
          const carsData = [...monthlyData.datasets[0].data];
          const accessoriesData = [...monthlyData.datasets[1].data];
          const totalData = carsData.map((val, idx) => val + (accessoriesData[idx] || 0));
          
          monthlyData.datasets.push({
            label: 'Total',
            data: totalData,
            borderColor: '#FF6B6B',
            tension: 0.4,
            fill: false
          });
        }
        
        salesTrendData.value = monthlyData;
        break;
      }
      
      case 'yearly': {
        const currentYear = new Date().getFullYear();
        const yearlyData: ChartData = {
          labels: [
            (currentYear - 3).toString(),
            (currentYear - 2).toString(),
            (currentYear - 1).toString(),
            currentYear.toString()
          ],
          datasets: [
            {
              label: 'Cars',
              data: generateRandomData(4, 160000, 220000),
              borderColor: '#7B66FF',
              tension: 0.4,
              fill: false
            },
            {
              label: 'Accessories',
              data: generateRandomData(4, 35000, 42000),
              borderColor: '#FFB534',
              tension: 0.4,
              fill: false
            }
          ]
        };
        
        if (yearlyData.datasets[0] && yearlyData.datasets[1]) {
          const carsData = [...yearlyData.datasets[0].data];
          const accessoriesData = [...yearlyData.datasets[1].data];
          const totalData = carsData.map((val, idx) => val + (accessoriesData[idx] || 0));
          
          yearlyData.datasets.push({
            label: 'Total',
            data: totalData,
            borderColor: '#FF6B6B',
            tension: 0.4,
            fill: false
          });
        }
        
        yearlySalesTrendData.value = yearlyData;
        break;
      }
    }
  } catch (error) {
    console.error('Error fetching time period data:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to fetch time period data'
    });
  } finally {
    isLoading.value = false;
  }
};

const generateRandomData = (length: number, min: number, max: number): number[] => {
  return Array.from({ length }, () => Math.floor(Math.random() * (max - min + 1)) + min);
};

onMounted(() => {
  void Promise.all([
    refreshData(),
    fetchTimeData() // Initial data fetch
  ]);
});

// Add watch for timePeriod changes
watch(
  () => timePeriod.value,
  (newValue) => {
    console.log('Time period changed to:', newValue);
    fetchTimeData();
  },
  { immediate: true }
);
</script>

<template>
  <q-page class="q-pa-md">
    <div class="q-px-lg">
      <div class="row items-center q-mb-md">
        <div class="col">
          <h1 class="text-h4 q-mb-none">Dashboard</h1>
          <p class="text-subtitle1 q-mt-xs">System Overview</p>
        </div>
        <div class="col-auto">
          <q-btn 
            :class="[
              $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white'
            ]" 
            @click="refreshData" 
            :loading="isLoading"
          >
            <q-icon :name="'refresh'" :color="$q.dark.isActive ? 'black' : 'white'" />
            <span :class="$q.dark.isActive ? 'text-black' : 'text-white'">Refresh</span>
          </q-btn>
        </div>
      </div>

      <!-- Key Metrics -->
      <div class="row q-col-gutter-md q-mb-lg">
        <!-- Total Inventory Value -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Total Inventory Value"
            icon="attach_money"
            :value="formatCurrency(totalInventoryValue)"
            :trend-percentage="16.2"
            ytd-value="$24.5M"
            :trend-data="[65, 68, 70, 72, 75, 78, 80]"
            :is-loading="isLoading"
          />
        </div>

        <!-- Total Sales -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Total Sales"
            icon="shopping_cart"
            :value="formatCurrency(totalSales)"
            :trend-percentage="12.5"
            ytd-value="$7.2M"
            :trend-data="[40, 42, 45, 48, 50, 52, 55]"
            :is-loading="isLoading"
          />
        </div>

        <!-- Inventory Items -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Inventory Items"
            icon="inventory_2"
            :value="inventoryItems"
            :trend-percentage="7.2"
            additional-info="85 new this week"
            :trend-data="[3200, 3300, 3400, 3500, 3600, 3700, 3845]"
            :is-loading="isLoading"
          />
        </div>

        <!-- Low Stock Items -->
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            class="metrics-card"
            title="Low Stock Items"
            icon="warning"
            :value="lowStockItems"
            :trend-percentage="-4"
            additional-info="5 critical"
            :trend-data="[28, 26, 25, 24, 25, 24, 24]"
            :is-loading="isLoading"
          />
        </div>
      </div>

      <!-- Charts Section - Using the new component -->
      <chart-section
        :sales-trend-data="currentSalesTrendData"
        :inventory-data="inventoryData"
        v-model:time-period="timePeriod"
        v-model:chart-type="chartType"
        :is-dark="$q.dark.isActive"
        :is-loading="isLoading"
      />

      <!-- Activity Feed Section -->
      <div class="row q-col-gutter-md q-mb-xl">
        <div class="col-12">
          <recent-activities-section :activities="recentActivities" :is-loading="isLoading" />
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
