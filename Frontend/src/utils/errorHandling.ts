import { Notify } from 'quasar';
import type { QNotifyCreateOptions } from 'quasar';
import { operationNotifications } from './notifications';

export type ErrorType = 'validation' | 'database' | 'inventory' | 'purchase' | 'network' | 'auth' | 'unknown';

export class AppError extends Error {
    public readonly type: ErrorType;
    public readonly userMessage: string;
    public readonly severity: 'info' | 'warning' | 'error' | 'critical';

    constructor(
        message: string,
        type: ErrorType = 'unknown',
        userMessage?: string,
        severity: 'info' | 'warning' | 'error' | 'critical' = 'error'
    ) {
        super(message);
        this.name = 'AppError';
        this.type = type;
        this.userMessage = userMessage || message;
        this.severity = severity;
    }
}

// Centralized error handling service
export const errorHandler = {
    handle(error: unknown, context: string = 'operation'): { message: string; type: ErrorType } {
        console.error(`Error in ${context}:`, error);

        let errorMessage = 'An unexpected error occurred';
        let errorType: ErrorType = 'unknown';

        if (error instanceof AppError) {
            errorMessage = error.userMessage;
            errorType = error.type;

            // Display notification based on error
            this.showErrorNotification(error);

            // For critical errors, log to monitoring/analytics system
            if (error.severity === 'critical') {
                // TODO: implement logging to external monitoring system
                console.error('CRITICAL ERROR:', error);
            }
        } else if (error instanceof Error) {
            errorMessage = error.message;

            // Basic notification for standard errors
            Notify.create({
                type: 'negative',
                message: errorMessage,
                timeout: 3000
            });
        }

        return { message: errorMessage, type: errorType };
    },

    showErrorNotification(error: AppError): void {
        const options: QNotifyCreateOptions = {
            message: error.userMessage,
            timeout: error.severity === 'critical' ? 0 : 5000,
        };

        switch (error.severity) {
            case 'info':
                options.type = 'info';
                break;
            case 'warning':
                options.type = 'warning';
                break;
            case 'error':
            case 'critical':
                options.type = 'negative';
                options.caption = `Error type: ${error.type}`;
                break;
        }

        Notify.create(options);
    },

    // Helper method for handling common operation errors
    handleOperation(
        error: unknown,
        operation: 'add' | 'update' | 'delete' | 'fetch',
        itemType: string
    ): { message: string; type: ErrorType } {
        const result = this.handle(error, `${operation} ${itemType}`);

        switch (operation) {
            case 'add':
                operationNotifications.add.error(result.message);
                break;
            case 'update':
                operationNotifications.update.error(result.message);
                break;
            case 'delete':
                operationNotifications.delete.error(result.message);
                break;
            case 'fetch':
                Notify.create({
                    type: 'negative',
                    message: `Failed to fetch ${itemType}: ${result.message}`,
                    timeout: 3000
                });
                break;
        }
        return result;
    },

    // Recovery methods for specific error types
    recoverFromInventoryError(
        refreshCallbacks: Array<() => Promise<void>>
    ): Promise<void> {
        Notify.create({
            type: 'info',
            message: 'Attempting to refresh inventory data...',
            timeout: 2000
        });

        return Promise.all(refreshCallbacks.map(cb => cb()))
            .then(() => {
                Notify.create({
                    type: 'positive',
                    message: 'Inventory data refreshed successfully',
                    timeout: 2000
                });
            })
            .catch(refreshError => {
                console.error('Error refreshing data:', refreshError);
                Notify.create({
                    type: 'negative',
                    message: 'Failed to refresh inventory data',
                    timeout: 3000
                });
            });
    }
}; 
