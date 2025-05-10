/**
 * Represents a customer in the system.
 */
export interface Customer {
    /** Unique identifier for the customer */
    id: string;

    /** Full name of the customer */
    fullName: string;

    /** Email address of the customer */
    email: string;

    /** Phone number of the customer */
    phone: string;

    /** Physical address of the customer */
    address: string;

    /** The date when the customer was registered, in ISO string format */
    dateRegistered: string;

    /** Timestamp of when the customer record was created, in ISO string format */
    createdAt: string;

    /** Timestamp of when the customer record was last updated, in ISO string format */
    updatedAt: string;
} 
