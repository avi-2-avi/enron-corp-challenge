<script setup lang="ts">
defineProps({
  orientation: {
    type: String,
    default: 'horizontal',
  },
});
</script>

<template>
   <div
    :class="[
      'bg-blue',
      orientation === 'horizontal'
        ? 'mx-auto max-w-full w-[calc(100%-4rem)] h-px'
        : 'my-auto max-h-full h-[calc(100%-2rem)] w-px',
    ]"
  />
</template>
