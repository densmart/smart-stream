<script setup>
import { ref, onMounted, computed } from 'vue';
import { playlistsApi, mediaApi, uploadApi } from '@/api';
import { useToast } from 'primevue/usetoast';
import { PlaylistType } from '@/types/api';

const toast = useToast();

// State
const playlists = ref([]);
const loading = ref(false);
const totalRecords = ref(0);
const displayDialog = ref(false);
const displayDeleteDialog = ref(false);
const displayMediaDialog = ref(false);
const displayPosterDialog = ref(false);
const isEditMode = ref(false);
const selectedPlaylist = ref(null);
const selectedPosterUrl = ref(null);
const uploadingPoster = ref(false);

// Hierarchy navigation
const currentParentId = ref(null);
const breadcrumbs = ref([{ label: 'All Series', id: null }]);

// Computed properties for hierarchy
const isInsideSeries = computed(() => currentParentId.value !== null);
const createButtonLabel = computed(() => (isInsideSeries.value ? 'Add Season' : 'Add Series'));

// Parent playlist autocomplete
const loadingPlaylists = ref(false);
const filteredPlaylists = ref([]);

// Form data
const formData = ref({
    name: '',
    type: PlaylistType.SERIES,
    poster: null,
    parent_id: null
});

// Poster preview
const posterFile = ref(null);
const posterPreviewUrl = ref(null);

// Media management
const playlistMedia = ref([]);
const availableMedia = ref([]);
const availableMediaTotal = ref(0);
const availableMediaSearch = ref('');
const availableMediaPage = ref(1);
const availableMediaLimit = ref(50);
const loadingMedia = ref(false);
const selectedMediaToAdd = ref([]);
const selectedMediaToRemove = ref([]);

// Upload Season functionality
const displayUploadSeasonDialog = ref(false);
const uploadFiles = ref([]);
const uploadFormat = ref('mp4');
const uploading = ref(false);
const uploadProgress = ref(0);
const uploadProgressText = ref('');
const browseMediaPath = ref('');
const browseMediaItems = ref([]);
const browseMediaBreadcrumbs = ref([]);
const loadingBrowse = ref(false);

// Pagination
const lazyParams = ref({
    page: 1,
    limit: 24,
    search: ''
});

// Search functionality
const searchQuery = ref('');

// Methods
const loadPlaylists = async () => {
    try {
        loading.value = true;
        const params = {
            offset: (lazyParams.value.page - 1) * lazyParams.value.limit,
            limit: lazyParams.value.limit,
            name: lazyParams.value.search || undefined,
            type: isInsideSeries.value ? PlaylistType.SEASON : PlaylistType.SERIES
        };

        // Add parent_id filter if we're viewing children
        if (currentParentId.value) {
            params.parent_id = currentParentId.value;
        }

        const response = await playlistsApi.getPlaylists(params);
        playlists.value = response.result;
        totalRecords.value = response.pagination.total;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to load series',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const onPage = (event) => {
    lazyParams.value.page = event.page + 1;
    loadPlaylists();
};

// Search handler
const onSearch = () => {
    lazyParams.value.search = searchQuery.value;
    lazyParams.value.page = 1;
    loadPlaylists();
};

// Hierarchy navigation methods
const viewChildren = (playlist) => {
    currentParentId.value = playlist.id;
    breadcrumbs.value.push({ label: playlist.name, id: playlist.id });
    lazyParams.value.page = 1;
    loadPlaylists();
};

const navigateToBreadcrumb = (index) => {
    const breadcrumb = breadcrumbs.value[index];
    currentParentId.value = breadcrumb.id;
    breadcrumbs.value = breadcrumbs.value.slice(0, index + 1);
    lazyParams.value.page = 1;
    loadPlaylists();
};

// Autocomplete search using new API endpoint
const searchPlaylists = async (event) => {
    try {
        loadingPlaylists.value = true;
        const query = event.query;
        // Call new search API endpoint with name parameter and type filter
        const results = await playlistsApi.searchPlaylists(query || undefined, PlaylistType.SERIES);
        filteredPlaylists.value = results;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: 'Failed to search series',
            life: 3000
        });
        filteredPlaylists.value = [];
    } finally {
        loadingPlaylists.value = false;
    }
};

const openCreateDialog = () => {
    isEditMode.value = false;
    formData.value = {
        name: '',
        type: isInsideSeries.value ? PlaylistType.SEASON : PlaylistType.SERIES,
        poster: null,
        parent_id: isInsideSeries.value ? currentParentId.value : null
    };
    posterFile.value = null;
    posterPreviewUrl.value = null;
    displayDialog.value = true;
};

const openEditDialog = async (playlist) => {
    isEditMode.value = true;
    selectedPlaylist.value = playlist;

    // Load parent playlist object if parent_id is set
    let parentPlaylist = null;
    if (playlist.parent_id) {
        try {
            parentPlaylist = await playlistsApi.getPlaylist(playlist.parent_id);
        } catch (error) {
            console.error('Failed to load parent playlist:', error);
        }
    }

    formData.value = {
        name: playlist.name,
        type: playlist.type,
        poster: playlist.poster,
        parent_id: parentPlaylist || null
    };
    posterFile.value = null;
    posterPreviewUrl.value = playlist.poster ? uploadApi.getPosterUrl(playlist.poster) : null;
    displayDialog.value = true;
};

const hideDialog = () => {
    displayDialog.value = false;
    selectedPlaylist.value = null;
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

const savePlaylist = async () => {
    if (!formData.value.name) {
        toast.add({
            severity: 'warn',
            summary: 'Validation Error',
            detail: 'Name is required',
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
                type: selectedPlaylist.value.type
            };

            // Only include poster if it was changed
            if (posterPath !== selectedPlaylist.value.poster) {
                updateData.poster = posterPath;
            }

            // Extract parent_id (can be object or null)
            updateData.parent_id = formData.value.parent_id?.id || null;

            await playlistsApi.updatePlaylist(selectedPlaylist.value.id, updateData);
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: selectedPlaylist.value.type === PlaylistType.SEASON ? 'Season updated successfully' : 'Series updated successfully',
                life: 3000
            });
        } else {
            await playlistsApi.createPlaylist({
                name: formData.value.name,
                type: isInsideSeries.value ? PlaylistType.SEASON : PlaylistType.SERIES,
                poster: posterPath,
                parent_id: isInsideSeries.value ? currentParentId.value : formData.value.parent_id?.id || null
            });
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: isInsideSeries.value ? 'Season created successfully' : 'Series created successfully',
                life: 3000
            });
        }

        hideDialog();
        loadPlaylists();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to save series',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const confirmDelete = (playlist) => {
    selectedPlaylist.value = playlist;
    displayDeleteDialog.value = true;
};

const deletePlaylist = async () => {
    try {
        loading.value = true;
        const itemType = selectedPlaylist.value.type === PlaylistType.SEASON ? 'Season' : 'Series';
        await playlistsApi.deletePlaylist(selectedPlaylist.value.id);
        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: `${itemType} deleted successfully`,
            life: 3000
        });
        displayDeleteDialog.value = false;
        selectedPlaylist.value = null;
        loadPlaylists();
    } catch (error) {
        const itemType = selectedPlaylist.value.type === PlaylistType.SEASON ? 'season' : 'series';
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || `Failed to delete ${itemType}`,
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const loadAvailableMedia = async () => {
    if (!selectedPlaylist.value) return;

    try {
        loadingMedia.value = true;

        // Load all media with pagination and search
        const params = {
            page: availableMediaPage.value,
            limit: availableMediaLimit.value,
            only_unassigned: true
        };

        if (availableMediaSearch.value) {
            params.name = availableMediaSearch.value;
        }

        const response = await mediaApi.getMedia(params);
        availableMedia.value = response.result;
        availableMediaTotal.value = response.pagination.total;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to load available media',
            life: 3000
        });
    } finally {
        loadingMedia.value = false;
    }
};

const openMediaDialog = async (playlist) => {
    // Check if playlist has children
    if (playlist.has_children) {
        toast.add({
            severity: 'warn',
            summary: 'Cannot Manage Episodes',
            detail: 'Episodes can only be added to series without seasons. This series contains seasons.',
            life: 5000
        });
        return;
    }

    selectedPlaylist.value = playlist;
    selectedMediaToAdd.value = [];
    selectedMediaToRemove.value = [];
    availableMediaSearch.value = '';
    availableMediaPage.value = 1;
    displayMediaDialog.value = true;

    // Load playlist media and available media
    try {
        loadingMedia.value = true;
        const mediaInPlaylist = await playlistsApi.getPlaylistMedia(playlist.id, { page: 1, limit: 1000 });
        playlistMedia.value = mediaInPlaylist.result;
        loadingMedia.value = false;

        // Load available media separately with filters
        await loadAvailableMedia();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to load media',
            life: 3000
        });
        loadingMedia.value = false;
    }
};

const onAvailableMediaSearch = async () => {
    availableMediaPage.value = 1;
    await loadAvailableMedia();
};

const onAvailableMediaPage = async (event) => {
    availableMediaPage.value = event.page + 1;
    await loadAvailableMedia();
};

const addMediaToPlaylist = async () => {
    if (selectedMediaToAdd.value.length === 0) {
        toast.add({
            severity: 'warn',
            summary: 'No Selection',
            detail: 'Please select episodes to add',
            life: 3000
        });
        return;
    }

    try {
        loadingMedia.value = true;

        // Generate order values starting from current playlist media count
        const startOrder = playlistMedia.value.length;

        // Use batch API for multiple media
        const mediaToAdd = selectedMediaToAdd.value.map((media, index) => ({
            mediaId: media.id,
            order: startOrder + index
        }));

        await playlistsApi.batchAddMediaToPlaylist(selectedPlaylist.value.id, mediaToAdd);

        const itemType = selectedPlaylist.value.type === PlaylistType.SEASON ? 'season' : 'series';
        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: `Added ${selectedMediaToAdd.value.length} episode(s) to ${itemType}`,
            life: 3000
        });

        // Clear selection
        selectedMediaToAdd.value = [];

        // Refresh both lists
        const mediaInPlaylist = await playlistsApi.getPlaylistMedia(selectedPlaylist.value.id, { page: 1, limit: 1000 });
        playlistMedia.value = mediaInPlaylist.result;
        await loadAvailableMedia();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to add episodes',
            life: 3000
        });
    } finally {
        loadingMedia.value = false;
    }
};

const removeMediaFromPlaylist = async () => {
    if (selectedMediaToRemove.value.length === 0) {
        toast.add({
            severity: 'warn',
            summary: 'No Selection',
            detail: 'Please select episodes to remove',
            life: 3000
        });
        return;
    }

    try {
        loadingMedia.value = true;

        // Use batch API for multiple media
        const mediaIds = selectedMediaToRemove.value.map((media) => media.id);
        await playlistsApi.batchRemoveMediaFromPlaylist(selectedPlaylist.value.id, mediaIds);

        const itemType = selectedPlaylist.value.type === PlaylistType.SEASON ? 'season' : 'series';
        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: `Removed ${selectedMediaToRemove.value.length} episode(s) from ${itemType}`,
            life: 3000
        });

        // Clear selection
        selectedMediaToRemove.value = [];

        // Refresh both lists
        const mediaInPlaylist = await playlistsApi.getPlaylistMedia(selectedPlaylist.value.id, { page: 1, limit: 1000 });
        playlistMedia.value = mediaInPlaylist.result;
        await loadAvailableMedia();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to remove episodes',
            life: 3000
        });
    } finally {
        loadingMedia.value = false;
    }
};

const onMediaReorder = async (event) => {
    try {
        loadingMedia.value = true;

        // Update local array with new order
        const reorderedMedia = event.value;

        // Collect all updates that need to be made
        const updates = reorderedMedia
            .map((media, index) => ({
                mediaId: media.id,
                order: index
            }))
            .filter((update, index) => reorderedMedia[index].order !== update.order);

        // Use batch API if there are any updates needed
        if (updates.length > 0) {
            await playlistsApi.batchUpdateMediaOrder(selectedPlaylist.value.id, updates);
        }

        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: 'Media order updated',
            life: 3000
        });

        // Refresh playlist media list
        const mediaInPlaylist = await playlistsApi.getPlaylistMedia(selectedPlaylist.value.id, { page: 1, limit: 1000 });
        playlistMedia.value = mediaInPlaylist.result;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to update order',
            life: 3000
        });
    } finally {
        loadingMedia.value = false;
    }
};

const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
};

const formatFileSize = (bytes) => {
    if (!bytes) return 'N/A';
    const kb = bytes / 1024;
    const mb = kb / 1024;
    const gb = mb / 1024;

    if (gb >= 1) return `${gb.toFixed(2)} GB`;
    if (mb >= 1) return `${mb.toFixed(2)} MB`;
    if (kb >= 1) return `${kb.toFixed(2)} KB`;
    return `${bytes} bytes`;
};

const getPosterUrl = (posterPath) => {
    return uploadApi.getPosterUrl(posterPath);
};

const openPosterDialog = (posterPath) => {
    if (posterPath) {
        selectedPosterUrl.value = getPosterUrl(posterPath);
        displayPosterDialog.value = true;
    }
};

// Upload Season methods
const openUploadSeasonDialog = () => {
    uploadFiles.value = [];
    uploadFormat.value = 'mp4';
    uploadProgress.value = 0;
    uploadProgressText.value = '';
    browseMediaPath.value = '';
    browseMediaItems.value = [];
    browseMediaBreadcrumbs.value = [];
    displayUploadSeasonDialog.value = true;

    // Load root directory
    loadBrowseMediaFiles('');
};

const loadBrowseMediaFiles = async (path) => {
    try {
        loadingBrowse.value = true;
        const response = await mediaApi.browseMediaFiles(path);
        browseMediaItems.value = response.result.items;
        browseMediaPath.value = response.result.current_path;

        // Build breadcrumbs
        if (!path || path === '') {
            browseMediaBreadcrumbs.value = [{ label: 'Root', path: '' }];
        } else {
            const parts = path.split('/').filter((p) => p);
            browseMediaBreadcrumbs.value = [{ label: 'Root', path: '' }];
            let currentPath = '';
            parts.forEach((part) => {
                currentPath += '/' + part;
                browseMediaBreadcrumbs.value.push({ label: part, path: currentPath });
            });
        }
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Browse Error',
            detail: error.response?.data?.error || 'Failed to browse files',
            life: 3000
        });
    } finally {
        loadingBrowse.value = false;
    }
};

const navigateToBrowsePath = (path) => {
    loadBrowseMediaFiles(path);
};

const toggleFileSelection = (item) => {
    const index = uploadFiles.value.findIndex((f) => f.path === item.path);
    if (index > -1) {
        uploadFiles.value.splice(index, 1);
    } else {
        uploadFiles.value.push(item);
    }
};

const isFileSelected = (item) => {
    return uploadFiles.value.some((f) => f.path === item.path);
};

const chooseAllFiles = () => {
    // Фильтруем только файлы (не папки)
    const allFiles = browseMediaItems.value.filter(item => item.type === 'file');

    // Добавляем все файлы, которых еще нет в выбранных
    allFiles.forEach(file => {
        if (!isFileSelected(file)) {
            uploadFiles.value.push(file);
        }
    });
};

const uploadSeason = async () => {
    // Validation
    if (uploadFiles.value.length === 0) {
        toast.add({
            severity: 'warn',
            summary: 'No Files',
            detail: 'Please select at least one video file',
            life: 3000
        });
        return;
    }

    if (!selectedPlaylist.value) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: 'No playlist selected',
            life: 3000
        });
        return;
    }

    uploading.value = true;
    const createdMediaIds = [];

    try {
        // Sort files by name to ensure consistent order
        const sortedFiles = [...uploadFiles.value].sort((a, b) => a.name.localeCompare(b.name));

        // Upload files one by one
        for (let i = 0; i < sortedFiles.length; i++) {
            const file = sortedFiles[i];
            uploadProgressText.value = `Creating Эпизод ${i + 1} of ${sortedFiles.length}...`;
            uploadProgress.value = ((i + 1) / sortedFiles.length) * 100;

            try {
                // Create media entry for this file
                const mediaData = {
                    name: `Эпизод ${i + 1}`,
                    format: uploadFormat.value,
                    path: file.path,
                    poster: null
                };

                const createdMedia = await mediaApi.createMedia(mediaData);
                createdMediaIds.push(createdMedia.id);
            } catch (error) {
                console.error(`Failed to create media for ${file.name}:`, error);
                toast.add({
                    severity: 'error',
                    summary: 'Error',
                    detail: `Failed to create Эпизод ${i + 1}: ${error.response?.data?.error || error.message}`,
                    life: 5000
                });
                // Continue with next file
            }
        }

        if (createdMediaIds.length === 0) {
            toast.add({
                severity: 'error',
                summary: 'Upload Failed',
                detail: 'No episodes were created successfully',
                life: 5000
            });
            return;
        }

        // Now add all created media to the playlist using batch API
        uploadProgressText.value = 'Adding episodes to playlist...';
        const mediaToAdd = createdMediaIds.map((mediaId, index) => ({
            mediaId: mediaId,
            order: index
        }));

        await playlistsApi.batchAddMediaToPlaylist(selectedPlaylist.value.id, mediaToAdd);

        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: `Successfully uploaded ${createdMediaIds.length} episodes`,
            life: 3000
        });

        displayUploadSeasonDialog.value = false;

        // Refresh playlist media
        const mediaInPlaylist = await playlistsApi.getPlaylistMedia(selectedPlaylist.value.id, { page: 1, limit: 1000 });
        playlistMedia.value = mediaInPlaylist.result;
        await loadAvailableMedia();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Upload Error',
            detail: error.response?.data?.error || 'Failed to upload episodes',
            life: 5000
        });
    } finally {
        uploading.value = false;
        uploadProgress.value = 0;
        uploadProgressText.value = '';
    }
};

onMounted(() => {
    loadPlaylists();
});
</script>

<template>
    <div class="card">
        <Toast />

        <div class="flex justify-between items-center mb-6">
            <h5 class="mb-0">Series Management</h5>
            <Button :label="createButtonLabel" icon="pi pi-plus" @click="openCreateDialog" />
        </div>

        <!-- Search Field -->
        <div class="mb-4">
            <IconField>
                <InputIcon class="pi pi-search" />
                <InputText v-model="searchQuery" placeholder="Search..." @keyup.enter="onSearch" @input="onSearch" class="w-full md:w-96" />
            </IconField>
        </div>

        <!-- Breadcrumb Navigation -->
        <div v-if="breadcrumbs.length > 1" class="mb-4">
            <Breadcrumb
                :model="
                    breadcrumbs.map((bc, idx) => ({
                        label: bc.label,
                        command: () => navigateToBreadcrumb(idx)
                    }))
                "
            />
        </div>

        <!-- Grid View -->
        <div v-if="!loading && playlists.length > 0" class="grid grid-cols-12 gap-4 mb-4">
            <div v-for="playlist in playlists" :key="playlist.id" class="col-span-12 sm:col-span-6 md:col-span-4 xl:col-span-3">
                <div class="card mb-0 p-4">
                    <div class="relative mb-4">
                        <img
                            v-if="playlist.poster"
                            :src="getPosterUrl(playlist.poster)"
                            :alt="playlist.name"
                            class="w-full h-48 object-cover rounded-border cursor-pointer hover:opacity-90 transition-opacity"
                            @click="openPosterDialog(playlist.poster)"
                        />
                        <div v-else class="w-full h-48 bg-surface-100 dark:bg-surface-700 rounded-border flex items-center justify-center">
                            <i class="pi pi-list text-6xl text-surface-400"></i>
                        </div>
                        <Tag :value="playlist.type === PlaylistType.SEASON ? 'Season' : 'Series'" severity="success" class="absolute top-2 right-2" />
                    </div>
                    <div class="mb-3">
                        <h6 class="mb-1">{{ playlist.name }}</h6>
                        <p class="text-xs text-muted-color">Created: {{ formatDate(playlist.created_at) }}</p>
                    </div>

                    <!-- Buttons for root level (Series) -->
                    <div v-if="!currentParentId" class="flex gap-2">
                        <Button icon="pi pi-pencil" outlined rounded size="small" @click="openEditDialog(playlist)" />
                        <Button v-if="playlist.has_children" icon="pi pi-folder-open" outlined rounded severity="success" size="small" @click="viewChildren(playlist)" title="Seasons" />
                        <template v-else>
                            <Button icon="pi pi-video" outlined rounded severity="info" size="small" @click="openMediaDialog(playlist)" title="Episodes" />
                            <Button icon="pi pi-folder-open" outlined rounded severity="success" size="small" @click="viewChildren(playlist)" title="Seasons" />
                        </template>
                        <Button icon="pi pi-trash" outlined rounded severity="danger" size="small" @click="confirmDelete(playlist)" />
                    </div>

                    <!-- Buttons for inside series (Seasons) -->
                    <div v-else class="flex gap-2">
                        <Button icon="pi pi-pencil" outlined rounded size="small" @click="openEditDialog(playlist)" />
                        <Button icon="pi pi-video" outlined rounded severity="info" size="small" @click="openMediaDialog(playlist)" title="Episodes" />
                        <Button icon="pi pi-trash" outlined rounded severity="danger" size="small" @click="confirmDelete(playlist)" />
                    </div>
                </div>
            </div>
        </div>

        <!-- Empty State -->
        <div v-if="!loading && playlists.length === 0" class="flex flex-col items-center justify-center py-12">
            <i class="pi pi-list text-6xl text-surface-400 mb-4"></i>
            <p class="text-xl text-muted-color mb-4">No series found</p>
            <Button label="Create Your First Series" icon="pi pi-plus" @click="openCreateDialog" />
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
            @page="onPage"
            template="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport"
            currentPageReportTemplate="Showing {first} to {last} of {totalRecords} series"
        ></Paginator>

        <!-- Create/Edit Dialog -->
        <Dialog
            v-model:visible="displayDialog"
            :header="isEditMode ? (selectedPlaylist?.type === PlaylistType.SEASON ? 'Edit Season' : 'Edit Series') : isInsideSeries ? 'Create Season' : 'Create Series'"
            :modal="true"
            :closable="true"
            class="p-fluid"
            style="width: 600px"
        >
            <div class="flex flex-col gap-6 py-4">
                <div class="flex flex-col gap-2">
                    <label for="name">Name *</label>
                    <InputText id="name" v-model="formData.name" required="true" autofocus :class="{ 'p-invalid': !formData.name }" />
                </div>

                <div v-if="!isInsideSeries" class="flex flex-col gap-2">
                    <label for="parent">Parent Playlist (optional)</label>
                    <AutoComplete id="parent" v-model="formData.parent_id" :suggestions="filteredPlaylists" @complete="searchPlaylists" optionLabel="name" placeholder="Search parent series..." :loading="loadingPlaylists">
                        <template #option="slotProps">
                            <div class="flex items-center gap-2">
                                <i class="pi pi-list text-surface-500"></i>
                                <span>{{ slotProps.option.name }}</span>
                            </div>
                        </template>
                    </AutoComplete>
                    <small class="text-muted-color">Leave empty to make this a root-level series</small>
                </div>

                <div class="flex flex-col gap-2">
                    <label>Poster Image</label>
                    <div v-if="posterPreviewUrl" class="mb-3">
                        <div class="relative inline-block">
                            <img :src="posterPreviewUrl" alt="Poster preview" class="max-w-full h-48 rounded-border" />
                            <Button icon="pi pi-times" rounded severity="danger" class="absolute top-2 right-2" @click="clearPoster" />
                        </div>
                    </div>
                    <FileUpload mode="basic" accept="image/*" :maxFileSize="5000000" :auto="false" chooseLabel="Choose Poster" @select="onPosterSelect" :disabled="uploadingPoster" />
                    <small class="text-muted-color">Max file size: 5MB. Supported: JPG, PNG, WEBP</small>
                </div>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="hideDialog" :disabled="loading || uploadingPoster" />
                <Button label="Save" icon="pi pi-check" @click="savePlaylist" :loading="loading || uploadingPoster" />
            </template>
        </Dialog>

        <!-- Delete Confirmation Dialog -->
        <Dialog v-model:visible="displayDeleteDialog" header="Confirm Delete" :modal="true" :closable="true" style="width: 450px">
            <div class="flex items-center gap-4">
                <i class="pi pi-exclamation-triangle !text-3xl text-orange-500" />
                <span v-if="selectedPlaylist">
                    Are you sure you want to delete {{ selectedPlaylist.type === PlaylistType.SEASON ? 'season' : 'series' }} <b>{{ selectedPlaylist.name }}</b
                    >?
                </span>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="displayDeleteDialog = false" />
                <Button label="Delete" icon="pi pi-trash" severity="danger" @click="deletePlaylist" :loading="loading" />
            </template>
        </Dialog>

        <!-- Media Management Dialog -->
        <Dialog v-model:visible="displayMediaDialog" :header="selectedPlaylist?.type === PlaylistType.SEASON ? 'Manage Season Episodes' : 'Manage Series Episodes'" :modal="true" :closable="true" class="p-fluid" style="width: 900px; max-width: 95vw">
            <div v-if="selectedPlaylist" class="mb-4">
                <h6>{{ selectedPlaylist.name }}</h6>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <!-- Available Media -->
                <div class="col-span-2 md:col-span-1">
                    <div class="flex justify-between items-center mb-3">
                        <h6 class="mb-0">Available Episodes ({{ availableMediaTotal }})</h6>
                        <Button label="Add Selected" icon="pi pi-plus" size="small" @click="addMediaToPlaylist" :disabled="loadingMedia || selectedMediaToAdd.length === 0" />
                    </div>

                    <!-- Search -->
                    <div class="mb-3">
                        <IconField>
                            <InputIcon class="pi pi-search" />
                            <InputText v-model="availableMediaSearch" placeholder="Search media..." @keyup.enter="onAvailableMediaSearch" :disabled="loadingMedia" class="w-full" />
                        </IconField>
                    </div>

                    <div v-if="loadingMedia" class="flex justify-center py-8">
                        <ProgressSpinner style="width: 50px; height: 50px" />
                    </div>
                    <DataTable v-else v-model:selection="selectedMediaToAdd" :value="availableMedia" scrollable scrollHeight="350px" :rowHover="true" class="text-sm">
                        <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>
                        <Column field="name" header="Name"></Column>
                        <Column field="format" header="Format">
                            <template #body="{ data }">
                                <Tag :value="data.format.toUpperCase()" severity="info" />
                            </template>
                        </Column>
                    </DataTable>
                    <p v-if="!loadingMedia && availableMedia.length === 0" class="text-center text-muted-color py-4">No available media found</p>

                    <!-- Pagination -->
                    <Paginator
                        v-if="!loadingMedia && availableMediaTotal > availableMediaLimit"
                        :rows="availableMediaLimit"
                        :totalRecords="availableMediaTotal"
                        :first="(availableMediaPage - 1) * availableMediaLimit"
                        @page="onAvailableMediaPage"
                        template="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport"
                        currentPageReportTemplate="Showing {first} to {last} of {totalRecords}"
                        class="mt-3"
                    ></Paginator>
                </div>

                <!-- Media in Playlist -->
                <div class="col-span-2 md:col-span-1">
                    <div class="flex justify-between items-center mb-3">
                        <div class="flex items-center gap-2">
                            <h6 class="mb-0">Media in {{ selectedPlaylist?.type === PlaylistType.SEASON ? 'Season' : 'Series' }} ({{ playlistMedia.length }})</h6>
                            <Button v-if="playlistMedia.length === 0" label="Upload Season" icon="pi pi-upload" size="small" @click="openUploadSeasonDialog" :disabled="loadingMedia" />
                        </div>
                        <Button label="Remove Selected" icon="pi pi-minus" size="small" severity="danger" @click="removeMediaFromPlaylist" :disabled="loadingMedia || selectedMediaToRemove.length === 0" />
                    </div>
                    <div v-if="loadingMedia" class="flex justify-center py-8">
                        <ProgressSpinner style="width: 50px; height: 50px" />
                    </div>
                    <DataTable v-else v-model:selection="selectedMediaToRemove" :value="playlistMedia" :reorderableRows="true" @row-reorder="onMediaReorder" scrollable scrollHeight="400px" :rowHover="true" class="text-sm">
                        <Column :rowReorder="true" headerStyle="width: 3rem" />
                        <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>
                        <Column field="order" header="Order" headerStyle="width: 4rem">
                            <template #body="{ data }">
                                <Tag :value="data.order.toString()" severity="secondary" />
                            </template>
                        </Column>
                        <Column field="name" header="Name"></Column>
                        <Column field="format" header="Format" headerStyle="width: 5rem">
                            <template #body="{ data }">
                                <Tag :value="data.format.toUpperCase()" severity="info" />
                            </template>
                        </Column>
                    </DataTable>
                    <p v-if="!loadingMedia && playlistMedia.length === 0" class="text-center text-muted-color py-4">No episodes in this {{ selectedPlaylist?.type === PlaylistType.SEASON ? 'season' : 'series' }}</p>
                </div>
            </div>

            <template #footer>
                <Button label="Close" icon="pi pi-times" @click="displayMediaDialog = false" />
            </template>
        </Dialog>

        <!-- Poster Preview Dialog -->
        <Dialog v-model:visible="displayPosterDialog" :modal="true" :closable="true" :dismissableMask="true" header="Poster Preview" class="poster-preview-dialog" style="width: auto; max-width: 90vw">
            <div class="flex justify-center items-center">
                <img v-if="selectedPosterUrl" :src="selectedPosterUrl" alt="Poster" class="max-w-full max-h-[80vh] object-contain rounded-border" />
            </div>
        </Dialog>

        <!-- Upload Season Dialog -->
        <Dialog v-model:visible="displayUploadSeasonDialog" header="Upload Season Episodes" :modal="true" :closable="!uploading" class="p-fluid" style="width: 800px; max-width: 95vw">
            <div class="flex flex-col gap-4 py-4">
                <!-- Format Selection -->
                <div class="flex flex-col gap-2">
                    <label for="format">Video Format (same for all episodes)</label>
                    <Select id="format" v-model="uploadFormat" :options="['mp4', 'mkv', 'avi', 'mov', 'webm', 'flv', 'wmv', 'm4v']" placeholder="Select format" :disabled="uploading" />
                </div>

                <!-- File Browser -->
                <div class="flex flex-col gap-2">
                    <label>Browse and Select Video Files</label>

                    <!-- Breadcrumbs -->
                    <div v-if="browseMediaBreadcrumbs.length > 0" class="mb-2">
                        <Breadcrumb
                            :model="
                                browseMediaBreadcrumbs.map((bc) => ({
                                    label: bc.label,
                                    command: () => navigateToBrowsePath(bc.path)
                                }))
                            "
                        />
                    </div>

                    <!-- Choose All Button -->
                    <div class="flex justify-between items-center mb-2">
                        <span class="text-sm text-muted-color">Click files to select</span>
                        <Button
                            label="Choose All"
                            icon="pi pi-check-square"
                            size="small"
                            text
                            @click="chooseAllFiles"
                            :disabled="browseMediaItems.filter(i => i.type === 'file').length === 0 || uploading"
                        />
                    </div>

                    <!-- File List -->
                    <div v-if="loadingBrowse" class="flex justify-center py-8">
                        <ProgressSpinner style="width: 50px; height: 50px" />
                    </div>
                    <div v-else class="border border-surface-200 dark:border-surface-700 rounded-border p-3 max-h-96 overflow-y-auto">
                        <div v-if="browseMediaItems.length === 0" class="text-center text-muted-color py-4">No files or folders found</div>
                        <div
                            v-for="item in browseMediaItems"
                            :key="item.path"
                            class="flex items-center gap-3 p-2 hover:bg-surface-50 dark:hover:bg-surface-700/50 rounded-border cursor-pointer"
                            @click="item.type === 'directory' ? navigateToBrowsePath(item.path) : toggleFileSelection(item)"
                        >
                            <Checkbox v-if="item.type === 'file'" :modelValue="isFileSelected(item)" :disabled="uploading" @click.stop="toggleFileSelection(item)" />
                            <i v-if="item.type === 'directory'" class="pi pi-folder text-xl text-yellow-500"></i>
                            <i v-else class="pi pi-file text-xl text-blue-500"></i>
                            <div class="flex-1">
                                <div class="font-medium">{{ item.name }}</div>
                                <div v-if="item.type === 'file'" class="text-xs text-muted-color">{{ item.format?.toUpperCase() }} - {{ formatFileSize(item.size) }}</div>
                            </div>
                            <i v-if="item.type === 'directory'" class="pi pi-chevron-right text-muted-color"></i>
                        </div>
                    </div>

                    <small class="text-muted-color"> Click folders to navigate, click checkboxes to select files. Selected: {{ uploadFiles.length }} file(s) </small>
                </div>

                <!-- Selected Files Preview -->
                <div v-if="uploadFiles.length > 0" class="flex flex-col gap-2">
                    <label>Episodes Preview ({{ uploadFiles.length }} files)</label>
                    <div class="border border-surface-200 dark:border-surface-700 rounded-border p-3 max-h-48 overflow-y-auto">
                        <div class="grid gap-2">
                            <div v-for="(file, index) in uploadFiles.slice().sort((a, b) => a.name.localeCompare(b.name))" :key="file.path" class="flex items-center gap-3 p-2 bg-surface-50 dark:bg-surface-800 rounded-border">
                                <Tag :value="`#${index + 1}`" severity="secondary" />
                                <div class="flex-1">
                                    <div class="font-medium">Эпизод {{ index + 1 }}</div>
                                    <div class="text-xs text-muted-color">{{ file.name }}</div>
                                </div>
                                <Tag :value="uploadFormat.toUpperCase()" severity="info" />
                                <Button icon="pi pi-times" text rounded size="small" severity="danger" @click="toggleFileSelection(file)" :disabled="uploading" />
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Upload Progress -->
                <div v-if="uploading" class="flex flex-col gap-2">
                    <label>{{ uploadProgressText }}</label>
                    <ProgressBar :value="uploadProgress" />
                </div>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="displayUploadSeasonDialog = false" :disabled="uploading" />
                <Button label="Upload & Add to Season" icon="pi pi-check" @click="uploadSeason" :loading="uploading" :disabled="uploadFiles.length === 0" />
            </template>
        </Dialog>
    </div>
</template>
