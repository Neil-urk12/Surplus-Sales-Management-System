export type ActionStatus = 'successful' | 'failed';
export type ActionType = 'All Actions' | 'Created' | 'Updated' | 'Deleted' | 'Login' | 'Logout';

export interface UserSnippet {
  id: string;
  fullName: string;
  role: 'admin' | 'staff';
}

export interface ActivityLogEntry {
  id: string;
  timestamp: Date;
  user: UserSnippet;
  actionType: ActionType;
  details: string;
  status: ActionStatus;
  isSystemAction: boolean;
}

export interface LogFilters {
  dateFrom: Date | null;
  dateTo: Date | null;
  showSystemActions: boolean;
  showSuccessful: boolean;
  showFailed: boolean;
  selectedActionType: ActionType;
  selectedUserId: string | null;
  searchQuery?: string | null;
}
