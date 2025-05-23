import { Notify } from 'quasar';
import type { QNotifyCreateOptions } from 'quasar';

interface NotifyOptions extends Partial<QNotifyCreateOptions> {
  message: string;
  color?: 'positive' | 'negative' | 'warning' | 'info';
  timeout?: number;
  position?: 'top' | 'bottom' | 'left' | 'right' | 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right' | 'center';
}

const defaultOptions: Partial<NotifyOptions> = {
  position: 'top',
  timeout: 2000,
};

export function showSuccessNotification(options: NotifyOptions) {
  Notify.create({
    ...defaultOptions,
    color: 'positive',
    ...options,
  });
}

export function showErrorNotification(options: NotifyOptions) {
  Notify.create({
    ...defaultOptions,
    color: 'negative',
    ...options,
  });
}

export function showWarningNotification(options: NotifyOptions) {
  Notify.create({
    ...defaultOptions,
    color: 'warning',
    ...options,
  });
}

export function showInfoNotification(options: NotifyOptions) {
  Notify.create({
    ...defaultOptions,
    color: 'info',
    ...options,
  });
}

type ItemType = string;

export const operationNotifications = {
  add: {
    success: (itemName: ItemType) => showSuccessNotification({
      message: `Added new ${itemName}`,
    }),
    error: (itemName: ItemType) => showErrorNotification({
      message: `Failed to add ${itemName}`,
    }),
  },
  update: {
    success: (itemName: ItemType) => showSuccessNotification({
      message: `Updated ${itemName}`,
    }),
    error: (itemName: ItemType) => showErrorNotification({
      message: `Failed to update ${itemName}`,
    }),
  },
  delete: {
    success: (itemName: ItemType) => showSuccessNotification({
      message: `Successfully deleted ${itemName}`,
    }),
    error: (itemName: ItemType) => showErrorNotification({
      message: `Failed to delete ${itemName}`,
    }),
  },
  validation: {
    error: (message: string) => showErrorNotification({
      message,
    }),
    warning: (message: string) => showWarningNotification({
      message,
    }),
  },
  filters: {
    success: () => showSuccessNotification({
      message: 'Filters applied successfully',
    }),
  },
}; 