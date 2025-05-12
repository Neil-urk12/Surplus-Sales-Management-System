<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">Inventory Alerts</div>
    </q-card-section>
    <q-separator />
    <q-card-section class="q-pa-none">
      <q-list>
        <q-item v-for="alert in alerts" :key="`${alert.category}-${alert.status}`" class="q-pa-md">
          <q-item-section avatar>
            <q-icon :name="getAlertIcon(alert.category)" :color="getAlertColor(alert.status)" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ alert.category }} - {{ alert.status }}</q-item-label>
            <q-item-label caption>{{ alert.count }} {{ alert.category.toLowerCase() }} {{ alert.status.toLowerCase() }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-btn
              flat
              round
              dense
              icon="visibility"
              @click="handleAlertAction(alert)"
              :color="getAlertColor(alert.status)"
            />
          </q-item-section>
        </q-item>
        <q-item v-if="alerts.length === 0" class="text-center q-pa-lg text-grey">
          <q-item-section>No inventory alerts</q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { alertsService, type InventoryAlert } from 'src/services/alertsService';

const $q = useQuasar();
const router = useRouter();
const alerts = ref<InventoryAlert[]>([]);
const isLoading = ref(false);

// Fetch alerts on component mount
onMounted(async () => {
  await fetchAlerts();
});

/**
 * Fetches inventory alerts from the service
 */
async function fetchAlerts() {
  try {
    isLoading.value = true;
    alerts.value = await alertsService.getInventoryAlerts();
  } catch (error) {
    console.error('Error fetching alerts:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to load inventory alerts',
      position: 'top',
      timeout: 3000
    });
  } finally {
    isLoading.value = false;
  }
}

/**
 * Map of category to icon mappings for better maintainability
 */
const alertIcons = new Map<string, string>([
  ['Cabs', 'directions_car'],
  ['Accessories', 'settings_input_component'],
  ['Materials', 'category'],
]);

/**
 * Gets the appropriate icon for each category using the Map
 */
function getAlertIcon(category: string): string {
  return alertIcons.get(category) || 'warning';
}

/**
 * Map of status to color mappings for better maintainability
 */
const alertColors = new Map<string, string>([
  ['Low Stock', 'warning'],
  ['Out of Stock', 'negative'],
]);

/**
 * Gets the appropriate color for each alert status using the Map
 */
function getAlertColor(status: string): string {
  return alertColors.get(status) || 'grey';
}

/**
 * Handles the action when an alert is clicked
 */
async function handleAlertAction(alert: InventoryAlert) {
  try {
    // Navigate to the appropriate inventory page with filter
    const category = alert.category.toLowerCase();
    const status = alert.status.replace(' ', '-').toLowerCase();

    // Use template literals for better readability
    await router.push({
      path: `/inventory/${category}`,
      query: { status }
    });

    // Show notification using the same variables for consistency
    $q.notify({
      type: 'info',
      message: `Viewing ${status.replace('-', ' ')} ${category}`,
      position: 'top',
      timeout: 2000
    });
  } catch (error) {
    console.error('Error handling alert action:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to process alert action',
      position: 'top',
      timeout: 3000
    });
  }
}
</script>
