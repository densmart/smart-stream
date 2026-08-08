<script setup>
import { ref, computed, watch } from 'vue';
import { mediaApi } from '@/api';
import { useToast } from 'primevue/usetoast';

const props = defineProps({
    modelValue: {
        type: String,
        default: ''
    },
    visible: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['update:modelValue', 'update:visible']);

const toast = useToast();

// State
const currentPath = ref('');
const items = ref([]);
const loading = ref(false);
const selectedFile = ref(null);

// Computed
const breadcrumbs = computed(() => {
    if (!currentPath.value) {
        return [{ label: 'Root', path: '' }];
    }

    const parts = currentPath.value.split('/').filter(Boolean);
    const crumbs = [{ label: 'Root', path: '' }];

    let accumulatedPath = '';
    for (const part of parts) {
        accumulatedPath = accumulatedPath ? `${accumulatedPath}/${part}` : part;
        crumbs.push({ label: part, path: accumulatedPath });
    }

    return crumbs;
});

const canSelectFile = computed(() => {
    return selectedFile.value !== null && selectedFile.value.type === 'file';
});

// Methods
const loadDirectory = async (path = '') => {
    try {
        loading.value = true;
        selectedFile.value = null;

        const response = await mediaApi.browseMediaFiles(path);
        items.value = response.result.items || [];
        currentPath.value = response.result.current_path || '';
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to load directory',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const handleItemClick = (item) => {
    if (item.type === 'directory') {
        // Для папок - сразу переходим
        loadDirectory(item.path);
    } else {
        // Для файлов - выбираем
        selectedFile.value = item;
    }
};

const navigateToBreadcrumb = (path) => {
    loadDirectory(path);
};

const confirmSelection = () => {
    if (canSelectFile.value) {
        // Используем путь как есть (относительный путь без начального слеша)
        emit('update:modelValue', selectedFile.value.path);
        closeDialog();
    }
};

const closeDialog = () => {
    emit('update:visible', false);
    selectedFile.value = null;
};

const formatFileSize = (bytes) => {
    if (!bytes) return 'N/A';
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    if (bytes === 0) return '0 Bytes';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round((bytes / Math.pow(1024, i)) * 100) / 100 + ' ' + sizes[i];
};

// Watch for dialog visibility
watch(
    () => props.visible,
    (newVal) => {
        if (newVal) {
            selectedFile.value = null;

            // Если есть текущий путь к файлу, открываем папку с этим файлом
            if (props.modelValue) {
                // Убираем начальный слеш если есть
                let path = props.modelValue.startsWith('/') ? props.modelValue.substring(1) : props.modelValue;

                // Извлекаем путь к директории (убираем имя файла)
                const parts = path.split('/');
                if (parts.length > 1) {
                    parts.pop(); // Убираем последний элемент (имя файла)
                    const dirPath = parts.join('/');
                    loadDirectory(dirPath);
                } else {
                    // Файл в корне
                    loadDirectory('');
                }
            } else {
                // Нет пути - открываем корень
                currentPath.value = '';
                loadDirectory('');
            }
        }
    }
);
</script>

<template>
    <Dialog :visible="visible" @update:visible="emit('update:visible', $event)" header="Select Media File" :modal="true" :closable="true" :style="{ width: '800px', maxHeight: '80vh' }" class="file-browser-dialog">
        <div class="flex flex-col gap-4">
            <!-- Breadcrumb Navigation -->
            <div class="breadcrumb-wrapper">
                <Breadcrumb :model="breadcrumbs" class="mb-0">
                    <template #item="{ item }">
                        <a class="cursor-pointer text-primary hover:underline" @click.prevent="navigateToBreadcrumb(item.path)">
                            {{ item.label }}
                        </a>
                    </template>
                </Breadcrumb>
            </div>

            <!-- File/Folder List -->
            <div v-if="loading" class="flex justify-center py-8">
                <ProgressSpinner />
            </div>

            <div v-else-if="items.length === 0" class="flex flex-col items-center justify-center py-8">
                <i class="pi pi-folder-open text-6xl text-surface-400 mb-4"></i>
                <p class="text-muted-color">No files or folders found</p>
            </div>

            <div v-else class="file-list-container" style="max-height: 400px; overflow-y: auto">
                <DataView :value="items" layout="list">
                    <template #list="{ items }">
                        <div class="flex flex-col gap-2">
                            <div
                                v-for="item in items"
                                :key="item.path"
                                class="file-item p-3 border rounded-border cursor-pointer transition-colors"
                                :class="{
                                    'bg-green-50 border-green-200 dark:bg-green-500/10 dark:border-green-500/30': selectedFile?.path === item.path,
                                    'border-surface-200 dark:border-surface-600 hover:bg-surface-50 dark:hover:bg-surface-700/50': selectedFile?.path !== item.path
                                }"
                                @click="handleItemClick(item)"
                            >
                                <div class="flex items-center gap-3">
                                    <i :class="item.type === 'directory' ? 'pi pi-folder text-yellow-500' : 'pi pi-file-video text-blue-500'" class="text-2xl"></i>
                                    <div class="flex-1">
                                        <div class="font-semibold">{{ item.name }}</div>
                                        <div class="text-sm text-muted-color flex gap-4">
                                            <span v-if="item.type === 'file' && item.format"> Format: {{ item.format.toUpperCase() }} </span>
                                            <span v-if="item.type === 'file' && item.size"> Size: {{ formatFileSize(item.size) }} </span>
                                        </div>
                                    </div>
                                    <i v-if="item.type === 'directory'" class="pi pi-chevron-right text-muted-color"></i>
                                </div>
                            </div>
                        </div>
                    </template>
                </DataView>
            </div>
        </div>

        <template #footer>
            <Button label="Cancel" icon="pi pi-times" text @click="closeDialog" />
            <Button label="Select" icon="pi pi-check" @click="confirmSelection" :disabled="!canSelectFile" />
        </template>
    </Dialog>
</template>

<style scoped>
.file-browser-dialog :deep(.p-dialog-content) {
    overflow: visible;
}

.file-item {
    user-select: none;
}

.file-list-container {
    scrollbar-width: thin;
}

.file-list-container::-webkit-scrollbar {
    width: 8px;
}

.file-list-container::-webkit-scrollbar-track {
    background: var(--surface-100);
}

.file-list-container::-webkit-scrollbar-thumb {
    background: var(--surface-300);
    border-radius: 4px;
}

.file-list-container::-webkit-scrollbar-thumb:hover {
    background: var(--surface-400);
}
</style>
