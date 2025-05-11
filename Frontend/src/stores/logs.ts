import { defineStore } from 'pinia';
import { useUsersStore } from './users';
import type { ActivityLogEntry, UserSnippet, ActionType, LogFilters } from '../types/logTypes';

export const actionTypeOptions: ActionType[] = ['All Actions', 'Created', 'Logged In', 'Updated', 'Deleted', 'Login', 'Logout'];

// For generating unique IDs for new logs, very basic for now.
let logIdCounter = 0;

export const useLogStore = defineStore('activityLogs', {
  state: () => ({
    logs: [] as ActivityLogEntry[],
    users: [] as UserSnippet[], // This will be populated from useUsersStore
    allActionTypes: [...actionTypeOptions] as ActionType[],
    filters: {
      dateFrom: null,
      dateTo: null,
      showSystemActions: true,
      showSuccessful: true,
      showFailed: true,
      selectedActionType: 'All Actions',
      selectedUserId: null,
      searchQuery: null, 
    } as LogFilters,
    isLoading: false,
    error: null as string | null,
    currentPage: 1,
    pageSize: 10,
  }),

  getters: {
    filteredLogs(state): ActivityLogEntry[] {
      let logsToFilter = [...state.logs];

      if (state.filters.searchQuery && state.filters.searchQuery.trim() !== '') {
        const query = state.filters.searchQuery.toLowerCase().trim();
        logsToFilter = logsToFilter.filter(log => 
          log.user.fullName.toLowerCase().includes(query) ||
          log.user.role.toLowerCase().includes(query) ||
          log.actionType.toLowerCase().includes(query) ||
          log.details.toLowerCase().includes(query)
        );
      }

      return logsToFilter.filter(log => {
        const dateFromMatch = !state.filters.dateFrom || new Date(log.timestamp) >= new Date(state.filters.dateFrom);
        const dateToMatch = !state.filters.dateTo || new Date(log.timestamp) <= new Date(state.filters.dateTo);
        const systemActionMatch = state.filters.showSystemActions || !log.isSystemAction;
        
        let statusMatch = false;
        const successfulPass = state.filters.showSuccessful ? log.status === 'successful' : false;
        const failedPass = state.filters.showFailed ? log.status === 'failed' : false;

        if(state.filters.showSuccessful && state.filters.showFailed) {
            statusMatch = successfulPass || failedPass;
        } else if (state.filters.showSuccessful) {
            statusMatch = successfulPass;
        } else if (state.filters.showFailed) {
            statusMatch = failedPass;
        } else {
            statusMatch = true; 
        }

        const actionTypeMatch = state.filters.selectedActionType === 'All Actions' || log.actionType === state.filters.selectedActionType;
        const userMatch = !state.filters.selectedUserId || log.user.id === state.filters.selectedUserId;

        return dateFromMatch && dateToMatch && systemActionMatch && statusMatch && actionTypeMatch && userMatch;
      });
    },
    paginatedLogs(state): ActivityLogEntry[] {
      const start = (state.currentPage - 1) * state.pageSize;
      const end = start + state.pageSize;
      const sortedLogs = [...this.filteredLogs].sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
      return sortedLogs.slice(start, end);
    },
    totalFilteredLogs(): number {
      return this.filteredLogs.length;
    },
  },

  actions: {
    async fetchLogs() {
      this.isLoading = true;
      this.error = null;
      try {
        // Fetch users from the users store
        const usersStore = useUsersStore();
        if (usersStore.users.length === 0) {
          await usersStore.fetchUsers();
        }
        
        // Map users to UserSnippet format
        this.users = usersStore.users.map(user => ({
          id: user.id,
          fullName: user.fullName,
          role: user.role
        }));
        
        // Add system user if it doesn't exist
        if (!this.users.some(u => u.id === 'system')) {
          this.users.push({
            id: 'system',
            fullName: 'System',
            role: 'admin'
          });
        }
        
        // Simulate API call delay if needed, but now initializes logs as empty.
        await new Promise(resolve => setTimeout(resolve, 200)); 
        this.logs = []; // Initialize logs as an empty array
      } catch (e) {
        this.error = 'Failed to initialize activity logs.';
        console.error(e);
        this.logs = []; // Ensure logs is empty on error too
      } finally {
        this.isLoading = false;
      }
    },
    addLogEntry(entryData: Omit<ActivityLogEntry, 'id' | 'timestamp'>) {
      const newEntry: ActivityLogEntry = {
        ...entryData,
        id: `log_${logIdCounter++}`,
        timestamp: new Date(),
      };
      this.logs.unshift(newEntry); // Add to the beginning of the array
    },
    updateFilters(newFilters: Partial<LogFilters>) {
      this.filters = { ...this.filters, ...newFilters };
      this.currentPage = 1; 
    },
    setCurrentPage(page: number) {
      this.currentPage = page;
    },
    setPageSize(size: number) {
      this.pageSize = size;
      this.currentPage = 1; 
    },
    refreshLogs() {
      // This will now clear current logs and prepare for new ones or fetched ones.
      // Refresh users as well to ensure up-to-date user data
      const usersStore = useUsersStore();
      void usersStore.fetchUsers().then(() => {
        void this.fetchLogs();
      });
    },
  },
});
