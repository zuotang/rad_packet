<script setup>
defineProps({
  title: { type: String, required: true },
  desc: { type: String, default: "" },
});
</script>

<template>
  <section class="card page-head">
    <h2>{{ title }}</h2>
    <p class="muted" v-if="desc">{{ desc }}</p>
  </section>
</template>
