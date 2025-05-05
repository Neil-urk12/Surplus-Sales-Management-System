/**
 * Utility functions for exporting data
 */

/**
 * Convert an array of objects to CSV format and trigger download
 * @param data Array of objects to convert
 * @param filename Filename for the download (without extension)
 * @param columns Optional column configuration for custom headers and field mapping
 */
export function exportToCsv<T extends Record<string, string | number | boolean | null | undefined>>(
    data: T[],
    filename: string,
    columns?: Array<{ header: string; field: keyof T }>
): void {
    if (!data || data.length === 0) {
        console.warn('No data to export');
        return;
    }

    // Determine fields to export
    let fieldsToExport: Array<{ header: string; field: keyof T }>;

    if (columns && columns.length > 0) {
        // Use provided columns
        fieldsToExport = columns;
    } else {
        // Auto-generate columns from first data item
        // We know data[0] exists because we checked data.length above
        // Use type assertion to tell TypeScript this is a valid object
        const keys = Object.keys(data[0] as object);
        fieldsToExport = keys.map(key => ({
            header: key.charAt(0).toUpperCase() + key.slice(1).replace(/_/g, ' '),
            field: key as keyof T
        }));
    }

    // Create CSV header row
    const csvHeader = fieldsToExport.map(col => `"${col.header}"`).join(',');

    // Convert data rows to CSV
    const csvRows = data.map(item => {
        const values = fieldsToExport.map(col => {
            const value = item[col.field];
            // Handle different data types
            if (value === null || value === undefined) return '""';
            if (typeof value === 'string') return `"${value.replace(/"/g, '""')}"`;
            if (typeof value === 'number' || typeof value === 'boolean') return value;
            // Handle objects, dates, etc.
            return `"${String(value).replace(/"/g, '""')}"`;
        });
        return values.join(',');
    });

    // Combine header and rows
    const csvContent = [csvHeader, ...csvRows].join('\n');

    // Create download link
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.setAttribute('href', url);
    link.setAttribute('download', `${filename}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
} 
