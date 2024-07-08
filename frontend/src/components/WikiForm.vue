<script setup>
import {marked} from 'marked'
import {debounce} from 'lodash-es'
import {ref, computed} from 'vue'

const page = ref('/')
const content = ref('# hello')

const output = computed(() => marked(content.value))
const updateContent = debounce((e) => {
  content.value = e.target.value
}, 100)
</script>

<template>
  <div class="editor">
    <div class="m-2">
      <label>
        <span class="p-2">Page</span>
        <input class="input lg:min-w-full bg-gray-700 rounded-sm p-2" :value="page"/>
      </label>
    </div>
    <div class="m-2">
      <label>
        <span class="p-2">Content</span>
        <textarea class="input lg:min-w-full lg:min-h-96 bg-gray-700 rounded-sm p-2" :value="content"
                  @input="updateContent"/>
      </label>
    </div>
    <div class="output" v-html="output"></div>
  </div>
</template>

<style scoped>

</style>