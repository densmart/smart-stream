<script setup>
import { ref, onMounted } from 'vue';
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

// Form data
const formData = ref({
    name: '',
    type: PlaylistType.SERIES,
    poster: null
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

// Pagination
const lazyParams = ref({
    page: 1,
    limit: 12,
    search: ''
});

// Playlist type options
const typeOptions = [
    { label: 'Series', value: PlaylistType.SERIES },
    { label: 'Franchise', value: PlaylistType.FRANCHISE }
];

// Methods
const loadPlaylists = async () => {
    try {
        loading.value = true;
        const response = await playlistsApi.getPlaylists({
            page: lazyParams.value.page,
            limit: lazyParams.value.limit,
            search: lazyParams.value.search || undefined
        });
        playlists.value = response.result;
        totalRecords.value = response.pagination.total;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to load playlists',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const onPage = (event) => {
    lazyParams.value.page = event.page + 1;
    lazyParams.value.limit = event.rows;
    loadPlaylists();
};

const openCreateDialog = () => {
    isEditMode.value = false;
    formData.value = {
        name: '',
        type: PlaylistType.SERIES,
        poster: null
    };
    posterFile.value = null;
    posterPreviewUrl.value = null;
    displayDialog.value = true;
};

const openEditDialog = (playlist) => {
    isEditMode.value = true;
    selectedPlaylist.value = playlist;
    formData.value = {
        name: playlist.name,
        type: playlist.type,
        poster: playlist.poster
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
    if (!formData.value.name || !formData.value.type) {
        toast.add({
            severity: 'warn',
            summary: 'Validation Error',
            detail: 'Name and type are required',
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
                type: formData.value.type
            };

            // Only include poster if it was changed
            if (posterPath !== selectedPlaylist.value.poster) {
                updateData.poster = posterPath;
            }

            await playlistsApi.updatePlaylist(selectedPlaylist.value.id, updateData);
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Playlist updated successfully',
                life: 3000
            });
        } else {
            await playlistsApi.createPlaylist({
                name: formData.value.name,
                type: formData.value.type,
                poster: posterPath
            });
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Playlist created successfully',
                life: 3000
            });
        }

        hideDialog();
        loadPlaylists();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to save playlist',
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
        await playlistsApi.deletePlaylist(selectedPlaylist.value.id);
        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: 'Playlist deleted successfully',
            life: 3000
        });
        displayDeleteDialog.value = false;
        selectedPlaylist.value = null;
        loadPlaylists();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to delete playlist',
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
            detail: 'Please select media to add',
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

        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: `Added ${selectedMediaToAdd.value.length} media to playlist`,
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
            detail: error.response?.data?.error || 'Failed to add media',
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
            detail: 'Please select media to remove',
            life: 3000
        });
        return;
    }

    try {
        loadingMedia.value = true;

        // Use batch API for multiple media
        const mediaIds = selectedMediaToRemove.value.map(media => media.id);
        await playlistsApi.batchRemoveMediaFromPlaylist(selectedPlaylist.value.id, mediaIds);

        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: `Removed ${selectedMediaToRemove.value.length} media from playlist`,
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
            detail: error.response?.data?.error || 'Failed to remove media',
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

const getPosterUrl = (posterPath) => {
    return uploadApi.getPosterUrl(posterPath);
};

const getTypeLabel = (type) => {
    const option = typeOptions.find(opt => opt.value === type);
    return option ? option.label : type;
};

const getTypeSeverity = (type) => {
    switch (type) {
        case PlaylistType.SERIES:
            return 'success';
        case PlaylistType.FRANCHISE:
            return 'info';
        default:
            return 'secondary';
    }
};

const openPosterDialog = (posterPath) => {
    if (posterPath) {
        selectedPosterUrl.value = getPosterUrl(posterPath);
        displayPosterDialog.value = true;
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
            <h5 class="mb-0">Playlists Management</h5>
            <Button label="Add Playlist" icon="pi pi-plus" @click="openCreateDialog" />
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
                        <div
                            v-else
                            class="w-full h-48 bg-surface-100 dark:bg-surface-700 rounded-border flex items-center justify-center"
                        >
                            <i class="pi pi-list text-6xl text-surface-400"></i>
                        </div>
                        <Tag
                            :value="getTypeLabel(playlist.type)"
                            :severity="getTypeSeverity(playlist.type)"
                            class="absolute top-2 right-2"
                        />
                    </div>
                    <div class="mb-3">
                        <h6 class="mb-1">{{ playlist.name }}</h6>
                        <p class="text-xs text-muted-color">Created: {{ formatDate(playlist.created_at) }}</p>
                    </div>
                    <div class="flex gap-2">
                        <Button
                            icon="pi pi-pencil"
                            outlined
                            rounded
                            size="small"
                            @click="openEditDialog(playlist)"
                        />
                        <Button
                            icon="pi pi-video"
                            outlined
                            rounded
                            severity="info"
                            size="small"
                            @click="openMediaDialog(playlist)"
                            title="Manage Media"
                        />
                        <Button
                            icon="pi pi-trash"
                            outlined
                            rounded
                            severity="danger"
                            size="small"
                            @click="confirmDelete(playlist)"
                        />
                    </div>
                </div>
            </div>
        </div>

        <!-- Empty State -->
        <div
            v-if="!loading && playlists.length === 0"
            class="flex flex-col items-center justify-center py-12"
        >
            <i class="pi pi-list text-6xl text-surface-400 mb-4"></i>
            <p class="text-xl text-muted-color mb-4">No playlists found</p>
            <Button label="Create Your First Playlist" icon="pi pi-plus" @click="openCreateDialog" />
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
            currentPageReportTemplate="Showing {first} to {last} of {totalRecords} playlists"
        ></Paginator>

        <!-- Create/Edit Dialog -->
        <Dialog
            v-model:visible="displayDialog"
            :header="isEditMode ? 'Edit Playlist' : 'Create Playlist'"
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
                    <label for="type">Type *</label>
                    <Select
                        id="type"
                        v-model="formData.type"
                        :options="typeOptions"
                        optionLabel="label"
                        optionValue="value"
                        placeholder="Select a type"
                        :class="{ 'p-invalid': !formData.type }"
                    />
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
                <Button label="Save" icon="pi pi-check" @click="savePlaylist" :loading="loading || uploadingPoster" />
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
                <span v-if="selectedPlaylist">
                    Are you sure you want to delete playlist <b>{{ selectedPlaylist.name }}</b>?
                </span>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="displayDeleteDialog = false" />
                <Button label="Delete" icon="pi pi-trash" severity="danger" @click="deletePlaylist" :loading="loading" />
            </template>
        </Dialog>

        <!-- Media Management Dialog -->
        <Dialog
            v-model:visible="displayMediaDialog"
            header="Manage Playlist Media"
            :modal="true"
            :closable="true"
            class="p-fluid"
            style="width: 900px; max-width: 95vw"
        >
            <div v-if="selectedPlaylist" class="mb-4">
                <h6>{{ selectedPlaylist.name }}</h6>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <!-- Available Media -->
                <div class="col-span-2 md:col-span-1">
                    <div class="flex justify-between items-center mb-3">
                        <h6 class="mb-0">Available Media ({{ availableMediaTotal }})</h6>
                        <Button
                            label="Add Selected"
                            icon="pi pi-plus"
                            size="small"
                            @click="addMediaToPlaylist"
                            :disabled="loadingMedia || selectedMediaToAdd.length === 0"
                        />
                    </div>

                    <!-- Search -->
                    <div class="mb-3">
                        <IconField>
                            <InputIcon class="pi pi-search" />
                            <InputText
                                v-model="availableMediaSearch"
                                placeholder="Search media..."
                                @keyup.enter="onAvailableMediaSearch"
                                :disabled="loadingMedia"
                                class="w-full"
                            />
                        </IconField>
                    </div>

                    <div v-if="loadingMedia" class="flex justify-center py-8">
                        <ProgressSpinner style="width: 50px; height: 50px" />
                    </div>
                    <DataTable
                        v-else
                        v-model:selection="selectedMediaToAdd"
                        :value="availableMedia"
                        scrollable
                        scrollHeight="350px"
                        :rowHover="true"
                        class="text-sm"
                    >
                        <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>
                        <Column field="name" header="Name"></Column>
                        <Column field="format" header="Format">
                            <template #body="{ data }">
                                <Tag :value="data.format.toUpperCase()" severity="info" />
                            </template>
                        </Column>
                    </DataTable>
                    <p v-if="!loadingMedia && availableMedia.length === 0" class="text-center text-muted-color py-4">
                        No available media found
                    </p>

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
                        <h6 class="mb-0">Media in Playlist ({{ playlistMedia.length }})</h6>
                        <Button
                            label="Remove Selected"
                            icon="pi pi-minus"
                            size="small"
                            severity="danger"
                            @click="removeMediaFromPlaylist"
                            :disabled="loadingMedia || selectedMediaToRemove.length === 0"
                        />
                    </div>
                    <div v-if="loadingMedia" class="flex justify-center py-8">
                        <ProgressSpinner style="width: 50px; height: 50px" />
                    </div>
                    <DataTable
                        v-else
                        v-model:selection="selectedMediaToRemove"
                        :value="playlistMedia"
                        :reorderableRows="true"
                        @row-reorder="onMediaReorder"
                        scrollable
                        scrollHeight="400px"
                        :rowHover="true"
                        class="text-sm"
                    >
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
                    <p v-if="!loadingMedia && playlistMedia.length === 0" class="text-center text-muted-color py-4">
                        No media in this playlist
                    </p>
                </div>
            </div>

            <template #footer>
                <Button label="Close" icon="pi pi-times" @click="displayMediaDialog = false" />
            </template>
        </Dialog>

        <!-- Poster Preview Dialog -->
        <Dialog
            v-model:visible="displayPosterDialog"
            :modal="true"
            :closable="true"
            :dismissableMask="true"
            header="Poster Preview"
            class="poster-preview-dialog"
            style="width: auto; max-width: 90vw"
        >
            <div class="flex justify-center items-center">
                <img
                    v-if="selectedPosterUrl"
                    :src="selectedPosterUrl"
                    alt="Poster"
                    class="max-w-full max-h-[80vh] object-contain rounded-border"
                />
            </div>
        </Dialog>
    </div>
</template>
