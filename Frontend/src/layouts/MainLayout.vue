<template>
  <q-layout view="lHh Lpr lFf">
    <q-header :elevated="!isDark"  class="q-ma-md q-mx-lg custom-nav-height-width custom-rounded flex" style="height: 62px; background-color: var(--header-nav-bg);">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
          class="text-soft-light"

        />
        <q-toolbar-title class="row flext items-center justify-between">
          <div>
            {{ activePage }}
          </div>
          <div class="flex row q-gutter-x-xs justify-center" v-if="route.path === '/inventory'">
            <q-btn dense round flat icon="download" ></q-btn>
            <q-btn class="text-white bg-primary">
              <q-icon color="white" name="add" />
              Add</q-btn>
          </div>

        </q-toolbar-title>
        <!-- DO NOT REMOVE YET -->
        <!-- <q-toolbar-title class="row flext items-center" v-if="route.path === '/sales'">
          <q-input
            dark
            borderless
            v-model="text"
            input-class="text-left custom-style"
            class="q-ml-md"
            placeholder="Search"
            debounce="500"
          >
            <template v-slot:prepend>
              <q-icon v-if="text === ''" name="search" class="text-soft-light" />
              <q-icon v-else name="clear" class="cursor-pointer text-soft-light"   @click="clearSearch"/>
            </template>
          </q-input>
          <q-btn class="col">
              <q-icon name="filter_alt" class="text-soft-light"/>
            <p>Filter</p>
          </q-btn>
        </q-toolbar-title> -->


      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
    >
      <q-list>
        <q-item class="q-px-none q-pt-md q-pb-sm q-hoverable ">
          <q-item-section class="flex flex-start">
            <div class="row items-center justify-start q-ml-xs q-gutter-x-md">
              <img
                src="../assets/logo.png"
                style="width: 64px; height: 64px;"
                class="q-mt-xs"
              >
              <div class="text-h6 q-pl-xs text-soft-light">Cortes Surplus</div>
            </div>
          </q-item-section>
        </q-item>
        <div class="flex-col ">
          <!-- menu -->
          <MenuItems
            class="text-soft-light"
            v-for="link in menuItemsList"
            :key="link.title"
            v-bind="link"
          />
            <q-item
              clickable
              v-bind="$attrs"
              class="text-soft-light q-hoverable absolute-bottom"
            >
              <!-- user -->
              <div class="flex flex-start row items-center justify-between">
                <div class="flex flex-start items-center q-gutter-md">
                  <q-avatar color="red user-profile-width user-profile-height" text-color="text-primary">
                      <img src="https://cdn.quasar.dev/img/avatar.png">
                  </q-avatar>
                  <q-item-section>
                    <q-item-label>{{ currentUser.name }}</q-item-label>
                  </q-item-section>
                </div>
              </div>
              <div class="q-pa-md col">
                <q-btn-dropdown
                flat
                class="absolute-right"
                >
                  <div class="row no-wrap q-pa-md">
                    <div class="column">
                      <div class="text-h6 q-mb-md">Preference</div>
                        <div class="flex-start items-center">
                          <q-btn dense flat class="text-soft-light row full-width flex justify-between" @click=toggleColorMode>
                            <div class="flex items-center col-auto q-pr-sm full-width">
                                <div class="flex q-pr-sm q-pl-sm ">
                                  <q-icon
                                    v-if="!isDark"
                                    name="light_mode"
                                    class="text-soft-light"
                                    size="24px"
                                  />
                                  <q-icon
                                    v-else
                                    name="dark_mode"
                                    class="text-soft-light"
                                    size="24px"
                                  />
                                </div>
                                <div class="flex items-center col">
                                  <p class="q-mb-none text-weight-medium">{{ themeMode }}</p>
                                </div>
                            </div>
                          </q-btn>
                          <q-btn dense flat class="text-soft-light row full-width" @click=alert>
                            <div class="flex flex-start items-center">
                                <div class="flex items-center col-auto q-pr-sm">
                                  <q-btn dense round flat icon="notifications" class="text-soft-light">
                                  <q-badge color="red" floating>
                                    4
                                  </q-badge>
                                </q-btn>
                                </div>
                                <div class="flex items-center col">
                                  <p class="q-mb-none text-weight-medium">Notifications</p>
                                </div>
                            </div>
                          </q-btn>
                        </div>
                    </div>

                    <q-separator vertical inset class="q-mx-lg" />

                    <div class="column items-center">
                      <q-avatar size="72px">
                        <img src="https://cdn.quasar.dev/img/boy-avatar.png">
                      </q-avatar>

                      <div class="text-subtitle1 q-mt-md q-mb-xs">{{ currentUser.name }}</div>

                      <q-btn
                        color="primary"
                        label="Logout"
                        push
                        flat
                        size="sm"
                        v-close-popup
                        @click="handleLogout"
                      />
                    </div>
                  </div>
                </q-btn-dropdown>
              </div>
            </q-item>
        </div>

      </q-list>
    </q-drawer>

    <q-page-container style="padding-top: 62px;">

      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed, defineAsyncComponent } from 'vue'
import { useRoute, useRouter } from 'vue-router';
import { useQuasar } from 'quasar'
import { useAuthStore } from '../stores/auth'
const MenuItems = defineAsyncComponent(() => import('../components/SideMenuItems.vue'))
import type { menuItemsProps } from '../types/menu-items'

const menuItemsList: menuItemsProps[] = [
  {
    title: 'Dashboard',
    icon: 'dashboard',
    to: '/',
    exact: true
  },
  {
    title: 'Inventory',
    icon: 'storage',
    hasSubmenu: true,
    children: [
      {
        title: 'Cabs',
        icon: 'directions_car',
        to: '/inventory/cabs'
      },
      {
        title: 'Materials',
        icon: 'category',
        to: '/inventory/materials'
      },
      {
        title: 'Accessories',
        icon: 'settings_input_component',
        to: '/inventory/accessories'
      }
    ]
  },
  {
    title: 'Sales',
    icon: 'trending_up',
    to: '/app/sales'
  },
  {
    title:"Contacts",
    icon:"contacts",
    to: '/app/contacts'
  },
]

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const activePage = ref('')

watch(() => route.path, (newPath) => {
  switch (newPath) {
    case '/':
      activePage.value = 'Dashboard'
      break
    case '/sales':
      activePage.value = 'Sales'
      break
    case '/inventory/cabs':
      activePage.value = 'Cabs'
      break
    case '/inventory/materials':
      activePage.value = 'Materials'
      break
    case '/inventory/accessories':
      activePage.value = 'Accessories'
      break
    case '/contacts':
      activePage.value = 'Contacts'
      break
    default:
      activePage.value = 'Dashboard'
  }
}, { immediate: true })

const leftDrawerOpen = ref(false)
const isDark = ref(false)
const $q = useQuasar()


const text = ref<string>('')
const timeoutId = ref<ReturnType<typeof setTimeout> | null>(null)

const performSearch = (searchTerm: string) => {
console.log('Searching for:', searchTerm)
//for later

};

watch(text, (newValue) => {
  if (timeoutId.value) clearTimeout(timeoutId.value);
  timeoutId.value = setTimeout(() => {
    performSearch(newValue)
  }, 300)
})

const themeMode = ref('')

onMounted(() => {
  const savedMode = localStorage.getItem('quasar-theme')
  if (savedMode) {
    isDark.value = savedMode === 'dark'
    themeMode.value = savedMode
  } else {
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  $q.dark.set(isDark.value)
//for the time out cleanup
})

const toggleColorMode = () => {
  isDark.value = !isDark.value
  $q.dark.set(isDark.value)
  localStorage.setItem('quasar-theme', isDark.value ? 'dark' : 'light')

  if(isDark.value){
    themeMode.value = 'Dark Mode'
  }else{
    themeMode.value = 'Light Mode'
  }
}

const currentUser = computed(() => {
  return authStore.user ? {
    name: authStore.user.fullName,
    email: authStore.user.email,
    avatar: 'https://cdn.quasar.dev/img/avatar.png',
  } : {
    name: 'Guest User',
    email: '',
    avatar: 'https://cdn.quasar.dev/img/avatar.png',
  }
})


function alert () {
  console.log('alert')
}

function toggleLeftDrawer () {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

async function handleLogout() {
  authStore.logout()
  await router.push('/login')
}
</script>
