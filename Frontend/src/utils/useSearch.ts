import { ref, watch } from 'vue';

interface UseSearchOptions {
    debounceTime?: number;
    initialValue?: string;
    onSearch?: (value: string) => void;
}

export function useSearch(options: UseSearchOptions = {}) {
    const {
        debounceTime = 300,
        initialValue = '',
        onSearch
    } = options;

    const searchInput = ref(initialValue);
    const searchValue = ref(initialValue);
    let debounceTimeout: ReturnType<typeof setTimeout> | null = null;

    function updateSearch(value: string) {
        if (debounceTimeout) {
            clearTimeout(debounceTimeout);
        }

        debounceTimeout = setTimeout(() => {
            searchValue.value = value;

            if (onSearch) {
                onSearch(value);
            }
        }, debounceTime);
    }

    function clearSearch() {
        searchInput.value = '';
        searchValue.value = '';

        if (onSearch) {
            onSearch('');
        }
    }

    watch(searchInput, (newValue) => {
        updateSearch(newValue);
    });

    return {
        searchInput,
        searchValue,
        updateSearch,
        clearSearch
    };
} 
