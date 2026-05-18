<script setup>
import { ref, onMounted, computed } from 'vue';
import { mediaApi, uploadApi, playlistsApi } from '@/api';
import { useToast } from 'primevue/usetoast';
import FileBrowser from '@/components/FileBrowser.vue';

const toast = useToast();

// State
const mediaList = ref([]);
const loading = ref(false);
const totalRecords = ref(0);
const displayDialog = ref(false);
const displayDeleteDialog = ref(false);
const isEditMode = ref(false);
const selectedMedia = ref(null);
const uploadingPoster = ref(false);
const showFileBrowser = ref(false);

// Filter by playlist
const filteredPlaylists = ref([]);
const selectedPlaylist = ref(null);
const loadingPlaylists = ref(false);

// Form data
const formData = ref({
    name: '',
    format: '',
    path: '',
    poster: null,
    duration: null,
    size: null
});

// Poster preview
const posterFile = ref(null);
const posterPreviewUrl = ref(null);

// Pagination
const lazyParams = ref({
    page: 1,
    limit: 12, // Using 12 for a nice grid layout
    search: ''
});

// Format options
const formatOptions = [
    { label: 'MP4', value: 'mp4' },
    { label: 'MKV', value: 'mkv' },
    { label: 'AVI', value: 'avi' },
    { label: 'MOV', value: 'mov' },
    { label: 'WEBM', value: 'webm' },
    { label: 'FLV', value: 'flv' },
    { label: 'WMV', value: 'wmv' },
    { label: 'M4V', value: 'm4v' },
    { label: '3GP', value: '3gp' },
    { label: 'TS', value: 'ts' }
];

// Methods
const loadMedia = async () => {
    try {
        loading.value = true;
        const params = {
            offset: (lazyParams.value.page - 1) * lazyParams.value.limit,
            limit: lazyParams.value.limit,
            search: lazyParams.value.search || undefined
        };

        // Filter by playlist or show only unassigned
        if (selectedPlaylist.value) {
            params.playlist_id = selectedPlaylist.value.id;
        } else {
            params.only_unassigned = true;
        }

        const response = await mediaApi.getMedia(params);
        mediaList.value = response.result;
        totalRecords.value = response.pagination.total;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to load media',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

// Autocomplete search for playlists using new API endpoint
const searchPlaylistsForFilter = async (event) => {
    try {
        loadingPlaylists.value = true;
        const query = event.query;
        // Call new search API endpoint with name parameter
        const results = await playlistsApi.searchPlaylists(query || undefined);
        filteredPlaylists.value = results;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: 'Failed to search playlists',
            life: 3000
        });
        filteredPlaylists.value = [];
    } finally {
        loadingPlaylists.value = false;
    }
};

// Apply playlist filter
const applyPlaylistFilter = () => {
    lazyParams.value.page = 1;
    loadMedia();
};

// Clear playlist filter
const clearPlaylistFilter = () => {
    selectedPlaylist.value = null;
    applyPlaylistFilter();
};

const onPage = (event) => {
    lazyParams.value.page = event.page + 1;
    lazyParams.value.limit = event.rows;
    loadMedia();
};

const openCreateDialog = () => {
    isEditMode.value = false;
    formData.value = {
        name: '',
        format: 'mp4',
        path: '',
        poster: null,
        duration: null,
        size: null
    };
    posterFile.value = null;
    posterPreviewUrl.value = null;
    displayDialog.value = true;
};

const openEditDialog = (media) => {
    isEditMode.value = true;
    selectedMedia.value = media;
    formData.value = {
        name: media.name,
        format: media.format,
        path: media.path,
        poster: media.poster,
        duration: media.duration,
        size: media.size
    };
    posterFile.value = null;
    posterPreviewUrl.value = media.poster ? uploadApi.getPosterUrl(media.poster) : null;
    displayDialog.value = true;
};

const hideDialog = () => {
    displayDialog.value = false;
    selectedMedia.value = null;
    posterFile.value = null;
    posterPreviewUrl.value = null;
};

const onPosterSelect = (event) => {
    const file = event.files[0];
    if (file) {
        posterFile.value = file;
        posterPreviewUrl.value = URL.createObjectURL(file);
    }
};

const clearPoster = () => {
    posterFile.value = null;
    posterPreviewUrl.value = null;
    formData.value.poster = null;
};

const openFileBrowser = () => {
    showFileBrowser.value = true;
};

const saveMedia = async () => {
    if (!formData.value.name || !formData.value.format || !formData.value.path) {
        toast.add({
            severity: 'warn',
            summary: 'Validation Error',
            detail: 'Name, format, and path are required',
            life: 3000
        });
        return;
    }

    try {
        loading.value = true;

        // Upload poster if a new file was selected
        let posterPath = formData.value.poster;
        if (posterFile.value) {
            try {
                uploadingPoster.value = true;
                const uploadResponse = await uploadApi.uploadPoster(posterFile.value);
                posterPath = uploadResponse.filename;
            } catch (error) {
                toast.add({
                    severity: 'error',
                    summary: 'Upload Error',
                    detail: error.response?.data?.error || 'Failed to upload poster',
                    life: 3000
                });
                return;
            } finally {
                uploadingPoster.value = false;
            }
        }

        if (isEditMode.value) {
            const updateData = {
                name: formData.value.name,
                format: formData.value.format,
                path: formData.value.path
            };

            // Only include duration if provided
            if (formData.value.duration !== null && formData.value.duration !== undefined) {
                updateData.duration = formData.value.duration;
            }

            // Only include size if provided
            if (formData.value.size !== null && formData.value.size !== undefined) {
                updateData.size = formData.value.size;
            }

            // Only include poster if it was changed
            if (posterPath !== selectedMedia.value.poster) {
                updateData.poster = posterPath;
            }

            await mediaApi.updateMedia(selectedMedia.value.id, updateData);
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Media updated successfully',
                life: 3000
            });
        } else {
            const createData = {
                name: formData.value.name,
                format: formData.value.format,
                path: formData.value.path,
                poster: posterPath
            };

            // Only include duration if provided
            if (formData.value.duration !== null && formData.value.duration !== undefined) {
                createData.duration = formData.value.duration;
            }

            // Only include size if provided
            if (formData.value.size !== null && formData.value.size !== undefined) {
                createData.size = formData.value.size;
            }

            await mediaApi.createMedia(createData);
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Media created successfully',
                life: 3000
            });
        }

        hideDialog();
        loadMedia();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to save media',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const confirmDelete = (media) => {
    selectedMedia.value = media;
    displayDeleteDialog.value = true;
};

const deleteMedia = async () => {
    try {
        loading.value = true;
        await mediaApi.deleteMedia(selectedMedia.value.id);
        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: 'Media deleted successfully',
            life: 3000
        });
        displayDeleteDialog.value = false;
        selectedMedia.value = null;
        loadMedia();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to delete media',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

const formatDuration = (seconds) => {
    if (!seconds || seconds === 0) return 'N/A';

    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    if (hours > 0) {
        return `${hours}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
    } else {
        return `${minutes}:${String(secs).padStart(2, '0')}`;
    }
};

const formatFileSize = (bytes) => {
    if (!bytes || bytes === 0) return 'N/A';
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round((bytes / Math.pow(1024, i)) * 100) / 100 + ' ' + sizes[i];
};

const getPosterUrl = (posterPath) => {
    return uploadApi.getPosterUrl(posterPath);
};

onMounted(() => {
    loadMedia();
});
</script>

<template>
    <div class="card">
        <Toast />

        <div class="flex justify-between items-center mb-6">
            <h5 class="mb-0">Media Management</h5>
            <Button label="Add Media" icon="pi pi-plus" @click="openCreateDialog" />
        </div>

        <!-- Filter by Playlist -->
        <div class="mb-4">
            <div class="flex items-center">
                <div class="relative inline-block">
                    <AutoComplete
                        id="playlistFilter"
                        v-model="selectedPlaylist"
                        :suggestions="filteredPlaylists"
                        @complete="searchPlaylistsForFilter"
                        optionLabel="name"
                        placeholder="Filter by playlist"
                        :loading="loadingPlaylists"
                        :inputStyle="{ width: '24rem', paddingRight: selectedPlaylist ? '2.5rem' : '0.75rem' }"
                    >
                        <template #option="slotProps">
                            <div class="flex items-center gap-2">
                                <i class="pi pi-list text-surface-500"></i>
                                <span>{{ slotProps.option.name }}</span>
                            </div>
                        </template>
                    </AutoComplete>
                    <Button
                        v-if="selectedPlaylist"
                        icon="pi pi-times"
                        text
                        rounded
                        severity="secondary"
                        @click="clearPlaylistFilter"
                        title="Clear filter"
                        class="!absolute"
                        style="right: 0.25rem; top: 50%; transform: translateY(-50%); width: 2rem; height: 2rem; z-index: 10;"
                    />
                </div>
                <Button
                    icon="pi pi-filter"
                    label="Filter"
                    @click="applyPlaylistFilter"
                    title="Apply filter"
                    class="ml-2"
                />
            </div>
        </div>

        <!-- Grid View -->
        <div v-if="!loading && mediaList.length > 0" class="grid grid-cols-12 gap-4 mb-4">
            <div v-for="media in mediaList" :key="media.id" class="col-span-12 sm:col-span-6 md:col-span-4 xl:col-span-3">
                <div class="card mb-0 p-4">
                    <div class="relative mb-4">
                        <img
                            v-if="media.poster"
                            :src="getPosterUrl(media.poster)"
                            :alt="media.name"
                            class="w-full h-48 object-cover rounded-border"
                        />
                        <div
                            v-else
                            class="w-full h-48 bg-surface-100 dark:bg-surface-700 rounded-border flex items-center justify-center"
                        >
                            <i class="pi pi-video text-6xl text-surface-400"></i>
                        </div>
                        <Tag :value="media.format.toUpperCase()" class="absolute top-2 right-2" severity="info" />
                    </div>
                    <div class="mb-3">
                        <h6 class="mb-1">{{ media.name }}</h6>
                        <p class="text-sm text-muted-color mb-1">{{ media.path }}</p>
                        <div class="flex gap-3 mb-1">
                            <p class="text-xs text-muted-color">
                                <i class="pi pi-clock mr-1"></i>{{ formatDuration(media.duration) }}
                            </p>
                            <p class="text-xs text-muted-color">
                                <i class="pi pi-database mr-1"></i>{{ formatFileSize(media.size) }}
                            </p>
                        </div>
                        <p class="text-xs text-muted-color">Created: {{ formatDate(media.created_at) }}</p>
                    </div>
                    <div class="flex gap-2">
                        <Button
                            icon="pi pi-pencil"
                            outlined
                            rounded
                            size="small"
                            @click="openEditDialog(media)"
                        />
                        <Button
                            icon="pi pi-trash"
                            outlined
                            rounded
                            severity="danger"
                            size="small"
                            @click="confirmDelete(media)"
                        />
                    </div>
                </div>
            </div>
        </div>

        <!-- Empty State -->
        <div
            v-if="!loading && mediaList.length === 0"
            class="flex flex-col items-center justify-center py-12"
        >
            <i class="pi pi-video text-6xl text-surface-400 mb-4"></i>
            <p class="text-xl text-muted-color mb-4">No media files found</p>
            <Button label="Add Your First Media" icon="pi pi-plus" @click="openCreateDialog" />
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="flex justify-center py-12">
            <ProgressSpinner />
        </div>

        <!-- Pagination -->
        <Paginator
            v-if="totalRecords > 0"
            :rows="lazyParams.limit"
            :totalRecords="totalRecords"
            :rowsPerPageOptions="[12, 24, 48]"
            @page="onPage"
            template="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
            currentPageReportTemplate="Showing {first} to {last} of {totalRecords} media files"
        ></Paginator>

        <!-- Create/Edit Dialog -->
        <Dialog
            v-model:visible="displayDialog"
            :header="isEditMode ? 'Edit Media' : 'Create Media'"
            :modal="true"
            :closable="true"
            class="p-fluid"
            style="width: 600px"
        >
            <div class="flex flex-col gap-6 py-4">
                <div class="flex flex-col gap-2">
                    <label for="name">Name *</label>
                    <InputText
                        id="name"
                        v-model="formData.name"
                        required="true"
                        autofocus
                        :class="{ 'p-invalid': !formData.name }"
                    />
                </div>

                <div class="flex flex-col gap-2">
                    <label for="format">Format *</label>
                    <Select
                        id="format"
                        v-model="formData.format"
                        :options="formatOptions"
                        optionLabel="label"
                        optionValue="value"
                        placeholder="Select a format"
                        :class="{ 'p-invalid': !formData.format }"
                    />
                </div>

                <div class="flex flex-col gap-2">
                    <label for="path">File Path *</label>
                    <div class="flex gap-2">
                        <InputText
                            id="path"
                            v-model="formData.path"
                            readonly
                            placeholder="Click 'Browse' to select a file"
                            :class="{ 'p-invalid': !formData.path }"
                            class="flex-1"
                        />
                        <Button
                            icon="pi pi-folder-open"
                            label="Browse"
                            @click="openFileBrowser"
                            outlined
                        />
                    </div>
                    <small class="text-muted-color">Select a video file from the server storage</small>
                </div>

                <div class="flex flex-col gap-2">
                    <label for="duration">Duration (seconds)</label>
                    <InputNumber
                        id="duration"
                        v-model="formData.duration"
                        :min="0"
                        placeholder="Auto-detected if left empty"
                    />
                    <small class="text-muted-color">Optional - will be auto-detected from file if not provided</small>
                </div>

                <div class="flex flex-col gap-2">
                    <label for="size">File Size (bytes)</label>
                    <InputNumber
                        id="size"
                        v-model="formData.size"
                        :min="0"
                        placeholder="Auto-detected if left empty"
                    />
                    <small class="text-muted-color">Optional - will be auto-detected from file if not provided</small>
                </div>

                <div class="flex flex-col gap-2">
                    <label>Poster Image</label>
                    <div v-if="posterPreviewUrl" class="mb-3">
                        <div class="relative inline-block">
                            <img :src="posterPreviewUrl" alt="Poster preview" class="max-w-full h-48 rounded-border" />
                            <Button
                                icon="pi pi-times"
                                rounded
                                severity="danger"
                                class="absolute top-2 right-2"
                                @click="clearPoster"
                            />
                        </div>
                    </div>
                    <FileUpload
                        mode="basic"
                        accept="image/*"
                        :maxFileSize="5000000"
                        :auto="false"
                        chooseLabel="Choose Poster"
                        @select="onPosterSelect"
                        :disabled="uploadingPoster"
                    />
                    <small class="text-muted-color">Max file size: 5MB. Supported: JPG, PNG, WEBP</small>
                </div>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="hideDialog" :disabled="loading || uploadingPoster" />
                <Button label="Save" icon="pi pi-check" @click="saveMedia" :loading="loading || uploadingPoster" />
            </template>
        </Dialog>

        <!-- Delete Confirmation Dialog -->
        <Dialog
            v-model:visible="displayDeleteDialog"
            header="Confirm Delete"
            :modal="true"
            :closable="true"
            style="width: 450px"
        >
            <div class="flex items-center gap-4">
                <i class="pi pi-exclamation-triangle !text-3xl text-orange-500" />
                <span v-if="selectedMedia">
                    Are you sure you want to delete media <b>{{ selectedMedia.name }}</b>?
                </span>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="displayDeleteDialog = false" />
                <Button label="Delete" icon="pi pi-trash" severity="danger" @click="deleteMedia" :loading="loading" />
            </template>
        </Dialog>

        <!-- File Browser Dialog -->
        <FileBrowser
            v-model="formData.path"
            v-model:visible="showFileBrowser"
        />
    </div>
</template>
