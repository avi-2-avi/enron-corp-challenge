<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { Email, GetEmailsParams } from "../types";
import { useEmailStore } from "../stores/emailStore";
import Spinner from "./Spinner.vue";

const emailStore = useEmailStore();

const emailsPerPage = 10;
const currentPage = ref(1);
const isLoading = ref(false);
const filterTerm = ref("");

const fetchEmailData = async (id: string) => {
  try {
    await emailStore.fetchEmail(id);
  } catch (error) {
    console.error('Failed to fetch email:', error);
  }
};

const fetchEmails = async () => {
  try {
    isLoading.value = true;
    const emailParams = <GetEmailsParams>({
      page: currentPage.value,
      size: emailsPerPage,
      filter: filterTerm.value
    });
    await emailStore.fetchEmails(emailParams);
  } catch (error) {
    console.error('Failed to fetch emails:', error);
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchEmails();
});

watch(currentPage, () => {
  console.log('currentPage:', currentPage.value);
  console.log('filterTerm:', filterTerm.value);
  fetchEmails();
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
    <input type="text" v-model="filterTerm" @keydown.enter="fetchEmails" :disabled="isLoading" placeholder="Search for content keyword..." class="w-full p-2 border border-blue rounded-md" />
    <div v-if="isLoading">
      <Spinner />
    </div>
    <div v-else-if="emailStore.emails.length > 0">
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
            <tr v-for="(email, index) in emailStore.emails" @click="fetchEmailData(email.id)" :key="email.id" :class="['hover:bg-blue/5 h-12',
              emailStore.selectedEmail !== null && emailStore.selectedEmail.id === email.id ? 'bg-blue/5' : '',
              index !== emailStore.emails.length - 1 ? 'border-black/5 border-b' : '']">
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
    <div v-else>
      <p class="mt-4">No emails found.</p>
    </div>
  </div>
</template>