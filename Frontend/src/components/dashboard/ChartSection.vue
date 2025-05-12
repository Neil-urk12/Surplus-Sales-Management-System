<script setup lang="ts">
import { defineProps } from 'vue';
import SalesTrendChart from '../charts/SalesTrendChart.vue';
import InventoryChart from './InventoryChart.vue';

interface ChartDataset {
  label: string;
  data: number[];
  borderColor: string;
  tension: number;
  fill: boolean;
}

interface ChartData {
  labels: string[];
  datasets: ChartDataset[];
}

interface InventoryChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string[];
  }[];
}

defineProps({
  salesTrendData: {
    type: Object as () => ChartData,
    required: true
  },
  inventoryData: {
    type: Object as () => InventoryChartData,
    required: true
  },
  timePeriod: {
    type: String as () => 'weekly' | 'monthly' | 'yearly',
    required: true
  },
  chartType: {
    type: String as () => 'line' | 'bar' | 'area',
    required: true
  },
  isDark: {
    type: Boolean,
    default: false
  },
  isLoading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:timePeriod', 'update:chartType']);

const updateTimePeriod = (value: 'weekly' | 'monthly' | 'yearly') => {
  emit('update:timePeriod', value);
};

const updateChartType = (type: 'line' | 'bar' | 'area') => {
  emit('update:chartType', type);
};
</script>

<template>
  <div class="row q-col-gutter-lg q-mb-xl">
    <div class="col-12 col-md-8">
      <q-card bordered class="chart-card">
        <q-card-section>
          <div class="row items-center justify-between q-mb-md">
            <div class="text-h6">Sales Trend</div>
            <div class="row items-center q-gutter-md">
              <!-- Time Period Dropdown -->
              <q-select
                :model-value="timePeriod"
                @update:model-value="updateTimePeriod"
                :options="[
                  { label: 'Weekly', value: 'weekly' },
                  { label: 'Monthly', value: 'monthly' },
                  { label: 'Yearly', value: 'yearly' }
                ]"
                dense
                outlined
                style="width: 150px"
                emit-value
                map-options
                :disable="isLoading"
              />

              <!-- Chart Type Button Group -->
              <q-btn-group flat>
                <q-btn 
                  :class="[
                    chartType === 'line' ? 
                      (isDark ? 'bg-white text-black' : 'bg-primary text-white') : 
                      'bg-white text-black'
                  ]" 
                  @click="updateChartType('line')" 
                  :disable="isLoading"
                >
                  <span>Line Chart</span>
                </q-btn>
                <q-btn 
                  :class="[
                    chartType === 'bar' ? 
                      (isDark ? 'bg-white text-black' : 'bg-primary text-white') : 
                      'bg-white text-black'
                  ]" 
                  @click="updateChartType('bar')" 
                  :disable="isLoading"
                >
                  <span>Bar Chart</span>
                </q-btn>
                <q-btn 
                  :class="[
                    chartType === 'area' ? 
                      (isDark ? 'bg-white text-black' : 'bg-primary text-white') : 
                      'bg-white text-black'
                  ]" 
                  @click="updateChartType('area')" 
                  :disable="isLoading"
                >
                  <span>Area Chart</span>
                </q-btn>
              </q-btn-group>
            </div>
          </div>
          <div v-if="isLoading" class="chart-skeleton">
            <q-skeleton type="rect" height="300px" />
          </div>
          <sales-trend-chart v-else :chart-data="salesTrendData" :chart-type="chartType" :isDark="isDark" />
        </q-card-section>
      </q-card>
    </div>
    <div class="col-12 col-md-4">
      <q-card bordered class="chart-card">
        <q-card-section>
          <div class="text-h6">Inventory Distribution</div>
          <div v-if="isLoading" class="chart-skeleton">
            <q-skeleton type="rect" height="300px" />
          </div>
          <inventory-chart v-else :chart-data="inventoryData as InventoryChartData" :isDark="isDark" />
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<style scoped>
.chart-card {
  height: 100%;
  min-height: 400px;
}

.chart-skeleton {
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>

<style scoped>
.bg-white.text-black {
  border: 1px solid black;
}
</style>
