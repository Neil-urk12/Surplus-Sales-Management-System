<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue';
import type { NewMaterialInput } from 'src/types/materials';
import type { PropType } from 'vue';
import { validateAndSanitizeBase64Image } from '../../utils/imageValidation';
import { operationNotifications } from '../../utils/notifications';

const props = defineProps({
    modelValue: {
        type: Boolean,
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
    (e: 'add-material', materialData: NewMaterialInput): void
}>();

// Constants for file validation
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const MAX_DIMENSION = 4096; // Maximum image dimension in pixels

// --- State --- 
const newMaterial = ref<NewMaterialInput>({
    name: '',
    category: '',
    supplier: '',
    quantity: 0,
    status: 'Out of Stock',
    image: props.defaultImageUrl,
});
const imageUrlValid = ref(true);
const fileInput = ref<HTMLInputElement | null>(null);
const previewUrl = ref<string>(props.defaultImageUrl);
const isUploadingImage = ref(false);
const isDragging = ref(false);
let currentAbortController: AbortController | null = null;

// --- Computed --- 
const capitalizedName = computed({
    get: () => newMaterial.value.name,
    set: (value: string) => {
        if (value) {
            newMaterial.value.name = value.charAt(0).toUpperCase() + value.slice(1);
        } else {
            newMaterial.value.name = value;
        }
    }
});

// --- Watchers --- 
watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        resetForm();
    }
});

// Add watcher for quantity to update status automatically
watch(() => newMaterial.value.quantity, (newQuantity) => {
    if (newQuantity === 0) {
        newMaterial.value.status = 'Out of Stock';
    } else if (newQuantity <= 10) {
        newMaterial.value.status = 'Low Stock';
    } else {
        newMaterial.value.status = 'In Stock';
    }
});

// --- Lifecycle Hooks ---
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

// --- Functions --- 
function resetForm() {
    newMaterial.value = {
        name: '',
        category: '',
        supplier: '',
        quantity: 0,
        status: 'Out of Stock',
        image: props.defaultImageUrl,
    };
    clearImageInput();
}

function handleAddNewMaterial() {
    if (!newMaterial.value.name) {
        operationNotifications.validation.error('Please enter a material name');
        return;
    }

    if (!newMaterial.value.category || !newMaterial.value.supplier) {
        operationNotifications.validation.error('Please fill in all required fields');
        return;
    }

    if (!imageUrlValid.value) {
        operationNotifications.validation.error('Please provide a valid image');
        return;
    }

    if (!newMaterial.value.image || newMaterial.value.image === '') {
        newMaterial.value.image = props.defaultImageUrl;
    }

    // Full sanitization right before sending to backend
    if (newMaterial.value.image.startsWith('data:image/')) {
        const sanitizationResult = validateAndSanitizeBase64Image(newMaterial.value.image, MAX_DIMENSION);
        if (!sanitizationResult.isValid) {
            operationNotifications.validation.error(sanitizationResult.error || 'Invalid image data');
            return;
        }
        newMaterial.value.image = sanitizationResult.sanitizedData!;
    }

    emit('add-material', { ...newMaterial.value });
}

function closeDialog() {
    emit('update:modelValue', false);
}

// --- Image Handling --- 
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
        reader.onload = (e) => {
            const result = e.target?.result as string;
            if (result) {
                const validationResult = validateAndSanitizeBase64Image(result, MAX_DIMENSION);
                if (validationResult.isValid) {
                    newMaterial.value.image = validationResult.sanitizedData!;
                    previewUrl.value = validationResult.sanitizedData!;
                    imageUrlValid.value = true;
                } else {
                    operationNotifications.validation.error(validationResult.error || 'Invalid image data');
                    clearImageInput();
                }
            }
        };
        reader.readAsDataURL(file);
    } catch (error) {
        console.error('Error processing image:', error);
        operationNotifications.validation.error('Error processing image');
        clearImageInput();
    } finally {
        isUploadingImage.value = false;
    }
}

function clearImageInput() {
    if (previewUrl.value && previewUrl.value !== props.defaultImageUrl) {
        URL.revokeObjectURL(previewUrl.value);
    }
    previewUrl.value = props.defaultImageUrl;
    newMaterial.value.image = props.defaultImageUrl;
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
        handleFile(file);
    }
}

// Add the removeImage function
function removeImage() {
    clearImageInput();
}
</script>

<template>
    <q-dialog :model-value="modelValue" persistent @update:model-value="closeDialog" @hide="resetForm">
        <q-card style="min-width: 400px; max-width: 95vw">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">New Material</div>
                <q-space />
                <q-btn icon="close" flat round dense @click="closeDialog" />
            </q-card-section>

            <q-card-section>
                <q-form @submit.prevent="handleAddNewMaterial" class="q-gutter-sm">
                    <q-input v-model="capitalizedName" label="Material Name" dense outlined required lazy-rules
                        :rules="[val => !!val || 'Name is required']">
                        <template v-slot:prepend>
                            <q-icon name="inventory_2" />
                        </template>
                    </q-input>

                    <div class="row q-col-gutter-sm">
                        <div class="col-12 col-sm-6">
                            <q-select v-model="newMaterial.category" :options="categories" label="Category" dense outlined
                                required emit-value map-options placeholder="Select a category" lazy-rules
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
                            <q-select v-model="newMaterial.supplier" :options="suppliers" label="Supplier" dense outlined
                                required emit-value map-options placeholder="Select a supplier" lazy-rules
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
                            <q-input v-model.number="newMaterial.quantity" type="number" min="0" label="Quantity" dense
                                outlined required lazy-rules
                                :rules="[val => val !== null && val !== undefined && val >= 0 || 'Quantity must be positive']">
                                <template v-slot:prepend>
                                    <q-icon name="numbers" />
                                </template>
                            </q-input>
                        </div>

                        <div class="col-12 col-sm-6">
                            <q-input v-model="newMaterial.status" label="Status" dense outlined readonly>
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
                                    <div class="text-caption text-grey q-mt-sm">JPG, PNG, GIF | Max 5MB | Max {{ MAX_DIMENSION }}px
                                    </div>
                                </div>
                                <div v-else class="row items-center justify-center">
                                    <div class="col-8 text-center">
                                        <img :src="previewUrl" class="preview-image"
                                            :alt="newMaterial.name || 'Preview image'"
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
                <q-btn unelevated color="primary" label="Add Material" @click="handleAddNewMaterial"
                    :disable="!newMaterial.name || !newMaterial.category || !newMaterial.supplier || newMaterial.quantity < 0 || !imageUrlValid || isUploadingImage" />
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