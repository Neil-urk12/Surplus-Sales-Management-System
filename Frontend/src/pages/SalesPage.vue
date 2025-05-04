<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuasar } from 'quasar';
import type { MultiCab, Sale, Customer, User } from '../types/models';
import SalesTrendChart from '../components/charts/SalesTrendChart.vue';
import SalesByModelChart from '../components/charts/SalesByModelChart.vue';

const $q = useQuasar();

// Generate historical sales data for the chart
function generateHistoricalSales(existingCarIds: string[]) {
  const months = [
    'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec', 'Jan', 'Feb', 'Mar', 'Apr'
  ];

  const historicalSales: Sale[] = [];

  // Generate monthly sales data to match the chart in the mockup
  const monthlyRevenues = [700000, 400000, 600000, 500000, 400000, 400000, 400000, 600000, 300000, 500000, 300000];

  let saleId = 100;

  months.forEach((month, index) => {
    const year = index < 7 ? 2024 : 2025;
    const monthIndex = index < 7 ? index + 5 : index - 7;

    // Generate multiple sales for each month to reach the target revenue
    const targetRevenue = monthlyRevenues[index] || 0; // Default to 0 if undefined
    let currentRevenue = 0;

    while (currentRevenue < targetRevenue) {
      const price = 15000 + Math.floor(Math.random() * 55000);
      currentRevenue += price;

      const day = Math.floor(Math.random() * 28) + 1;
      const date = new Date(year, monthIndex, day);
      const randomCarId = existingCarIds[Math.floor(Math.random() * existingCarIds.length)] || 'unknown'; // Assign a random car ID

      historicalSales.push({
        id: `hist-${saleId++}`,
        customerId: String(Math.floor(Math.random() * 3) + 1),
        soldBy: String(Math.floor(Math.random() * 4) + 1),
        saleDate: date.toISOString(),
        totalPrice: price,
        createdAt: date.toISOString(),
        updatedAt: date.toISOString(),
        multiCabId: randomCarId // Add the multiCabId
      });
    }
  });

  return historicalSales;
}

// Generate additional cars for the statistics
function generateAdditionalCars() {
  const models = ['Coupe B', 'Truck Z', 'Hatchback J', 'Sedan X', 'SUV Y', 'Minivan C'];
  const colors = ['Red', 'Blue', 'Black', 'White', 'Silver', 'Gray'];
  const makes = ['BMW', 'Toyota', 'Honda', 'Nissan', 'Suzuki', 'Ford'];

  const additionalCars: MultiCab[] = [];

  // Generate more cars to match the statistics in the mockup
  for (let i = 0; i < 193; i++) { // To reach a total of 200 cars
    const modelIndex = i % models.length;
    const model = models[modelIndex]!;
    const make = makes[modelIndex]!;

    additionalCars.push({
      id: `additional-${i + 1}`,
      make,
      model,
      year: 2023 + Math.floor(Math.random() * 2),
      color: colors[Math.floor(Math.random() * colors.length)]!,
      condition: Math.random() > 0.3 ? 'New' : 'Used',
      price: 15000 + Math.floor(Math.random() * 55000),
      status: 'Sold',
      dateAdded: new Date(2024, Math.floor(Math.random() * 12), Math.floor(Math.random() * 28) + 1).toISOString(),
      serialNumber: `${model.substring(0, 2)}-${2023 + Math.floor(Math.random() * 2)}-${1000 + i}`,
      createdAt: new Date(2024, Math.floor(Math.random() * 12), Math.floor(Math.random() * 28) + 1).toISOString(),
      updatedAt: new Date(2025, Math.floor(Math.random() * 4), Math.floor(Math.random() * 28) + 1).toISOString()
    });
  }

  return additionalCars;
}

// Mock data for multicabs
const mockMultiCabs = ref<MultiCab[]>([
  {
    id: '1',
    make: 'Suzuki',
    model: 'Sedan X',
    year: 2024,
    color: 'Silver',
    condition: 'New',
    price: 20541,
    status: 'Sold',
    dateAdded: '2025-03-15T00:00:00Z',
    serialNumber: 'SX-2024-001',
    createdAt: '2025-03-15T00:00:00Z',
    updatedAt: '2025-04-28T10:30:00Z'
  },
  {
    id: '2',
    make: 'Toyota',
    model: 'Truck Z',
    year: 2024,
    color: 'White',
    condition: 'New',
    price: 47432,
    status: 'Sold',
    dateAdded: '2025-03-10T00:00:00Z',
    serialNumber: 'TZ-2024-001',
    createdAt: '2025-03-10T00:00:00Z',
    updatedAt: '2025-04-26T14:15:00Z'
  },
  {
    id: '3',
    make: 'Honda',
    model: 'Minivan C',
    year: 2024,
    color: 'Black',
    condition: 'New',
    price: 68118,
    status: 'Sold',
    dateAdded: '2025-03-05T00:00:00Z',
    serialNumber: 'MC-2024-001',
    createdAt: '2025-03-05T00:00:00Z',
    updatedAt: '2025-04-24T09:45:00Z'
  },
  {
    id: '4',
    make: 'Nissan',
    model: 'SUV Y',
    year: 2024,
    color: 'Gray',
    condition: 'New',
    price: 30526,
    status: 'Sold',
    dateAdded: '2025-03-01T00:00:00Z',
    serialNumber: 'SY-2024-001',
    createdAt: '2025-03-01T00:00:00Z',
    updatedAt: '2025-04-23T16:20:00Z'
  },
  {
    id: '5',
    make: 'BMW',
    model: 'Coupe B',
    year: 2024,
    color: 'Gray',
    condition: 'New',
    price: 47182,
    status: 'Sold',
    dateAdded: '2025-02-25T00:00:00Z',
    serialNumber: 'CB-2024-001',
    createdAt: '2025-02-25T00:00:00Z',
    updatedAt: '2025-04-23T11:10:00Z'
  },
  {
    id: '6',
    make: 'Suzuki',
    model: 'Sedan X',
    year: 2023,
    color: 'Silver',
    condition: 'Used',
    price: 15186,
    status: 'Sold',
    dateAdded: '2025-02-20T00:00:00Z',
    serialNumber: 'SX-2023-001',
    createdAt: '2025-02-20T00:00:00Z',
    updatedAt: '2025-04-20T13:30:00Z'
  },
  {
    id: '7',
    make: 'Toyota',
    model: 'Truck Z',
    year: 2023,
    color: 'Silver',
    condition: 'Used',
    price: 49451,
    status: 'Sold',
    dateAdded: '2025-02-15T00:00:00Z',
    serialNumber: 'TZ-2023-001',
    createdAt: '2025-02-15T00:00:00Z',
    updatedAt: '2025-04-18T10:00:00Z'
  },
  // More cars for the model statistics
  ...generateAdditionalCars()
]);

// Get all car IDs for historical sales
const allCarIds = [...mockMultiCabs.value.map(car => car.id), ...generateAdditionalCars().map(car => car.id)];

// Mock data for sales
const mockSales = ref<Sale[]>([
  {
    id: '1',
    customerId: '1',
    soldBy: '1',
    saleDate: '2025-04-28T10:30:00Z',
    totalPrice: 20541,
    createdAt: '2025-04-28T10:30:00Z',
    updatedAt: '2025-04-28T10:30:00Z',
    multiCabId: '1' // Linked to mockMultiCabs[0]
  },
  {
    id: '2',
    customerId: '1',
    soldBy: '1',
    saleDate: '2025-04-26T14:15:00Z',
    totalPrice: 47432,
    createdAt: '2025-04-26T14:15:00Z',
    updatedAt: '2025-04-26T14:15:00Z',
    multiCabId: '2' // Linked to mockMultiCabs[1]
  },
  {
    id: '3',
    customerId: '2',
    soldBy: '2',
    saleDate: '2025-04-24T09:45:00Z',
    totalPrice: 68118,
    createdAt: '2025-04-24T09:45:00Z',
    updatedAt: '2025-04-24T09:45:00Z',
    multiCabId: '3' // Linked to mockMultiCabs[2]
  },
  {
    id: '4',
    customerId: '1',
    soldBy: '3',
    saleDate: '2025-04-23T16:20:00Z',
    totalPrice: 30526,
    createdAt: '2025-04-23T16:20:00Z',
    updatedAt: '2025-04-23T16:20:00Z',
    multiCabId: '4' // Linked to mockMultiCabs[3]
  },
  {
    id: '5',
    customerId: '2',
    soldBy: '2',
    saleDate: '2025-04-23T11:10:00Z',
    totalPrice: 47182,
    createdAt: '2025-04-23T11:10:00Z',
    updatedAt: '2025-04-23T11:10:00Z',
    multiCabId: '5' // Linked to mockMultiCabs[4]
  },
  {
    id: '6',
    customerId: '3',
    soldBy: '2',
    saleDate: '2025-04-20T13:30:00Z',
    totalPrice: 15186,
    createdAt: '2025-04-20T13:30:00Z',
    updatedAt: '2025-04-20T13:30:00Z',
    multiCabId: '6' // Linked to mockMultiCabs[5]
  },
  {
    id: '7',
    customerId: '1',
    soldBy: '4',
    saleDate: '2025-04-18T10:00:00Z',
    totalPrice: 49451,
    createdAt: '2025-04-18T10:00:00Z',
    updatedAt: '2025-04-18T10:00:00Z',
    multiCabId: '7' // Linked to mockMultiCabs[6]
  },
  // Additional historical data for the chart
  ...generateHistoricalSales(allCarIds)
]);

// Mock data for customers
const mockCustomers = ref<Customer[]>([
  {
    id: '1',
    fullName: 'Jane Roe',
    email: 'jane.roe@example.com',
    phone: '555-123-4567',
    address: '123 Main St',
    dateRegistered: '2025-01-15T00:00:00Z',
    createdAt: '2025-01-15T00:00:00Z',
    updatedAt: '2025-01-15T00:00:00Z'
  },
  {
    id: '2',
    fullName: 'Emily Stone',
    email: 'emily.stone@example.com',
    phone: '555-987-6543',
    address: '456 Oak Ave',
    dateRegistered: '2025-01-20T00:00:00Z',
    createdAt: '2025-01-20T00:00:00Z',
    updatedAt: '2025-01-20T00:00:00Z'
  },
  {
    id: '3',
    fullName: 'Chris Evans',
    email: 'chris.evans@example.com',
    phone: '555-456-7890',
    address: '789 Pine Rd',
    dateRegistered: '2025-01-25T00:00:00Z',
    createdAt: '2025-01-25T00:00:00Z',
    updatedAt: '2025-01-25T00:00:00Z'
  }
]);

// Mock data for users (salespeople)
const mockUsers = ref<User[]>([
  {
    id: '1',
    fullName: 'Edward Allen',
    email: 'edward.allen@company.com',
    role: 'salesperson',
    createdAt: '2024-12-01T00:00:00Z',
    updatedAt: '2024-12-01T00:00:00Z',
    isActive: true
  },
  {
    id: '2',
    fullName: 'Charlie Brown',
    email: 'charlie.brown@company.com',
    role: 'salesperson',
    createdAt: '2024-12-05T00:00:00Z',
    updatedAt: '2024-12-05T00:00:00Z',
    isActive: true
  },
  {
    id: '3',
    fullName: 'Bob Johnson',
    email: 'bob.johnson@company.com',
    role: 'salesperson',
    createdAt: '2024-12-10T00:00:00Z',
    updatedAt: '2024-12-10T00:00:00Z',
    isActive: true
  },
  {
    id: '4',
    fullName: 'Alice Smith',
    email: 'alice.smith@company.com',
    role: 'salesperson',
    createdAt: '2024-12-15T00:00:00Z',
    updatedAt: '2024-12-15T00:00:00Z',
    isActive: true
  }
]);

// Table columns
const columns = [
  { name: 'carModel', required: true, label: 'CAR MODEL', align: 'left' as const, field: 'carModel', sortable: true },
  { name: 'date', required: true, label: 'DATE', align: 'left' as const, field: 'date', sortable: true },
  { name: 'price', required: true, label: 'PRICE', align: 'left' as const, field: 'price', sortable: true },
  { name: 'customer', required: true, label: 'CUSTOMER', align: 'left' as const, field: 'customer', sortable: true },
  { name: 'salesperson', required: true, label: 'SALESPERSON', align: 'left' as const, field: 'salesperson', sortable: true },
  { name: 'color', required: true, label: 'COLOR', align: 'left' as const, field: 'color', sortable: true }
];

// Calculate total revenue
const totalRevenue = computed(() => {
  return mockSales.value.reduce((sum, sale) => sum + sale.totalPrice, 0);
});

// Calculate total sales
const totalSales = computed(() => mockSales.value.length);

// Calculate average sale price
const averageSalePrice = computed(() => {
  return totalRevenue.value / totalSales.value;
});

// Find top selling model
const topSellingModel = computed(() => {
  const modelCounts = mockMultiCabs.value.reduce((counts, car) => {
    counts[car.model] = (counts[car.model] || 0) + 1;
    return counts;
  }, {} as Record<string, number>);

  let topModel = '';
  let topCount = 0;

  Object.entries(modelCounts).forEach(([model, count]) => {
    if (count > topCount) {
      topModel = model;
      topCount = count;
    }
  });

  return { model: topModel, count: topCount };
});

// Prepare data for the sales trend chart
const monthlySalesData = computed(() => {
  const months = ['Jun 24', 'Jul 24', 'Aug 24', 'Sep 24', 'Oct 24', 'Nov 24', 'Dec 24', 'Jan 25', 'Feb 25', 'Mar 25', 'Apr 25', 'May 25'];
  const monthlySales = new Array(12).fill(0);
  const monthlyRevenue = new Array(12).fill(0);

  mockSales.value.forEach(sale => {
    const date = new Date(sale.saleDate);
    const month = date.getMonth();
    const year = date.getFullYear();

    let index;
    if (year === 2024 && month >= 5) {
      index = month - 5; // June 2024 (month 5) is index 0
    } else if (year === 2025 && month <= 4) {
      index = month + 7; // January 2025 (month 0) is index 7
    } else {
      return; // Outside our 12-month window
    }

    monthlySales[index]++;
    monthlyRevenue[index] += sale.totalPrice;
  });

  return {
    labels: months,
    datasets: [
      {
        label: 'Total Revenue',
        data: monthlyRevenue,
        borderColor: '#36A2EB',
        tension: 0.4,
        fill: false
      },
      {
        label: 'Units Sold',
        data: monthlySales,
        borderColor: '#4BC0C0',
        tension: 0.4,
        fill: false
      }
    ]
  };
});

// Prepare data for the sales by model chart
const salesByModelData = computed(() => {
  const modelCounts = mockMultiCabs.value.reduce((counts, car) => {
    counts[car.model] = (counts[car.model] || 0) + 1;
    return counts;
  }, {} as Record<string, number>);

  const sortedModels = Object.entries(modelCounts)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6);

  return {
    labels: sortedModels.map(([model]) => model),
    datasets: [
      {
        label: 'Units Sold',
        data: sortedModels.map(([, count]) => count),
        backgroundColor: '#4BC0C0',
        borderWidth: 0
      }
    ]
  };
});

// Format recent sales data for the table
const recentSales = computed(() => {
  return mockSales.value
    .slice(0, 7) // Get only the first 7 sales (non-historical)
    .map(sale => {
      const car = mockMultiCabs.value.find(car => car.id === sale.multiCabId); // Link using multiCabId
      const customer = mockCustomers.value.find(customer => customer.id === sale.customerId);
      const salesperson = mockUsers.value.find(user => user.id === sale.soldBy);

      return {
        id: sale.id,
        carModel: car?.model || 'Unknown Model',
        date: sale.saleDate,
        price: sale.totalPrice,
        customer: customer?.fullName || 'Unknown Customer', // Use fullName
        salesperson: salesperson?.fullName || 'Unknown Salesperson',
        color: car?.color || 'Unknown'
      };
    });
});

// Format number with commas
function formatNumber(num: number): string {
  return num.toLocaleString('en-US', { maximumFractionDigits: 0 });
}

// Format date
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
}
</script>

<template>
  <q-page class="q-pa-md">
    <div class="q-px-lg">
      <h1 class="text-h4 q-mb-none">Car Sales Analytics</h1>
      <p class="text-subtitle1 q-mt-xs">Overview of sales performance and trends.</p>

      <!-- Key Metrics Cards -->
      <div class="row q-col-gutter-md q-mt-md">
        <!-- Total Revenue -->
        <div class="col-md-3 col-sm-6 col-xs-12">
          <q-card bordered>
            <q-card-section>
              <div class="text-caption text-uppercase">TOTAL REVENUE</div>
              <div class="row items-center">
                <q-icon name="attach_money" size="sm" class="q-mr-sm" />
              </div>
              <div class="text-h4">${{ formatNumber(totalRevenue) }}</div>
              <div class="text-caption">Across {{ totalSales }} sales</div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Units Sold -->
        <div class="col-md-3 col-sm-6 col-xs-12">
          <q-card bordered>
            <q-card-section>
              <div class="text-caption text-uppercase">UNITS SOLD</div>
              <div class="row items-center">
                <q-icon name="shopping_cart" size="sm" class="q-mr-sm" />
              </div>
              <div class="text-h4">{{ totalSales }}</div>
              <div class="text-caption">Total cars sold</div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Average Sale Price -->
        <div class="col-md-3 col-sm-6 col-xs-12">
          <q-card bordered>
            <q-card-section>
              <div class="text-caption text-uppercase">AVG. SALE PRICE</div>
              <div class="row items-center">
                <q-icon name="attach_money" size="sm" class="q-mr-sm" />
              </div>
              <div class="text-h4">${{ formatNumber(averageSalePrice) }}</div>
              <div class="text-caption">Average price per unit</div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Top Selling Model -->
        <div class="col-md-3 col-sm-6 col-xs-12">
          <q-card bordered>
            <q-card-section>
              <div class="text-caption text-uppercase">TOP SELLING MODEL</div>
              <div class="row items-center">
                <q-icon name="directions_car" size="sm" class="q-mr-sm" />
              </div>
              <div class="text-h4">{{ topSellingModel.model }}</div>
              <div class="text-caption">{{ topSellingModel.count }} units sold</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Charts Section -->
      <div class="row q-col-gutter-md q-mt-md">
        <!-- Sales Trend Chart -->
        <div class="col-md-6 col-sm-12">
          <q-card bordered>
            <q-card-section>
              <div class="text-h6">Sales Trend (Last 12 Months)</div>
              <sales-trend-chart :chart-data="monthlySalesData" :isDark="$q.dark.isActive" />
            </q-card-section>
          </q-card>
        </div>

        <!-- Sales by Model Chart -->
        <div class="col-md-6 col-sm-12">
          <q-card bordered>
            <q-card-section>
              <div class="text-h6">Sales by Model</div>
              <sales-by-model-chart :chart-data="salesByModelData" :isDark="$q.dark.isActive" />
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Recent Sales Table -->
      <div class="q-mt-md">
        <q-card bordered>
          <q-card-section>
            <div class="text-h6">Recent Sales</div>
          </q-card-section>
          <q-separator />
          <q-table
            :rows="recentSales"
            :columns="columns"
            row-key="id"
            :pagination="{ rowsPerPage: 10 }"
          >
            <template v-slot:body-cell-carModel="props">
              <q-td :props="props">
                <div class="row items-center">
                  <q-icon name="directions_car" size="xs" class="q-mr-sm" />
                  {{ props.value }}
                </div>
              </q-td>
            </template>
            <template v-slot:body-cell-date="props">
              <q-td :props="props">
                <div class="row items-center">
                  <q-icon name="event" size="xs" class="q-mr-sm" />
                  {{ formatDate(props.value) }}
                </div>
              </q-td>
            </template>
            <template v-slot:body-cell-price="props">
              <q-td :props="props">
                <div class="row items-center">
                  <q-icon name="attach_money" size="xs" class="q-mr-sm" />
                  ${{ formatNumber(props.value) }}
                </div>
              </q-td>
            </template>
            <template v-slot:body-cell-customer="props">
              <q-td :props="props">
                <div class="row items-center">
                  <q-icon name="person" size="xs" class="q-mr-sm" />
                  {{ props.value }}
                </div>
              </q-td>
            </template>
            <template v-slot:body-cell-salesperson="props">
              <q-td :props="props">
                <div class="row items-center">
                  <q-icon name="badge" size="xs" class="q-mr-sm" />
                  {{ props.value }}
                </div>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </div>
    </div>
  </q-page>
</template>
