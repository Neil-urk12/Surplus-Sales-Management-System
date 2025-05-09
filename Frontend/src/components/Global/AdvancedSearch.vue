<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue';

const props = defineProps({
    modelValue: {
        type: String,
        default: ''
    },
    placeholder: {
        type: String,
        default: 'Search'
    },
    debounceTime: {
        type: Number,
        default: 300
    },
    showClearButton: {
        type: Boolean,
        default: true
    },
    outlined: {
        type: Boolean,
        default: true
    },
    dense: {
        type: Boolean,
        default: true
    },
    color: {
        type: String,
        default: 'primary'
    }
});

const emit = defineEmits(['update:modelValue', 'search', 'clear']);

const inputValue = ref(props.modelValue);
let debounceTimeout: ReturnType<typeof setTimeout> | null = null;

watch(() => props.modelValue, (newValue) => {
    inputValue.value = newValue;
});

watch(inputValue, (newValue) => {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
    }

    debounceTimeout = setTimeout(() => {
        emit('update:modelValue', newValue);
        emit('search', newValue);
    }, props.debounceTime);
});

function clearSearch() {
    inputValue.value = '';
    emit('update:modelValue', '');
    emit('search', '');
    emit('clear');
}
</script>

<template>
    <div class="advanced-search-container">
        <q-input v-model="inputValue" :outlined="outlined" :dense="dense" :placeholder="placeholder" class="full-width"
            :color="color">
            <template v-slot:prepend>
                <q-icon name="search" />
            </template>
            <template v-slot:append v-if="showClearButton && inputValue">
                <q-icon name="close" class="cursor-pointer" @click="clearSearch" aria-label="Clear search" />
            </template>
        </q-input>
    </div>
</template>

<style lang="sass" scoped>
.advanced-search-container
  width: 100%
</style>
