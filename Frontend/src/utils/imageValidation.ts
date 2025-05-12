export const MAX_IMAGE_SIZE = 5 * 1024 * 1024;
export const DEFAULT_MAX_DIMENSION = 4096; // Default maximum dimension

// Valid image MIME types
export const VALID_IMAGE_TYPES = [
  'image/jpeg',
  'image/png',
  'image/gif',
  'image/webp',
  'image/svg+xml'
];

interface ImageValidationResult {
  isValid: boolean;
  error?: string;
  sanitizedData?: string;
}

/**
 * Validates and sanitizes a base64 image string
 * @param base64String The base64 image string to validate
 * @returns Validation result with sanitized data if valid
 */
export function validateAndSanitizeBase64Image(
  base64String: string
): ImageValidationResult {
  // Check if the string is empty
  if (!base64String) {
    return { isValid: false, error: 'No image data provided' };
  }

  // Check if it's a valid base64 image format
  const base64Regex = /^data:image\/[a-zA-Z+]+;base64,/;
  if (!base64Regex.test(base64String)) {
    return { isValid: false, error: 'Invalid image format' };
  }

  try {
    // Extract MIME type safely
    const parts = base64String.split(';');
    if (!parts?.[0]) {
      return { isValid: false, error: 'Invalid image format' };
    }

    const mimeParts = parts[0].split(':');
    if (!mimeParts?.[1]) {
      return { isValid: false, error: 'Invalid image format' };
    }

    const mimeType = mimeParts[1];
    if (!VALID_IMAGE_TYPES.includes(mimeType)) {
      return { isValid: false, error: 'Unsupported image type' };
    }

    // Remove the data URL prefix to get the actual base64 data
    const base64Parts = base64String.split(',');
    if (!base64Parts?.[1]) {
      return { isValid: false, error: 'Invalid base64 data' };
    }
    const base64Data = base64Parts[1];

    // Try to decode a small portion of the base64 data to verify it's valid
    try {
      const sampleSize = Math.min(100, base64Data.length);
      const sample = base64Data.slice(0, sampleSize);
      atob(sample);
    } catch {
      return { isValid: false, error: 'Invalid base64 encoding' };
    }

    // If we got here, the base64 data appears valid
    // Return the original data to avoid potential corruption from re-encoding
    return {
      isValid: true,
      sanitizedData: base64String
    };
  } catch {
    return {
      isValid: false,
      error: 'Invalid base64 encoding'
    };
  }
}

/**
 * Asynchronously validates image dimensions from a base64 string
 * @param base64String The base64 image string to validate
 * @param maxDimension The maximum allowed dimension (width/height) for the image
 * @returns Promise resolving to validation result
 */
export async function validateImageDimensions(
  base64String: string,
  maxDimension: number = DEFAULT_MAX_DIMENSION
): Promise<ImageValidationResult> {
  if (!maxDimension) {
    return { isValid: true };
  }
  
  return new Promise((resolve) => {
    const img = new Image();
    
    img.onload = () => {
      if (img.width > maxDimension || img.height > maxDimension) {
        resolve({
          isValid: false,
          error: `Image dimensions exceed maximum limit of ${maxDimension}px`
        });
      } else {
        resolve({ isValid: true });
      }
    };
    
    img.onerror = () => {
      resolve({
        isValid: false,
        error: 'Failed to load image for dimension validation'
      });
    };
    
    img.src = base64String;
  });
}

/**
 * Validates a File object for image upload
 * @param file The File object to validate
 * @returns Validation result
 */
export function validateImageFile(file: File): ImageValidationResult {
  if (!file) {
    return { isValid: false, error: 'No file provided' };
  }

  if (!VALID_IMAGE_TYPES.includes(file.type)) {
    return { isValid: false, error: 'Unsupported image type' };
  }

  if (file.size > MAX_IMAGE_SIZE) {
    return {
      isValid: false,
      error: `File size exceeds maximum limit of ${MAX_IMAGE_SIZE / (1024 * 1024)}MB`
    };
  }

  return { isValid: true };
}

/**
 * Comprehensive validation for file uploads
 * Validates file type, size, and optionally dimensions in one function
 * 
 * @param file The File object to validate
 * @param options Optional configuration parameters
 * @returns Promise resolving to validation result with sanitized data if valid
 */
export async function validateFileUpload(
  file: File,
  options: {
    maxFileSize?: number,
    maxDimension?: number,
    validateDimensions?: boolean
  } = {}
): Promise<ImageValidationResult> {
  // Set default options
  const maxFileSize = options.maxFileSize || MAX_IMAGE_SIZE;
  const maxDimension = options.maxDimension || DEFAULT_MAX_DIMENSION;
  const validateDimensions = options.validateDimensions !== undefined ? options.validateDimensions : true;
  
  // First validate the file itself
  const fileValidation = validateImageFile(file);
  if (!fileValidation.isValid) {
    return fileValidation;
  }
  
  // If file size exceeds the limit
  if (file.size > maxFileSize) {
    return {
      isValid: false,
      error: `File size exceeds maximum limit of ${maxFileSize / (1024 * 1024)}MB`
    };
  }
  
  try {
    // Read the file as data URL
    const base64String = await readFileAsDataURL(file);
    
    // Validate and sanitize the base64 data
    const validationResult = validateAndSanitizeBase64Image(base64String);
    if (!validationResult.isValid || !validationResult.sanitizedData) {
      return validationResult;
    }
    
    // Optionally validate dimensions
    if (validateDimensions) {
      const dimensionResult = await validateImageDimensions(validationResult.sanitizedData, maxDimension);
      if (!dimensionResult.isValid) {
        return dimensionResult;
      }
    }
    
    // All validations passed
    return {
      isValid: true,
      sanitizedData: validationResult.sanitizedData
    };
  } catch (error) {
    console.error('Error during file validation:', error);
    return {
      isValid: false,
      error: error instanceof Error ? error.message : 'Unknown error during file validation'
    };
  }
}

/**
 * Reads a File object as a data URL
 * @param file The File object to read
 * @returns Promise resolving to the data URL string
 */
function readFileAsDataURL(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const result = e.target?.result as string;
      if (result) {
        resolve(result);
      } else {
        reject(new Error('Failed to read file'));
      }
    };
    reader.onerror = () => reject(new Error('Error reading file'));
    reader.readAsDataURL(file);
  });
}
