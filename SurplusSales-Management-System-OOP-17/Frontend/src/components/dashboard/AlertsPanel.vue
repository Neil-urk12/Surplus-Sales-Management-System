<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">System Alerts</div>
    </q-card-section>
    <q-separator />
    <q-card-section class="q-pa-none">
      <q-list>
        <q-item v-for="alert in alerts" :key="alert.id" class="q-pa-md">
          <q-item-section avatar>
            <q-icon :name="alert.icon" :color="getAlertColor(alert.severity)" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ alert.title }}</q-item-label>
            <q-item-label caption>{{ alert.message }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-btn 
              flat 
              round 
              dense 
              :icon="alert.actionIcon" 
              @click="handleAlertAction(alert.id)"
              :color="getAlertColor(alert.severity)"
            />
          </q-item-section>
        </q-item>
        <q-item v-if="alerts.length === 0" class="text-center q-pa-lg text-grey">
          <q-item-section>No active alerts</q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar';

/**
 * Defines the possible severity levels for system alerts
 * - error: Critical issues requiring immediate attention
 * - warning: Important issues that need attention but aren't critical
 * - success: Positive outcomes or completed actions
 * - info: General information or updates
 */
export type AlertSeverity = 
  | 'error'    // Critical system issues
  | 'warning'  // Important but not critical issues
  | 'success'  // Positive outcomes
  | 'info';    // Informational updates

/**
 * Maps AlertSeverity to Quasar color classes for consistent styling
 */
const severityColorMap: Record<AlertSeverity, string> = {
  error: 'negative',
  warning: 'warning',
  success: 'positive',
  info: 'info'
} as const;

/**
 * Represents a system alert notification
 */
export interface Alert {
  /** Unique identifier for the alert */
  id: string;
  /** Alert header text */
  title: string;
  /** Detailed alert message */
  message: string;
  /** Alert severity level */
  severity: AlertSeverity;
  /** Icon name from the icon set */
  icon: string;
  /** Icon for the action button */
  actionIcon: string;
}

const $q = useQuasar();
const props = defineProps<{
  alerts: Alert[];
}>();

const emit = defineEmits<{
  (e: 'alertAction', alertId: string): void;
}>();

/**
 * Get the appropriate color class for an alert's severity
 */
function getAlertColor(severity: AlertSeverity): string {
  return severityColorMap[severity];
}

function handleAlertAction(alertId: string) {
  try {
    // Validate alertId
    if (!alertId?.trim()) {
      throw new Error('Invalid alert ID provided');
    }

    // Check if the alert exists before emitting
    const alert = props.alerts.find(a => a.id === alertId);
    if (!alert) {
      throw new Error(`Alert with ID ${alertId} not found`);
    }

    // Log the action attempt
    console.log(`Processing ${alert.severity} alert action for: ${alert.title} (${alertId})`);

    // Emit the event
    emit('alertAction', alertId);

    // Log successful emission
    console.log(`Successfully processed ${alert.severity} alert action for: ${alert.title} (${alertId})`);
  } catch (error) {
    // Log the error with details
    console.error('Error handling alert action:', {
      error: error instanceof Error ? error.message : 'Unknown error',
      alertId,
      timestamp: new Date().toISOString()
    });

    // Show user-friendly notification
    $q.notify({
      type: 'negative',
      message: 'Failed to process alert action. Please try again.',
      position: 'top',
      timeout: 3000,
      actions: [
        { label: 'Dismiss', color: 'white' }
      ]
    });
  }
}
</script>

<style scoped>
.q-item {
  min-height: 72px;
}
</style> 