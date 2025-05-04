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

<<<<<<< HEAD
type ItemType = string;

export const operationNotifications = {
  add: {
    success: (itemName: ItemType) => showSuccessNotification({
      message: `Added new ${itemName}`,
    }),
    error: (itemName: ItemType) => showErrorNotification({
=======
// Predefined notifications for common operations
export const operationNotifications = {
  add: {
    success: (itemName: string) => showSuccessNotification({
      message: `Added new ${itemName}`,
    }),
    error: (itemName: string) => showErrorNotification({
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      message: `Failed to add ${itemName}`,
    }),
  },
  update: {
<<<<<<< HEAD
    success: (itemName: ItemType) => showSuccessNotification({
      message: `Updated ${itemName}`,
    }),
    error: (itemName: ItemType) => showErrorNotification({
=======
    success: (itemName: string) => showSuccessNotification({
      message: `Updated ${itemName}`,
    }),
    error: (itemName: string) => showErrorNotification({
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      message: `Failed to update ${itemName}`,
    }),
  },
  delete: {
<<<<<<< HEAD
    success: (itemName: ItemType) => showSuccessNotification({
      message: `Successfully deleted ${itemName}`,
    }),
    error: (itemName: ItemType) => showErrorNotification({
=======
    success: (itemName: string) => showSuccessNotification({
      message: `Successfully deleted ${itemName}`,
    }),
    error: (itemName: string) => showErrorNotification({
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
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