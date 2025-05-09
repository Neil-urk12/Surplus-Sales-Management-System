<template>
  <q-layout view="lHh Lpr lFf">
    <q-header :elevated="!isDark" class="custom-nav-height-width custom-rounded flex" style="height: 62px; background-color: var(--header-nav-bg);">
      <q-toolbar class="row items-center justify-between">
        <!-- Left section -->
        <div class="row items-center">
          <q-btn
            flat
            dense
            round
            icon="menu"
            aria-label="Menu"
            @click="toggleLeftDrawer"
            class="text-soft-light q-mr-sm"
          />
          <div class="text-soft-light text-weight-medium">{{ activePage }}</div>
        </div>

        <!-- Right section -->
        <div class="row items-center q-gutter-sm">
          <q-btn dense round flat :icon="$q.dark.isActive ? 'dark_mode' : 'light_mode'" @click="toggleColorMode" class="text-soft-light">
            <q-tooltip>{{ themeMode }}</q-tooltip>
          </q-btn>
          <div class="notification-wrapper">
            <q-btn dense round flat icon="notifications" class="text-soft-light">
            <q-badge
              v-if="unreadAlertsCount > 0"
              color="red"
              floating
            >
              {{ unreadAlertsCount }}
            </q-badge>
              <q-menu anchor="bottom right" self="top right" @hide="resetReadAlertsLimit">
                <q-list style="min-width: 350px">
                  <q-item-label header class="row items-center justify-between">
                    <span>{{ alertMessages.sections.systemAlerts }}</span>
                    <q-btn flat round dense icon="done_all" size="sm" @click="markAllAsRead" v-if="systemAlerts.length">
                      <q-tooltip>{{ alertMessages.actions.markAllAsRead }}</q-tooltip>
                    </q-btn>
                  </q-item-label>
                  <q-separator />
                  <template v-if="systemAlerts.length">
                    <q-scroll-area style="height: 300px;">
                      <!-- Unread Notifications Section -->
                      <div>
                        <q-item-label header class="text-primary q-px-md">
                          {{ alertMessages.sections.unreadNotifications }}
                        </q-item-label>

                        <template v-if="unreadAlerts.length > 0">
                          <q-item
                            v-for="alert in unreadAlerts"
                            :key="alert.id"
                            clickable
                            v-close-popup
                            @click="handleAlertAction(alert.id)"
                            class="notification-item"
                            :class="`notification-unread notification-unread-${getAlertColor(alert.severity)}`"
                          >
                            <q-item-section avatar>
                              <q-icon
                                :name="alert.icon"
                                :color="getAlertColor(alert.severity)"
                                class="notification-icon-unread"
                              />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="notification-title-unread">
                                {{ alert.title }}
                              </q-item-label>
                              <q-item-label caption>{{ alert.message }}</q-item-label>
                              <q-item-label caption class="text-grey-7">{{ formatAlertTime(alert.timestamp) }}</q-item-label>
                            </q-item-section>
                            <q-item-section side>
                              <q-btn
                                flat
                                round
                                dense
                                icon="done"
                                size="sm"
                                :color="getAlertColor(alert.severity)"
                                @click.stop="markAlertAsRead(alert.id)"
                              >
                                <q-tooltip>{{ alertMessages.actions.markAsRead }}</q-tooltip>
                              </q-btn>
                            </q-item-section>
                          </q-item>
                        </template>

                        <q-item v-else>
                          <q-item-section class="text-center text-grey q-py-sm">
                            <div>{{ alertMessages.empty.noUnreadNotifications }}</div>
                          </q-item-section>
                        </q-item>
                      </div>

                      <q-separator v-if="unreadAlerts.length > 0 && readAlerts.length > 0" class="q-my-sm" />

                      <!-- Read Notifications Section -->
                      <div>
                        <q-item-label header class="text-grey q-px-md q-pt-md">
                          {{ alertMessages.sections.readNotifications }}
                        </q-item-label>

                        <template v-if="readAlerts.length > 0">
                          <q-item
                            v-for="alert in readAlerts"
                            :key="alert.id"
                            clickable
                            v-close-popup
                            @click="handleAlertAction(alert.id)"
                            class="notification-item notification-read"
                          >
                            <q-item-section avatar>
                              <q-icon
                                :name="alert.icon"
                                :color="getAlertColor(alert.severity)"
                                class="notification-icon-read"
                              />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="notification-title-read">
                                {{ alert.title }}
                              </q-item-label>
                              <q-item-label caption>{{ alert.message }}</q-item-label>
                              <q-item-label caption class="text-grey-7">{{ formatAlertTime(alert.timestamp) }}</q-item-label>
                            </q-item-section>
                          </q-item>
                        </template>

                        <q-item v-else>
                          <q-item-section class="text-center text-grey q-py-sm">
                            <div>{{ alertMessages.empty.noReadNotifications }}</div>
                          </q-item-section>
                        </q-item>

                        <!-- Load More Button - Improves performance by loading read alerts incrementally -->
                        <q-item
                          v-if="hasMoreReadAlerts"
                          clickable
                          class="text-primary"
                          @click="loadMoreReadAlerts"
                        >
                          <q-item-section class="text-center">
                            <q-item-label>
                              {{ alertMessages.actions.loadMore(remainingReadAlerts) }}
                            </q-item-label>
                          </q-item-section>
                        </q-item>

                        <!-- View All Button -->
                        <q-item
                          v-if="allReadAlerts.length > 0"
                          clickable
                          class="text-primary"
                          @click="showAllReadAlerts"
                        >
                          <q-item-section class="text-center">
                            <q-item-label>
                              {{ alertMessages.actions.viewAll(allReadAlerts.length) }}
                            </q-item-label>
                          </q-item-section>
                        </q-item>
                      </div>
                    </q-scroll-area>
                  </template>
                  <q-item v-else>
                    <q-item-section class="text-center text-grey q-py-md flex content-center row">
                      <q-icon class="col" name="notifications_none" size="48px" color="grey-4" />
                      <div class="q-mt-sm col">{{ alertMessages.empty.noActiveNotifications }}</div>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>
          <q-btn v-if="route.path === '/inventory'" dense round flat icon="download" class="text-soft-light"></q-btn>
          <q-btn v-if="route.path === '/inventory'" class="text-white bg-primary">
            <q-icon color="white" name="add" />
            Add
          </q-btn>
        </div>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
    >
      <q-list>
        <q-item class="q-px-none q-pt-md q-pb-sm q-hoverable ">
          <q-item-section class="flex flex-start">
            <div class="row items-center justify-start q-ml-xs q-gutter-x-md">
              <img
                src="../assets/logo.png"
                style="width: 64px; height: 64px;"
                class="q-mt-xs"
              >
              <div class="text-h6 q-pl-xs text-soft-light">Cortes Surplus</div>
            </div>
          </q-item-section>
        </q-item>
        <div class="flex-col ">
          <!-- menu -->
          <MenuItems
            class="text-soft-light"
            v-for="link in filteredMenuItems"
            :key="link.title"
            v-bind="link"
            :isDark="isDark"
          />
            <q-item
              clickable
              v-bind="$attrs"
              class="text-soft-light q-hoverable absolute-bottom"
            >
              <!-- user -->
              <div class="flex row items-center justify-between q-px-sm">
                <div class="flex flex-start items-center q-gutter-md">
                  <q-avatar color="red user-profile-width user-profile-height" text-color="text-primary">
                      <img src="https://cdn.quasar.dev/img/avatar.png">
                  </q-avatar>
                  <q-item-section>
                    <q-item-label>{{ currentUser.name }}</q-item-label>
                  </q-item-section>
                </div>
              </div>
              <div class="q-pa-md col">
                <q-btn
                  label="Logout"
                  size="sm"
                  :class="[
                    $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white',
                    'q-px-md q-py-sm'
                  ]"
                  @click="handleLogout"
                />
              </div>
            </q-item>
        </div>

      </q-list>
    </q-drawer>

    <q-page-container style="padding-top: 62px;">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed, defineAsyncComponent } from 'vue'
import { useRoute, useRouter } from 'vue-router';
import { useQuasar } from 'quasar'
import { useAuthStore } from '../stores/auth'
const MenuItems = defineAsyncComponent(() => import('../components/SideMenuItems.vue'))
import type { menuItemsProps } from '../types/menu-items'
import type { AlertSeverity } from '../components/dashboard/AlertsPanel.vue';
import { alertsService, type InventoryAlert } from '../services/alertsService';
import { alertMessages } from '../constants/uiMessages';

const menuItemsList: menuItemsProps[] = [
  {
    title: 'Dashboard',
    icon: 'dashboard',
    to: '/',
    exact: true
  },
  {
    title: 'Inventory',
    icon: 'storage',
    hasSubmenu: true,
    children: [
      {
        title: 'Cabs',
        icon: 'directions_car',
        to: '/inventory/cabs'
      },
      {
        title: 'Materials',
        icon: 'category',
        to: '/inventory/materials'
      },
      {
        title: 'Accessories',
        icon: 'settings_input_component',
        to: '/inventory/accessories'
      }
    ]
  },
  {
    title: 'Sales',
    icon: 'trending_up',
    to: '/sales'
  },
  {
    title:"Contacts",
    icon:"contacts",
    to: '/contacts'
  },
];

const userManagementMenuItem: menuItemsProps = {
  title: 'User Management',
  icon: 'manage_accounts',
  to: '/user-management'
};

const filteredMenuItems = computed(() => {
  const baseItems = [...menuItemsList]; // Create a mutable copy
  const userRole = authStore.user?.role;
  if (userRole === 'admin' || userRole === 'staff') {
    baseItems.splice(1, 0, userManagementMenuItem);
  }
  return baseItems;
});

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const $q = useQuasar();
const activePage = ref('')
const leftDrawerOpen = ref(false)
const isDark = computed(() => $q.dark.isActive)

// Update the Alert interface
interface Alert {
  id: string;
  title: string;
  message: string;
  severity: AlertSeverity;
  icon: string;
  actionIcon: string;
  timestamp: Date;
  read?: boolean;
}

// Add severity color mapping
const severityColorMap: Record<AlertSeverity, string> = {
  error: 'negative',
  warning: 'warning',
  success: 'positive',
  info: 'info'
};

// System alerts array
const systemAlerts = ref<Alert[]>([]);

// Inventory alerts from the service
const inventoryAlerts = ref<InventoryAlert[]>([]);

// Function to fetch inventory alerts and convert them to system alerts
async function fetchInventoryAlerts() {
  try {
    inventoryAlerts.value = await alertsService.getInventoryAlerts();

    // Convert inventory alerts to system alerts
    const newSystemAlerts: Alert[] = inventoryAlerts.value.map((alert) => {
      // Create a unique ID for each alert based on category and status
      // This ensures consistent IDs across refreshes and prevents duplicate IDs
      const alertId = `inventory-${alert.category.toLowerCase()}-${alert.status.toLowerCase()}`;

      // Determine severity based on status
      const severity: AlertSeverity = alert.status === 'Out of Stock' ? 'error' : 'warning';

      // Determine icon based on category
      let icon = 'warning';
      if (alert.category === 'Cabs') icon = 'directions_car';
      else if (alert.category === 'Accessories') icon = 'settings_input_component';
      else if (alert.category === 'Materials') icon = 'category';

      // Create the alert message
      const title = `${alert.category} ${alert.status}`;
      const message = `${alert.count} ${alert.category.toLowerCase()} ${alert.status.toLowerCase()}`;

      return {
        id: alertId,
        title,
        message,
        severity,
        icon,
        actionIcon: 'visibility',
        timestamp: new Date(),
        read: false
      };
    });

    // Add inventory check reminder if we're within 5 days of the end of the month
    const daysUntilEndOfMonth = getDaysUntilEndOfMonth();
    if (shouldShowInventoryCheckReminder()) {
      // Get the name of the current month
      const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
                          'July', 'August', 'September', 'October', 'November', 'December'];
      const currentMonth = monthNames[new Date().getMonth()];

      // Determine severity based on how close we are to the end of the month
      let severity: AlertSeverity = 'info';
      let title = 'Monthly Inventory Check Reminder';
      let message = `${currentMonth} inventory check due in ${daysUntilEndOfMonth} day${daysUntilEndOfMonth === 1 ? '' : 's'}`;

      if (daysUntilEndOfMonth <= 2) {
        severity = 'warning';
        title = 'Urgent: Inventory Check Due Soon';
      }

      if (daysUntilEndOfMonth === 0) {
        severity = 'error';
        title = 'CRITICAL: Inventory Check Due Today';
        message = `${currentMonth} inventory check must be completed TODAY!`;
      }

      newSystemAlerts.push({
        id: 'inventory-check',
        title: title,
        message: message,
        severity: severity,
        icon: 'inventory',
        actionIcon: 'event',
        timestamp: new Date(),
        read: false
      });
    }

    // Update system alerts
    systemAlerts.value = newSystemAlerts;
  } catch (error) {
    console.error('Error fetching inventory alerts:', error);
    // Add a fallback alert if there's an error
    systemAlerts.value = [{
      id: 'error',
      title: 'Error Loading Notifications',
      message: 'Failed to load inventory notifications',
      severity: 'error',
      icon: 'error',
      actionIcon: 'refresh',
      timestamp: new Date(),
      read: false
    }];
  }
}

watch(() => route.path, (newPath) => {
  switch (newPath) {
    case '/':
      activePage.value = 'Dashboard'
      break
    case '/sales':
      activePage.value = 'Sales'
      break
    case '/inventory/cabs':
      activePage.value = 'Cabs'
      break
    case '/inventory/materials':
      activePage.value = 'Materials'
      break
    case '/inventory/accessories':
      activePage.value = 'Accessories'
      break
    case '/contacts':
      activePage.value = 'Contacts'
      break
    default:
      activePage.value = 'Dashboard'
  }
}, { immediate: true })

const themeMode = computed(() => ($q.dark.isActive ? 'Dark Mode' : 'Light Mode'))

// Computed property for unread alerts
const unreadAlerts = computed(() => {
  return systemAlerts.value.filter(alert => !alert.read)
    .sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime()); // Sort by timestamp, newest first
})

// Get all read alerts (for counting)
const allReadAlerts = computed(() => {
  return systemAlerts.value.filter(alert => alert.read)
    .sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime()); // Sort by timestamp, newest first
})

// State variable to track how many read alerts to display
const readAlertsLimit = ref(5);

/**
 * Computed property for read alerts with dynamic limit
 * This implements a pagination-like approach where we initially show a limited number
 * of read alerts and allow the user to load more as needed, improving performance
 * when there are many read alerts in the system.
 */
const readAlerts = computed(() => {
  return allReadAlerts.value.slice(0, readAlertsLimit.value); // Limit to the current readAlertsLimit
})

// Check if there are more read alerts than we're showing
const hasMoreReadAlerts = computed(() => {
  return allReadAlerts.value.length > readAlertsLimit.value;
})

// Calculate the number of remaining read alerts
const remainingReadAlerts = computed(() => {
  return allReadAlerts.value.length - readAlertsLimit.value;
})

/**
 * Function to load more read alerts
 * Increases the limit by 5 each time, up to a maximum of 20 at a time
 * to prevent performance issues with loading too many alerts at once
 */
function loadMoreReadAlerts() {
  // Define a maximum number of alerts to show at once to prevent performance issues
  const MAX_ALERTS_TO_SHOW = 20;

  // Increase the limit by 5 each time, but don't exceed the maximum
  readAlertsLimit.value = Math.min(readAlertsLimit.value + 5, MAX_ALERTS_TO_SHOW);

  // If we've reached the maximum but there are still more alerts,
  // show a notification to the user
  if (readAlertsLimit.value === MAX_ALERTS_TO_SHOW && remainingReadAlerts.value > 0) {
    $q.notify({
      type: 'info',
      message: `Showing ${MAX_ALERTS_TO_SHOW} most recent read alerts. Use "View All" to see all ${allReadAlerts.value.length} alerts.`,
      position: 'top',
      timeout: 3000
    });
  }
}

// Function to reset the read alerts limit when the menu is closed
function resetReadAlertsLimit() {
  // Reset to the initial limit
  readAlertsLimit.value = 5;
}

// Computed property to count unread alerts
const unreadAlertsCount = computed(() => {
  return unreadAlerts.value.length;
})

onMounted(async () => {
  // Set theme mode
  const savedMode = localStorage.getItem('quasar-theme')
  if (savedMode) {
    $q.dark.set(savedMode === 'dark')
  } else {
    $q.dark.set(window.matchMedia('(prefers-color-scheme: dark)').matches)
  }

  // Fetch inventory alerts
  try {
    await fetchInventoryAlerts();
  } catch (error) {
    console.error('Failed to fetch inventory alerts during initialization:', error);
    $q.notify({
      type: 'negative',
      message: 'Failed to load inventory alerts',
      caption: 'Please refresh the page or contact support if the issue persists',
      position: 'top',
      timeout: 5000
    });
  }

  // Set up interval to refresh alerts every 5 minutes
  const alertsInterval = setInterval(() => {
    // Wrap in async IIFE with try-catch for proper error handling
    void (async () => {
      try {
        await fetchInventoryAlerts();
      } catch (error) {
        console.error('Failed to fetch inventory alerts during refresh interval:', error);
        // Don't show notification for background refresh errors to avoid spamming the user
        // The error alert will already be shown by the fetchInventoryAlerts function
      }
    })();
  }, 5 * 60 * 1000)

  // Clean up interval on component unmount
  return () => {
    clearInterval(alertsInterval)
  }
})

const toggleColorMode = () => {
  $q.dark.toggle()
  localStorage.setItem('quasar-theme', $q.dark.isActive ? 'dark' : 'light')
}

const currentUser = computed(() => {
  return authStore.user ? {
    name: authStore.user.fullName,
    email: authStore.user.email,
    role: authStore.user.role
  } : { name: 'Guest', email: '', role: 'guest' }
})


function toggleLeftDrawer () {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

async function handleLogout() {
  authStore.logout()
  await router.push('/login')
}

// Add helper functions
function getAlertColor(severity: AlertSeverity): string {
  return severityColorMap[severity];
}

function formatAlertTime(date: Date): string {
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes} minutes ago`;
  if (hours < 24) return `${hours} hours ago`;
  if (days === 1) return 'yesterday';
  return `${days} days ago`;
}

/**
 * Calculates the number of days until the end of the current month
 * @returns The number of days until the end of the month
 */
function getDaysUntilEndOfMonth(): number {
  const today = new Date();
  // Create a date for the first day of the next month
  const nextMonth = new Date(today.getFullYear(), today.getMonth() + 1, 1);
  // Subtract one day to get the last day of the current month
  const lastDayOfMonth = new Date(nextMonth.getTime() - 86400000);
  // Calculate the difference in days
  const diffTime = lastDayOfMonth.getTime() - today.getTime();
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return diffDays;
}

/**
 * Checks if we're within the reminder period for the monthly inventory check
 * and if the user hasn't already scheduled the check for this month
 * @param daysBeforeEndOfMonth Number of days before the end of the month to show the reminder
 * @returns True if we should show the reminder, false otherwise
 */
function shouldShowInventoryCheckReminder(daysBeforeEndOfMonth: number = 5): boolean {
  const daysRemaining = getDaysUntilEndOfMonth();

  // Check if the user has already scheduled the inventory check for this month
  const today = new Date();
  const currentMonthYear = `${today.getMonth()}-${today.getFullYear()}`;
  const scheduledCheck = localStorage.getItem('scheduledInventoryCheck');

  // If the user has scheduled a check for the current month, don't show the reminder
  if (scheduledCheck === currentMonthYear) {
    return false;
  }

  // Show the reminder if we're within the specified number of days of the end of the month
  return daysRemaining <= daysBeforeEndOfMonth;
}

/**
 * Marks the inventory check as scheduled for the current month
 */
function markInventoryCheckScheduled(): void {
  const today = new Date();
  const currentMonthYear = `${today.getMonth()}-${today.getFullYear()}`;
  localStorage.setItem('scheduledInventoryCheck', currentMonthYear);
}

/**
 * Marks a specific notification as read
 * @param alertId The ID of the notification to mark as read
 */
function markAlertAsRead(alertId: string) {
  const alertIndex = systemAlerts.value.findIndex(a => a.id === alertId);
  if (alertIndex !== -1) {
    const alert = systemAlerts.value[alertIndex];
    if (alert) {
      alert.read = true;
    }
  }
}

/**
 * Shows a dialog with all read notifications
 */
function showAllReadAlerts() {
  // Create a list of all read notifications for the dialog
  const alertsList = allReadAlerts.value.map(alert => {
    return `<div class="q-mb-md">
      <div class="text-weight-medium">${alert.title}</div>
      <div class="text-caption">${alert.message}</div>
      <div class="text-caption text-grey">${formatAlertTime(alert.timestamp)}</div>
    </div>`;
  }).join('');

  // Show dialog with all read notifications
  $q.dialog({
    title: alertMessages.sections.readNotifications,
    message: `<div class="q-pa-md">${alertsList}</div>`,
    html: true,
    style: 'width: 500px; max-width: 80vw',
    ok: {
      label: 'Close',
      flat: true,
      color: 'primary'
    }
  });
}

/**
 * Marks all notifications as read
 */
function markAllAsRead() {
  systemAlerts.value = systemAlerts.value.map(alert => ({
    ...alert,
    read: true
  }));
}

// Update handleAlertAction to handle notification clicks
async function handleAlertAction(alertId: string) {
  try {
    // Mark the notification as read
    const alertIndex = systemAlerts.value.findIndex(a => a.id === alertId);
    if (alertIndex !== -1) {
      const alert = systemAlerts.value[alertIndex];
      if (alert) {
        alert.read = true;
      }
    }

    // Handle inventory alerts
    if (alertId.startsWith('inventory-')) {
      const parts = alertId.split('-');
      const category = parts[1];
      const status = parts[2];

      // Navigate to the appropriate inventory page with filter
      if (category && ['cabs', 'accessories', 'materials'].includes(category)) {
        await router.push({
          path: `/inventory/${category}`,
          query: { status: status }
        });

        $q.notify({
          type: 'info',
          message: `Viewing ${status} ${category}`,
          position: 'top',
          timeout: 2000
        });
      }
    }
    // Handle inventory check alert
    else if (alertId === 'inventory-check') {
      // Get the last day of the month for the inventory check date
      const today = new Date();
      const nextMonth = new Date(today.getFullYear(), today.getMonth() + 1, 1);
      const lastDayOfMonth = new Date(nextMonth.getTime() - 86400000);
      const formattedDate = lastDayOfMonth.toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });

      // Get the name of the current month
      const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
                          'July', 'August', 'September', 'October', 'November', 'December'];
      const currentMonth = monthNames[today.getMonth()];

      $q.dialog({
        title: 'Schedule Inventory Check',
        message: `The ${currentMonth} inventory check is due by ${formattedDate}. Would you like to schedule it now?`,
        html: true,
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
        // Mark the inventory check as scheduled for this month
        markInventoryCheckScheduled();

        // Remove the notification from the list
        const alertIndex = systemAlerts.value.findIndex(a => a.id === alertId);
        if (alertIndex !== -1) {
          systemAlerts.value.splice(alertIndex, 1);
        }

        $q.notify({
          type: 'positive',
          message: `${currentMonth} inventory check scheduled for ${formattedDate}`,
          timeout: 3000
        });
      });
    }
    // Handle error alert
    else if (alertId === 'error') {
      // Refresh notifications
      try {
        $q.notify({
          type: 'info',
          message: 'Refreshing notifications',
          position: 'top',
          timeout: 2000
        });

        await fetchInventoryAlerts();

        $q.notify({
          type: 'positive',
          message: 'Notifications refreshed successfully',
          position: 'top',
          timeout: 2000
        });
      } catch (error) {
        console.error('Failed to refresh notifications:', error);
        $q.notify({
          type: 'negative',
          message: 'Failed to refresh notifications',
          caption: 'Please try again or contact support if the issue persists',
          position: 'top',
          timeout: 3000
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
</script>

<style>
.notification-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.q-badge {
  font-size: 0.75rem;
  padding: 2px 6px;
}

/* Remove old styles */
.notification-container,
.notification-btn,
.notif-badge,
.notif-count {
  display: none !important;
}

/* Notification styling */
.notification-item {
  transition: all 0.3s ease;
}

.notification-unread {
  border-left: 4px solid;
  background-color: rgba(var(--q-primary-rgb), 0.05);
}

.notification-unread-error {
  border-left-color: var(--q-negative);
}

.notification-unread-warning {
  border-left-color: var(--q-warning);
}

.notification-unread-info {
  border-left-color: var(--q-info);
}

.notification-unread-success {
  border-left-color: var(--q-positive);
}

.notification-read {
  opacity: 0.7;
  background-color: transparent;
  border-left: 4px solid transparent;
}

.q-dark .notification-unread {
  background-color: rgba(255, 255, 255, 0.1);
}

.q-dark .notification-read {
  background-color: transparent;
  opacity: 0.5;
}

.notification-title-unread {
  font-weight: 600;
}

.notification-title-read {
  font-weight: 400;
}

.notification-icon-unread {
  transform: scale(1.1);
}

.notification-icon-read {
  opacity: 0.7;
}
</style>
