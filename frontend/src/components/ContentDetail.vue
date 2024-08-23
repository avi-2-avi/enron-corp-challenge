<script setup lang="ts">
import { computed } from "vue";
import { useEmailStore } from "../stores/emailStore";

const emailStore = useEmailStore();
const email = computed(() => emailStore.selectedEmail);
</script>

<template>
  <div class="h-1/2 py-2 md:p-6 xl:h-full mb-4 xl:mb-0">
    <div v-if="email" class="h-1/2 xl:h-full overflow-auto hover:overflow-scroll">
      <div class="flex flex-col md:flex-row mb-2">
        <p class="font-bold w-[15%]">Date:</p>
        <p>{{ email.date }}</p>
      </div>
      <div class="flex flex-col md:flex-row mb-2">
        <p class="font-bold w-[15%]">From:</p>
        <p>{{ email.from }}</p>
      </div>
      <div class="flex flex-col md:flex-row mb-2">
        <p class="font-bold w-[15%]">To:</p>
        <p>{{ email.to }}</p>
      </div>
      <div class="flex flex-col md:flex-row">
        <p class="font-bold w-[15%]">Subject:</p>
        <p>{{ email.subject }}</p>
      </div>
      <div class="mt-4">
        <p class="font-bold">Content:</p>
        <p class="font-roboto whitespace-pre-line" v-html="email.content"></p>
      </div>
    </div>
    <div v-else>
      <p>No email selected.</p>
    </div>
  </div>
</template>
