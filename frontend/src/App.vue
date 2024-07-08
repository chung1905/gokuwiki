<script setup>
import {RouterLink, RouterView} from 'vue-router'
</script>

<template>
  <header>
    <nav>
      <RouterLink to="/"
                  class="mb-4 text-4xl font-extrabold leading-none tracking-tight text-gray-900 md:text-5xl lg:text-6xl dark:text-white">
        <span class="text-transparent bg-clip-text bg-gradient-to-r to-emerald-600 from-sky-400">Wiki</span>
      </RouterLink>
    </nav>
  </header>

  <RouterView class="pt-2"/>
</template>

<style scoped>
</style>
