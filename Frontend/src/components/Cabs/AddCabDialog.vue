<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';
import type { NewCabInput } from 'src/types/cabs';
import type { PropType } from 'vue';
import { validateFileUpload, validateAndSanitizeBase64Image, MAX_IMAGE_SIZE, DEFAULT_MAX_DIMENSION } from '../../utils/imageValidation'; // Updated import
import { operationNotifications } from '../../utils/notifications'; // Adjusted path
// Import getNextFallbackImage or handle fallbacks differently if needed
// import { getNextFallbackImage } from 'src/config/defaultImages'; 

const props = defineProps({
    modelValue: {
        type: Boolean,
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
    (e: 'add-cab', cabData: NewCabInput): void
}>();

const $q = useQuasar();

// Constants for file validation
const MAX_FILE_SIZE = MAX_IMAGE_SIZE;
const MAX_DIMENSION = DEFAULT_MAX_DIMENSION;

// --- State --- 
const newCab = ref<NewCabInput>({
    name: '',
    make: '',
    quantity: 0,
    price: 0,
    unit_color: '',
    status: 'Out of Stock',
    image: props.defaultImageUrl, // Initialize with prop
});
const imageUrlValid = ref(true);
const validatingImage = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref<string>(props.defaultImageUrl); // Initialize with prop
const isUploadingImage = ref(false);
const isDragging = ref(false);
let currentAbortController: AbortController | null = null;

// --- Computed --- 
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

// --- Watchers --- 
watch(() => newCab.value.quantity, (newQuantity) => {
    if (newQuantity === 0) {
        newCab.value.status = 'Out of Stock';
    } else if (newQuantity <= 7) {
        newCab.value.status = 'Low Stock';
    } else {
        newCab.value.status = 'In Stock';
    }
});

watch(() => newCab.value.image, (newUrl: string) => {
    if (!newUrl || newUrl === props.defaultImageUrl) {
        previewUrl.value = props.defaultImageUrl;
        imageUrlValid.value = true;
        return;
    }
    if (newUrl.startsWith('data:image/')) {
        const validationResult = validateAndSanitizeBase64Image(newUrl);
        if (validationResult.isValid) {
            newCab.value.image = validationResult.sanitizedData!;
            previewUrl.value = validationResult.sanitizedData!;
            imageUrlValid.value = true;
        } else {
            $q.notify({
                color: 'negative',
                message: validationResult.error || 'Invalid image data',
                position: 'top',
            });
            previewUrl.value = props.defaultImageUrl;
            imageUrlValid.value = false;
        }
        return;
    }

    // Validate HTTP/HTTPS URLs
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

// Reset state when dialog opens
watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        resetForm();
    }
});

// --- Lifecycle Hooks ---
onBeforeUnmount(() => {
    // Abort any ongoing image URL validation when the component is unmounted
    if (currentAbortController) {
        console.log('Aborting ongoing image validation on unmount');
        currentAbortController.abort();
        currentAbortController = null;
    }
});

// --- Functions --- 

function resetForm() {
    newCab.value = {
        name: '',
        make: '',
        quantity: 0,
        price: 0,
        unit_color: '',
        status: 'Out of Stock',
        image: props.defaultImageUrl,
    };
    clearImageInput(); // Also clears preview, file input etc.
}

function handleAddNewCab() {
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

    // Ensure image is set, fallback to default if necessary
    if (!newCab.value.image || newCab.value.image === '') {
        newCab.value.image = props.defaultImageUrl;
    }

    // Emit the data to the parent
    emit('add-cab', { ...newCab.value });
}

function closeDialog() {
    emit('update:modelValue', false);
}

// --- Image Handling --- 

async function validateImageUrl(url: string): Promise<boolean> {
    if (!url || !url.startsWith('http')) {
        imageUrlValid.value = false;
        return false;
    }

    validatingImage.value = true;

    if (currentAbortController) {
        currentAbortController.abort();
    }
    currentAbortController = new AbortController();
    const signal = currentAbortController.signal;

    try {
        const result = await new Promise<boolean>((resolve) => {
            const img = new Image();
            const cleanup = () => {
                img.onload = null;
                img.onerror = null;
                if (currentAbortController?.signal === signal) {
                    currentAbortController = null;
                }
            };

            signal.addEventListener('abort', () => { cleanup(); resolve(false); });
            img.onload = () => { cleanup(); imageUrlValid.value = true; validatingImage.value = false; resolve(true); };
            img.onerror = () => {
                cleanup();
                // Basic fallback: if validation fails, mark as invalid
                // Fallback logic (getNextFallbackImage) could be added here if needed
                imageUrlValid.value = false;
                validatingImage.value = false;
                resolve(false);
            };

            const timeoutId = setTimeout(() => {
                if (!signal.aborted) {
                    currentAbortController?.abort();
                    imageUrlValid.value = false;
                    validatingImage.value = false;
                    resolve(false);
                }
            }, 5000);

            signal.addEventListener('abort', () => { clearTimeout(timeoutId); });
            img.src = url;
        });
        return result;
    } catch (error) {
        console.error('Error validating image URL:', error);
        imageUrlValid.value = false;
        validatingImage.value = false;
        return false;
    } finally {
        if (validatingImage.value) {
            validatingImage.value = false; // Ensure this is always reset
        }
    }
}

function handleFile(file: File) {
    isUploadingImage.value = true;
    
    validateFileUpload(file, {
        maxFileSize: MAX_FILE_SIZE,
        maxDimension: MAX_DIMENSION
    })
        .then(result => {
            if (result.isValid && result.sanitizedData) {
                newCab.value.image = result.sanitizedData;
                previewUrl.value = result.sanitizedData;
                imageUrlValid.value = true;
            } else {
                $q.notify({
                    color: 'negative',
                    message: result.error || 'Invalid image',
                    position: 'top',
                });
                clearImageInput();
            }
        })
        .catch(error => {
            console.error('Error processing image:', error);
            $q.notify({
                color: 'negative',
                message: 'Error processing image',
                position: 'top',
            });
            clearImageInput();
        })
        .finally(() => {
            isUploadingImage.value = false;
        });
}

function removeImage(event?: Event) {
    if (event) event.stopPropagation();
    clearImageInput();
}

function clearImageInput() {
    if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
        URL.revokeObjectURL(previewUrl.value);
    }
    previewUrl.value = props.defaultImageUrl;
    newCab.value.image = props.defaultImageUrl; // Reset model image too
    imageUrlValid.value = true;
    if (fileInput.value) {
        fileInput.value.value = '';
    }
    isUploadingImage.value = false;
}

function handleFileSelect(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
        const file = input.files[0];
        void handleFile(file);
    }
    if (input) {
        input.value = '';
    }
}

function triggerFileInput() {
    fileInput.value?.click();
}

// --- Drag & Drop --- 
function handleDragLeave(event: DragEvent) {
    const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
    const x = event.clientX;
    const y = event.clientY;
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
                <div class="text-h6">New Cab</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section>
                <q-form @submit.prevent="handleAddNewCab" class="q-gutter-sm">
                    <q-input v-model="capitalizedName" label="Cab Name" dense outlined required lazy-rules
                        :rules="[val => !!val || 'Name is required']">
                        <template v-slot:prepend>
                            <q-icon name="directions_car" />
                        </template>
                    </q-input>

                    <div class="row q-col-gutter-sm">
                        <div class="col-12 col-sm-6">
                            <q-select v-model="newCab.make" :options="makes" label="Make" dense outlined required
                                emit-value map-options placeholder="Select a make" lazy-rules
                                :rules="[val => !!val || 'Make is required']">
                                <template v-slot:prepend>
                                    <q-icon name="business" />
                                </template>
                                <template v-slot:no-option>
                                    <q-item>
                                        <q-item-section class="text-grey">
                                            No results
                                        </q-item-section>
                                    </q-item>
                                </template>
                            </q-select>
                        </div>

                        <div class="col-12 col-sm-6">
                            <q-select v-model="newCab.unit_color" :options="colors" label="Color" dense outlined
                                required emit-value map-options placeholder="Select a color" lazy-rules
                                :rules="[val => !!val || 'Color is required']">
                                <template v-slot:prepend>
                                    <q-icon name="palette" />
                                </template>
                                <template v-slot:no-option>
                                    <q-item>
                                        <q-item-section class="text-grey">
                                            No results
                                        </q-item-section>
                                    </q-item>
                                </template>
                            </q-select>
                        </div>
                    </div>

                    <div class="row q-col-gutter-sm">
                        <div class="col-12 col-sm-6">
                            <q-input v-model.number="newCab.quantity" type="number" min="0" label="Quantity" dense
                                outlined required lazy-rules
                                :rules="[val => val !== null && val !== undefined && val >= 0 || 'Quantity must be positive']">
                                <template v-slot:prepend>
                                    <q-icon name="numbers" />
                                </template>
                            </q-input>
                        </div>

                        <div class="col-12 col-sm-6">
                            <q-input v-model.number="newCab.price" type="number" label="Price" dense outlined required
                                lazy-rules
                                :rules="[val => val !== null && val !== undefined && val > 0 || 'Price must be greater than 0']">
                                <template v-slot:prepend>
                                    <q-icon name="attach_money" />
                                </template>
                            </q-input>
                        </div>
                    </div>

                    <div class="row q-col-gutter-sm">
                        <div class="col-12">
                            <q-input v-model="newCab.status" label="Status" dense outlined readonly>
                                <template v-slot:prepend>
                                    <q-icon name="info" />
                                </template>
                            </q-input>
                        </div>
                    </div>

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
                                    <div class="text-body1 q-mt-sm">
                                        Processing image...
                                    </div>
                                </div>
                                <div v-else-if="!previewUrl || previewUrl === defaultImageUrl" class="text-center">
                                    <q-icon name="cloud_upload" size="48px" color="primary" />
                                    <div class="text-body1 q-mt-sm">
                                        Drag and drop an image here or click to select
                                    </div>
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
                                        <img :src="previewUrl" class="preview-image"
                                            :alt="newCab.name || 'Preview image'"
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
                <q-btn unelevated color="primary" label="Add Cab" @click="handleAddNewCab"
                    :disable="!newCab.name || !newCab.make || !newCab.unit_color || newCab.quantity < 0 || newCab.price <= 0 || !imageUrlValid || isUploadingImage" />
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
