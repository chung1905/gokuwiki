<script setup>
import {ref, onMounted} from 'vue';

const items = ref([]);

const fetchData = async () => {
  try {
    const response = await fetch('/api/wiki/list');
    const data = await response.json();
    items.value = data.result;
  } catch (error) {
    console.error('Error fetching data:', error);
  }
};

onMounted(fetchData);
</script>

<template>
  <div class="flex flex-col">
    <RouterLink :to="`/wiki/${value}`" v-for="(value, key) in items" :key="key"
                class="border-solid border border-gray-500 p-2 rounded-lg m-1">{{ value }}
    </RouterLink>
  </div>
</template>

<style scoped>

</style>