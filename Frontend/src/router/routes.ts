import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [ { path: '', component: () => import('pages/DashboardPage.vue')},
                { path: 'inventory', component: () => import('pages/InventoryPage.vue') },
                { path: 'sales', component: () => import('pages/SalesPage.vue') },
                { path: 'contacts', component: () => import('pages/ContactsPage.vue') }
    ],
  },
  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
