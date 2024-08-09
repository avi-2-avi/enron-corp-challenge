<script setup lang="ts">
import Header from "./components/Header.vue";
import Footer from "./components/Footer.vue";
import Divider from "./components/Divider.vue";
import ContentDetail from "./components/ContentDetail.vue";
import Table from "./components/Table.vue";
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <Header />
    <Divider />
    <div class="flex flex-row flex-grow space-x-4 h-80 mx-8">
      <div class="basis-2/5">
        <Table />
      </div>
      <Divider orientation="vertical" />
      <div class="basis-3/5">
        <ContentDetail />
      </div>
    </div>
    <Divider />
    <Footer />
  </div>

</template>
