<template>
  <q-page class="q-pa-md">
    <div class="q-px-lg">
      <div class="row items-center q-mb-lg">
        <div class="col">
          <h1 class="text-h4 q-mb-none">Dashboard</h1>
          <p class="text-subtitle1 q-mt-xs">System Overview</p>
        </div>
        <div class="col-auto">
          <q-btn color="primary" icon="refresh" label="Refresh" @click="refreshData" :loading="isLoading" />
        </div>
      </div>

      <!-- Key Metrics -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            title="Total Inventory"
            icon="inventory_2"
            :value="totalInventory"
            :subtitle="`${lowStockItems} items low in stock`"
          />
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            title="Monthly Sales"
            icon="trending_up"
            :value="formatCurrency(monthlySales)"
            subtitle="This month"
          />
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            title="Sold Cars"
            icon="directions_car"
            :value="recentOrders"
            subtitle="Last 7 days"
          />
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <metrics-card
            title="Total Revenue"
            icon="payments"
            :value="formatCurrency(totalRevenue)"
            subtitle="This month"
          />
        </div>
      </div>

      <!-- Charts Section -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-md-8">
          <q-card bordered>
            <q-card-section>
              <div class="text-h6">Sales Trend</div>
              <sales-trend-chart :chart-data="salesTrendData" :isDark="$q.dark.isActive" />
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-4">
          <q-card bordered>
            <q-card-section>
              <div class="text-h6">Inventory Distribution</div>
              <inventory-chart :chart-data="inventoryData" :isDark="$q.dark.isActive" />
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Activity and Alerts Section -->
      <div class="row q-col-gutter-md">
        <div class="col-12 col-md-8">
          <activity-feed :activities="recentActivities" />
        </div>
        <div class="col-12 col-md-4">
          <alerts-panel :alerts="systemAlerts" @alert-action="handleAlertAction" />
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useRouter } from 'vue-router';
import MetricsCard from '../components/dashboard/MetricsCard.vue';
import InventoryChart from '../components/dashboard/InventoryChart.vue';
import SalesTrendChart from '../components/charts/SalesTrendChart.vue';
import ActivityFeed from '../components/dashboard/ActivityFeed.vue';
import AlertsPanel, { type Alert, type AlertSeverity } from '../components/dashboard/AlertsPanel.vue';

const $q = useQuasar();
const router = useRouter();

// State
const isLoading = ref(false);
const totalInventory = ref(0);
const lowStockItems = ref(0);
const monthlySales = ref(0);
const recentOrders = ref(0);
const totalRevenue = ref(0);

// Mock data for demonstration
const salesTrendData = ref({
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
  datasets: [
    {
      label: 'Sales',
      data: [30, 45, 35, 50, 40, 60],
      borderColor: '#36A2EB',
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

const recentActivities = ref([
  {
    id: '1',
    title: 'New Car Sold',
    description: 'Multicab X1 sold to customer',
    timestamp: new Date(),
    icon: 'directions_car',
    color: 'positive'
  },
  {
    id: '2',
    title: 'Low Stock Alert',
    description: 'Multicab X2 stock is below minimum threshold',
    timestamp: new Date(Date.now() - 3600000),
    icon: 'warning',
    color: 'warning'
  },
  {
    id: '3',
    title: 'Inventory Updated',
    description: 'Added 5 units of Multicab X3',
    timestamp: new Date(Date.now() - 7200000),
    icon: 'inventory',
    color: 'info'
  }
]);

/**
 * System alerts configuration
 * Each alert represents a different type of system notification that requires user attention
 */
const systemAlerts = ref<Alert[]>([
  {
    id: '1',
    title: 'Low Stock Warning',
    message: '5 items are running low on stock',
    severity: 'warning' as AlertSeverity,
    icon: 'warning',
    actionIcon: 'visibility',
    // This alert will navigate to the inventory page to show low stock items
  },
  {
    id: '2',
    title: 'Inventory Check Required',
    message: 'Monthly inventory check due in 2 days',
    severity: 'info' as AlertSeverity,
    icon: 'inventory',
    actionIcon: 'event',
    // This alert will open the inventory check scheduling modal
  }
]);

// Methods
function formatCurrency(value: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'PHP'
  }).format(value);
}

async function refreshData() {
  isLoading.value = true;
  try {
    // TODO: Implement actual data fetching
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Mock data update
    totalInventory.value = 550;
    lowStockItems.value = 5;
    monthlySales.value = 250000;
    recentOrders.value = 12;
    totalRevenue.value = 1250000;
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

/**
 * Handles actions triggered from the AlertsPanel component
 * Each alert has a specific action associated with it
 * 
 * @param alertId - The unique identifier of the alert being acted upon
 */
async function handleAlertAction(alertId: string) {
  try {
    switch (alertId) {
      case '1': {
        // Navigate to inventory page with low stock filter
        // TODO: Implement navigation with proper filter state
        console.log('Navigating to inventory page for low stock items...');
        await router.push({
          path: '/inventory/cabs',
          query: { filter: 'low-stock' }
        });
        break;
      }
      case '2': {
        // Show inventory check scheduling modal
        // TODO: Implement inventory check modal/page
        console.log('Opening inventory check scheduling modal...');
        $q.dialog({
          title: 'Schedule Inventory Check',
          message: 'Would you like to schedule the monthly inventory check?',
          ok: {
            label: 'Schedule Now',
            color: 'primary'
          },
          cancel: {
            label: 'Remind Me Later',
            color: 'grey'
          },
          persistent: true
        }).onOk(() => {
          try {
            // TODO: Implement actual async scheduling functionality
            console.log('Scheduling inventory check...');
            $q.notify({
              type: 'positive',
              message: 'Inventory check scheduled successfully'
            });
          } catch (scheduleError) {
            console.error('Error scheduling inventory check:', scheduleError);
            $q.notify({
              type: 'negative',
              message: 'Failed to schedule inventory check',
              caption: 'Please try again later'
            });
          }
        });
        break;
      }
      default: {
        console.warn(`Unhandled alert action for alert ID: ${alertId}`);
        $q.notify({
          type: 'warning',
          message: 'This action is not yet implemented'
        });
      }
    }
  } catch (error) {
    console.error('Error handling alert action:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to process alert action',
      caption: 'Please try again or contact support if the issue persists'
    });
  }
}

onMounted(async () => {
  await refreshData();
});
</script>

<style scoped>
.text-h4 {
  font-weight: 600;
}
</style>

