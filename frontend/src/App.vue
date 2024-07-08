<script setup>
import {RouterLink, RouterView} from 'vue-router'
</script>

<template>
  <header>
    <nav class="flex justify-between mb-4">
      <div>
        <RouterLink to="/"
                    class="text-4xl font-extrabold leading-none tracking-tight text-gray-900 md:text-5xl lg:text-6xl dark:text-white">
          <span class="text-transparent bg-clip-text bg-gradient-to-r to-emerald-600 from-sky-400">Wiki</span>
        </RouterLink>
      </div>
      <div class="flex items-end">
        <RouterLink to="/new"
                    class="text-2xl font-extrabold mb-0 dark:text-white">
          <span>Create</span>
        </RouterLink>
      </div>
    </nav>
  </header>

  <RouterView class="pt-2"/>
</template>

<style scoped>
</style>
