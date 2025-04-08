<template>
  <q-layout view="lHh Lpr lFf">
    <q-header :elevated="!isDark"  class="q-ma-md q-mx-lg custom-nav-height-width custom-rounded flex " style="height: 62px; background-color: var(--header-nav-bg);">
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
        <q-toolbar-title>
          <q-input
            dark
            borderless
            v-model="text"
            input-class="text-left custom-style"
            class="q-ml-md"
            placeholder="Search"
          >
            <template v-slot:prepend>
              <q-icon v-if="text === ''" name="search" class="text-soft-light" />
              <q-icon v-else name="clear" class="cursor-pointer text-soft-light"   @click="clearSearch"/>
            </template>
          </q-input>
        </q-toolbar-title>
        <div class="flex row q-gutter-x-xs justify-center">
          <q-toolbar-title class="flex items-center">
            <q-btn dense round flat class="text-soft-light" @click=toggleColorMode>
              <q-icon v-if="isDark != true" name="light_mode" class="text-soft-light"/>
              <q-icon v-else name="dark_mode" class="text-soft-light"/>
            </q-btn>
            <q-btn dense round flat icon="notifications" class="text-soft-light">
              <q-badge color="red" floating>
                4
              </q-badge>
            </q-btn>
          </q-toolbar-title>
          <q-avatar color="red user-profile-width user-profile-height" text-color="text-primary">
            <img src="https://cdn.quasar.dev/img/avatar.png">
          </q-avatar>
        </div>

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
        <MenuItems
        class="text-soft-light"
        v-for="link in menuItemsList"
        :key="link.title"
        v-bind="link"
        />

      </q-list>
    </q-drawer>

    <q-page-container style="padding-top: 62px;">

      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import MenuItems from 'components/SideMenuItems.vue'
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
    to: '/inventory'
  },
  {
    title: 'Sales',
    icon: 'trending_up',
    to: '/sales'
  },
  {
    title:"Contacts",
    icon:"contacts",
    to: '/contacts'

  },
]

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

const clearSearch = () => {
  text.value = ''
  if (timeoutId.value) {
    clearTimeout(timeoutId.value)
    timeoutId.value = null
  }
};

onMounted(() => {
  const savedMode = localStorage.getItem('quasar-theme')
  if (savedMode) {
    isDark.value = savedMode === 'dark'
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
}



function toggleLeftDrawer () {
  leftDrawerOpen.value = !leftDrawerOpen.value
}
</script>


