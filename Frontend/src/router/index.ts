import { defineRouter } from '#q-app/wrappers';
import {
  createMemoryHistory,
  createRouter,
  createWebHashHistory,
  createWebHistory,
} from 'vue-router';
import routes from './routes';
import { useAuthStore } from 'stores/auth';
import { showErrorNotification } from 'src/utils/notifications';

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

export default defineRouter(function (/* { store, ssrContext } */) {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : (process.env.VUE_ROUTER_MODE === 'history' ? createWebHistory : createWebHashHistory);

  const Router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,

    // Leave this as is and make changes in quasar.conf.js instead!
    // quasar.conf.js -> build -> vueRouterMode
    // quasar.conf.js -> build -> publicPath
    history: createHistory(process.env.VUE_ROUTER_BASE),
  });

  // Navigation guards
  Router.beforeEach((to, _from, next) => {
    // Check if the route requires authentication
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
    // Find the deepest matched route that defines roles
    const routeRoles = to.matched.slice().reverse().find(record => record.meta.roles)?.meta.roles as string[] | undefined;

    // Need to get store instance inside the guard
    const authStore = useAuthStore();
    const isAuthenticated = authStore.isAuthenticated; // Use store getter
    const userRole = authStore.user?.role;

    // Explicitly handle case where no routes are matched
    if (!to.matched || to.matched.length === 0) {
      next();
      return;
    }

    if (requiresAuth && !isAuthenticated) {
      // Needs auth, but user is not logged in -> redirect to login
      next('/login');
    } else if (to.path === '/login' && isAuthenticated) {
      // Already logged in, trying to access login -> redirect to dashboard
      next('/');
    } else if (routeRoles && isAuthenticated) {
      // Explicitly check for undefined userRole when route requires roles
      if (!userRole) {
        console.warn(`User role is undefined. Redirecting unauthorized.`);
        showErrorNotification({
          message: `Your user role is not defined. Access denied to ${to.path}.`,
        });
        next('/unauthorized');
        return;
      }
      // Route requires specific roles AND user is logged in
      if (userRole && routeRoles.includes(userRole)) {
        // User has the required role -> allow access
        next();
      } else {
        // User does not have the required role -> redirect to dashboard
        console.warn(`Unauthorized access attempt to ${to.path} by user with role ${userRole}`);
        showErrorNotification({
          message: `You do not have the necessary permissions to access ${to.path}.`,
        });
        next('/unauthorized'); // Redirect to unauthorized page as they don't have permission
      }
    } else {
      // Route either doesn't require auth, or doesn't require specific roles,
      // or user is not authenticated but route doesn't require auth.
      // -> allow access
      next();
    }
  });

  return Router;
});
