<script setup lang="ts">
<<<<<<< HEAD
<<<<<<< HEAD
import { ref, watch, computed, onMounted, onUnmounted } from 'vue';
=======
import { ref, watch, computed } from 'vue';
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
import { ref, watch, computed, onMounted, onUnmounted } from 'vue';
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
import type { QTableColumn, QTableProps } from 'quasar';
import ProductCardModal from 'src/components/Global/ProductModal.vue'
import { useQuasar } from 'quasar';
import { useCabsStore } from 'src/stores/cabs';
<<<<<<< HEAD
<<<<<<< HEAD
import type { CabsRow, NewCabInput } from 'src/types/cabs';
import { getFirstFallbackImage, getNextFallbackImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
=======
import type { CabsRow } from 'src/stores/cabs';
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
import type { CabsRow, NewCabInput } from 'src/types/cabs';
import { getDefaultImage, getNextFallbackImage } from 'src/config/defaultImages';
import { validateAndSanitizeBase64Image } from '../utils/imageValidation';
import { operationNotifications } from '../utils/notifications';
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)

const $q = useQuasar();
const store = useCabsStore();
const showFilterDialog = ref(false);
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const cabToDelete = ref<CabsRow | null>(null);
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
const showProductCardModal = ref(false);
const isDragging = ref(false);

// Add global drag event handlers
onMounted(() => {
  const handleGlobalDragEnd = () => {
    isDragging.value = false;
  };

  document.addEventListener('dragend', handleGlobalDragEnd);

  // Clean up on unmount
  onUnmounted(() => {
    document.removeEventListener('dragend', handleGlobalDragEnd);
  });
});

// Enhanced drag event handlers
function handleDragLeave(event: DragEvent) {
  // Check if the mouse left the container (not just moved between child elements)
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;

  // Check if the mouse is outside the container's bounds
  if (x <= rect.left || x >= rect.right || y <= rect.top || y >= rect.bottom) {
    isDragging.value = false;
  }
}

function handleDrop(event: DragEvent) {
  event.preventDefault();
  isDragging.value = false;

  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    const file = event.dataTransfer.files[0];
    void handleFile(file);
  }
}
<<<<<<< HEAD
=======
const show = ref(false);
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)

const selected = ref<CabsRow>({
  name: '',
  id: 0,
  make: 'Mazda',
<<<<<<< HEAD
  quantity: 0,
  price: 0,
  unit_color: 'Black',
  status: 'Out of Stock',
  image: '',
})

const newCab = ref<NewCabInput>({
  name: '',
  make: '',
  quantity: 0,
  price: 0,
  unit_color: '',
  status: 'Out of Stock',
  image: 'https://loremflickr.com/600/400/car',
})

<<<<<<< HEAD
// Image validation
const imageUrlValid = ref(true);
const validatingImage = ref(false);
const defaultImageUrl = getFirstFallbackImage('cab');
=======
const newCab = ref<Omit<CabsRow, 'id'>>({
=======
  quantity: 0,
  price: 0,
  unit_color: 'Black',
  status: 'Out of Stock',
  image: '',
})

const newCab = ref<NewCabInput>({
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  name: '',
  make: '',
  quantity: 0,
  price: 0,
  unit_color: '',
  status: 'Out of Stock',
  image: 'https://loremflickr.com/600/400/car',
})

// Image validation
const imageUrlValid = ref(true);
const validatingImage = ref(false);
<<<<<<< HEAD
const defaultImageUrl = 'https://loremflickr.com/600/400/car';
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
const defaultImageUrl = getDefaultImage('cab');
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)

// Available options from store
const { makes, colors, statuses } = store;

const capitalizedName = computed({
  get: () => newCab.value.name,
  set: (value: string) => {
    if (value) {
      newCab.value.name = value.charAt(0).toUpperCase() + value.slice(1);
    } else {
      newCab.value.name = value;
    }
  }
});

const columns: QTableColumn[] = [
  { name: 'id', align: 'center', label: 'ID', field: 'id', sortable: true },
  {
    name: 'unitName',
    required: true,
    label: 'Unit Name',
    align: 'left',
    field: 'name',
    sortable: true
  },
  { name: 'make', label: 'Make', field: 'make' },
  { name: 'quantity', label: 'Quantity', field: 'quantity', sortable: true },
  { name: 'price', label: 'Price', field: 'price', sortable: true, format: (val: number) =>
        `₱ ${val.toLocaleString('en-PH', {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2
        })}`
  },
  { name: 'status', label: 'Status', field: 'status' },
  { name: 'color', label: 'Color', field: 'unit_color' },
  {
    name: 'actions',
    label: 'Actions',
    field: 'actions',
<<<<<<< HEAD
<<<<<<< HEAD
    align: 'center',
    sortable: false
=======
    align: 'center'
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
    align: 'center',
    sortable: false
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  }
];

const onRowClick: QTableProps['onRowClick'] = (evt, row) => {
  // Check if the click originated from the action button or its menu
  const target = evt.target as HTMLElement;
  if (target.closest('.action-button') || target.closest('.action-menu')) {
    return; // Do nothing if clicked on action button or its menu
  }
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  
  // Update selected with a proper copy of the row data
  selected.value = { ...row as CabsRow };
  
  // Validate and set the image
  if (selected.value.image) {
    if (selected.value.image.startsWith('data:image/')) {
      // For base64 images, validate but preserve the original if it's a valid format
      const validationResult = validateAndSanitizeBase64Image(selected.value.image);
      if (validationResult.isValid && validationResult.sanitizedData) {
        selected.value.image = validationResult.sanitizedData;
      }
      // Even if validation fails, we'll let the ProductCardModal handle the fallback
    }
  } else {
    // If no image, use default
    selected.value.image = defaultImageUrl;
  }
  
  showProductCardModal.value = true;
<<<<<<< HEAD
}

function addToCart() {
  if (selected.value) {
    console.log('added to cart for', selected.value.name);
  }
  showProductCardModal.value = false;
}

function openAddDialog() {
  newCab.value = {
    name: '',
    make: '',
    quantity: 0,
    price: 0,
    unit_color: '',
    status: 'Out of Stock',
    image: defaultImageUrl
  };
  previewUrl.value = defaultImageUrl;
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = ''; // Clear the file input
  }
  showAddDialog.value = true;
}

async function addNewCab() {
  try {
    // Validate required fields
    if (!newCab.value.make || !newCab.value.unit_color) {
      operationNotifications.validation.error('Please fill in all required fields');
      return;
    }

    // Validate image
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image');
      return;
    }

    // If no image is uploaded, use default
    if (!newCab.value.image || newCab.value.image === '') {
=======
  selected.value = row as CabsRow
  show.value = true
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
}

function addToCart() {
  if (selected.value) {
    console.log('added to cart for', selected.value.name);
  }
  showProductCardModal.value = false;
}

function openAddDialog() {
  newCab.value = {
    name: '',
    make: '',
    quantity: 0,
    price: 0,
    unit_color: '',
    status: 'Out of Stock',
    image: defaultImageUrl
  };
  previewUrl.value = defaultImageUrl;
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = ''; // Clear the file input
  }
  showAddDialog.value = true;
}

async function addNewCab() {
  try {
    // Validate required fields
    if (!newCab.value.make || !newCab.value.unit_color) {
      operationNotifications.validation.error('Please fill in all required fields');
      return;
    }

<<<<<<< HEAD
    // If image URL is empty, use default
    if (!newCab.value.image) {
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
    // Validate image
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image');
      return;
    }

    // If no image is uploaded, use default
    if (!newCab.value.image || newCab.value.image === '') {
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      newCab.value.image = defaultImageUrl;
    }

    // Execute the store action and await its completion
    const result = await store.addCab(newCab.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showAddDialog.value = false;
<<<<<<< HEAD
<<<<<<< HEAD
      clearImageInput(); // Clear the image input state
      operationNotifications.add.success(`cab: ${newCab.value.name}`);
    }
  } catch (error) {
    console.error('Error adding cab:', error);
    operationNotifications.add.error('cab');
=======

      $q.notify({
        color: 'positive',
        message: `Added new cab: ${newCab.value.name}`,
        position: 'top',
        timeout: 2000
      });
    }
  } catch (error) {
    console.error('Error adding cab:', error);
    $q.notify({
      color: 'negative',
      message: 'Failed to add cab',
      position: 'top',
      timeout: 2000
    });
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
      clearImageInput(); // Clear the image input state
      operationNotifications.add.success(`cab: ${newCab.value.name}`);
    }
  } catch (error) {
    console.error('Error adding cab:', error);
    operationNotifications.add.error('cab');
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  }
}

function applyFilters() {
  showFilterDialog.value = false;
<<<<<<< HEAD
<<<<<<< HEAD
  operationNotifications.filters.success();
=======

  $q.notify({
    color: 'positive',
    message: 'Filters applied successfully',
    position: 'top',
    timeout: 2000
  });
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
  operationNotifications.filters.success();
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
}

// Add watch for quantity changes
watch(() => newCab.value.quantity, (newQuantity) => {
  if (newQuantity === 0) {
    newCab.value.status = 'Out of Stock';
  } else if (newQuantity <= 2) {
    newCab.value.status = 'Low Stock';
  } else if (newQuantity <= 5) {
    newCab.value.status = 'In Stock';
  } else {
    newCab.value.status = 'Available';
  }
});

// Function to validate if URL is a valid image
<<<<<<< HEAD
<<<<<<< HEAD
let currentAbortController: AbortController | null = null;

=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
let currentAbortController: AbortController | null = null;

>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
async function validateImageUrl(url: string): Promise<boolean> {
  if (!url) {
    imageUrlValid.value = false;
    return false;
  }

  if (!url.startsWith('http')) {
    imageUrlValid.value = false;
    return false;
  }

  validatingImage.value = true;

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  // Abort any existing validation
  if (currentAbortController) {
    currentAbortController.abort();
  }

  // Create new abort controller for this validation
  currentAbortController = new AbortController();
  const signal = currentAbortController.signal;

<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  try {
    const result = await new Promise<boolean>((resolve) => {
      const img = new Image();

      const cleanup = () => {
        img.onload = null;
        img.onerror = null;
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
        if (currentAbortController?.signal === signal) {
          currentAbortController = null;
        }
      };

      // Handle abort signal
      signal.addEventListener('abort', () => {
        cleanup();
        resolve(false);
      });

<<<<<<< HEAD
=======
      };

>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      img.onload = () => {
        cleanup();
        imageUrlValid.value = true;
        validatingImage.value = false;
        resolve(true);
      };

      img.onerror = () => {
        cleanup();
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
        // Try next fallback image if current one fails
        if (url === newCab.value.image) {
          const nextFallback = getNextFallbackImage(url, 'cab');
          newCab.value.image = nextFallback;
          // Handle the promise without making the handler async
          void validateImageUrl(nextFallback).catch(error => {
            console.error('Error validating fallback image:', error);
            imageUrlValid.value = false;
          });
        }
<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
        imageUrlValid.value = false;
        validatingImage.value = false;
        resolve(false);
      };

      // Set a timeout to avoid hanging
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      const timeoutId = setTimeout(() => {
        if (!signal.aborted) {
          currentAbortController?.abort();
          // Try next fallback image if current one times out
          if (url === newCab.value.image) {
            const nextFallback = getNextFallbackImage(url, 'cab');
            newCab.value.image = nextFallback;
            // Handle the promise without making the timeout callback async
            void validateImageUrl(nextFallback).catch(error => {
              console.error('Error validating fallback image:', error);
              imageUrlValid.value = false;
            });
          }
          imageUrlValid.value = false;
          validatingImage.value = false;
          resolve(false);
        }
<<<<<<< HEAD
      }, 5000);

      // Clean up timeout if aborted
      signal.addEventListener('abort', () => {
        clearTimeout(timeoutId);
      });

=======
      setTimeout(() => {
        cleanup();
        imageUrlValid.value = false;
        validatingImage.value = false;
        resolve(false);
      }, 5000);

>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
      }, 5000);

      // Clean up timeout if aborted
      signal.addEventListener('abort', () => {
        clearTimeout(timeoutId);
      });

>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      img.src = url;
    });

    return result;
  } catch (error) {
    console.error('Error validating image URL:', error);
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
    // Try next fallback image if current one errors
    if (url === newCab.value.image) {
      const nextFallback = getNextFallbackImage(url, 'cab');
      newCab.value.image = nextFallback;
      return validateImageUrl(nextFallback).catch(error => {
        console.error('Error validating fallback image:', error);
        imageUrlValid.value = false;
        return false;
      });
    }
<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
    imageUrlValid.value = false;
    validatingImage.value = false;
    return false;
  } finally {
    if (validatingImage.value) {
      validatingImage.value = false;
    }
  }
}

// Modify the watch for image URL changes to handle the default image case
<<<<<<< HEAD
<<<<<<< HEAD
watch(() => newCab.value.image, (newUrl: string) => {
=======
watch(() => newCab.value.image, async (newUrl: string) => {
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
watch(() => newCab.value.image, (newUrl: string) => {
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  if (!newUrl || newUrl === defaultImageUrl) {
    imageUrlValid.value = true; // Default image or empty should be valid
    return;
  }
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  if (newUrl.startsWith('data:image/')) {
    const validationResult = validateAndSanitizeBase64Image(newUrl);
    if (validationResult.isValid) {
      newCab.value.image = validationResult.sanitizedData!;
      imageUrlValid.value = true;
<<<<<<< HEAD
    } else {
      $q.notify({
        color: 'negative',
        message: validationResult.error || 'Invalid image data',
        position: 'top',
      });
      imageUrlValid.value = false;
    }
    return;
  }

  // Handle the promise in the watcher
  void validateImageUrl(newUrl).catch(error => {
    console.error('Error in image URL watcher:', error);
    imageUrlValid.value = false;
  });
=======
  try {
    if (newUrl.startsWith('data:image/')) {
      imageUrlValid.value = true; // Base64 image data is valid
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
    } else {
      $q.notify({
        color: 'negative',
        message: validationResult.error || 'Invalid image data',
        position: 'top',
      });
      imageUrlValid.value = false;
    }
    return;
  }

  // Handle the promise in the watcher
  void validateImageUrl(newUrl).catch(error => {
    console.error('Error in image URL watcher:', error);
    imageUrlValid.value = false;
<<<<<<< HEAD
  }
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
  });
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
});

// Add new refs for file handling
const fileInput = ref<HTMLInputElement | null>(null);
<<<<<<< HEAD
<<<<<<< HEAD
const previewUrl = ref('');

// Add these constants and types at the top of the script
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels

type AllowedMimeType = typeof ALLOWED_TYPES[number];

// Update the validateFile function with stronger file type validation
async function validateFile(file: File): Promise<{ isValid: boolean; error?: string }> {
  try {
    console.log('Starting file validation:', {
      name: file.name,
      type: file.type,
      size: file.size
    });

    // Step 1: Basic file validation
    if (!file) {
      console.error('Validation failed: No file provided');
      return { isValid: false, error: 'No file provided.' };
    }

    // Step 2: Size validation
    if (file.size > MAX_FILE_SIZE) {
      const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
      console.error(`Validation failed: File size ${sizeMB}MB exceeds limit`);
      return {
        isValid: false,
        error: `File size (${sizeMB}MB) exceeds 5MB limit. Please choose a smaller file.`
      };
    }

    // Step 3: Enhanced MIME type validation with file signature check
    const validMimeTypes = {
      'image/jpeg': [0xFF, 0xD8, 0xFF],
      'image/png': [0x89, 0x50, 0x4E, 0x47],
      'image/gif': [0x47, 0x49, 0x46, 0x38]
    };

    if (!Object.keys(validMimeTypes).includes(file.type)) {
      console.error(`Validation failed: Invalid file type ${file.type}`);
      return {
        isValid: false,
        error: `Invalid file type: ${file.type}. Please upload a JPG, PNG, or GIF image.`
      };
    }

    // Step 4: File signature validation
    const arrayBuffer = await file.slice(0, 4).arrayBuffer();
    const bytes = new Uint8Array(arrayBuffer);
    const expectedSignature = validMimeTypes[file.type as keyof typeof validMimeTypes];

    const isValidSignature = expectedSignature.every((byte, i) => byte === bytes[i]);
    if (!isValidSignature) {
      console.error('Validation failed: File signature mismatch');
      return {
        isValid: false,
        error: 'Invalid image format. The file content does not match its extension.'
      };
    }

    // Step 5: Validate image dimensions
    try {
      const dimensionValidation = await validateImageDimensions(file);
      if (!dimensionValidation.isValid) {
        console.error('Validation failed:', dimensionValidation.error);
        return dimensionValidation;
      }
    } catch (error) {
      console.error('Error validating image dimensions:', error);
      return {
        isValid: false,
        error: 'Error validating image dimensions. Please try again.'
      };
    }

    console.log('File validation passed successfully');
    return { isValid: true };
  } catch (error) {
    console.error('Unexpected error during file validation:', error);
    return {
      isValid: false,
      error: 'An unexpected error occurred during validation. Please try again.'
    };
  }
}

// Function to validate image dimensions with better error handling
function validateImageDimensions(file: File): Promise<{ isValid: boolean; error?: string }> {
  return new Promise((resolve) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    const cleanup = () => {
      URL.revokeObjectURL(objectUrl);
    };

    img.onload = () => {
      cleanup();
      console.log('Image dimensions:', {
        width: img.width,
        height: img.height,
        maxAllowed: MAX_DIMENSION
      });

      if (img.width > MAX_DIMENSION || img.height > MAX_DIMENSION) {
        resolve({
          isValid: false,
          error: `Image dimensions (${img.width}x${img.height}) cannot exceed ${MAX_DIMENSION}x${MAX_DIMENSION} pixels.`
        });
      } else if (img.width === 0 || img.height === 0) {
        resolve({
          isValid: false,
          error: 'Invalid image dimensions.'
        });
      } else {
        resolve({ isValid: true });
      }
    };

    img.onerror = () => {
      cleanup();
      console.error('Error loading image for dimension validation');
      resolve({
        isValid: false,
        error: 'Error loading image. Please ensure it is a valid image file.'
      });
    };

    // Set a timeout to avoid hanging
    const timeout = setTimeout(() => {
      cleanup();
      console.error('Dimension validation timed out');
      resolve({
        isValid: false,
        error: 'Image validation timed out. Please try again.'
      });
    }, 10000); // 10 second timeout

    img.src = objectUrl;

    // Clear timeout when image loads or errors
    img.onload = () => {
      clearTimeout(timeout);
      cleanup();
      if (img.width > MAX_DIMENSION || img.height > MAX_DIMENSION) {
        resolve({
          isValid: false,
          error: `Image dimensions (${img.width}x${img.height}) cannot exceed ${MAX_DIMENSION}x${MAX_DIMENSION} pixels.`
        });
      } else if (img.width === 0 || img.height === 0) {
        resolve({
          isValid: false,
          error: 'Invalid image dimensions.'
        });
      } else {
        resolve({ isValid: true });
      }
    };

    img.onerror = () => {
      clearTimeout(timeout);
      cleanup();
      resolve({
        isValid: false,
        error: 'Error loading image. Please ensure it is a valid image file.'
      });
    };
  });
}

// Add new ref for upload loading state
const isUploadingImage = ref(false);

// Update handleFile function to focus on file processing
async function handleFile(file: File) {
  try {
    isUploadingImage.value = true;

    console.log('Starting file validation for:', file.name);
    const validation = await validateFile(file);
    if (!validation.isValid) {
      console.error('File validation failed:', validation.error);
      $q.notify({
        type: 'negative',
        message: validation.error || 'Invalid file',
        position: 'top',
        timeout: 3000
      });
      return;
    }
    console.log('File validation passed');

    // Create a temporary URL for preview
    console.log('Creating preview URL');
    const tempPreviewUrl = URL.createObjectURL(file);
    previewUrl.value = tempPreviewUrl;
    console.log('Preview URL set:', previewUrl.value);

    console.log('Starting FileReader');
    const reader = new FileReader();

    reader.onload = (e) => {
      console.log('FileReader loaded');
      if (e.target?.result) {
        const base64String = e.target.result as string;
        console.log('Processing base64 data');
        const base64ValidationResult = validateAndSanitizeBase64Image(base64String);

        if (!base64ValidationResult.isValid) {
          console.error('Base64 validation failed:', base64ValidationResult.error);
          $q.notify({
            type: 'negative',
            message: base64ValidationResult.error || 'Invalid image data',
            position: 'top',
            timeout: 3000
          });
          previewUrl.value = defaultImageUrl;
          return;
        }

        console.log('Base64 validation passed, updating image');
        newCab.value.image = base64ValidationResult.sanitizedData!;
        imageUrlValid.value = true;

        $q.notify({
          type: 'positive',
          message: 'Image uploaded successfully',
          position: 'top',
          timeout: 2000
        });
      }
    };

    reader.onerror = (error) => {
      console.error('FileReader error:', error);
      previewUrl.value = defaultImageUrl;
      $q.notify({
        type: 'negative',
        message: 'Error reading file. Please try again.',
        position: 'top',
        timeout: 3000
      });
    };

    console.log('Starting file read');
    reader.readAsDataURL(file);
  } catch (error) {
    console.error('Error in handleFile:', error);
    previewUrl.value = defaultImageUrl;
    $q.notify({
      type: 'negative',
      message: 'An unexpected error occurred. Please try again.',
      position: 'top',
      timeout: 3000
    });
  } finally {
    isUploadingImage.value = false;
  }
}

// Make sure to only revoke the URL when necessary
function removeImage(event?: Event) {
  if (event) {
    event.stopPropagation();
  }
  clearImageInput();
}

// Update the clearImageInput function
function clearImageInput() {
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = defaultImageUrl;
  newCab.value.image = defaultImageUrl;
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  isUploadingImage.value = false;
=======
const isDragging = ref(false);
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
const previewUrl = ref('');

// Add these constants and types at the top of the script
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels

type AllowedMimeType = typeof ALLOWED_TYPES[number];

// Update the validateFile function with stronger file type validation
async function validateFile(file: File): Promise<{ isValid: boolean; error?: string }> {
  try {
    console.log('Starting file validation:', {
      name: file.name,
      type: file.type,
      size: file.size
    });

    // Step 1: Basic file validation
    if (!file) {
      console.error('Validation failed: No file provided');
      return { isValid: false, error: 'No file provided.' };
    }

    // Step 2: Size validation
    if (file.size > MAX_FILE_SIZE) {
      const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
      console.error(`Validation failed: File size ${sizeMB}MB exceeds limit`);
      return {
        isValid: false,
        error: `File size (${sizeMB}MB) exceeds 5MB limit. Please choose a smaller file.`
      };
    }

    // Step 3: Enhanced MIME type validation with file signature check
    const validMimeTypes = {
      'image/jpeg': [0xFF, 0xD8, 0xFF],
      'image/png': [0x89, 0x50, 0x4E, 0x47],
      'image/gif': [0x47, 0x49, 0x46, 0x38]
    };

    if (!Object.keys(validMimeTypes).includes(file.type)) {
      console.error(`Validation failed: Invalid file type ${file.type}`);
      return {
        isValid: false,
        error: `Invalid file type: ${file.type}. Please upload a JPG, PNG, or GIF image.`
      };
    }

    // Step 4: File signature validation
    const arrayBuffer = await file.slice(0, 4).arrayBuffer();
    const bytes = new Uint8Array(arrayBuffer);
    const expectedSignature = validMimeTypes[file.type as keyof typeof validMimeTypes];

    const isValidSignature = expectedSignature.every((byte, i) => byte === bytes[i]);
    if (!isValidSignature) {
      console.error('Validation failed: File signature mismatch');
      return {
        isValid: false,
        error: 'Invalid image format. The file content does not match its extension.'
      };
    }

    // Step 5: Validate image dimensions
    try {
      const dimensionValidation = await validateImageDimensions(file);
      if (!dimensionValidation.isValid) {
        console.error('Validation failed:', dimensionValidation.error);
        return dimensionValidation;
      }
    } catch (error) {
      console.error('Error validating image dimensions:', error);
      return {
        isValid: false,
        error: 'Error validating image dimensions. Please try again.'
      };
    }

    console.log('File validation passed successfully');
    return { isValid: true };
  } catch (error) {
    console.error('Unexpected error during file validation:', error);
    return {
      isValid: false,
      error: 'An unexpected error occurred during validation. Please try again.'
    };
  }
}

// Function to validate image dimensions with better error handling
function validateImageDimensions(file: File): Promise<{ isValid: boolean; error?: string }> {
  return new Promise((resolve) => {
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    const cleanup = () => {
      URL.revokeObjectURL(objectUrl);
    };

    img.onload = () => {
      cleanup();
      console.log('Image dimensions:', {
        width: img.width,
        height: img.height,
        maxAllowed: MAX_DIMENSION
      });

      if (img.width > MAX_DIMENSION || img.height > MAX_DIMENSION) {
        resolve({
          isValid: false,
          error: `Image dimensions (${img.width}x${img.height}) cannot exceed ${MAX_DIMENSION}x${MAX_DIMENSION} pixels.`
        });
      } else if (img.width === 0 || img.height === 0) {
        resolve({
          isValid: false,
          error: 'Invalid image dimensions.'
        });
      } else {
        resolve({ isValid: true });
      }
    };

    img.onerror = () => {
      cleanup();
      console.error('Error loading image for dimension validation');
      resolve({
        isValid: false,
        error: 'Error loading image. Please ensure it is a valid image file.'
      });
    };

    // Set a timeout to avoid hanging
    const timeout = setTimeout(() => {
      cleanup();
      console.error('Dimension validation timed out');
      resolve({
        isValid: false,
        error: 'Image validation timed out. Please try again.'
      });
    }, 10000); // 10 second timeout

    img.src = objectUrl;

    // Clear timeout when image loads or errors
    img.onload = () => {
      clearTimeout(timeout);
      cleanup();
      if (img.width > MAX_DIMENSION || img.height > MAX_DIMENSION) {
        resolve({
          isValid: false,
          error: `Image dimensions (${img.width}x${img.height}) cannot exceed ${MAX_DIMENSION}x${MAX_DIMENSION} pixels.`
        });
      } else if (img.width === 0 || img.height === 0) {
        resolve({
          isValid: false,
          error: 'Invalid image dimensions.'
        });
      } else {
        resolve({ isValid: true });
      }
    };

    img.onerror = () => {
      clearTimeout(timeout);
      cleanup();
      resolve({
        isValid: false,
        error: 'Error loading image. Please ensure it is a valid image file.'
      });
    };
  });
}

// Add new ref for upload loading state
const isUploadingImage = ref(false);

// Update handleFile function to focus on file processing
async function handleFile(file: File) {
  try {
    isUploadingImage.value = true;

    console.log('Starting file validation for:', file.name);
    const validation = await validateFile(file);
    if (!validation.isValid) {
      console.error('File validation failed:', validation.error);
      $q.notify({
        type: 'negative',
        message: validation.error || 'Invalid file',
        position: 'top',
        timeout: 3000
      });
      return;
    }
    console.log('File validation passed');

    // Create a temporary URL for preview
    console.log('Creating preview URL');
    const tempPreviewUrl = URL.createObjectURL(file);
    previewUrl.value = tempPreviewUrl;
    console.log('Preview URL set:', previewUrl.value);

    console.log('Starting FileReader');
    const reader = new FileReader();

    reader.onload = (e) => {
      console.log('FileReader loaded');
      if (e.target?.result) {
        const base64String = e.target.result as string;
        console.log('Processing base64 data');
        const base64ValidationResult = validateAndSanitizeBase64Image(base64String);

        if (!base64ValidationResult.isValid) {
          console.error('Base64 validation failed:', base64ValidationResult.error);
          $q.notify({
            type: 'negative',
            message: base64ValidationResult.error || 'Invalid image data',
            position: 'top',
            timeout: 3000
          });
          previewUrl.value = defaultImageUrl;
          return;
        }

        console.log('Base64 validation passed, updating image');
        newCab.value.image = base64ValidationResult.sanitizedData!;
        imageUrlValid.value = true;

        $q.notify({
          type: 'positive',
          message: 'Image uploaded successfully',
          position: 'top',
          timeout: 2000
        });
      }
    };

    reader.onerror = (error) => {
      console.error('FileReader error:', error);
      previewUrl.value = defaultImageUrl;
      $q.notify({
        type: 'negative',
        message: 'Error reading file. Please try again.',
        position: 'top',
        timeout: 3000
      });
    };

    console.log('Starting file read');
    reader.readAsDataURL(file);
  } catch (error) {
    console.error('Error in handleFile:', error);
    previewUrl.value = defaultImageUrl;
    $q.notify({
      type: 'negative',
      message: 'An unexpected error occurred. Please try again.',
      position: 'top',
      timeout: 3000
    });
  } finally {
    isUploadingImage.value = false;
  }
}

// Make sure to only revoke the URL when necessary
function removeImage(event?: Event) {
  if (event) {
    event.stopPropagation();
  }
  clearImageInput();
}

// Update the clearImageInput function
function clearImageInput() {
  if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = defaultImageUrl;
  newCab.value.image = defaultImageUrl;
  imageUrlValid.value = true;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
<<<<<<< HEAD
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
  isUploadingImage.value = false;
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
}

// Function to handle edit cab
function editCab(cab: CabsRow) {
  selected.value = { ...cab };
  newCab.value = {
    name: cab.name,
    make: cab.make,
    quantity: cab.quantity,
    price: cab.price,
    unit_color: cab.unit_color,
    status: cab.status,
    image: cab.image
  };
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)

  // Handle the image preview for base64 images
  if (cab.image.startsWith('data:image/')) {
    const validationResult = validateAndSanitizeBase64Image(cab.image);
    if (validationResult.isValid) {
      previewUrl.value = validationResult.sanitizedData!;
      newCab.value.image = validationResult.sanitizedData!;
      imageUrlValid.value = true;
    } else {
      previewUrl.value = defaultImageUrl;
      newCab.value.image = defaultImageUrl;
      imageUrlValid.value = true;
      operationNotifications.validation.warning('Invalid image data, using default image');
    }
  } else {
    // For any other case, use default image
    previewUrl.value = defaultImageUrl;
    newCab.value.image = defaultImageUrl;
    imageUrlValid.value = true;
  }
<<<<<<< HEAD
=======
  previewUrl.value = cab.image; // Set the preview URL for the existing image
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  showEditDialog.value = true;
}

// Function to handle update cab
async function updateCab() {
  try {
<<<<<<< HEAD
<<<<<<< HEAD
    // Validate required fields
    if (!newCab.value.make || !newCab.value.unit_color) {
      operationNotifications.validation.error('Please fill in all required fields');
      return;
    }

    // Validate image
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image');
      return;
    }

    // If no image is uploaded, use default
    if (!newCab.value.image || newCab.value.image === '') {
      newCab.value.image = defaultImageUrl;
    }

    if (!selected.value) {
      throw new Error('No cab selected for update');
    }

=======
    // Validate image URL before proceeding
    if (!imageUrlValid.value) {
      $q.notify({
        color: 'negative',
        message: 'Please provide a valid image URL',
        position: 'top',
        timeout: 2000
      });
=======
    // Validate required fields
    if (!newCab.value.make || !newCab.value.unit_color) {
      operationNotifications.validation.error('Please fill in all required fields');
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      return;
    }

    // Validate image
    if (!imageUrlValid.value) {
      operationNotifications.validation.error('Please provide a valid image');
      return;
    }

    // If no image is uploaded, use default
    if (!newCab.value.image || newCab.value.image === '') {
      newCab.value.image = defaultImageUrl;
    }

<<<<<<< HEAD
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
    if (!selected.value) {
      throw new Error('No cab selected for update');
    }

>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
    // Execute the store action and await its completion
    const result = await store.updateCab(selected.value.id, newCab.value);

    // Only close dialog and show notification after operation successfully completes
    if (result.success) {
      showEditDialog.value = false;
<<<<<<< HEAD
<<<<<<< HEAD
      clearImageInput(); // Clear the image input state
      operationNotifications.update.success(`cab: ${newCab.value.name}`);
    }
  } catch (error) {
    console.error('Error updating cab:', error);
    operationNotifications.update.error('cab');
=======
      clearImageInput();

      $q.notify({
        color: 'positive',
        message: `Updated cab: ${newCab.value.name}`,
        position: 'top',
        timeout: 2000
      });
    }
  } catch (error) {
    console.error('Error updating cab:', error);
    $q.notify({
      color: 'negative',
      message: 'Failed to update cab',
      position: 'top',
      timeout: 2000
    });
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
      clearImageInput(); // Clear the image input state
      operationNotifications.update.success(`cab: ${newCab.value.name}`);
    }
  } catch (error) {
    console.error('Error updating cab:', error);
    operationNotifications.update.error('cab');
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  }
}

// Function to handle delete cab
function deleteCab(cab: CabsRow) {
  cabToDelete.value = cab;
  showDeleteDialog.value = true;
}

// Function to confirm and execute delete
async function confirmDelete() {
  try {
    if (!cabToDelete.value) return;

    await store.deleteCab(cabToDelete.value.id);
    showDeleteDialog.value = false;
    cabToDelete.value = null;
<<<<<<< HEAD
<<<<<<< HEAD
    operationNotifications.delete.success('cab');
  } catch (error) {
    console.error('Error deleting cab:', error);
    operationNotifications.delete.error('cab');
  }
}

// Function to handle file selection
async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    console.log('Selected file:', {
      name: file.name,
      type: file.type,
      size: file.size,
      lastModified: file.lastModified
    });

    // Check file type
    if (!ALLOWED_TYPES.includes(file.type as AllowedMimeType)) {
      $q.notify({
        type: 'negative',
        message: `Invalid file type: ${file.type}. Allowed types are: JPEG, PNG, and GIF`,
        position: 'top',
        timeout: 3000
      });
      return;
    }

    // Check file size
    if (file.size > MAX_FILE_SIZE) {
      $q.notify({
        type: 'negative',
        message: `File size (${(file.size / 1024 / 1024).toFixed(2)}MB) exceeds the 5MB limit`,
        position: 'top',
        timeout: 3000
      });
      return;
    }

    try {
      await handleFile(file);
    } catch (error) {
      console.error('Error handling file:', error);
      $q.notify({
        type: 'negative',
        message: 'Error processing the image. Please try again.',
        position: 'top',
        timeout: 3000
      });
    }
  }
}

// Function to trigger file input
function triggerFileInput() {
  fileInput.value?.click();
}

// Initialize data on component mount
onMounted(async () => {
  await store.initializeCabs();
});

=======

    $q.notify({
      color: 'positive',
      message: `Successfully deleted cab`,
      position: 'top',
      timeout: 2000
    });
=======
    operationNotifications.delete.success('cab');
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  } catch (error) {
    console.error('Error deleting cab:', error);
    operationNotifications.delete.error('cab');
  }
}

<<<<<<< HEAD
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
// Function to handle file selection
async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files[0]) {
    const file = input.files[0];
    console.log('Selected file:', {
      name: file.name,
      type: file.type,
      size: file.size,
      lastModified: file.lastModified
    });

    // Check file type
    if (!ALLOWED_TYPES.includes(file.type as AllowedMimeType)) {
      $q.notify({
        type: 'negative',
        message: `Invalid file type: ${file.type}. Allowed types are: JPEG, PNG, and GIF`,
        position: 'top',
        timeout: 3000
      });
      return;
    }

    // Check file size
    if (file.size > MAX_FILE_SIZE) {
      $q.notify({
        type: 'negative',
        message: `File size (${(file.size / 1024 / 1024).toFixed(2)}MB) exceeds the 5MB limit`,
        position: 'top',
        timeout: 3000
      });
      return;
    }

    try {
      await handleFile(file);
    } catch (error) {
      console.error('Error handling file:', error);
      $q.notify({
        type: 'negative',
        message: 'Error processing the image. Please try again.',
        position: 'top',
        timeout: 3000
      });
    }
  }
}

// Function to trigger file input
function triggerFileInput() {
  fileInput.value?.click();
}

// Initialize data on component mount
onMounted(async () => {
  await store.initializeCabs();
});

>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
</script>

<template>
  <q-page class="flex q-pa-md">
    <div class="q-pa-sm full-width">
      <div class="flex row q-my-sm">
        <div class="flex full-width col">
          <div class="flex col q-mr-sm">
            <q-input
              v-model="store.rawCabSearch"
              outlined
              dense
              placeholder="Search"
              class="full-width"
            >
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
          <div class="flex col">
            <q-btn
              outline
              icon="filter_list"
              label="Filters"
              @click="showFilterDialog = true"
            />
          </div>
        </div>

        <div class="flex row q-gutter-x-sm">
          <q-btn class="text-white bg-primary" unelevated @click="openAddDialog">
            <q-icon name="add" color="white" />
            Add
          </q-btn>
          <div class="flex row">
            <q-btn dense flat class="bg-primary text-white q-pa-sm">
              <q-icon name="download" color="white" />
              Download CSV
            </q-btn>
          </div>
        </div>
      </div>

<<<<<<< HEAD
<<<<<<< HEAD
      <!--CABS TABLE-->
=======
      <!--INVENTORY TABLE-->
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
      <!--CABS TABLE-->
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      <q-table
        class="my-sticky-column-table"
        flat
        bordered
        title="Cabs"
        :rows="store.filteredCabRows"
        :columns="columns"
        row-key="id"
        :filter="store.cabSearch"
        @row-click="onRowClick"
        :pagination="{ rowsPerPage: 5 }"
<<<<<<< HEAD
<<<<<<< HEAD
        :loading="store.isLoading"
      >
        <template v-slot:loading>
          <q-inner-loading showing color="primary">
            <q-spinner-gears size="50px" color="primary" />
          </q-inner-loading>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width :key="props.row.id">
            <q-btn
              flat
              round
              dense
              color="grey"
              icon="more_vert"
              class="action-button"
              :aria-label="'Actions for ' + props.row.name"
            >
              <q-menu class="action-menu" :aria-label="'Available actions for ' + props.row.name">
                <q-list style="min-width: 100px">
                  <q-item
                    clickable
                    v-close-popup
                    @click.stop="editCab(props.row)"
                    role="button"
                    :aria-label="'Edit ' + props.row.name"
                  >
                    <q-item-section>
                      <q-item-label>
                        <q-icon name="edit" size="xs" class="q-mr-sm" aria-hidden="true" />
=======
=======
        :loading="store.isLoading"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
      >
        <template v-slot:loading>
          <q-inner-loading showing color="primary">
            <q-spinner-gears size="50px" color="primary" />
          </q-inner-loading>
        </template>
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width :key="props.row.id">
            <q-btn
              flat
              round
              dense
              color="grey"
              icon="more_vert"
              class="action-button"
              :aria-label="'Actions for ' + props.row.name"
            >
              <q-menu class="action-menu" :aria-label="'Available actions for ' + props.row.name">
                <q-list style="min-width: 100px">
                  <q-item
                    clickable
                    v-close-popup
                    @click.stop="editCab(props.row)"
                    role="button"
                    :aria-label="'Edit ' + props.row.name"
                  >
                    <q-item-section>
                      <q-item-label>
<<<<<<< HEAD
                        <q-icon name="edit" size="xs" class="q-mr-sm" />
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                        <q-icon name="edit" size="xs" class="q-mr-sm" aria-hidden="true" />
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                        Edit
                      </q-item-label>
                    </q-item-section>
                  </q-item>
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                  <q-item
                    clickable
                    v-close-popup
                    @click.stop="deleteCab(props.row)"
                    role="button"
                    :aria-label="'Delete ' + props.row.name"
                    class="text-negative"
                  >
<<<<<<< HEAD
                    <q-item-section>
                      <q-item-label class="text-negative">
                        <q-icon name="delete" size="xs" class="q-mr-sm" aria-hidden="true" />
=======
                  <q-item clickable v-close-popup @click.stop="deleteCab(props.row)">
                    <q-item-section>
                      <q-item-label class="text-negative">
                        <q-icon name="delete" size="xs" class="q-mr-sm" />
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    <q-item-section>
                      <q-item-label class="text-negative">
                        <q-icon name="delete" size="xs" class="q-mr-sm" aria-hidden="true" />
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                        Delete
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </q-td>
        </template>
      </q-table>

      <!-- Existing Cab Modal -->
      <ProductCardModal
        v-model="showProductCardModal"
        :image="selected?.image || ''"
        :title="selected?.name || ''"
        :unit_color="selected?.unit_color || ''"
        :price="selected?.price || 0"
        :quantity="selected?.quantity || 0"
        :details="`${selected?.make}`"
        :status="selected?.status || ''"
        @add="addToCart"
      />

      <!-- Add Cab Dialog -->
      <q-dialog
        v-model="showAddDialog"
        persistent
        @hide="clearImageInput"
      >
        <q-card style="min-width: 400px; max-width: 95vw">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">New Cab</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section>
            <q-form @submit.prevent="addNewCab" class="q-gutter-sm">
              <q-input
                v-model="capitalizedName"
                label="Cab Name"
                dense
                outlined
                required
                :rules="[val => !!val || 'Name is required']"
              >
                <template v-slot:prepend>
                  <q-icon name="directions_car" />
                </template>
              </q-input>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.make"
                    :options="makes"
                    label="Make"
                    dense
                    outlined
                    required
<<<<<<< HEAD
<<<<<<< HEAD
                    emit-value
                    map-options
                    placeholder="Select a make"
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    emit-value
                    map-options
                    placeholder="Select a make"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    :rules="[val => !!val || 'Make is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="business" />
                    </template>
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    <template v-slot:no-option>
                      <q-item>
                        <q-item-section class="text-grey">
                          No results
                        </q-item-section>
                      </q-item>
                    </template>
<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                  </q-select>
                </div>

                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.unit_color"
                    :options="colors"
                    label="Color"
                    dense
                    outlined
                    required
<<<<<<< HEAD
<<<<<<< HEAD
                    emit-value
                    map-options
                    placeholder="Select a color"
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    emit-value
                    map-options
                    placeholder="Select a color"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    :rules="[val => !!val || 'Color is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="palette" />
                    </template>
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    <template v-slot:no-option>
                      <q-item>
                        <q-item-section class="text-grey">
                          No results
                        </q-item-section>
                      </q-item>
                    </template>
<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                  </q-select>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-input
                    v-model.number="newCab.quantity"
                    type="number"
<<<<<<< HEAD
<<<<<<< HEAD
                    min="0"
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    min="0"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    label="Quantity"
                    dense
                    outlined
                    required
                    :rules="[val => val >= 0 || 'Quantity must be positive']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="numbers" />
                    </template>
                  </q-input>
                </div>

                <div class="col-12 col-sm-6">
                  <q-input
                    v-model.number="newCab.price"
                    type="number"
                    label="Price"
                    dense
                    outlined
                    required
                    :rules="[val => val > 0 || 'Price must be greater than 0']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="attach_money" />
                    </template>
                  </q-input>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12">
                  <q-input
                    v-model="newCab.status"
                    label="Status"
                    dense
                    outlined
                    readonly
<<<<<<< HEAD
<<<<<<< HEAD
=======
                    disable
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                  >
                    <template v-slot:prepend>
                      <q-icon name="info" />
                    </template>
                  </q-input>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12">
                  <div
                    class="upload-container q-pa-md"
                    :class="{ 'dragging': isDragging }"
                    @dragenter.prevent="isDragging = true"
                    @dragover.prevent="isDragging = true"
<<<<<<< HEAD
<<<<<<< HEAD
                    @dragleave.prevent="handleDragLeave"
                    @drop.prevent="handleDrop"
                    @dragend.prevent="isDragging = false"
=======
                    @dragleave.prevent="isDragging = false"
                    @drop.prevent="handleDrop"
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    @dragleave.prevent="handleDragLeave"
                    @drop.prevent="handleDrop"
                    @dragend.prevent="isDragging = false"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    @click="triggerFileInput"
                  >
                    <input
                      type="file"
                      ref="fileInput"
<<<<<<< HEAD
<<<<<<< HEAD
                      accept="image/jpeg,image/png,image/gif"
                      class="hidden"
                      @change="handleFileSelect"
                    >
                    <div v-if="isUploadingImage" class="text-center">
                      <q-spinner-dots color="primary" size="40px" />
                      <div class="text-body1 q-mt-sm">
                        Processing image...
                      </div>
                    </div>
                    <div v-else-if="!previewUrl" class="text-center">
=======
                      accept="image/*"
                      class="hidden"
                      @change="handleFileSelect"
                    >
                    <div class="text-center" v-if="!previewUrl">
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                      accept="image/jpeg,image/png,image/gif"
                      class="hidden"
                      @change="handleFileSelect"
                    >
                    <div v-if="isUploadingImage" class="text-center">
                      <q-spinner-dots color="primary" size="40px" />
                      <div class="text-body1 q-mt-sm">
                        Processing image...
                      </div>
                    </div>
                    <div v-else-if="!previewUrl" class="text-center">
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                      <q-icon name="cloud_upload" size="48px" color="primary" />
                      <div class="text-body1 q-mt-sm">
                        Drag and drop an image here or click to select
                      </div>
<<<<<<< HEAD
<<<<<<< HEAD
                      <div class="text-caption text-grey q-mt-sm">
                        Supported formats: JPG, PNG, GIF
                        <br>
                        Maximum file size: 5MB
                        <br>
                        Maximum dimensions: 4096x4096px
                      </div>
                    </div>
                    <div v-else class="row items-center justify-center">
                      <div class="col-8 text-center">
                        <img
                          :src="previewUrl"
                          class="preview-image"
                          :alt="newCab.name || 'Preview image'"
                        />
=======
                      <div class="text-caption text-grey">
=======
                      <div class="text-caption text-grey q-mt-sm">
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                        Supported formats: JPG, PNG, GIF
                        <br>
                        Maximum file size: 5MB
                        <br>
                        Maximum dimensions: 4096x4096px
                      </div>
                    </div>
<<<<<<< HEAD
                    <div v-else class="row items-center">
                      <div class="col-8">
                        <img :src="previewUrl" class="preview-image" />
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    <div v-else class="row items-center justify-center">
                      <div class="col-8 text-center">
                        <img
                          :src="previewUrl"
                          class="preview-image"
                          :alt="newCab.name || 'Preview image'"
                        />
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                      </div>
                      <div class="col-4 text-center">
                        <q-btn
                          flat
                          round
                          color="negative"
                          icon="close"
<<<<<<< HEAD
<<<<<<< HEAD
                          @click.stop="removeImage"
=======
                          @click.stop="removeImage($event)"
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                          @click.stop="removeImage"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                        >
                          <q-tooltip>Remove Image</q-tooltip>
                        </q-btn>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </q-form>
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md">
            <q-btn flat label="Cancel" v-close-popup />
            <q-btn
              unelevated
              color="primary"
              label="Add Cab"
              @click="addNewCab"
              :disable="!newCab.name || !newCab.make || !newCab.unit_color || newCab.quantity < 0 || newCab.price <= 0"
            />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Filter Dialog -->
      <q-dialog v-model="showFilterDialog">
        <q-card style="min-width: 350px">
          <q-card-section class="row items-center">
            <div class="text-h6">Filter Cabs</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-select
              v-model="store.filterMake"
              :options="makes"
              label="Make"
              clearable
              outlined
              class="q-mb-md"
            />

            <q-select
              v-model="store.filterColor"
              :options="colors"
              label="Color"
              clearable
              outlined
              class="q-mb-md"
            />

            <q-select
              v-model="store.filterStatus"
              :options="statuses"
              label="Status"
              clearable
              outlined
              class="q-mb-md"
            />
          </q-card-section>

          <q-card-actions align="right" class="text-primary">
            <q-btn flat label="Reset" color="negative" @click="store.resetFilters" />
            <q-btn flat label="Apply Filters" @click="applyFilters" />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Edit Cab Dialog -->
      <q-dialog
        v-model="showEditDialog"
        persistent
        @hide="clearImageInput"
      >
        <q-card style="min-width: 400px; max-width: 95vw">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">Edit Cab</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section>
            <q-form @submit.prevent="updateCab" class="q-gutter-sm">
              <q-input
                v-model="capitalizedName"
                label="Cab Name"
                dense
                outlined
                required
                :rules="[val => !!val || 'Name is required']"
              >
                <template v-slot:prepend>
                  <q-icon name="directions_car" />
                </template>
              </q-input>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.make"
                    :options="makes"
                    label="Make"
                    dense
                    outlined
                    required
<<<<<<< HEAD
<<<<<<< HEAD
                    emit-value
                    map-options
                    placeholder="Select a make"
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    emit-value
                    map-options
                    placeholder="Select a make"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    :rules="[val => !!val || 'Make is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="business" />
                    </template>
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    <template v-slot:no-option>
                      <q-item>
                        <q-item-section class="text-grey">
                          No results
                        </q-item-section>
                      </q-item>
                    </template>
<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                  </q-select>
                </div>

                <div class="col-12 col-sm-6">
                  <q-select
                    v-model="newCab.unit_color"
                    :options="colors"
                    label="Color"
                    dense
                    outlined
                    required
<<<<<<< HEAD
<<<<<<< HEAD
                    emit-value
                    map-options
                    placeholder="Select a color"
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    emit-value
                    map-options
                    placeholder="Select a color"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    :rules="[val => !!val || 'Color is required']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="palette" />
                    </template>
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    <template v-slot:no-option>
                      <q-item>
                        <q-item-section class="text-grey">
                          No results
                        </q-item-section>
                      </q-item>
                    </template>
<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                  </q-select>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12 col-sm-6">
                  <q-input
                    v-model.number="newCab.quantity"
                    type="number"
<<<<<<< HEAD
<<<<<<< HEAD
                    min="0"
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    min="0"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    label="Quantity"
                    dense
                    outlined
                    required
                    :rules="[val => val >= 0 || 'Quantity must be positive']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="numbers" />
                    </template>
                  </q-input>
                </div>

                <div class="col-12 col-sm-6">
                  <q-input
                    v-model.number="newCab.price"
                    type="number"
                    label="Price"
                    dense
                    outlined
                    required
                    :rules="[val => val > 0 || 'Price must be greater than 0']"
                  >
                    <template v-slot:prepend>
                      <q-icon name="attach_money" />
                    </template>
                  </q-input>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12">
                  <q-input
                    v-model="newCab.status"
                    label="Status"
                    dense
                    outlined
                    readonly
<<<<<<< HEAD
<<<<<<< HEAD
=======
                    disable
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                  >
                    <template v-slot:prepend>
                      <q-icon name="info" />
                    </template>
                  </q-input>
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-12">
                  <div
                    class="upload-container q-pa-md"
                    :class="{ 'dragging': isDragging }"
                    @dragenter.prevent="isDragging = true"
                    @dragover.prevent="isDragging = true"
<<<<<<< HEAD
<<<<<<< HEAD
                    @dragleave.prevent="handleDragLeave"
                    @drop.prevent="handleDrop"
                    @dragend.prevent="isDragging = false"
=======
                    @dragleave.prevent="isDragging = false"
                    @drop.prevent="handleDrop"
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    @dragleave.prevent="handleDragLeave"
                    @drop.prevent="handleDrop"
                    @dragend.prevent="isDragging = false"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                    @click="triggerFileInput"
                  >
                    <input
                      type="file"
                      ref="fileInput"
<<<<<<< HEAD
<<<<<<< HEAD
                      accept="image/jpeg,image/png,image/gif"
                      class="hidden"
                      @change="handleFileSelect"
                    >
                    <div v-if="isUploadingImage" class="text-center">
                      <q-spinner-dots color="primary" size="40px" />
                      <div class="text-body1 q-mt-sm">
                        Processing image...
                      </div>
                    </div>
                    <div v-else-if="!previewUrl" class="text-center">
=======
                      accept="image/*"
                      class="hidden"
                      @change="handleFileSelect"
                    >
                    <div v-if="!previewUrl && !newCab.image" class="text-center">
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                      accept="image/jpeg,image/png,image/gif"
                      class="hidden"
                      @change="handleFileSelect"
                    >
                    <div v-if="isUploadingImage" class="text-center">
                      <q-spinner-dots color="primary" size="40px" />
                      <div class="text-body1 q-mt-sm">
                        Processing image...
                      </div>
                    </div>
                    <div v-else-if="!previewUrl" class="text-center">
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                      <q-icon name="cloud_upload" size="48px" color="primary" />
                      <div class="text-body1 q-mt-sm">
                        Drag and drop an image here or click to select
                      </div>
<<<<<<< HEAD
<<<<<<< HEAD
                      <div class="text-caption text-grey q-mt-sm">
                        Supported formats: JPG, PNG, GIF
                        <br>
                        Maximum file size: 5MB
                        <br>
                        Maximum dimensions: 4096x4096px
                      </div>
                    </div>
                    <div v-else class="row items-center justify-center">
                      <div class="col-8 text-center">
                        <img
                          :src="previewUrl"
                          class="preview-image"
                          :alt="newCab.name || 'Preview image'"
                        />
=======
                      <div class="text-caption text-grey">
=======
                      <div class="text-caption text-grey q-mt-sm">
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                        Supported formats: JPG, PNG, GIF
                        <br>
                        Maximum file size: 5MB
                        <br>
                        Maximum dimensions: 4096x4096px
                      </div>
                    </div>
<<<<<<< HEAD
                    <div v-else class="row items-center">
                      <div class="col-8">
                        <img :src="previewUrl || newCab.image" class="preview-image" />
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                    <div v-else class="row items-center justify-center">
                      <div class="col-8 text-center">
                        <img
                          :src="previewUrl"
                          class="preview-image"
                          :alt="newCab.name || 'Preview image'"
                        />
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                      </div>
                      <div class="col-4 text-center">
                        <q-btn
                          flat
                          round
                          color="negative"
                          icon="close"
<<<<<<< HEAD
<<<<<<< HEAD
                          @click.stop="removeImage"
=======
                          @click.stop="removeImage($event)"
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
                          @click.stop="removeImage"
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
                        >
                          <q-tooltip>Remove Image</q-tooltip>
                        </q-btn>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </q-form>
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md">
            <q-btn flat label="Cancel" v-close-popup />
            <q-btn
              unelevated
              color="primary"
              label="Update Cab"
              @click="updateCab"
              :disable="!newCab.name || !newCab.make || !newCab.unit_color || newCab.quantity < 0 || newCab.price <= 0"
            />
          </q-card-actions>
        </q-card>
      </q-dialog>

      <!-- Delete Confirmation Dialog -->
      <q-dialog v-model="showDeleteDialog" persistent>
        <q-card>
          <q-card-section class="row items-center">
            <q-avatar icon="warning" color="negative" text-color="white" />
            <span class="q-ml-sm text-h6">Delete Cab</span>
          </q-card-section>

          <q-card-section>
            Are you sure you want to delete {{ cabToDelete?.name }}? This action cannot be undone.
          </q-card-section>

          <q-card-actions align="right">
            <q-btn flat label="Cancel" v-close-popup />
            <q-btn flat label="Delete" color="negative" @click="confirmDelete" />
          </q-card-actions>
        </q-card>
      </q-dialog>
    </div>
  </q-page>
</template>

<style lang="sass">
.my-sticky-column-table
  max-width: 100%

  thead tr:first-child th:nth-child(2)
    background-color: var(--sticky-column-bg)

  td:nth-child(2)
    background-color: var(--sticky-column-bg)

  th:nth-child(2),
  td:nth-child(2)
    position: sticky
    left: 0
    z-index: 1
    color: white

    .body--dark &
      color: black

<<<<<<< HEAD
<<<<<<< HEAD
.cab-page-z-top
  z-index: 1000

.upload-container
  position: relative
=======
.z-top
  z-index: 1000

.upload-container
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
.cab-page-z-top
  z-index: 1000

.upload-container
  position: relative
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  border: 2px dashed #ccc
  border-radius: 8px
  cursor: pointer
  transition: all 0.3s ease
  min-height: 200px
  display: flex
  align-items: center
  justify-content: center
<<<<<<< HEAD
<<<<<<< HEAD
  background-color: transparent
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05)

  &:hover
    border-color: var(--q-primary)
    background-color: rgba(var(--q-primary-rgb), 0.04)
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1)
    transform: translateY(-1px)

  &.dragging
    border-color: var(--q-primary)
    background-color: rgba(var(--q-primary-rgb), 0.08)
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15)
    transform: translateY(-2px)

  .q-spinner-dots
    margin: 0 auto
=======
=======
  background-color: transparent
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05)
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)

  &:hover
    border-color: var(--q-primary)
    background-color: rgba(var(--q-primary-rgb), 0.04)
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1)
    transform: translateY(-1px)

  &.dragging
<<<<<<< HEAD
    border-color: #00b4ff
    background: rgba(0, 180, 255, 0.1)
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
    border-color: var(--q-primary)
    background-color: rgba(var(--q-primary-rgb), 0.08)
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15)
    transform: translateY(-2px)

  .q-spinner-dots
    margin: 0 auto
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)

.preview-image
  width: 100%
  max-height: 180px
  object-fit: contain
  border-radius: 4px
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1)
  transition: all 0.3s ease

  &:hover
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15)
<<<<<<< HEAD
=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======
>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)

.hidden
  display: none

.action-button
  position: relative
  z-index: 1

.action-menu
  z-index: 1001 !important
<<<<<<< HEAD
<<<<<<< HEAD

=======
>>>>>>> dc75c8f (feat(CabsPage) Enhance CabsPage functionality and UI)
=======

>>>>>>> 52c0309 (feat(ProductModal, CabsPage, MaterialsPage) Enhance image handling and validation)
</style>
