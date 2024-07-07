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
  <div>
    <h1>Home Page</h1>
    <ul>
      <li v-for="(value, key) in items" :key="key">
        <RouterLink :to="`/wiki/${value}`">{{ value }}</RouterLink>
      </li>
    </ul>
  </div>
</template>

<style scoped>

</style>