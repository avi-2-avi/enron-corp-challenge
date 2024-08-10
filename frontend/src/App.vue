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
    <div class="flex flex-col xl:flex-row xl:flex-grow xl:space-x-4 xl:h-80 mx-8">
      <div class="xl:basis-7/12">
        <Table />
      </div>
      <div class="xl:hidden my-2 xl:my-0">
        <Divider />
      </div>
      <Divider orientation="vertical" />
      <div class="xl:basis-5/12">
        <ContentDetail />
      </div>
    </div>
    <Divider />
    <Footer />
  </div>

</template>
