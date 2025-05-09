<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted, computed } from 'vue';
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

// Validate chart data to ensure it's properly formatted
const validatedChartData = computed(() => {
  console.log('Computing validatedChartData with input:', props.chartData);
  
  if (!props.chartData || !props.chartData.datasets || !props.chartData.labels) {
    console.error('Invalid chart data format');
    return {
      labels: [],
      datasets: []
    };
  }
  
  // Ensure datasets all have the same length as labels
  const datasets = props.chartData.datasets.map(dataset => {
    if (!dataset.data || dataset.data.length !== props.chartData.labels.length) {
      console.warn(`Dataset ${dataset.label} has mismatched data length. Expected: ${props.chartData.labels.length}, Got: ${dataset.data?.length || 0}`);
      // Fill missing values with zeros if needed
      const data = Array(props.chartData.labels.length).fill(0);
      if (dataset.data) {
        for (let i = 0; i < Math.min(dataset.data.length, data.length); i++) {
          data[i] = dataset.data[i];
        }
      }
      return { ...dataset, data };
    }
    return dataset;
  });

  const result = {
    labels: props.chartData.labels,
    datasets
  };
  
  console.log('Validated chart data:', result);
  return result;
});

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
  
  // Use validated chart data
  const datasets = validatedChartData.value.datasets.map((dataset) => ({
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
      labels: validatedChartData.value.labels,
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
          },
          title: {
            display: true,
            text: 'Quantity Sold',
            color: props.isDark ? 'rgba(255, 255, 255, 0.7)' : 'rgba(0, 0, 0, 0.7)',
            font: {
              weight: 'bold'
            }
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
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || '';
              if (label) {
                label += ': ';
              }
              if (context.parsed.y !== null) {
                // Check if the dataset is tracking quantity or revenue
                if (context.dataset.label?.toLowerCase().includes('revenue')) {
                  label += new Intl.NumberFormat('en-US', { 
                    style: 'currency',
                    currency: 'USD',
                    maximumFractionDigits: 0 
                  }).format(context.parsed.y);
                } else {
                  // For quantities (Cabs, Accessories, Total items)
                  label += new Intl.NumberFormat('en-US', { 
                    style: 'decimal',
                    maximumFractionDigits: 0 
                  }).format(context.parsed.y) + ' units';
                }
              }
              return label;
            }
          }
        }
      }
    }
  });
};

onMounted(() => {
  createChart();
});

// Watch for changes in chartData, chartType, and isDark
watch([() => props.chartData, () => props.chartType, () => props.isDark], (newValues) => {
  console.log('Chart data, type, or theme changed');
  console.log('New chart data:', newValues[0]);
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
    <div v-if="!validatedChartData.datasets.length" class="no-data">
      No data available to display
    </div>
  </div>
</template>

<style scoped>
.chart-container {
  position: relative;
  height: 300px;
  width: 100%;
}

.no-data {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #888;
  font-style: italic;
}
</style>
