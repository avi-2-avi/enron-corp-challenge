<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { GetEmailsParams } from "../types/email";
import { useEmailStore } from "../stores/emailStore";
import { formatNumber } from "../utils/formatNumber";
import Spinner from "./Spinner.vue";

const emailStore = useEmailStore();

const emailsPerPage = 10;
const maxPageNum = 100;
const currentPage = ref(1);
const isLoading = ref(false);
const filterTerm = ref("");
const sortColumn = ref<string | null>(null);
const sortOrder = ref<'asc' | 'desc'>('desc');

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
      filter: filterTerm.value,
      sort: sortColumn.value || undefined,
      order: sortOrder.value
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
  fetchEmails();
});

const changePage = (page: number) => {
  if (page > 0 && page <= emailStore.totalPages) {
    currentPage.value = page;
  }
};

const updatePage = (event: KeyboardEvent) => {
  const input = event.target as HTMLInputElement;
  const newPage = parseInt(input.value);

  if (!isNaN(newPage) && newPage > 0) {
    currentPage.value = newPage;
    fetchEmails();
  }
};

const pages = computed(() => {
  const totalPages = Math.min(emailStore.totalPages, maxPageNum);
  const current = currentPage.value;

  let start = Math.max(1, current - 1);
  let end = Math.min(totalPages, start + 3);

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

const toggleSort = (column: string) => {
  if (sortColumn.value === column) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortColumn.value = column;
    sortOrder.value = 'asc';
  }
  fetchEmails();
};

</script>

<template>
  <div class="py-4 md:p-6 flex flex-col h-full">
    <div class="flex flex-row w-full">
      <input type="text" v-model="filterTerm" @keydown.enter="fetchEmails" :disabled="isLoading"
        placeholder="Type a content keyword..."
        class="flex-grow p-2 border border-r-0 border-blue rounded-md rounded-r-none" />
      <button @click="fetchEmails" :disabled="isLoading"
        class="h-full p-2 border border-l-0 border-blue rounded-md rounded-l-none bg-blue text-white">
        Search
      </button>
    </div>
    <div v-if="isLoading">
      <Spinner />
    </div>
    <div v-else-if="emailStore.emails.length > 0">
      <div class="mt-8 flex-grow rounded-xl">
        <table class="table-fixed w-full shadow-inner rounded-xl">
          <thead class="rounded-xl border-black/5 border-b">
            <tr>
              <th class="py-3 px-4 text-left w-1/3 hover:cursor-pointer" 
              @click="toggleSort('subject')">
                Subject
                <span class="text-blue pl-1" v-if="sortColumn === 'subject'">
                  {{ sortOrder === 'asc' ? '▲' : '▼' }}
                </span>
              </th>
              <th class="hidden sm:table-cell py-3 px-4 text-left w-1/3 hover:cursor-pointer"
                @click="toggleSort('from')">
                From
                <span class="text-blue pl-1" v-if="sortColumn === 'from'">
                  {{ sortOrder === 'asc' ? '▲' : '▼' }}
                </span>
              </th>
              <th class="hidden md:table-cell py-3 px-4 text-left w-1/3 hover:cursor-pointer"
                @click="toggleSort('to')">
                To
                <span class="text-blue pl-1" v-if="sortColumn === 'to'">
                  {{ sortOrder === 'asc' ? '▲' : '▼' }}
                </span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(email, index) in emailStore.emails" @click="fetchEmailData(email.id)" :key="email.id" :class="['hover:bg-blue/5 h-12',
              emailStore.selectedEmail !== null && emailStore.selectedEmail.id === email.id ? 'bg-blue/5' : '',
              index !== emailStore.emails.length - 1 ? 'border-black/5 border-b' : '']">
              <td class="py-3 px-4 truncate">{{ email.subject }}</td>
              <td class="hidden sm:table-cell py-3 px-4 text-ellipsis truncate">{{ email.from }}</td>
              <td class="hidden md:table-cell py-3 px-4 text-ellipsis truncate">{{ email.to }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="flex flex-row justify-center sm:justify-between items-center mt-4">
        <div class="hidden sm:flex">
          <small>Total Found: <span class="font-bold">{{ formatNumber(emailStore.totalElements) }}</span></small>
        </div>
        <div class="flex justify-center space-x-2 text-blue">
          <button @click="changePage(1)" :disabled="currentPage === 1" class="p-1 disabled:text-blue/40">
            &laquo;
          </button>
          <button @click="previousPage" :disabled="currentPage === 1" class="p-1 disabled:text-blue/40">
            &lsaquo;
          </button>
          <button v-for="page in pages" :key="page" @click="changePage(page)"
            :class="['px-4', { 'font-bold': currentPage === page }]"
            class="border border-blue rounded hover:bg-blue/5 text-black">
            {{ page }}
          </button>
          <button @click="nextPage" :disabled="currentPage === emailStore.totalPages || currentPage === maxPageNum"
            class="p-1 disabled:text-blue/40">
            &rsaquo;
          </button>
          <button @click="changePage(maxPageNum)" :disabled="currentPage === maxPageNum"
            class="p-1 disabled:text-blue/40">
            &raquo;
          </button>
        </div>
        <div class="hidden sm:flex">
          <input class="py-1 px-0 border border-blue rounded-md text-center w-16" type="number"
            :placeholder="currentPage.toString()" @keydown.enter="updatePage" />
        </div>
      </div>
    </div>
    <div v-else>
      <p class="mt-4">No emails found.</p>
    </div>
  </div>
</template>