<script setup lang="ts">
import { computed } from "vue";
import { Email } from "../types";
import { useEmailStore } from "../stores/emailStore";

const emailStore = useEmailStore();
const email = computed(() => emailStore.selectedEmail);

</script>

<template>
  <div class="p-6 h-full">
    <div v-if="email" class="h-full overflow-auto hover:overflow-scroll">
      <div class="flex flex-row mb-2">
        <p class="font-bold w-[12%]">Date:</p>
        <p>{{ email.date }}</p>
      </div>
      <div class="flex flex-row mb-2">
        <p class="font-bold w-[12%]">From:</p>
        <p>{{ email.from }}</p>
      </div>
      <div class="flex flex-row mb-2">
        <p class="font-bold w-[12%]">To:</p>
        <p>{{ email.to }}</p>
      </div>
      <div class="flex flex-row">
        <p class="font-bold w-[12%]">Subject:</p>
        <p>{{ email.subject }}</p>
      </div>
      <div class="mt-4">
        <p class="font-bold">Content:</p>
        <pre class="font-roboto">{{ email.content }}</pre>
      </div>
    </div>
    <div v-else>
      <p>No email selected</p>
    </div>
  </div>
</template>
