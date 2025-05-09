/**
 * UI text messages used throughout the application
 * Centralizing these messages makes it easier to:
 * - Maintain consistent wording
 * - Update text in one place
 * - Implement internationalization in the future
 */

/**
 * Alert-related messages
 */
export const alertMessages = {
  /**
   * Messages for notification sections
   */
  sections: {
    unreadNotifications: 'Unread Notifications',
    readNotifications: 'Read Notifications',
    systemAlerts: 'System Alerts',
  },
  /**
   * Empty state messages
   */
  empty: {
    noUnreadNotifications: 'No unread notifications',
    noReadNotifications: 'No read notifications',
    noActiveNotifications: 'No active notifications',
  },
  /**
   * Action messages
   */
  actions: {
    markAllAsRead: 'Mark all as read',
    markAsRead: 'Mark as read',
    viewAll: (count: number) => `View all ${count} read notifications`,
  },
};
