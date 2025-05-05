<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue';
import { Chart, registerables } from 'chart.js';

// Register all Chart.js components
Chart.register(...registerables);

const props = defineProps<{
  chartData: {
    labels: string[];
    datasets: {
      label: string;
      data: number[];
      borderColor: string;
      tension: number;
      fill: boolean;
    }[];
  };
  isDark: boolean;
  chartType: 'line' | 'bar' | 'area';
}>();

const chartCanvas = ref<HTMLCanvasElement | null>(null);
let chart: Chart | null = null;

const createChart = () => {
  if (!chartCanvas.value) return;

  const ctx = chartCanvas.value.getContext('2d');
  if (!ctx) return;

  // Destroy existing chart if it exists
  if (chart) {
    chart.destroy();
  }

  const isArea = props.chartType === 'area';
  const isBar = props.chartType === 'bar';
  
  const datasets = props.chartData.datasets.map((dataset) => ({
    ...dataset,
    backgroundColor: isBar 
      ? dataset.borderColor + '80'  // Add 50% transparency
      : isArea
        ? dataset.borderColor + '20' // Add 87.5% transparency
        : undefined,
    fill: isArea ? true : dataset.fill,
  }));

  chart = new Chart(ctx, {
    type: isArea ? 'line' : props.chartType,
    data: {
      labels: props.chartData.labels,
      datasets
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: {
          beginAtZero: true,
          grid: {
            color: props.isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)',
          },
          ticks: {
            color: props.isDark ? 'rgba(255, 255, 255, 0.7)' : 'rgba(0, 0, 0, 0.7)',
          }
        },
        x: {
          grid: {
            color: props.isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)',
          },
          ticks: {
            color: props.isDark ? 'rgba(255, 255, 255, 0.7)' : 'rgba(0, 0, 0, 0.7)',
          }
        }
      },
      plugins: {
        legend: {
          labels: {
            color: props.isDark ? 'rgba(255, 255, 255, 0.7)' : 'rgba(0, 0, 0, 0.7)',
          }
        },
        tooltip: {
          backgroundColor: props.isDark ? 'rgba(0, 0, 0, 0.7)' : 'rgba(255, 255, 255, 0.7)',
          titleColor: props.isDark ? 'rgba(255, 255, 255, 0.9)' : 'rgba(0, 0, 0, 0.9)',
          bodyColor: props.isDark ? 'rgba(255, 255, 255, 0.9)' : 'rgba(0, 0, 0, 0.9)',
        }
      }
    }
  });
};

onMounted(() => {
  createChart();
});

// Watch for changes in chartData, chartType, and isDark
watch([() => props.chartData, () => props.chartType, () => props.isDark], () => {
  createChart();
}, { deep: true });

onUnmounted(() => {
  if (chart) {
    chart.destroy();
  }
});
</script>

<template>
  <div class="chart-container">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<style scoped>
.chart-container {
  position: relative;
  height: 300px;
  width: 100%;
}
</style>
