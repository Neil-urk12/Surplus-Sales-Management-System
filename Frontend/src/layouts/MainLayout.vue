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
              <q-badge color="red" floating>
                {{ systemAlerts.length }}
              </q-badge>
              <q-menu anchor="bottom right" self="top right">
                <q-list style="min-width: 350px">
                  <q-item-label header class="row items-center justify-between">
                    <span>System Alerts</span>
                    <q-btn flat round dense icon="done_all" size="sm" @click="markAllAsRead" v-if="systemAlerts.length">
                      <q-tooltip>Mark all as read</q-tooltip>
                    </q-btn>
                  </q-item-label>
                  <q-separator />
                  <template v-if="systemAlerts.length">
                    <q-scroll-area style="height: 300px;">
                      <q-item v-for="alert in systemAlerts" :key="alert.id" clickable v-close-popup @click="handleAlertAction(alert.id)" :class="{ 'bg-grey-2': alert.read }">
                        <q-item-section avatar>
                          <q-icon :name="alert.icon" :color="getAlertColor(alert.severity)" />
                        </q-item-section>
                        <q-item-section>
                          <q-item-label class="text-weight-medium">{{ alert.title }}</q-item-label>
                          <q-item-label caption>{{ alert.message }}</q-item-label>
                          <q-item-label caption class="text-grey-7">{{ formatAlertTime(alert.timestamp) }}</q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-scroll-area>
                  </template>
                  <q-item v-else>
                    <q-item-section class="text-center text-grey q-py-md">
                      <q-icon name="notifications_none" size="48px" color="grey-4" />
                      <div class="q-mt-sm">No active alerts</div>
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
              <div class="flex flex-start row items-center justify-between">
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
                <q-btn-dropdown
                flat
                class="absolute-right"
                >
                  <div class="row no-wrap q-pa-md">
                    <div class="column">
                      <div class="text-h6 q-mb-md">Preference</div>
                        <div class="flex-start items-center">
                          <q-btn dense flat class="text-soft-light row full-width flex justify-between" @click=toggleColorMode>
                            <div class="flex items-center col-auto q-pr-sm full-width">
                                <div class="flex q-pr-sm q-pl-sm ">
                                  <q-icon
                                    v-if="!isDark"
                                    name="light_mode"
                                    class="text-soft-light"
                                    size="24px"
                                  />
                                  <q-icon
                                    v-else
                                    name="dark_mode"
                                    class="text-soft-light"
                                    size="24px"
                                  />
                                </div>
                                <div class="flex items-center col">
                                  <p class="q-mb-none text-weight-medium">{{ themeMode }}</p>
                                </div>
                            </div>
                          </q-btn>
                          <q-btn dense flat class="text-soft-light row full-width" @click=alert>
                            <div class="flex flex-start items-center">
                                <div class="flex items-center col-auto q-pr-sm">
                                  <q-btn dense round flat icon="notifications" class="text-soft-light">
                                  <q-badge color="red" floating class="notif-badge">
                                    {{ systemAlerts.length }}
                                  </q-badge>
                                </q-btn>
                                </div>
                                <div class="flex items-center col">
                                  <p class="q-mb-none text-weight-medium">Notifications</p>
                                </div>
                            </div>
                          </q-btn>
                        </div>
                    </div>

                    <q-separator vertical inset class="q-mx-lg" />

                    <div class="column items-center">
                      <q-avatar size="72px">
                        <img src="https://cdn.quasar.dev/img/boy-avatar.png">
                      </q-avatar>

                      <div class="text-subtitle1 q-mt-md q-mb-xs">{{ currentUser.name }}</div>

                      <q-btn
                        color="primary"
                        label="Logout"
                        push
                        flat
                        size="sm"
                        v-close-popup
                        @click="handleLogout"
                      />
                    </div>
                  </div>
                </q-btn-dropdown>
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

// Update systemAlerts with timestamps
const systemAlerts = ref<Alert[]>([
  {
    id: '1',
    title: 'Low Stock Warning',
    message: '5 items are running low on stock',
    severity: 'warning' as AlertSeverity,
    icon: 'warning',
    actionIcon: 'visibility',
    timestamp: new Date(),
    read: false
  },
  {
    id: '2',
    title: 'Inventory Check Required',
    message: 'Monthly inventory check due in 2 days',
    severity: 'info' as AlertSeverity,
    icon: 'inventory',
    actionIcon: 'event',
    timestamp: new Date(Date.now() - 3600000),
    read: false
  }
]);

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

onMounted(() => {
  const savedMode = localStorage.getItem('quasar-theme')
  if (savedMode) {
    $q.dark.set(savedMode === 'dark')
  } else {
    $q.dark.set(window.matchMedia('(prefers-color-scheme: dark)').matches)
  }
//for the time out cleanup
})

const toggleColorMode = () => {
  $q.dark.toggle()
}

const currentUser = computed(() => {
  return authStore.user ? {
    name: authStore.user.fullName,
    email: authStore.user.email,
    role: authStore.user.role
  } : { name: 'Guest', email: '', role: 'guest' }
})

function alert () {
  console.log('alert')
}

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

function markAllAsRead() {
  systemAlerts.value = systemAlerts.value.map(alert => ({
    ...alert,
    read: true
  }));
}

// Update handleAlertAction to mark alerts as read with type safety
async function handleAlertAction(alertId: string) {
  try {
    // Mark the alert as read
    const alertIndex = systemAlerts.value.findIndex(a => a.id === alertId);
    if (alertIndex !== -1) {
      const alert = systemAlerts.value[alertIndex];
      if (alert) {
        alert.read = true;
      }
    }

    switch (alertId) {
      case '1': {
        await router.push({
          path: '/inventory/cabs',
          query: { filter: 'low-stock' }
        });
        break;
      }
      case '2': {
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
          $q.notify({
            type: 'positive',
            message: 'Inventory check scheduled successfully'
          });
        });
        break;
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
</style>
