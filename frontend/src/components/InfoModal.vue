<script setup lang="ts">
defineProps<{
  show: boolean;
  close: () => void;
}>();
</script>

<template>
    <div v-if="show" class="fixed inset-0 flex items-center justify-center z-50">
        <div class="fixed inset-0 bg-black opacity-50" @click="close"/>
        <div class="flex flex-col bg-white p-10 rounded-lg shadow-lg z-10 max-w-sm lg:max-w-lg w-full">
            <p class="font-bold text-xl mb-4 text-center">How to Use</p>
            <p class="mb-2">Search for any word within the email content using a single term.</p>
            <p class="mb-2">Sort the data in ascending or descending order by clicking on the column headers.</p>
            <p class="mb-2">Navigate through the pages using the navigation panel below.</p>
            <p class="mb-2">For quicker navigation, enter the page number in the input field at the bottom right and
                press
                "Enter" to go directly to that page.</p>
            <p class="mb-4">This site is optimized for use on various devices, including mobile, tablet, and desktop.
            </p>
            <button @click="close" class="mt-4 px-4 py-2 bg-black text-white rounded-md w-40 mx-auto hover:bg-black/90">
                Close
            </button>
        </div>
    </div>
</template>
