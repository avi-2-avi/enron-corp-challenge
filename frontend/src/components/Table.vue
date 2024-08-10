<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { Email } from "../types";
import { useEmailStore } from "../stores/emailStore";
import Spinner from "./Spinner.vue";

const emailStore = useEmailStore();

const emailsPerPage = 10;
const currentPage = ref(1);

const fetchEmailData = async (id: string) => {
  try {
    await emailStore.fetchEmail(id);
  } catch (error) {
    console.error('Failed to fetch email:', error);
  }
};

const fetchEmails = async () => {
  try {
    await emailStore.fetchEmails();
  } catch (error) {
    console.error('Failed to fetch emails:', error);
  }
};

onMounted(() => {
  fetchEmails();
});

const paginatedEmails = computed(() => {
  const start = (currentPage.value - 1) * emailsPerPage;
  return emailStore.emails.slice(start, start + emailsPerPage);
});

const changePage = (page: number) => {
  if (page > 0 && page <= emailStore.totalPages) {
    currentPage.value = page;
  }
};

const pages = computed(() => {
  const current = currentPage.value;

  let start = Math.max(1, current - 1);
  let end = Math.min(emailStore.totalPages, start + 3);

  if (end - start < 3) {
    start = Math.max(1, end - 3);
  }

  return Array.from({ length: end - start + 1 }, (_, i) => i + start);
});

const nextPage = () => {
  if (currentPage.value < emailStore.totalPages) {
    currentPage.value++;
  }
};

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};

</script>

<template>
  <div class="p-6 flex flex-col h-full">
    <input type="text" placeholder="Search" class="w-full p-2 border border-blue rounded-md" />
    <div v-if="paginatedEmails.length === 0">
 <Spinner />
    </div>
    <div v-else>
      <div class="mt-8 flex-grow rounded-xl">
        <table class="table-fixed w-full shadow-inner rounded-xl">
          <thead class="rounded-xl border-black/5 border-b">
            <tr>
              <th class="py-3 px-4 text-left w-1/3">Subject</th>
              <th class="py-3 px-4 text-left w-1/3">From</th>
              <th class="py-3 px-4 text-left w-1/3">To</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(email, index) in paginatedEmails" @click="fetchEmailData(email.id)" :key="email.id" :class="['hover:bg-blue/5 h-12',
              emailStore.selectedEmail !== null && emailStore.selectedEmail.id === email.id ? 'bg-blue/5' : '',
              index !== paginatedEmails.length - 1 ? 'border-black/5 border-b' : '']">
              <td class="py-3 px-4 truncate">{{ email.subject }}</td>
              <td class="py-3 px-4 text-ellipsis truncate">{{ email.from }}</td>
              <td class="py-3 px-4 text-ellipsis truncate">{{ email.to }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="mt-4 flex justify-center space-x-2 text-blue">
        <button @click="previousPage" :disabled="currentPage === 1" class="p-2">
          &lt;
        </button>
        <button v-for="page in pages" :key="page" @click="changePage(page)"
          :class="['px-4', { 'font-bold': currentPage === page }]"
          class="border border-blue rounded hover:bg-blue/5 text-black">
          {{ page }}
        </button>
        <button @click="nextPage" :disabled="currentPage === emailStore.totalPages" class="p-2">
          &gt;
        </button>
      </div>
    </div>

  </div>
</template>