import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: () => import('pages/LoginPage.vue'),
  },
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    meta: { requiresAuth: false },
    children: [
      { path: '', component: () => import('pages/DashboardPage.vue') },
      { path: 'cabs', component: () => import('pages/CabsPage.vue') },
      { path: 'materials', component: () => import('pages/MaterialsPage.vue') },
      // { path: 'accessories', component: () => import('pages/AccessoriesPage.vue') },
      { path: 'sales', component: () => import('pages/SalesPage.vue') },
      { path: 'contacts', component: () => import('pages/ContactsPage.vue') },
    ],
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
]

export default routes;
