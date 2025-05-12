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

// Refs for two-step date-time selection process
const dateTimeFromInput = ref<string | null>(null);
const dateTimeToInput = ref<string | null>(null);

// Refs to track the current step in date-time selection
const showFromTimeSelection = ref(false);
const showToTimeSelection = ref(false);

// Dialog visibility state
const fromDateDialog = ref(false);
const toDateDialog = ref(false);

// Sync localFilters with store filters initially and when store filters might change externally
// Also, pre-fill date/time inputs if store has values
watch(() => logStore.filters, (newStoreFilters) => {
  localFilters.value = {
    ...newStoreFilters,
    searchQuery: localFilters.value.searchQuery
   };
  setDateTimeInputsFromFilters();
}, { deep: true, immediate: true });

onMounted(async () => {
  try {
    await logStore.fetchLogs();
    setDateTimeInputsFromFilters(); // Initialize inputs from store or defaults
  } catch (error) {
    console.error('Error fetching logs:', error);
  }
});

// Watchers for the new combined date-time inputs
watch(dateTimeFromInput, (newVal) => {
  if (newVal) {
    try {
      const parsedDate = quasarDateUtils.extractDate(newVal, 'YYYY-MM-DD HH:mm');
      localFilters.value.dateFrom = parsedDate;
    } catch (e) {
      console.error('Error parsing date from:', e);
      localFilters.value.dateFrom = null;
    }
  } else {
    localFilters.value.dateFrom = null;
  }
});

watch(dateTimeToInput, (newVal) => {
  if (newVal) {
    try {
      const parsedDate = quasarDateUtils.extractDate(newVal, 'YYYY-MM-DD HH:mm');
      localFilters.value.dateTo = parsedDate;
    } catch (e) {
      console.error('Error parsing date to:', e);
      localFilters.value.dateTo = null;
    }
  } else {
    localFilters.value.dateTo = null;
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
  console.log('applyFilters called. Local filters:', localFilters.value);
  // The watchers for dateTimeFromInput and dateTimeToInput already update localFilters.dateFrom and localFilters.dateTo
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
  setDateTimeInputsFromFilters(); // This was the intended call
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
    
    // Update the sorting in the store
    if (sortBy) logStore.setSorting(sortBy, descending ?? true);
  }
}

function exportLogs(){
  console.log("Export logs clicked");
}

function setDateTimeInputsFromFilters() {
  // Set date-time inputs from filters (or defaults)
  if (localFilters.value.dateFrom) {
    try {
      // Format the date to match the mask
      dateTimeFromInput.value = quasarDateUtils.formatDate(
        localFilters.value.dateFrom,
        'YYYY-MM-DD HH:mm'
      );
      // Reset time selection flag when loading from filters
      showFromTimeSelection.value = false;
    } catch (e) {
      console.error('Error formatting dateFrom:', e);
      dateTimeFromInput.value = null;
    }
  } else {
    dateTimeFromInput.value = null;
  }

  if (localFilters.value.dateTo) {
    try {
      // Format the date to match the mask
      dateTimeToInput.value = quasarDateUtils.formatDate(
        localFilters.value.dateTo,
        'YYYY-MM-DD HH:mm'
      );
      // Reset time selection flag when loading from filters
      showToTimeSelection.value = false;
    } catch (e) {
      console.error('Error formatting dateTo:', e);
      dateTimeToInput.value = null;
    }
  } else {
    dateTimeToInput.value = null;
  }
}

function handleFromDateTimeDone() {
  showFromTimeSelection.value = false; // Reset to show date picker next time
  fromDateDialog.value = false;      // Close the dialog
}

function handleToDateTimeDone() {
  showToTimeSelection.value = false;   // Reset to show date picker next time
  toDateDialog.value = false;        // Close the dialog
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
        <q-btn unelevated icon="file_download" label="Export" @click="exportLogs" aria-label="Export Logs" 
        :class="[
        $q.dark.isActive ? 'bg-transparent text-white' : 'bg-primary text-white'
        ]" 
        style="border:1px solid white"
        />
      </div>
    </div>

    <q-card class="q-mb-md shadow-1">
      <q-card-section class="q-pa-md">
        <!-- First Row of Filters -->
        <div class="row items-start q-col-gutter-md q-mb-sm">
          <div class="col-12 col-sm-6 col-md-3">
            <q-input
              outlined dense
              v-model="dateTimeFromInput"
              label="Date From"
              placeholder="YYYY-MM-DD HH:mm"
              clearable
            >
              <template v-slot:prepend>
                <q-icon name="event" class="cursor-pointer" @click="fromDateDialog = true" />
              </template>
            </q-input>
          </div>

          <div class="col-12 col-sm-6 col-md-3">
            <q-input
              outlined dense
              v-model="dateTimeToInput"
              label="Date To"
              placeholder="YYYY-MM-DD HH:mm"
              clearable
            >
              <template v-slot:prepend>
                <q-icon name="event" class="cursor-pointer" @click="toDateDialog = true" />
              </template>
            </q-input>
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
            <q-btn unelevated :class="[
              $q.dark.isActive ? 'bg-transparent text-white' : 'bg-primary text-white'
              ]" label="Apply Filters" @click="applyFilters" icon="filter_alt" 
              style="border:1px solid white"
              />
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

    <!-- Date From Dialog -->
    <q-dialog v-model="fromDateDialog">
      <q-card style="min-width: 580px;">
        <q-card-section class="q-pa-md">
          <div v-if="!showFromTimeSelection">
            <h6 class="q-my-none q-mb-md text-center">Select Date</h6>
            <q-date 
              :landscape="true"
              v-model="dateTimeFromInput" 
              mask="YYYY-MM-DD HH:mm"
              today-btn
              calendar-type="gregorian"
              class="q-mx-auto"
            />
            <div class="row items-center justify-end q-gutter-sm q-mt-md">
              <q-btn flat label="Cancel" :color="$q.dark.isActive ? 'text-black' : 'text-white'" v-close-popup @click="showFromTimeSelection = false" />
              <q-btn @click="showFromTimeSelection = true" label="Next" :class="[
              $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white'
              ]" />
            </div>
          </div>
          <div v-else>
            <h6 class="q-my-none q-mb-md text-center">Select Time</h6>
            <q-time 
              v-model="dateTimeFromInput"
              mask="YYYY-MM-DD HH:mm"
              format24h
              :landscape="true"
              now-btn
              class="q-mx-auto"
            />
            <div class="row items-center justify-between q-mt-md">
              <q-btn @click="showFromTimeSelection = false" label="Back" :color="$q.dark.isActive ? 'text-black' : 'text-white'" flat />
              <q-btn label="Done" 
              :class="[
              $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white'
              ]"
              @click="handleFromDateTimeDone" />
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Date To Dialog -->
    <q-dialog v-model="toDateDialog">
      <q-card style="min-width: 580px;">
        <q-card-section class="q-pa-md">
          <div v-if="!showToTimeSelection">
            <h6 class="q-my-none q-mb-md text-center">Select Date</h6>
            <q-date 
              :landscape="true"
              v-model="dateTimeToInput" 
              mask="YYYY-MM-DD HH:mm"
              today-btn
              calendar-type="gregorian"
              class="q-mx-auto"
            />
            <div class="row items-center justify-end q-gutter-sm q-mt-md">
              <q-btn flat label="Cancel" :color="$q.dark.isActive ? 'text-black' : 'text-white'" v-close-popup @click="showToTimeSelection = false" />
              <q-btn @click="showToTimeSelection = true" label="Next" :class="[
              $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white'
              ]" />
            </div>
          </div>
          <div v-else>
            <h6 class="q-my-none q-mb-md text-center">Select Time</h6>
            <q-time 
              v-model="dateTimeToInput"
              mask="YYYY-MM-DD HH:mm"
              format24h
              :landscape="true"
              now-btn
              class="q-mx-auto"
            />
            <div class="row items-center justify-between q-mt-md">
              <q-btn @click="showToTimeSelection = false" label="Back" color="grey-7" flat />
              <q-btn label="Done" :class="[
              $q.dark.isActive ? 'bg-white text-black' : 'bg-primary text-white'
              ]" @click="handleToDateTimeDone" />
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

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

</style>
