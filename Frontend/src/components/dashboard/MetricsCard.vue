<template>
  <q-card bordered class="metrics-card">
    <q-card-section class="q-pa-md">
      <!-- Title and Icon -->
      <div class="row items-center q-mb-md">
        <q-icon :name="icon" size="xs" class="q-mr-sm" />
        <div class="text-caption text-grey-7">{{ title }}</div>
      </div>

      <!-- Main Value -->
      <div class="text-h3 q-mb-md">{{ value }}</div>

      <!-- Growth and YTD -->
      <div class="row justify-between items-end q-mb-lg">
        <div class="column">
          <div class="row items-center q-mb-xs" v-if="trendPercentage !== undefined">
            <q-icon
              :name="trendPercentage >= 0 ? 'arrow_upward' : 'arrow_downward'"
              :class="trendPercentage >= 0 ? 'text-positive' : 'text-negative'"
              size="xs"
            />
            <span :class="[
              'text-caption q-ml-xs',
              trendPercentage >= 0 ? 'text-positive' : 'text-negative'
            ]">
              {{ Math.abs(trendPercentage).toFixed(1) }}%
            </span>
            <span class="text-caption text-grey-7 q-ml-xs">from last month</span>
          </div>
          <div class="text-caption text-grey-7 q-mb-xs" v-if="ytdValue">
            YTD: {{ ytdValue }}
          </div>
          <div class="text-caption text-grey-7" v-if="additionalInfo">
            {{ additionalInfo }}
          </div>
        </div>
      </div>

      <!-- Trend Graph -->
      <div class="trend-graph" v-if="trendData && trendData.length > 0">
        <canvas ref="chartCanvas"></canvas>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { Chart, registerables } from 'chart.js';

Chart.register(...registerables);

const props = defineProps<{
  title: string;
  icon: string;
  value: string | number;
  trendData?: number[];
  trendPercentage?: number;
  ytdValue?: string;
  additionalInfo?: string;
}>();

const chartCanvas = ref<HTMLCanvasElement | null>(null);
let chart: Chart | null = null;

const createChart = () => {
  if (!chartCanvas.value || !props.trendData) return;

  const ctx = chartCanvas.value.getContext('2d');
  if (!ctx) return;

  if (chart) {
    chart.destroy();
  }

  const isPositive = (props.trendPercentage || 0) >= 0;
  const color = isPositive ? '#21BA45' : '#C10015';
  
  // Create gradient
  const gradient = ctx.createLinearGradient(0, 0, 0, 60);
  gradient.addColorStop(0, `${color}15`);
  gradient.addColorStop(1, `${color}05`);

  chart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: new Array(props.trendData.length).fill(''),
      datasets: [{
        data: props.trendData,
        borderColor: color,
        borderWidth: 1.5,
        tension: 0.4,
        fill: true,
        backgroundColor: gradient,
        pointRadius: 0,
        pointHoverRadius: 0
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        tooltip: {
          enabled: false
        }
      },
      scales: {
        x: {
          display: false,
          grid: {
            display: false
          }
        },
        y: {
          display: false,
          grid: {
            display: false
          }
        }
      }
    }
  });
};

onMounted(() => {
  if (props.trendData) {
    createChart();
  }
});

watch(() => props.trendData, () => {
  if (props.trendData) {
    createChart();
  }
}, { deep: true });
</script>

<style scoped>
.metrics-card {
  transition: all 0.3s;
  height: 100%;
  min-height: 200px;
}

.metrics-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12);
}

.trend-graph {
  height: 60px;
  margin: 0 -16px -16px -16px;
  position: relative;
  bottom: -8px;
}

.text-h3 {
  font-weight: 600;
  font-size: 1.85rem;
  line-height: 1.2;
}

.q-card__section {
  padding-bottom: 8px;
}
</style> 