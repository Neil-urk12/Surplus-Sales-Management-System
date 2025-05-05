<script setup lang="ts">
import { ref, computed, watch, PropType, nextTick } from 'vue';
import { useQuasar } from 'quasar';
import type { CabsRow, NewCabInput, CabStatus } from 'src/types/cabs';
import { validateAndSanitizeBase64Image } from '../../utils/imageValidation';
import { operationNotifications } from '../../utils/notifications';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    cabData: {
        type: Object as PropType<CabsRow | null>,
        required: true,
    },
    makes: {
        type: Array as PropType<string[]>,
        required: true,
    },
    colors: {
        type: Array as PropType<string[]>,
        required: true,
    },
    defaultImageUrl: {
        type: String,
        required: true,
    },
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'update-cab', cabData: NewCabInput): void
}>();

const $q = useQuasar();

// Constants for file validation
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif'] as const;
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels
type AllowedMimeType = typeof ALLOWED_TYPES[number];

// --- Local State --- 
// Use a local ref to store editable data, initialized from prop
const localCabData = ref<NewCabInput>({
    name: '',
    make: '',
    quantity: 0,
    price: 0,
    unit_color: '',
    status: 'Out of Stock',
    image: props.defaultImageUrl,
});

const imageUrlValid = ref(true);
const validatingImage = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref<string>(props.defaultImageUrl);
const isUploadingImage = ref(false);
const isDragging = ref(false);
let currentAbortController: AbortController | null = null;

// --- Computed --- 
const capitalizedName = computed({
    get: () => localCabData.value.name,
    set: (value: string) => {
        if (value) {
            localCabData.value.name = value.charAt(0).toUpperCase() + value.slice(1);
        } else {
            localCabData.value.name = value;
        }
    }
});

// --- Watchers --- 
watch(() => localCabData.value.quantity, (newQuantity) => {
    if (newQuantity === 0) {
        localCabData.value.status = 'Out of Stock';
    } else if (newQuantity <= 2) {
        localCabData.value.status = 'Low Stock';
    } else if (newQuantity <= 5) {
        localCabData.value.status = 'In Stock';
    } else {
        localCabData.value.status = 'Available';
    }
});

watch(() => localCabData.value.image, (newUrl: string) => {
    if (!newUrl || newUrl === props.defaultImageUrl) {
        previewUrl.value = props.defaultImageUrl;
        imageUrlValid.value = true;
        return;
    }
    if (newUrl.startsWith('data:image/')) {
        // Only do basic validation here, not full sanitization
        const validationResult = validateAndSanitizeBase64Image(newUrl, true);
        if (validationResult.isValid) {
            // Don't update the data yet, just show the preview
            previewUrl.value = newUrl;
            imageUrlValid.value = true;
        } else {
            $q.notify({ color: 'negative', message: validationResult.error || 'Invalid image data', position: 'top' });
            previewUrl.value = props.defaultImageUrl;
            imageUrlValid.value = false;
        }
        return;
    }

    void validateImageUrl(newUrl).then(isValid => {
        if (isValid) {
            previewUrl.value = newUrl;
        } else {
            previewUrl.value = props.defaultImageUrl;
        }
    }).catch(error => {
        console.error('Error in image URL watcher:', error);
        previewUrl.value = props.defaultImageUrl;
        imageUrlValid.value = false;
    });
});

// Watch the incoming prop to update local state when the dialog opens/changes cab
watch(() => props.cabData, (newCab) => {
    if (newCab) {
        // Reset local state based on the new cab data
        localCabData.value = {
            name: newCab.name,
            make: newCab.make,
            quantity: newCab.quantity,
            price: newCab.price,
            unit_color: newCab.unit_color,
            status: newCab.status,
            image: newCab.image || props.defaultImageUrl // Use default if image is missing
        };
        // Reset image validation/preview state
        clearImageInput(); // This sets previewUrl and imageUrlValid correctly based on localCabData.image

        // Ensure the watcher for localCabData.image runs if the image is already set
        nextTick(() => {
            const imageWatcher = watch(() => localCabData.value.image, () => { }, { immediate: true });
            imageWatcher(); // Immediately unwatch after running once
        });

    } else {
        // Reset if cabData becomes null (though this might not happen with persistent dialog)
        resetForm();
    }
}, { immediate: true }); // Immediate run on component mount

// --- Functions --- 

function resetForm() {
    // Reset local data to empty/default state, matching AddCabDialog's reset
    localCabData.value = {
        name: '',
        make: '',
        quantity: 0,
        price: 0,
        unit_color: '',
        status: 'Out of Stock',
        image: props.defaultImageUrl,
    };
    clearImageInput();
}

function handleUpdateCab() {
    if (!localCabData.value.make || !localCabData.value.unit_color) {
        operationNotifications.validation.error('Please fill in all required fields');
        return;
    }
    if (!imageUrlValid.value) {
        operationNotifications.validation.error('Please provide a valid image');
        return;
    }
    // Ensure image field is not empty; use default if necessary
    if (!localCabData.value.image || localCabData.value.image === '') {
        localCabData.value.image = props.defaultImageUrl;
    }

    // Full sanitization right before sending to backend
    if (localCabData.value.image.startsWith('data:image/')) {
        const sanitizationResult = validateAndSanitizeBase64Image(localCabData.value.image, false);
        if (!sanitizationResult.isValid) {
            operationNotifications.validation.error(sanitizationResult.error || 'Invalid image data');
            return;
        }
        localCabData.value.image = sanitizationResult.sanitizedData!;
    }

    emit('update-cab', { ...localCabData.value });
}

function closeDialog() {
    emit('update:modelValue', false);
}


async function validateImageUrl(url: string): Promise<boolean> {
    if (!url || !url.startsWith('http')) {
        imageUrlValid.value = false;
        return false;
    }
    validatingImage.value = true;
    if (currentAbortController) currentAbortController.abort();
    currentAbortController = new AbortController();
    const signal = currentAbortController.signal;
    try {
        const result = await new Promise<boolean>((resolve) => {
            const img = new Image();
            const cleanup = () => { /* ... */ };
            signal.addEventListener('abort', () => { cleanup(); resolve(false); });
            img.onload = () => { cleanup(); imageUrlValid.value = true; validatingImage.value = false; resolve(true); };
            img.onerror = () => { cleanup(); imageUrlValid.value = false; validatingImage.value = false; resolve(false); };
            const timeoutId = setTimeout(() => { /* ... */ }, 5000);
            signal.addEventListener('abort', () => { clearTimeout(timeoutId); });
            img.src = url;
        });
        return result;
    } catch (error) {
        console.error('Error validating image URL:', error);
        imageUrlValid.value = false; validatingImage.value = false; return false;
    } finally {
        if (validatingImage.value) validatingImage.value = false;
    }
}

async function validateFile(file: File): Promise<{ isValid: boolean; error?: string }> {
    try {
        if (!file) return { isValid: false, error: 'No file provided.' };
        if (file.size > MAX_FILE_SIZE) return { isValid: false, error: `File size exceeds 5MB.` };
        const validMimeTypes = { 'image/jpeg': [0xFF, 0xD8, 0xFF], 'image/png': [0x89, 0x50, 0x4E, 0x47], 'image/gif': [0x47, 0x49, 0x46, 0x38] };
        if (!Object.keys(validMimeTypes).includes(file.type)) return { isValid: false, error: `Invalid file type.` };
        const arrayBuffer = await file.slice(0, 4).arrayBuffer();
        const bytes = new Uint8Array(arrayBuffer);
        const expectedSignature = validMimeTypes[file.type as keyof typeof validMimeTypes];
        if (!expectedSignature.every((byte, i) => byte === bytes[i])) return { isValid: false, error: 'File content mismatch.' };
        const dimensionValidation = await validateImageDimensions(file);
        if (!dimensionValidation.isValid) return dimensionValidation;
        return { isValid: true };
    } catch (error) { console.error('File validation error:', error); return { isValid: false, error: 'Validation error.' }; }

}

function validateImageDimensions(file: File): Promise<{ isValid: boolean; error?: string }> {
    return new Promise((resolve) => {
        const img = new Image();
        const objectUrl = URL.createObjectURL(file);
        const cleanup = () => URL.revokeObjectURL(objectUrl);
        let timeout: NodeJS.Timeout | null = null;

        const resolveClean = (result: { isValid: boolean; error?: string }) => {
            if (timeout) clearTimeout(timeout);
            cleanup();
            resolve(result);
        };

        img.onload = () => {
            if (img.width > MAX_DIMENSION || img.height > MAX_DIMENSION) {
                resolveClean({
                    isValid: false,
                    error: `Image dimensions exceed maximum of ${MAX_DIMENSION}px`
                });
            } else {
                resolveClean({ isValid: true });
            }
        };

        img.onerror = () => resolveClean({ isValid: false, error: 'Error loading image.' });
        timeout = setTimeout(() => resolveClean({ isValid: false, error: 'Validation timed out.' }), 10000);
        img.src = objectUrl;
    });
}

async function handleFile(file: File) {
    isUploadingImage.value = true;
    previewUrl.value = ''; // Clear preview during processing
    try {
        const validation = await validateFile(file);
        if (!validation.isValid) {
            $q.notify({ type: 'negative', message: validation.error || 'Invalid file', position: 'top' });
            clearImageInput();
            isUploadingImage.value = false;
            return;
        }
        const reader = new FileReader();

        reader.onload = (e) => {
            try {
                if (e.target?.result) {
                    const base64String = e.target.result as string;
                    // Only do basic validation during upload
                    const base64ValidationResult = validateAndSanitizeBase64Image(base64String, true);
                    if (!base64ValidationResult.isValid) {
                        $q.notify({ type: 'negative', message: base64ValidationResult.error || 'Invalid image data', position: 'top' });
                        clearImageInput();
                    } else {
                        // Store the original base64 string for later sanitization
                        localCabData.value.image = base64String;
                        previewUrl.value = base64String;
                        imageUrlValid.value = true;
                        $q.notify({ type: 'positive', message: 'Image uploaded successfully', position: 'top', timeout: 2000 });
                    }
                } else {
                    clearImageInput();
                }
            } catch (onloadError) {
                console.error('Error inside FileReader onload:', onloadError);
                $q.notify({ type: 'negative', message: 'Error processing image data.', position: 'top' });
                clearImageInput();
            } finally {
                isUploadingImage.value = false;
            }
        }

        reader.onerror = (error) => {
            console.error('FileReader error:', error);
            $q.notify({ type: 'negative', message: 'Error reading file.', position: 'top' });
            clearImageInput();
            isUploadingImage.value = false;
        };

        reader.readAsDataURL(file);

    } catch (error) {
        console.error('Error in handleFile:', error);
        $q.notify({ type: 'negative', message: 'An unexpected error occurred handling the file.', position: 'top' });
        clearImageInput();
        isUploadingImage.value = false;
    }
}

function removeImage(event?: Event) {
    if (event) event.stopPropagation();
    clearImageInput();
}

function clearImageInput() {
    if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
        URL.revokeObjectURL(previewUrl.value);
    }
    // Reset preview and image data based on the potentially existing image in localCabData or default
    const initialImage = props.cabData?.image || props.defaultImageUrl;
    localCabData.value.image = initialImage;
    previewUrl.value = initialImage;
    imageUrlValid.value = true; // Assume original/default is valid initially
    if (fileInput.value) {
        fileInput.value.value = '';
    }
}

async function handleFileSelect(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
        const file = input.files[0];
        if (!ALLOWED_TYPES.includes(file.type as AllowedMimeType)) { clearImageInput(); return; }
        if (file.size > MAX_FILE_SIZE) { clearImageInput(); return; }
        await handleFile(file);
    }
}

function triggerFileInput() {
    fileInput.value?.click();
}

// --- Drag & Drop --- 
function handleDragLeave(event: DragEvent) {
    const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
    const x = event.clientX; const y = event.clientY;
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

</script>

<template>
    <q-dialog :model-value="modelValue" persistent @update:model-value="closeDialog" @hide="resetForm">
        <q-card style="min-width: 400px; max-width: 95vw">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Edit Cab</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section>
                <!-- Use localCabData for v-model -->
                <q-form @submit.prevent="handleUpdateCab" class="q-gutter-sm">
                    <q-input v-model="capitalizedName" label="Cab Name" dense outlined required lazy-rules
                        :rules="[val => !!val || 'Name is required']">
                        <template v-slot:prepend>
                            <q-icon name="directions_car" />
                        </template>
                    </q-input>
                    <div class="row q-col-gutter-sm">
                        <div class="col-12 col-sm-6">
                            <q-select v-model="localCabData.make" :options="makes" label="Make" dense outlined required
                                emit-value map-options placeholder="Select a make" lazy-rules
                                :rules="[val => !!val || 'Make is required']">
                                <template v-slot:prepend>
                                    <q-icon name="business" />
                                </template>
                                <template v-slot:no-option>
                                    <q-item><q-item-section class="text-grey">No results</q-item-section></q-item>
                                </template>
                            </q-select>
                        </div>
                        <div class="col-12 col-sm-6">
                            <q-select v-model="localCabData.unit_color" :options="colors" label="Color" dense outlined
                                required emit-value map-options placeholder="Select a color" lazy-rules
                                :rules="[val => !!val || 'Color is required']">
                                <template v-slot:prepend>
                                    <q-icon name="palette" />
                                </template>
                                <template v-slot:no-option>
                                    <q-item><q-item-section class="text-grey">No results</q-item-section></q-item>
                                </template>
                            </q-select>
                        </div>
                    </div>
                    <div class="row q-col-gutter-sm">
                        <div class="col-12 col-sm-6">
                            <q-input v-model.number="localCabData.quantity" type="number" min="0" label="Quantity" dense
                                outlined required lazy-rules
                                :rules="[val => val !== null && val !== undefined && val >= 0 || 'Quantity must be positive']">
                                <template v-slot:prepend>
                                    <q-icon name="numbers" />
                                </template>
                            </q-input>
                        </div>
                        <div class="col-12 col-sm-6">
                            <q-input v-model.number="localCabData.price" type="number" label="Price" dense outlined
                                required lazy-rules
                                :rules="[val => val !== null && val !== undefined && val > 0 || 'Price must be greater than 0']">
                                <template v-slot:prepend>
                                    <q-icon name="attach_money" />
                                </template>
                            </q-input>
                        </div>
                    </div>
                    <div class="row q-col-gutter-sm">
                        <div class="col-12">
                            <q-input v-model="localCabData.status" label="Status" dense outlined readonly>
                                <template v-slot:prepend>
                                    <q-icon name="info" />
                                </template>
                            </q-input>
                        </div>
                    </div>
                    <!-- Edit Image Upload Section -->
                    <div class="row q-col-gutter-sm">
                        <div class="col-12">
                            <div class="upload-container q-pa-md" :class="{ 'dragging': isDragging }"
                                @dragenter.prevent="isDragging = true" @dragover.prevent="isDragging = true"
                                @dragleave.prevent="handleDragLeave" @drop.prevent="handleDrop"
                                @dragend.prevent="isDragging = false" @click="triggerFileInput">
                                <input type="file" ref="fileInput" accept="image/jpeg,image/png,image/gif"
                                    class="hidden" @change="handleFileSelect">
                                <div v-if="isUploadingImage" class="text-center">
                                    <q-spinner-dots color="primary" size="40px" />
                                    <div class="text-body1 q-mt-sm">Processing image...</div>
                                </div>
                                <div v-else-if="!previewUrl || previewUrl === defaultImageUrl" class="text-center">
                                    <q-icon name="cloud_upload" size="48px" color="primary" />
                                    <div class="text-body1 q-mt-sm">Drag/drop or click to select image</div>
                                    <div class="text-caption text-grey q-mt-sm">JPG, PNG, GIF | Max 5MB | Max 4096px
                                    </div>
                                </div>
                                <div v-else class="row items-center justify-center">
                                    <div class="col-8 text-center">
                                        <img :src="previewUrl" class="preview-image"
                                            :alt="localCabData.name || 'Preview'"
                                            @error="previewUrl = defaultImageUrl" />
                                    </div>
                                    <div class="col-4 text-center">
                                        <q-btn flat round color="negative" icon="close" @click.stop="removeImage">
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
                <q-btn flat label="Cancel" @click="closeDialog" />
                <q-btn unelevated color="primary" label="Update Cab" @click="handleUpdateCab"
                    :disable="!localCabData.name || !localCabData.make || !localCabData.unit_color || localCabData.quantity < 0 || localCabData.price <= 0 || !imageUrlValid || isUploadingImage" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<style scoped lang="sass">
/* Styles for upload container etc. will be moved here */
.upload-container
  position: relative
  border: 2px dashed #ccc
  border-radius: 8px
  cursor: pointer
  transition: all 0.3s ease
  min-height: 200px
  display: flex
  align-items: center
  justify-content: center
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

.preview-image
  width: 100%
  max-height: 180px
  object-fit: contain
  border-radius: 4px
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1)
  transition: all 0.3s ease

  &:hover
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15)

.hidden
  display: none
</style>
