<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted, computed } from 'vue';
import { Chart, registerables } from 'chart.js';

// Register all Chart.js components
Chart.register(...registerables);

// Import or define a simple logger utility to replace console.log statements
const logger = {
  debug: (process.env.NODE_ENV === 'development') ? console.debug : () => {},
  info: (process.env.NODE_ENV === 'development') ? console.info : () => {},
  warn: (process.env.NODE_ENV === 'development') ? console.warn : () => {},
  error: console.error // Keep error logging in all environments
};

// Define dataset type with dataType property for better formatting control
type DataType = 'revenue' | 'quantity' | (string & {});

interface ChartDataset {
  label: string;
  data: number[];
  borderColor: string;
  tension: number;
  fill: boolean;
  dataType?: DataType; // New property to explicitly set data type
}

interface ChartDataType {
  labels: string[];
  datasets: ChartDataset[];
}

// Create formatters once and reuse them
const currencyFormatter = new Intl.NumberFormat('en-US', { 
  style: 'currency',
  currency: 'USD',
  maximumFractionDigits: 0 
});

const quantityFormatter = new Intl.NumberFormat('en-US', { 
  style: 'decimal',
  maximumFractionDigits: 0 
});

const props = defineProps<{
  chartData: ChartDataType;
  isDark: boolean;
  chartType: 'line' | 'bar' | 'area';
  missingDataStrategy?: 'zero' | 'skip' | 'interpolate'; // Configuration option for handling missing data
}>();

const chartCanvas = ref<HTMLCanvasElement | null>(null);
let chart: Chart | null = null;

// Validate chart data structure to ensure it's properly formatted
const isValidChartDataStructure = (data: unknown): boolean => {
  if (!data || typeof data !== 'object') return false;
  
  const chartData = data as Partial<ChartDataType>;
  
  return (
    Array.isArray(chartData.labels) && 
    Array.isArray(chartData.datasets) &&
    chartData.datasets.every((dataset) => (
      typeof dataset === 'object' &&
      typeof dataset.label === 'string' &&
      Array.isArray(dataset.data) &&
      typeof dataset.borderColor === 'string'
    ))
  );
};

// Safe number conversion helper
const safeNumber = (value: unknown): number => {
  if (value === null || value === undefined) return 0;
  const num = Number(value);
  return isNaN(num) ? 0 : num;
};

// Handle missing or invalid values based on the selected strategy
const processMissingValues = (dataset: { data: unknown[] }, labelLength: number): number[] => {
  const strategy = props.missingDataStrategy || 'zero';
  
  // If data has exactly the right length, just convert values to safe numbers
  if (dataset.data.length === labelLength) {
    return dataset.data.map(val => safeNumber(val));
  }
  
  // Handle case where dataset length doesn't match labels length
  switch (strategy) {
    case 'zero': {
      // Fill with zeros
      const zeroFilledData = Array(labelLength).fill(0);
      for (let i = 0; i < Math.min(dataset.data.length, labelLength); i++) {
        zeroFilledData[i] = safeNumber(dataset.data[i]);
      }
      return zeroFilledData;
    }
      
    case 'skip':
      // Keep original data, even if mismatched (Chart.js will handle this)
      return dataset.data.map(val => safeNumber(val));
      
    case 'interpolate': {
      // Simple linear interpolation for missing values
      const interpolatedData = Array(labelLength).fill(0);
      const originalLength = dataset.data.length;
      
      for (let i = 0; i < labelLength; i++) {
        if (i < originalLength) {
          interpolatedData[i] = safeNumber(dataset.data[i]);
        } else {
          // Simple interpolation using the last valid value
          interpolatedData[i] = interpolatedData[i-1] || 0;
        }
      }
      return interpolatedData;
    }
      
    default:
      return Array(labelLength).fill(0);
  }
};

// Determine the data type for a dataset
const getDataType = (dataset: ChartDataset): DataType => {
  // First check explicit dataType property
  if (dataset.dataType) {
    return dataset.dataType;
  }
  
  // Fall back to checking the label (for backward compatibility)
  if (dataset.label?.toLowerCase().includes('revenue')) {
    return 'revenue';
  }
  
  return 'quantity';
};

// Validate and process chart data
const validatedChartData = computed(() => {
  if (!isValidChartDataStructure(props.chartData)) {
    logger.error('Invalid chart data structure. Expected format: { labels: string[], datasets: [{ label: string, data: number[], borderColor: string, ... }] }');
    return {
      labels: [],
      datasets: []
    };
  }
  
  const labelsLength = props.chartData.labels.length;
  
  // Process each dataset
  const datasets = props.chartData.datasets.map(dataset => {
    if (!dataset.data || dataset.data.length !== labelsLength) {
      logger.warn(
        `Dataset "${dataset.label}" has mismatched data length. ` +
        `Expected: ${labelsLength}, Got: ${dataset.data?.length || 0}. ` +
        `Using "${props.missingDataStrategy || 'zero'}" strategy.`
      );
      
      // Ensure dataset has data property before processing
      const datasetWithData = {
        ...dataset,
        data: Array.isArray(dataset.data) ? dataset.data : []
      };
      
      const processedData = processMissingValues(datasetWithData, labelsLength);
      return { ...dataset, data: processedData };
    }
    return dataset;
  });

  return {
    labels: props.chartData.labels,
    datasets
  };
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
                // Get the dataset with proper type information
                const typedDataset = context.dataset as unknown as ChartDataset;
                const dataType = getDataType(typedDataset);
                
                // Format according to data type
                if (dataType === 'revenue') {
                  label += currencyFormatter.format(context.parsed.y);
                } else {
                  label += quantityFormatter.format(context.parsed.y) + ' units';
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
  // Delay chart creation slightly to ensure data is fully loaded and canvas is ready
  setTimeout(() => {
    createChart();
  }, 50);
});

// Watch for changes in chartData, chartType, and isDark
watch([() => props.chartData, () => props.chartType, () => props.isDark, () => props.missingDataStrategy], () => {
  logger.debug('Chart configuration changed, recreating chart');
  // Use setTimeout to ensure the DOM has had time to update
  setTimeout(() => {
    createChart();
  }, 50);
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
