<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useQuasar } from 'quasar'
import type { QTableProps } from 'quasar'
import { useUsersStore } from '../stores/users'
import { useAuthStore } from '../stores/auth'
import type { User } from '../types/models'
import type { UserCreateData, UserUpdateData } from '../stores/users'

const $q = useQuasar()
const usersStore = useUsersStore()
const authStore = useAuthStore()

// Table data
const loading = ref(false)
const filter = ref('')
const pagination = ref({
  rowsPerPage: 10
})

// Form data
const createDialog = ref(false)
const editDialog = ref(false)
const selectedUser = ref<User | null>(null)

// Create user form
const newUser = ref<UserCreateData>({
  fullName: '',
  email: '',
  password: '',
  role: 'staff'
})

// Edit user form
const editedUser = ref<UserUpdateData>({
  role: 'staff',
  isActive: true
})

// Role options
const roleOptions = [
  { label: 'Admin', value: 'admin' },
  { label: 'Staff', value: 'staff' }
]

// Table columns
const columns: QTableProps['columns'] = [
  { name: 'fullName', align: 'left', label: 'Full Name', field: 'fullName', sortable: true },
  { name: 'email', align: 'left', label: 'Email', field: 'email', sortable: true },
  { name: 'role', align: 'left', label: 'Role', field: 'role', sortable: true },
  { name: 'isActive', align: 'left', label: 'Status', field: 'isActive', sortable: true },
  { name: 'actions', align: 'center', label: 'Actions', field: 'actions' }
]

// Current user role check
const isAdminOrStaff = computed(() => {
  const role = authStore.user?.role
  return role === 'admin' || role === 'staff'
})

// Highlight current user row
const rowClass = (row: User) => row.id === authStore.user?.id ? 'bg-grey-2' : ''

// Fetch users on component mount
onMounted(async () => {
  loading.value = true
  try {
    await usersStore.fetchUsers()
  } catch (error) {
    console.error('Failed to fetch users:', error)
    $q.notify({
      color: 'negative',
      message: 'Failed to load users',
      icon: 'error'
    })
  } finally {
    loading.value = false
  }
})

// Open create user dialog
function openCreateDialog() {
  // Reset form
  newUser.value = {
    fullName: '',
    email: '',
    password: '',
    role: 'staff'
  }
  createDialog.value = true
}

// Open edit user dialog
function openEditDialog(user: User) {
  selectedUser.value = user
  editedUser.value = {
    fullName: user.fullName,
    email: user.email,
    role: user.role,
    isActive: user.isActive
  }
  editDialog.value = true
}

// Submit create user form
async function submitCreateUser() {
  try {
    const created = await usersStore.createUser(newUser.value)
    if (!created) {
      throw new Error('User creation failed')
    }
    $q.notify({
      color: 'positive',
      message: 'User created successfully',
      icon: 'check'
    })
    createDialog.value = false
  } catch (error) {
    console.error('Failed to create user:', error)
    $q.notify({
      color: 'negative',
      message: 'Failed to create user',
      icon: 'error'
    })
  }
}

// Submit edit user form
async function submitEditUser() {
  if (!selectedUser.value) return

  try {
    const success = await usersStore.updateUser(selectedUser.value.id, {
      fullName: editedUser.value.fullName,
      email: editedUser.value.email,
      role: editedUser.value.role,
      isActive: editedUser.value.isActive
    } as UserUpdateData) // Explicitly cast to UserUpdateData
    if (!success) {
      throw new Error('User update failed')
    }
    $q.notify({
      color: 'positive',
      message: 'User updated successfully',
      icon: 'check'
    })
    editDialog.value = false
  } catch (error) {
    console.error('Failed to update user:', error)
    $q.notify({
      color: 'negative',
      message: 'Failed to update user',
      icon: 'error'
    })
  }
}

// Format status for display
function formatStatus(status: boolean) {
  return status ? 'Active' : 'Inactive'
}
</script>

<template>
  <q-page class="q-pa-md flex">
    <div class="q-pa-sm full-width">
      <!-- Page Header -->
      <div class="row items-center justify-between q-mb-md">
        <h1 class="text-h4 text-soft-light q-my-none">User Management</h1>
        <q-btn
          v-if="isAdminOrStaff"
          color="primary"
          icon="add"
          label="Create User"
          class="text-soft-light"
          @click="openCreateDialog"
        />
      </div>

      <!-- Users Table -->
      <q-card class="user-table-card">
        <q-card-section>
          <q-table
            :rows="usersStore.users"
            :columns="columns"
            :loading="loading || usersStore.loading"
            :pagination="pagination"
            row-key="id"
            :filter="filter"
            :row-class="rowClass"
            flat
            bordered
          >
            <!-- Table Top -->
            <template v-slot:top>
              <div class="row full-width items-center q-pb-sm">
                <q-input
                  dense
                  outlined
                  v-model="filter"
                  placeholder="Search"
                  class="col-grow"
                >
                  <template v-slot:append>
                    <q-icon name="search" />
                  </template>
                </q-input>
              </div>
            </template>

            <!-- Role Column -->
            <template v-slot:body-cell-role="props">
              <q-td :props="props">
                <q-badge :color="props.value === 'admin' ? 'accent' : props.value === 'staff' ? 'secondary' : 'grey-7'">
                  {{ props.value.toUpperCase() }}
                </q-badge>
              </q-td>
            </template>

            <!-- Status Column -->
            <template v-slot:body-cell-isActive="props">
              <q-td :props="props">
                <q-badge :color="props.value ? 'positive' : 'negative'">
                  {{ formatStatus(props.value) }}
                </q-badge>
              </q-td>
            </template>

            <!-- Actions Column -->
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" class="text-center">
                <q-btn
                  v-if="isAdminOrStaff"
                  flat
                  round
                  color="primary"
                  icon="edit"
                  @click="openEditDialog(props.row)"
                  size="sm"
                >
                  <q-tooltip>Edit User</q-tooltip>
                </q-btn>
              </q-td>
            </template>

            <!-- No Data Message -->
            <template v-slot:no-data>
              <div class="full-width row justify-center q-pa-md" style="color: var(--muted-foreground)">
                No users found
              </div>
            </template>
          </q-table>
        </q-card-section>
      </q-card>

      <!-- Create User Dialog -->
      <q-dialog v-model="createDialog" persistent>
        <q-card style="min-width: 350px">
          <q-card-section class="row items-center">
            <div class="text-h6">Create New User</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section>
            <q-form @submit="submitCreateUser">
              <q-input
                v-model="newUser.fullName"
                label="Full Name *"
                filled
                :rules="[val => !!val || 'Full name is required']"
                class="q-mb-md"
              />

              <q-input
                v-model="newUser.email"
                label="Email *"
                filled
                type="email"
                :rules="[
                  val => !!val || 'Email is required',
                  val => /^\S+@\S+\.\S+$/.test(val) || 'Please enter a valid email'
                ]"
                class="q-mb-md"
              />

              <q-input
                v-model="newUser.password"
                label="Password *"
                filled
                type="password"
                :rules="[val => !!val || 'Password is required', val => val.length >= 6 || 'Password must be at least 6 characters']"
                class="q-mb-md"
              />

              <q-select
                v-model="newUser.role"
                :options="roleOptions"
                option-value="value"
                option-label="label"
                label="Role *"
                filled
                emit-value
                map-options
                class="q-mb-md"
              />

              <div class="row justify-end q-mt-md">
                <q-btn label="Cancel" flat v-close-popup class="q-mr-sm" />
                <q-btn label="Create" type="submit" color="primary" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-dialog>

      <!-- Edit User Dialog -->
      <q-dialog v-model="editDialog" persistent>
        <q-card style="min-width: 350px">
          <q-card-section class="row items-center">
            <div class="text-h6">Edit User</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>

          <q-card-section>
            <q-form @submit="submitEditUser">
              <q-input
                v-model="editedUser.fullName"
                label="Full Name"
                filled
                class="q-mb-md"
              />

              <q-input
                v-model="editedUser.email"
                label="Email"
                filled
                type="email"
                class="q-mb-md"
              />

              <q-select
                v-model="editedUser.role"
                :options="roleOptions"
                option-value="value"
                option-label="label"
                label="Role *"
                filled
                emit-value
                map-options
                class="q-mb-md"
              />

              <q-toggle
                v-model="editedUser.isActive"
                label="Active Status"
                color="primary"
                class="q-mb-md"
              />

              <div class="row justify-end q-mt-md">
                <q-btn label="Cancel" flat v-close-popup class="q-mr-sm" />
                <q-btn label="Save" type="submit" color="primary" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-dialog>
    </div>
  </q-page>
</template>

<style scoped>
.user-table-card {
  border-radius: 8px;
  /* box-shadow removed for theme consistency */
}

.highlighted-row {
  background-color: var(--hover-bg);
}
</style>
