import { defineRouter } from '#q-app/wrappers';
import {
  createMemoryHistory,
  createRouter,
  createWebHashHistory,
  createWebHistory,
} from 'vue-router';
import routes from './routes';

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

    // Check if user is authenticated
    const isAuthenticated = localStorage.getItem('authToken');

    if (requiresAuth && !isAuthenticated) {
      // If route requires auth and user is not authenticated, redirect to login
      next('/login');
    } else if (to.path === '/login' && isAuthenticated) {
      // If user is already authenticated and tries to access login page, redirect to home
      next('/');
    } else {
      // Otherwise proceed as normal
      next();
    }
  });

  return Router;
});
