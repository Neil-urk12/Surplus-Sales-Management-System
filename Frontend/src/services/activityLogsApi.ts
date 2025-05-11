/**
 * API service for activity logs operations
 */
import { api } from 'src/boot/axios';
import { apiService } from './api';
import type { ActivityLogEntry } from 'src/types/logTypes';

interface PaginatedLogsResponse {
  data: ActivityLogEntry[];
  total: number;
  page: number;
  last_page: number;
}

/**
 * Service for handling activity logs-related API calls
 */
export const activityLogsService = {
  /**
   * Get paginated activity logs
   * @param page - Page number for pagination
   * @param limit - Number of logs per page
   * @returns Promise with paginated logs response
   */
  async getLogs(page = 1, limit = 10): Promise<PaginatedLogsResponse> {
    try {
      return await apiService.get<PaginatedLogsResponse>(
        '/api/activity-logs',
        { page, limit }
      );
    } catch (error) {
      console.error('Error fetching activity logs:', error);
      return {
        data: [],
        total: 0,
        page: 1,
        last_page: 1
      };
    }
  },

  /**
   * Get filtered activity logs based on various criteria
   * @param filters - Object containing filter parameters
   * @param page - Page number for pagination
   * @param limit - Number of logs per page
   * @returns Promise with filtered paginated logs response
   */
  async getFilteredLogs(
    filters: {
      user?: string;
      action?: string;
      status?: string;
      startDate?: string;
      endDate?: string;
    },
    page = 1,
    limit = 10
  ): Promise<PaginatedLogsResponse> {
    try {
      const params = {
        page,
        limit,
        ...filters
      };
      
      return await apiService.get<PaginatedLogsResponse>(
        '/api/activity-logs/filter',
        params
      );
    } catch (error) {
      console.error('Error fetching filtered activity logs:', error);
      return {
        data: [],
        total: 0,
        page: 1,
        last_page: 1
      };
    }
  },

  /**
   * Create a new activity log entry
   * @param logData - Data for the new log entry
   * @returns Promise resolving to boolean indicating success
   */
  async createLog(logData: {
    user: string;
    action: string;
    details: string;
    status: string;
    isSystemAction?: boolean;
  }): Promise<boolean> {
    try {
      await api.post('/api/activity-logs', logData);
      return true;
    } catch (error) {
      console.error('Error creating activity log:', error);
      return false;
    }
  }
};
