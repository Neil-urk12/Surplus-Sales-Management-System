<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';
import type { PropType } from 'vue';
import type { MaterialRow, NewMaterialInput } from 'src/types/materials';
import { validateAndSanitizeBase64Image, validateImageDimensions } from '../../utils/imageValidation';
import { operationNotifications } from '../../utils/notifications';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    materialData: {
        type: Object as PropType<MaterialRow | null>,
        required: true,
    },
    categories: {
        type: Array as PropType<string[]>,
        required: true,
    },
    suppliers: {
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
    (e: 'update-material', materialData: NewMaterialInput): void
}>();

const $q = useQuasar();

// Constants for file validation
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB

// --- Local State --- 
const localMaterialData = ref<NewMaterialInput>({
    name: '',
    category: '',
    supplier: '',
    quantity: 0,
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
    get: () => localMaterialData.value.name,
    set: (value: string) => {
        if (value) {
            localMaterialData.value.name = value.charAt(0).toUpperCase() + value.slice(1);
        } else {
            localMaterialData.value.name = value;
        }
    }
});

// Clean up any pending AbortController when component is unmounted
onBeforeUnmount(() => {
    if (currentAbortController) {
        currentAbortController.abort();
        currentAbortController = null;
    }
    // Cleanup any object URLs to prevent memory leaks
    if (previewUrl.value && previewUrl.value.startsWith('blob:')) {
        URL.revokeObjectURL(previewUrl.value);
    }
});

// Watch for changes in the modelValue prop
watch(() => props.modelValue, (newValue) => {
    console.log('modelValue changed to:', newValue);
    // We don't need to do anything when dialog opens as materialData watcher handles that
});

// Watch for changes in materialData
watch(() => props.materialData, async (newMaterial) => {
    console.log('materialData changed:', newMaterial);
    if (newMaterial) {
        // Reset local state based on the new material data
        const initialImage = newMaterial.image || props.defaultImageUrl;
        
        localMaterialData.value = {
            name: newMaterial.name,
            category: newMaterial.category,
            supplier: newMaterial.supplier,
            quantity: newMaterial.quantity,
            status: newMaterial.status,
            image: initialImage
        };

        // Always update the preview URL to match the current material
        previewUrl.value = initialImage;

        // Directly validate the initial image
        if (initialImage.startsWith('data:image/')) {
            const validationResult = validateAndSanitizeBase64Image(initialImage);
            if (validationResult.isValid) {
                // Check dimensions
                const dimensionResult = await validateImageDimensions(initialImage);
                if (dimensionResult.isValid) {
                    previewUrl.value = initialImage;
                    imageUrlValid.value = true;
                } else {
                    previewUrl.value = props.defaultImageUrl;
                    localMaterialData.value.image = props.defaultImageUrl;
                    imageUrlValid.value = false;
                    $q.notify({
                        type: 'negative',
                        message: dimensionResult.error || 'Image dimensions exceed limits',
                        position: 'top'
                    });
                }
            } else {
                previewUrl.value = props.defaultImageUrl;
                localMaterialData.value.image = props.defaultImageUrl;
                imageUrlValid.value = false;
                $q.notify({
                    type: 'negative',
                    message: validationResult.error || 'Invalid image data',
                    position: 'top'
                });
            }
        } else if (initialImage !== props.defaultImageUrl) {
            void validateImageUrl(initialImage);
        } else {
            previewUrl.value = props.defaultImageUrl;
            imageUrlValid.value = true;
        }
    }
}, { immediate: true });

// Add watcher for quantity to update status automatically
watch(() => localMaterialData.value.quantity, (newQuantity) => {
    if (newQuantity === 0) {
        localMaterialData.value.status = 'Out of Stock';
    } else if (newQuantity <= 10) {
        localMaterialData.value.status = 'Low Stock';
    } else {
        localMaterialData.value.status = 'In Stock';
    }
});

// --- Functions --- 
function resetForm() {
    console.log('Resetting form in EditMaterialDialog');
    localMaterialData.value = {
        name: '',
        category: '',
        supplier: '',
        quantity: 0,
        status: 'Out of Stock',
        image: props.defaultImageUrl,
    };
    clearImageInput();
}

function handleUpdateMaterial() {
    console.log('handleUpdateMaterial called in EditMaterialDialog');
    if (!localMaterialData.value.name) {
        operationNotifications.validation.error('Please enter a material name');
        return;
    }
    
    if (!localMaterialData.value.category || !localMaterialData.value.supplier) {
        operationNotifications.validation.error('Please fill in all required fields');
        return;
    }
    
    if (!imageUrlValid.value) {
        operationNotifications.validation.error('Please provide a valid image');
        return;
    }
    
    // Ensure image field is not empty; use default if necessary
    if (!localMaterialData.value.image || localMaterialData.value.image === '') {
        localMaterialData.value.image = props.defaultImageUrl;
    }

    // Full sanitization right before sending to backend
    if (localMaterialData.value.image.startsWith('data:image/')) {
        const sanitizationResult = validateAndSanitizeBase64Image(localMaterialData.value.image);
        if (!sanitizationResult.isValid) {
            operationNotifications.validation.error(sanitizationResult.error || 'Invalid image data');
            return;
        }
        localMaterialData.value.image = sanitizationResult.sanitizedData!;
    }

    emit('update-material', { ...localMaterialData.value });
    // Close the dialog immediately after emitting update
    emit('update:modelValue', false);
}

function closeDialog() {
    console.log('closeDialog called in EditMaterialDialog');
    emit('update:modelValue', false);
}

// --- Image Handling --- 
async function validateImageUrl(url: string): Promise<boolean> {
    if (!url || !url.startsWith('http')) {
        imageUrlValid.value = false;
        previewUrl.value = props.defaultImageUrl;
        return false;
    }

    validatingImage.value = true;

    if (currentAbortController) {
        currentAbortController.abort();
    }
    currentAbortController = new AbortController();
    const signal = currentAbortController.signal;

    try {
        const response = await fetch(url, { signal });
        if (!response.ok) {
            throw new Error('Failed to fetch image');
        }
        const blob = await response.blob();
        if (!blob.type.startsWith('image/')) {
            throw new Error('Not an image file');
        }
        if (blob.size > MAX_FILE_SIZE) {
            throw new Error('Image too large');
        }
        
        // Create a temporary URL to validate dimensions
        const blobUrl = URL.createObjectURL(blob);
        try {
            const dimensionResult = await validateImageDimensions(blobUrl);
            if (!dimensionResult.isValid) {
                throw new Error(dimensionResult.error || 'Image dimensions exceed limits');
            }
        } finally {
            URL.revokeObjectURL(blobUrl);
        }
        
        imageUrlValid.value = true;
        previewUrl.value = url;
        return true;
    } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') {
            console.log('Image validation aborted');
            return false;
        }
        imageUrlValid.value = false;
        previewUrl.value = props.defaultImageUrl;
        $q.notify({
            type: 'negative',
            message: error instanceof Error ? error.message : 'Failed to validate image',
            position: 'top'
        });
        return false;
    } finally {
        validatingImage.value = false;
        if (currentAbortController?.signal === signal) {
            currentAbortController = null;
        }
    }
}

function handleFile(file: File) {
    if (!file.type.startsWith('image/')) {
        operationNotifications.validation.error('Please select an image file');
        return;
    }

    if (file.size > MAX_FILE_SIZE) {
        operationNotifications.validation.error('Image size must be less than 5MB');
        return;
    }

    isUploadingImage.value = true;

    try {
        const reader = new FileReader();
        reader.onload = async (e) => {
            const result = e.target?.result as string;
            if (result) {
                const validationResult = validateAndSanitizeBase64Image(result);
                if (validationResult.isValid) {
                    // Validate dimensions
                    const dimensionResult = await validateImageDimensions(validationResult.sanitizedData!);
                    if (dimensionResult.isValid) {
                        localMaterialData.value.image = validationResult.sanitizedData!;
                        previewUrl.value = validationResult.sanitizedData!;
                        imageUrlValid.value = true;
                    } else {
                        operationNotifications.validation.error(dimensionResult.error || 'Image dimensions exceed limits');
                        clearImageInput();
                    }
                } else {
                    operationNotifications.validation.error(validationResult.error || 'Invalid image data');
                    clearImageInput();
                }
            }
            isUploadingImage.value = false;
        };
        reader.readAsDataURL(file);
    } catch (error) {
        console.error('Error processing image:', error);
        operationNotifications.validation.error('Error processing image');
        clearImageInput();
        isUploadingImage.value = false;
    }
}

function clearImageInput() {
    if (previewUrl.value && previewUrl.value !== props.defaultImageUrl) {
        URL.revokeObjectURL(previewUrl.value);
    }
    previewUrl.value = props.defaultImageUrl;
    localMaterialData.value.image = props.defaultImageUrl;
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
        handleFile(file);
    }
    if (input) {
        input.value = '';
    }
}

function triggerFileInput() {
    fileInput.value?.click();
}

// Add a function to remove the image (use clearImageInput)
function removeImage() {
    clearImageInput();
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
                <div class="text-h6">Edit Material</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section>
                <q-form @submit.prevent="handleUpdateMaterial" class="q-gutter-sm">
                    <q-input v-model="capitalizedName" label="Material Name" dense outlined required lazy-rules
                        :rules="[val => !!val || 'Name is required']">
                        <template v-slot:prepend>
                            <q-icon name="inventory_2" />
                        </template>
                    </q-input>

                    <div class="row q-col-gutter-sm">
                        <div class="col-12 col-sm-6">
                            <q-select v-model="localMaterialData.category" :options="categories" label="Category" dense
                                outlined required emit-value map-options placeholder="Select a category" lazy-rules
                                :rules="[val => !!val || 'Category is required']">
                                <template v-slot:prepend>
                                    <q-icon name="category" />
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
                            <q-select v-model="localMaterialData.supplier" :options="suppliers" label="Supplier" dense
                                outlined required emit-value map-options placeholder="Select a supplier" lazy-rules
                                :rules="[val => !!val || 'Supplier is required']">
                                <template v-slot:prepend>
                                    <q-icon name="local_shipping" />
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
                            <q-input v-model.number="localMaterialData.quantity" type="number" min="0" label="Quantity"
                                dense outlined required lazy-rules
                                :rules="[val => val !== null && val !== undefined && val >= 0 || 'Quantity must be positive']">
                                <template v-slot:prepend>
                                    <q-icon name="numbers" />
                                </template>
                            </q-input>
                        </div>

                        <div class="col-12 col-sm-6">
                            <q-input v-model="localMaterialData.status" label="Status" dense outlined readonly>
                                <template v-slot:prepend>
                                    <q-icon name="info" />
                                </template>
                            </q-input>
                        </div>
                    </div>

                    <div class="row q-col-gutter-sm">
                        <div class="col-12">
                            <div class="upload-container" :class="{ dragging: isDragging }"
                                @dragover.prevent="isDragging = true" @dragleave="handleDragLeave"
                                @drop.prevent="handleDrop" @click="triggerFileInput">
                                <input ref="fileInput" type="file" accept="image/*" class="hidden"
                                    @change="handleFileSelect" />
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
                                            :alt="localMaterialData.name || 'Preview image'"
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
                <q-btn unelevated color="primary" label="Update Material" @click="handleUpdateMaterial"
                    :disable="!localMaterialData.name || !localMaterialData.category || !localMaterialData.supplier || localMaterialData.quantity < 0 || !imageUrlValid || isUploadingImage" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<style scoped lang="sass">
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