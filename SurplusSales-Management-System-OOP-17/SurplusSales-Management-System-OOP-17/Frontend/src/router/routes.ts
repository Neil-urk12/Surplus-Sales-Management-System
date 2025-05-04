import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: () => import('pages/LoginPage.vue'),
  },
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/DashboardPage.vue') },
      { path: 'inventory/cabs', component: () => import('pages/CabsPage.vue') },
      { path: 'inventory/materials', component: () => import('pages/MaterialsPage.vue') },
      { path: 'inventory/accessories', component: () => import('pages/AccessoriesPage.vue') },
      { path: 'sales', component: () => import('pages/SalesPage.vue') },
      { path: 'contacts', component: () => import('pages/ContactsPage.vue') },
    ],
    meta: { requiresAuth: true },
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
]

export default routes;
