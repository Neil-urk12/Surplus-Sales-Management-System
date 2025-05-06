<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">Recent Activity</div>
      <div class="text-caption text-grey">Last 5 activities</div>
    </q-card-section>
    <q-separator />
    <q-card-section class="q-pa-none">
      <q-list separator class="activity-feed">
        <q-item v-for="activity in formattedActivities" :key="activity.id" class="activity-item q-pa-md">
          <q-item-section avatar>
            <q-icon :name="activity.icon" :color="activity.color" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ activity.title }}</q-item-label>
            <q-item-label caption>{{ activity.description }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-item-label caption class="timestamp">{{ activity.formattedTime }}</q-item-label>
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

function formatTimeAgo(date: Date): string {
  const now = new Date();
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return 'just now';
  } else if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60);
    return `${minutes}m ago`;
  } else if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600);
    return `${hours}h ago`;
  } else {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric'
    });
  }
}

const formattedActivities = computed(() => {
  // Ensure we only show the last 5 activities
  return props.activities.slice(0, 5).map(activity => ({
    ...activity,
    formattedTime: formatTimeAgo(new Date(activity.timestamp))
  }));
});
</script>

<style scoped>
.q-item {
  min-height: 72px;
}

.activity-feed {
  max-height: 400px;
  overflow-y: auto;
}

.activity-feed::-webkit-scrollbar {
  width: 8px;
}

.activity-feed::-webkit-scrollbar-track {
  background: transparent;
}

.activity-feed::-webkit-scrollbar-thumb {
  background: #e0e0e0;
  border-radius: 4px;
}

.activity-feed::-webkit-scrollbar-thumb:hover {
  background: #bdbdbd;
}

.dark .activity-feed::-webkit-scrollbar-thumb {
  background: #424242;
}

.dark .activity-feed::-webkit-scrollbar-thumb:hover {
  background: #616161;
}

.activity-item {
  transition: background-color 0.3s;
}

.activity-item:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

.dark .activity-item:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.timestamp {
  font-size: 0.8rem;
  opacity: 0.8;
}
</style> 