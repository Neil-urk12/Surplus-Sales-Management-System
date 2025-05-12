import { defineStore } from 'pinia';
import { useUsersStore } from './users';
import { activityLogsService } from '../services/activityLogsApi';
import type { ActivityLogEntry, UserSnippet, ActionType, LogFilters, ActionStatus } from '../types/logTypes';

export const actionTypeOptions: ActionType[] = ['All Actions', 'Created', 'Updated', 'Deleted', 'Login', 'Logout'];

// Interface for backend activity log format
interface BackendActivityLog {
  id: string;
  timestamp: string;
  user: string;
  action: string;
  details: string;
  status: string;
  isSystemAction: boolean;
  createdAt: string;
  updatedAt: string;
}

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
    totalLogs: 0,
    lastPage: 1,
    sortBy: 'timestamp',
    descending: true,
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
        // Date, user, and actionType filters are primarily handled by the backend
        // when fetchFilteredLogs is called with these filters.
        // Client-side filters applied here are those not passed to the backend,
        // or for additional refinement on the fetched dataset.

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
            // If neither showSuccessful nor showFailed is true, show all statuses
            statusMatch = true; 
        }

        // const actionTypeMatch = state.filters.selectedActionType === 'All Actions' || log.actionType === state.filters.selectedActionType;
        // const userMatch = !state.filters.selectedUserId || log.user.id === state.filters.selectedUserId;

        return systemActionMatch && statusMatch;
      });
    },
    paginatedLogs(state): ActivityLogEntry[] {
      const start = (state.currentPage - 1) * state.pageSize;
      const end = start + state.pageSize;
      
      // Get filtered logs
      const sortedLogs = [...this.filteredLogs];
      
      // Apply sorting based on q-table parameters
      const sortBy = state.sortBy || 'timestamp';
      const descending = state.descending !== undefined ? state.descending : true;
      
      sortedLogs.sort((a, b) => {
        let valueA, valueB;
        
        // Handle field accessor functions for complex fields
        if (sortBy === 'user') {
          valueA = a.user.fullName;
          valueB = b.user.fullName;
        } else if (sortBy === 'role') {
          valueA = a.user.role;
          valueB = b.user.role;
        } else if (sortBy === 'timestamp') {
          valueA = new Date(a.timestamp).getTime();
          valueB = new Date(b.timestamp).getTime();
        } else {
          // For other fields, access directly
          valueA = a[sortBy as keyof ActivityLogEntry];
          valueB = b[sortBy as keyof ActivityLogEntry];
        }
        
        // Compare values based on their types
        if (typeof valueA === 'string' && typeof valueB === 'string') {
          return descending ? valueB.localeCompare(valueA) : valueA.localeCompare(valueB);
        } else {
          // For numbers, dates, etc.
          return descending ? (valueB as number) - (valueA as number) : (valueA as number) - (valueB as number);
        }
      });
      
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
            role: 'admin' as const
          });
        }
        
        // Fetch logs from API
        const response = await activityLogsService.getLogs(this.currentPage, this.pageSize);
        
        // Map backend logs to frontend format if needed
        this.logs = this.mapBackendLogsToFrontend(response.data as unknown as BackendActivityLog[]);
        
        this.totalLogs = response.total;
        this.lastPage = response.last_page;
      } catch (e) {
        this.error = 'Failed to fetch activity logs.';
        console.error(e);
        this.logs = []; 
      } finally {
        this.isLoading = false;
      }
    },
    
    async fetchFilteredLogs() {
      this.isLoading = true;
      this.error = null;
      
      try {
        const filters: Record<string, string> = {};
        
        if (this.filters.selectedUserId) {
          filters.user = this.filters.selectedUserId;
        }
        
        if (this.filters.selectedActionType !== 'All Actions') {
          filters.action = this.filters.selectedActionType;
        }
        
        if (this.filters.dateFrom) {
          // Format with full ISO string to include time component
          filters.startDate = this.filters.dateFrom.toISOString();
        }
        
        if (this.filters.dateTo) {
          // Format with full ISO string to include time component
          filters.endDate = this.filters.dateTo.toISOString();
        }
        
        // We'll handle showSuccessful/showFailed logic on frontend since
        // backend might not have this exact filter format
        
        const response = await activityLogsService.getFilteredLogs(
          filters,
          this.currentPage,
          this.pageSize
        );
        
        // Map backend logs to frontend format
        this.logs = this.mapBackendLogsToFrontend(response.data as unknown as BackendActivityLog[]);
        
        this.totalLogs = response.total;
        this.lastPage = response.last_page;
      } catch (e) {
        this.error = 'Failed to fetch filtered activity logs.';
        console.error(e);
      } finally {
        this.isLoading = false;
      }
    },
    
    // Helper method to map backend log format to frontend format
    mapBackendLogsToFrontend(backendLogs: BackendActivityLog[]): ActivityLogEntry[] {
      return backendLogs.map(log => {
        // Find user in users array or create default user
        const userMatch = this.users.find(u => u.id === log.user);
        const user: UserSnippet = userMatch || {
          id: log.user,
          fullName: log.user,
          role: 'admin' as const
        };
        
        return {
          id: log.id,
          timestamp: new Date(log.timestamp),
          user,
          actionType: this.mapActionToActionType(log.action),
          details: log.details,
          status: log.status as ActionStatus,
          isSystemAction: log.isSystemAction
        };
      });
    },
    
    // Helper method to map backend action string to frontend ActionType
    mapActionToActionType(action: string): ActionType {
      // Check if the action is one of our predefined ActionTypes
      if (actionTypeOptions.includes(action as ActionType)) {
        return action as ActionType;
      }
      // Default to the first non-"All Actions" type if not recognized
      return actionTypeOptions.find(type => type !== 'All Actions') || 'Created';
    },
    
    async addLogEntry(entryData: Omit<ActivityLogEntry, 'id' | 'timestamp'>) {
      try {
        // Format log data for the API
        const logData = {
          user: entryData.user.id,
          action: entryData.actionType,
          details: entryData.details,
          status: entryData.status,
          isSystemAction: entryData.isSystemAction
        };
        
        // Send to API
        const success = await activityLogsService.createLog(logData);
        
        if (success) {
          // Optionally refresh logs after adding new entry
          await this.fetchLogs();
        }
      } catch (error) {
        console.error('Error adding log entry:', error);
        this.error = 'Failed to add log entry.';
      }
    },
    
    updateFilters(newFilters: Partial<LogFilters>) {
      this.filters = { ...this.filters, ...newFilters };
      this.currentPage = 1; 
      // When filters change, fetch filtered logs
      void this.fetchFilteredLogs();
    },
    
    async setCurrentPage(page: number) {
      this.currentPage = page;
      // Refetch logs for the new page
      if (Object.values(this.filters).some(v => v !== null && v !== undefined && v !== 'All Actions')) {
        await this.fetchFilteredLogs();
      } else {
        await this.fetchLogs();
      }
    },
    
    async setPageSize(size: number) {
      this.pageSize = size;
      this.currentPage = 1;
      // Refetch logs with new page size
      if (Object.values(this.filters).some(v => v !== null && v !== undefined && v !== 'All Actions')) {
        await this.fetchFilteredLogs();
      } else {
        await this.fetchLogs();
      }
    },
    
    refreshLogs() {
      // Refresh users and logs
      const usersStore = useUsersStore();
      void usersStore.fetchUsers().then(() => {
        if (Object.values(this.filters).some(v => v !== null && v !== undefined && v !== 'All Actions')) {
          void this.fetchFilteredLogs();
        } else {
          void this.fetchLogs();
        }
      });
    },
    
    setSorting(sortBy: string, descending: boolean) {
      this.sortBy = sortBy;
      this.descending = descending;
    },
  },
});
