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
            <q-icon :name="alert.icon" :color="alert.severity" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ alert.title }}</q-item-label>
            <q-item-label caption>{{ alert.message }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-btn flat round dense :icon="alert.actionIcon" @click="handleAlertAction(alert.id)" />
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
export interface Alert {
  id: string;
  title: string;
  message: string;
  severity: 'negative' | 'warning' | 'positive' | 'info';
  icon: string;
  actionIcon: string;
}

defineProps<{
  alerts: Alert[];
}>();

const emit = defineEmits<{
  (e: 'alertAction', alertId: string): void;
}>();

function handleAlertAction(alertId: string) {
  emit('alertAction', alertId);
}
</script>

<style scoped>
.q-item {
  min-height: 72px;
}
</style> 