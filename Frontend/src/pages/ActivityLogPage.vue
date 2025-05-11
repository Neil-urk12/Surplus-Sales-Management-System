<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useLogStore } from '../stores/logs';
import type { LogFilters, ActivityLogEntry, ActionStatus } from '../types/logTypes';
import type { QTableProps } from 'quasar';
import { date as quasarDateUtils } from 'quasar';

const logStore = useLogStore();

// Local state for filters to enable 'Apply Filters' button functionality
const localFilters = ref<LogFilters & { searchQuery: string | null }>({
  dateFrom: null,
  dateTo: null,
  showSystemActions: true,
  showSuccessful: true,
  showFailed: true,
  selectedActionType: 'All Actions',
  selectedUserId: null,
  searchQuery: null,
});

// Date and Time refs for q-input and q-popup-proxy
const dateFromProxy = ref<string | null>(null); // Format: YYYY/MM/DD
const timeFromProxy = ref<string | null>(null); // Format: HH:mm
const dateToProxy = ref<string | null>(null);   // Format: YYYY/MM/DD
const timeToProxy = ref<string | null>(null);   // Format: HH:mm

// Computed property for displaying the date range on the button
const dateRangeDisplay = computed(() => {
  if (localFilters.value.dateFrom && localFilters.value.dateTo) {
    const from = quasarDateUtils.formatDate(localFilters.value.dateFrom, 'YYYY/MM/DD HH:mm');
    const to = quasarDateUtils.formatDate(localFilters.value.dateTo, 'YYYY/MM/DD HH:mm');
    return `${from} - ${to}`;
  }
  if (localFilters.value.dateFrom) {
    return `From ${quasarDateUtils.formatDate(localFilters.value.dateFrom, 'YYYY/MM/DD HH:mm')}`;
  }
  if (localFilters.value.dateTo) {
    return `To ${quasarDateUtils.formatDate(localFilters.value.dateTo, 'YYYY/MM/DD HH:mm')}`;
  }
  return 'Custom Range';
});

// Function to combine date and time and update localFilters
function updateLocalFilterDate(type: 'dateFrom' | 'dateTo') {
  const dateVal = type === 'dateFrom' ? dateFromProxy.value : dateToProxy.value;
  const timeVal = type === 'dateFrom' ? timeFromProxy.value : timeToProxy.value;

  if (dateVal) {
    const dateParts = dateVal.split('/').map(part => parseInt(part, 10));
    const year = dateParts[0];
    const month = dateParts[1]; // This is 1-based from q-date
    const day = dateParts[2];

    // Check if all parts are valid numbers
    if (typeof year === 'number' && !isNaN(year) &&
        typeof month === 'number' && !isNaN(month) &&
        typeof day === 'number' && !isNaN(day)) {
      const combinedDate = new Date(year, month - 1, day); // Changed to const

      if (timeVal) {
        const timeParts = timeVal.split(':').map(part => parseInt(part, 10));
        const hours = timeParts[0];
        const minutes = timeParts[1];
        if (typeof hours === 'number' && !isNaN(hours) &&
            typeof minutes === 'number' && !isNaN(minutes)) {
          combinedDate.setHours(hours, minutes, 0, 0);
        } else {
          // If time is invalid or partially set, just use the date part (already set)
          // Or, decide on a default time, e.g., 00:00 for from, 23:59 for to
           if (type === 'dateFrom') combinedDate.setHours(0,0,0,0);
           else combinedDate.setHours(23,59,59,999);
        }
      } else {
        if (type === 'dateFrom') {
          combinedDate.setHours(0, 0, 0, 0);
        } else {
          combinedDate.setHours(23, 59, 59, 999);
        }
      }
      localFilters.value[type] = combinedDate;
    } else {
      localFilters.value[type] = null; // Invalid date string parts
    }
  } else {
    localFilters.value[type] = null;
  }
}

// Watchers to update localFilters when proxy models change
watch(dateFromProxy, () => updateLocalFilterDate('dateFrom'));
watch(timeFromProxy, () => updateLocalFilterDate('dateFrom'));
watch(dateToProxy, () => updateLocalFilterDate('dateTo'));
watch(timeToProxy, () => updateLocalFilterDate('dateTo'));

// Sync localFilters with store filters initially and when store filters might change externally
// Also, pre-fill date/time proxies if store has values
watch(() => logStore.filters, (newStoreFilters) => {
  localFilters.value = { 
    ...newStoreFilters,
    searchQuery: localFilters.value.searchQuery 
   };
  setDateTimeProxiesFromFilters();
}, { deep: true, immediate: true });

onMounted(async () => {
  try {
    await logStore.fetchLogs();
    setDateTimeProxiesFromFilters();
  } catch (error) {
    console.error('Error fetching logs:', error);
  }
});

const isLoading = computed(() => logStore.isLoading);
const baseUsersForFilter = computed<Array<{ id: string | null; fullName: string; role?: string }>>(() => [
  { id: null, fullName: 'All Users' }, 
  ...logStore.users.map(user => ({
    id: user.id,
    fullName: `${user.fullName} (${user.role})`,
    role: user.role
  }))
]);
const actionTypesForFilter = computed(() => logStore.allActionTypes);

// For q-select user filtering
const displayUserOptions = ref<Array<{ id: string | null; fullName: string; role?: string }>>([]);
watch(baseUsersForFilter,
  (newUsers) => {
    displayUserOptions.value = [...newUsers];
  },
  { immediate: true, deep: true }
);

function filterUserOptions(val: string, update: (callbackFn: () => void) => void) {
  update(() => {
    if (val === '') {
      displayUserOptions.value = [...baseUsersForFilter.value];
    } else {
      const needle = val.toLowerCase();
      displayUserOptions.value = baseUsersForFilter.value.filter(
        (v) => v.fullName.toLowerCase().indexOf(needle) > -1 || 
               (v.role && v.role.toLowerCase().indexOf(needle) > -1)
      );
    }
  });
}

const qTablePagination = ref<QTableProps['pagination']>({
  sortBy: 'timestamp',
  descending: true,
  page: 1,
  rowsPerPage: 10,
  rowsNumber: 0 // This will be updated by totalLogs
});

const paginatedLogs = computed(() => logStore.paginatedLogs);
const totalLogs = computed(() => logStore.totalFilteredLogs);

watch(totalLogs, (newTotal) => {
  if(qTablePagination.value) {
    qTablePagination.value.rowsNumber = newTotal;
  }
});

watch(qTablePagination, async (newVal) => {
  if (newVal) {
    try {
      await logStore.setCurrentPage(newVal.page ?? 1);
      await logStore.setPageSize(newVal.rowsPerPage ?? 10);
      // Sorting is handled by q-table client-side from the `paginatedLogs` array for now.
      // If server-side sorting is needed, this is where you'd trigger it.
    } catch (error) {
      console.error('Error updating pagination in store:', error);
      // Optionally, you could add a user-facing error message here
      // For example, using a toast notification.
    }
  }
}, { deep: true });

watch(() => logStore.currentPage, (newPage) => {
  if (qTablePagination.value && qTablePagination.value.page !== newPage) {
    qTablePagination.value.page = newPage;
  }
});
watch(() => logStore.pageSize, (newSize) => {
  if (qTablePagination.value && qTablePagination.value.rowsPerPage !== newSize) {
    qTablePagination.value.rowsPerPage = newSize;
  }
});

const columns: QTableProps['columns'] = [
  { name: 'timestamp', required: true, label: 'Date & Time', field: 'timestamp', sortable: true, align: 'left',
    format: (val) => formatDate(val) },
  { name: 'user', required: true, label: 'User', field: (row: ActivityLogEntry) => row.user.fullName, sortable: true, align: 'left' },
  { name: 'role', required: true, label: 'Role', field: (row: ActivityLogEntry) => row.user.role, sortable: true, align: 'left' },
  { name: 'actionType', required: true, label: 'Action Type', field: 'actionType', sortable: true, align: 'left' },
  { name: 'details', required: true, label: 'Details', field: 'details', align: 'left', style: 'min-width: 250px; white-space: normal;word-break:break-word;' },
  { name: 'status', required: true, label: 'Status', field: 'status', sortable: true, align: 'center' },
  { name: 'isSystemAction', required: true, label: 'System', field: 'isSystemAction', align: 'center' },
];

function applyFilters() {
  updateLocalFilterDate('dateFrom');
  updateLocalFilterDate('dateTo');
  if (qTablePagination.value) qTablePagination.value.page = 1;
  logStore.updateFilters({ 
    dateFrom: localFilters.value.dateFrom,
    dateTo: localFilters.value.dateTo,
    showSystemActions: localFilters.value.showSystemActions,
    showSuccessful: localFilters.value.showSuccessful,
    showFailed: localFilters.value.showFailed,
    selectedActionType: localFilters.value.selectedActionType,
    selectedUserId: localFilters.value.selectedUserId,
    searchQuery: localFilters.value.searchQuery
  });
}

function refreshData() {
  if (qTablePagination.value) qTablePagination.value.page = 1;
  logStore.refreshLogs();
  const currentStoreFilters = logStore.filters;
  localFilters.value = {
     ...currentStoreFilters,
     searchQuery: null // Reset search query on refresh as well
  };
  setDateTimeProxiesFromFilters(); // This was the intended call
}

function formatDate(val: Date | string | number | undefined | null): string {
  if (!val) return 'N/A';
  // Ensure val is a suitable type for Quasar's formatDate if it's not already a Date object
  const dateToFormat = typeof val === 'string' || typeof val === 'number' ? new Date(val) : val;
  return quasarDateUtils.formatDate(dateToFormat, 'YYYY-MM-DD HH:mm:ss'); // Corrected usage
}

function getStatusChipColor(status: ActionStatus): string {
  return status === 'successful' ? 'positive' : 'negative';
}

function onRequest(props: { pagination: QTableProps['pagination'] }) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination ?? {};
  if (qTablePagination.value) {
    qTablePagination.value.page = page ?? 1;
    qTablePagination.value.rowsPerPage = rowsPerPage ?? 10;
    qTablePagination.value.sortBy = sortBy ?? 'timestamp';
    qTablePagination.value.descending = descending ?? true;
    // The watcher on qTablePagination will sync these to the Pinia store
  }
}

// Placeholder for export functionality
function exportLogs(){
  console.log("Export logs clicked");
  // Actual export logic would go here
}

function setDateTimeProxiesFromFilters() {
  if (localFilters.value.dateFrom) {
    const dFrom = new Date(localFilters.value.dateFrom);
    dateFromProxy.value = quasarDateUtils.formatDate(dFrom, 'YYYY/MM/DD');
    timeFromProxy.value = quasarDateUtils.formatDate(dFrom, 'HH:mm');
  } else {
    dateFromProxy.value = null;
    timeFromProxy.value = null;
  }
  if (localFilters.value.dateTo) {
    const dTo = new Date(localFilters.value.dateTo);
    dateToProxy.value = quasarDateUtils.formatDate(dTo, 'YYYY/MM/DD');
    timeToProxy.value = quasarDateUtils.formatDate(dTo, 'HH:mm');
  } else {
    dateToProxy.value = null;
    timeToProxy.value = null;
  }
}

</script>

<template>
  <q-page padding>
    <!-- Page Header -->
    <div class="row items-center justify-between q-mb-xs">
      <h1 class="text-h5 text-weight-medium">Activity Log</h1>
      <div class="row q-gutter-xs items-center">
        <q-btn flat round dense icon="refresh" @click="refreshData" aria-label="Refresh Logs" >
          <q-tooltip>Refresh Logs</q-tooltip>
        </q-btn>
        <q-btn unelevated color="primary" icon="file_download" label="Export" @click="exportLogs" aria-label="Export Logs" />
      </div>
    </div>

    <q-card class="q-mb-md shadow-1">
      <q-card-section class="q-pa-md">
        <!-- First Row of Filters -->
        <div class="row items-start q-col-gutter-md q-mb-sm">
          <div class="col-12 col-sm-6 col-md-3">
            <q-btn outline color="grey-7" class="fit text-body2" icon="event" :label="dateRangeDisplay" no-caps>
              <q-tooltip>Select Date Range</q-tooltip>
              <q-popup-proxy>
                <q-card class="q-pa-sm" style="min-width: 300px;">
                  <q-card-section class="q-gutter-y-sm">
                    <div class="text-subtitle2 q-mb-xs">Date From:</div>
                    <q-input filled dense v-model="dateFromProxy" mask="YYYY/MM/DD" label="Date" clearable>
                      <template v-slot:append>
                        <q-icon name="event" class="cursor-pointer">
                          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                            <q-date v-model="dateFromProxy" mask="YYYY/MM/DD" />
                          </q-popup-proxy>
                        </q-icon>
                      </template>
                    </q-input>
                    <q-input filled dense v-model="timeFromProxy" mask="HH:mm" label="Time" clearable format24h>
                      <template v-slot:append>
                        <q-icon name="access_time" class="cursor-pointer">
                          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                            <q-time v-model="timeFromProxy" mask="HH:mm" format24h />
                          </q-popup-proxy>
                        </q-icon>
                      </template>
                    </q-input>
                    <q-separator class="q-my-md"/>
                    <div class="text-subtitle2 q-mb-xs">Date To:</div>
                     <q-input filled dense v-model="dateToProxy" mask="YYYY/MM/DD" label="Date" clearable>
                      <template v-slot:append>
                        <q-icon name="event" class="cursor-pointer">
                          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                            <q-date v-model="dateToProxy" mask="YYYY/MM/DD" />
                          </q-popup-proxy>
                        </q-icon>
                      </template>
                    </q-input>
                    <q-input filled dense v-model="timeToProxy" mask="HH:mm" label="Time" clearable format24h>
                      <template v-slot:append>
                        <q-icon name="access_time" class="cursor-pointer">
                          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                            <q-time v-model="timeToProxy" mask="HH:mm" format24h />
                          </q-popup-proxy>
                        </q-icon>
                      </template>
                    </q-input>
                  </q-card-section>
                  <q-card-actions align="right">
                     <q-btn flat label="Clear Dates" @click="() => { localFilters.dateFrom = null; localFilters.dateTo = null; setDateTimeProxiesFromFilters(); }" v-close-popup />
                     <q-btn flat label="Done" color="primary" v-close-popup />
                  </q-card-actions>
                </q-card>
              </q-popup-proxy>
            </q-btn>
          </div>

          <div class="col-12 col-sm-6 col-md-3">
            <q-select
              outlined dense
              v-model="localFilters.selectedUserId"
              :options="displayUserOptions"
              label="User"
              option-value="id"
              option-label="fullName"
              emit-value map-options use-input clearable
              @filter="filterUserOptions"
              options-dense
            />
          </div>

          <div class="col-12 col-sm-6 col-md-3">
            <q-select
              outlined dense
              v-model="localFilters.selectedActionType"
              :options="actionTypesForFilter"
              label="Action Type"
              emit-value map-options
              options-dense
            />
          </div>
          
          <div class="col-12 col-sm-6 col-md-3">
            <q-input
              outlined dense
              v-model="localFilters.searchQuery"
              label="Search activities..."
              clearable
            >
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
        </div>

        <!-- Second Row of Filters -->
        <div class="row items-center justify-between q-col-gutter-x-md q-gutter-y-sm">
          <div class="col-auto row items-center q-gutter-x-md">
            <q-checkbox dense v-model="localFilters.showSystemActions" label="System Actions" />
            <q-checkbox dense v-model="localFilters.showSuccessful" label="Successful" />
            <q-checkbox dense v-model="localFilters.showFailed" label="Failed" />
          </div>
          <div class="col-auto">
            <q-btn unelevated color="primary" label="Apply Filters" @click="applyFilters" icon="filter_alt" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card class="shadow-1">
      <q-table
        flat
        :rows="paginatedLogs" 
        :columns="columns"
        row-key="id"
        v-model:pagination="qTablePagination"
        :loading="isLoading"
        :rows-per-page-options="[10, 20, 50, 100]"
        @request="onRequest"
        binary-state-sort
        class="w-full"
      >
        <!-- <template v-slot:top-left>
          <div class="text-h6 q-table__title">Logs</div>
        </template> -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props" auto-width>
            <q-chip
              :color="getStatusChipColor(props.row.status)"
              text-color="white"
              size="sm"
              square dense
            >
              {{ props.row.status }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-isSystemAction="props">
          <q-td :props="props" auto-width style="text-align: center;">
            <q-chip v-if="props.row.isSystemAction" color="grey-7" text-color="white" size="sm" square dense>System</q-chip>
            <span v-else>-</span>
          </q-td>
        </template>

        <template v-slot:no-data="{ icon, message }">
            <div class="full-width row flex-center text-grey-7 q-gutter-sm q-my-xl">
                <q-icon size="2em" :name="paginatedLogs.length === 0 && (localFilters.dateFrom || localFilters.dateTo || localFilters.searchQuery || localFilters.selectedActionType !== 'All Actions' || localFilters.selectedUserId || !localFilters.showSystemActions || !localFilters.showSuccessful || !localFilters.showFailed ) ? 'filter_b_and_w' : icon" /> 
                <span>
                {{ paginatedLogs.length === 0 && (localFilters.dateFrom || localFilters.dateTo || localFilters.searchQuery || localFilters.selectedActionType !== 'All Actions' || localFilters.selectedUserId || !localFilters.showSystemActions || !localFilters.showSuccessful || !localFilters.showFailed ) ? 'No logs match your filter criteria.' : message }}
                </span>
            </div>
        </template>
         <template v-slot:loading>
            <q-inner-loading showing color="primary" />
        </template>

      </q-table>
       <div v-if="logStore.error && !isLoading" class="text-center q-pa-md text-negative">
        Error: {{ logStore.error }}
  </div>
    </q-card>

  </q-page>
</template>

<style scoped>
.w-full {
  width: 100%;
}

/* Custom styling for q-table title if needed via slot */
.q-table__title {
    font-size: 1.15rem; /* Adjusted from text-h6 default for a bit more control */
    font-weight: 500;
}

.q-card .q-select {
  min-width: 180px; 
}

.q-btn.fit {
  width: 100%;
}

/* Ensure the popup for date range is not overly wide on small screens */
.q-popup-proxy .q-card {
  min-width: 280px; /* Adjust as needed */
  max-width: 350px;
}

</style>
