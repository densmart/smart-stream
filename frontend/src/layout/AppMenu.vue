<script setup>
import { ref } from 'vue';
import AppMenuItem from './AppMenuItem.vue';

const model = ref([
    {
        label: 'Home',
        items: [
            {
                label: 'Dashboard',
                icon: 'pi pi-fw pi-home',
                to: '/'
            }
        ]
    },
    {
        label: 'Management',
        items: [
            {
                label: 'Clients',
                icon: 'pi pi-fw pi-users',
                to: '/clients'
            },
            {
                label: 'Movies',
                icon: 'pi pi-fw pi-video',
                to: '/media'
            },
            {
                label: 'Series',
                icon: 'pi pi-fw pi-list',
                to: '/series'
            },
            {
                label: 'Franchises',
                icon: 'pi pi-fw pi-folder',
                to: '/franchises'
            },
            {
                label: 'Cartoons',
                icon: 'pi pi-fw pi-star',
                to: '/cartoons'
            }
        ]
    }
]);
</script>

<template>
    <ul class="layout-menu">
        <template v-for="(item, i) in model" :key="item">
            <app-menu-item v-if="!item.separator" :item="item" :index="i"></app-menu-item>
            <li v-if="item.separator" class="menu-separator"></li>
        </template>
    </ul>
</template>

<style lang="scss" scoped></style>
