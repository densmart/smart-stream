<script setup>
import { ref, onMounted } from 'vue';
import { clientsApi } from '@/api';
import { useToast } from 'primevue/usetoast';

const toast = useToast();

// State
const clients = ref([]);
const loading = ref(false);
const totalRecords = ref(0);
const displayDialog = ref(false);
const displayDeleteDialog = ref(false);
const isEditMode = ref(false);
const selectedClient = ref(null);

// Form data
const formData = ref({
    login: '',
    password: '',
    email: '',
    is_active: true
});

// Email validation
const isValidEmail = (email) => {
    if (!email) return true; // Email is optional
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
};

// Pagination
const lazyParams = ref({
    page: 1,
    limit: 10,
    search: ''
});

// Methods
const loadClients = async () => {
    try {
        loading.value = true;
        const response = await clientsApi.getClients({
            offset: (lazyParams.value.page - 1) * lazyParams.value.limit,
            limit: lazyParams.value.limit,
            search: lazyParams.value.search || undefined
        });
        clients.value = response.result;
        totalRecords.value = response.pagination.total;
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to load clients',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const onPage = (event) => {
    lazyParams.value.page = event.page + 1;
    lazyParams.value.limit = event.rows;
    loadClients();
};

const openCreateDialog = () => {
    isEditMode.value = false;
    formData.value = {
        login: '',
        password: '',
        email: '',
        is_active: true
    };
    displayDialog.value = true;
};

const openEditDialog = (client) => {
    isEditMode.value = true;
    selectedClient.value = client;
    formData.value = {
        login: client.login,
        password: '', // Don't pre-fill password for security
        email: client.email || '',
        is_active: client.is_active
    };
    displayDialog.value = true;
};

const hideDialog = () => {
    displayDialog.value = false;
    selectedClient.value = null;
};

const saveClient = async () => {
    if (!formData.value.login || (!isEditMode.value && !formData.value.password)) {
        toast.add({
            severity: 'warn',
            summary: 'Validation Error',
            detail: 'Login and password are required',
            life: 3000
        });
        return;
    }

    if (formData.value.email && !isValidEmail(formData.value.email)) {
        toast.add({
            severity: 'warn',
            summary: 'Validation Error',
            detail: 'Please enter a valid email address',
            life: 3000
        });
        return;
    }

    try {
        loading.value = true;

        if (isEditMode.value) {
            const updateData = {
                login: formData.value.login,
                is_active: formData.value.is_active
            };

            // Only include email if it was provided
            if (formData.value.email) {
                updateData.email = formData.value.email;
            }

            // Only include password if it was changed
            if (formData.value.password) {
                updateData.password = formData.value.password;
            }

            await clientsApi.updateClient(selectedClient.value.id, updateData);
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Client updated successfully',
                life: 3000
            });
        } else {
            await clientsApi.createClient(formData.value);
            toast.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Client created successfully',
                life: 3000
            });
        }

        hideDialog();
        loadClients();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to save client',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const confirmDelete = (client) => {
    selectedClient.value = client;
    displayDeleteDialog.value = true;
};

const deleteClient = async () => {
    try {
        loading.value = true;
        await clientsApi.deleteClient(selectedClient.value.id);
        toast.add({
            severity: 'success',
            summary: 'Success',
            detail: 'Client deleted successfully',
            life: 3000
        });
        displayDeleteDialog.value = false;
        selectedClient.value = null;
        loadClients();
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: 'Error',
            detail: error.response?.data?.error || 'Failed to delete client',
            life: 3000
        });
    } finally {
        loading.value = false;
    }
};

const formatDate = (dateString) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

onMounted(() => {
    loadClients();
});
</script>

<template>
    <div class="card">
        <Toast />

        <div class="flex justify-between items-center mb-6">
            <h5 class="mb-0">Clients Management</h5>
            <Button label="Add Client" icon="pi pi-plus" @click="openCreateDialog" />
        </div>

        <DataTable
            :value="clients"
            :lazy="true"
            :paginator="true"
            :rows="lazyParams.limit"
            :totalRecords="totalRecords"
            :loading="loading"
            @page="onPage"
            paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
            :rowsPerPageOptions="[5, 10, 25, 50]"
            currentPageReportTemplate="Showing {first} to {last} of {totalRecords} clients"
            responsiveLayout="scroll"
            class="p-datatable-sm"
        >
            <Column field="id" header="ID" :sortable="false" style="min-width: 5rem"></Column>
            <Column field="login" header="Login" :sortable="false" style="min-width: 12rem"></Column>
            <Column field="is_active" header="Status" :sortable="false" style="min-width: 8rem">
                <template #body="{ data }">
                    <Tag :value="data.is_active ? 'Active' : 'Inactive'" :severity="data.is_active ? 'success' : 'danger'" />
                </template>
            </Column>
            <Column field="created_at" header="Created At" :sortable="false" style="min-width: 12rem">
                <template #body="{ data }">
                    {{ formatDate(data.created_at) }}
                </template>
            </Column>
            <Column field="last_login_at" header="Last Login" :sortable="false" style="min-width: 12rem">
                <template #body="{ data }">
                    {{ formatDate(data.last_login_at) }}
                </template>
            </Column>
            <Column header="Actions" :exportable="false" style="min-width: 10rem">
                <template #body="{ data }">
                    <Button icon="pi pi-pencil" outlined rounded class="mr-2" @click="openEditDialog(data)" />
                    <Button icon="pi pi-trash" outlined rounded severity="danger" @click="confirmDelete(data)" />
                </template>
            </Column>
        </DataTable>

        <!-- Create/Edit Dialog -->
        <Dialog v-model:visible="displayDialog" :header="isEditMode ? 'Edit Client' : 'Create Client'" :modal="true" :closable="true" class="p-fluid" style="width: 450px">
            <div class="flex flex-col gap-6 py-4">
                <div class="flex flex-col gap-2">
                    <label for="login">Login</label>
                    <InputText id="login" v-model="formData.login" required="true" autofocus :class="{ 'p-invalid': !formData.login }" />
                </div>

                <div class="flex flex-col gap-2">
                    <label for="email">Email (optional)</label>
                    <InputText id="email" v-model="formData.email" type="email" :class="{ 'p-invalid': formData.email && !isValidEmail(formData.email) }" />
                    <small v-if="formData.email && !isValidEmail(formData.email)" class="p-error">Please enter a valid email address</small>
                </div>

                <div class="flex flex-col gap-2">
                    <label for="password">Password {{ isEditMode ? '(leave empty to keep current)' : '' }}</label>
                    <Password id="password" v-model="formData.password" toggleMask :required="!isEditMode" :class="{ 'p-invalid': !isEditMode && !formData.password }" :feedback="false" />
                </div>

                <div class="flex items-center gap-2">
                    <Checkbox id="is_active" v-model="formData.is_active" :binary="true" />
                    <label for="is_active">Active</label>
                </div>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="hideDialog" />
                <Button label="Save" icon="pi pi-check" @click="saveClient" :loading="loading" />
            </template>
        </Dialog>

        <!-- Delete Confirmation Dialog -->
        <Dialog v-model:visible="displayDeleteDialog" header="Confirm Delete" :modal="true" :closable="true" style="width: 450px">
            <div class="flex items-center gap-4">
                <i class="pi pi-exclamation-triangle !text-3xl text-orange-500" />
                <span v-if="selectedClient">
                    Are you sure you want to delete client <b>{{ selectedClient.login }}</b
                    >?
                </span>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="displayDeleteDialog = false" />
                <Button label="Delete" icon="pi pi-trash" severity="danger" @click="deleteClient" :loading="loading" />
            </template>
        </Dialog>
    </div>
</template>
