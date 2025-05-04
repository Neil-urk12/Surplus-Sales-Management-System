<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">Recent Activity</div>
    </q-card-section>
    <q-separator />
    <q-card-section class="q-pa-none">
      <q-list separator>
        <q-item v-for="activity in formattedActivities" :key="activity.id" class="q-pa-md">
          <q-item-section avatar>
            <q-icon :name="activity.icon" :color="activity.color" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ activity.title }}</q-item-label>
            <q-item-label caption>{{ activity.description }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-item-label caption>{{ activity.formattedTime }}</q-item-label>
          </q-item-section>
        </q-item>
        <q-item v-if="activities.length === 0" class="text-center q-pa-lg text-grey">
          <q-item-section>No recent activities</q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Activity {
  id: string;
  title: string;
  description: string;
  timestamp: Date;
  icon: string;
  color: string;
}

const props = defineProps<{
  activities: Activity[];
}>();

const formatTimeOptions = {
  hour: 'numeric',
  minute: 'numeric',
  hour12: true,
  month: 'short',
  day: 'numeric'
} as const;

const formattedActivities = computed(() => {
  return props.activities.map(activity => ({
    ...activity,
    formattedTime: new Date(activity.timestamp).toLocaleString('en-US', formatTimeOptions)
  }));
});
</script>

<style scoped>
.q-item {
  min-height: 72px;
}
</style> 