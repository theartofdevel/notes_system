<template>
  <router-view></router-view>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  beforeCreate: function () {
    const token = localStorage.getItem("ns_token");
    if (!token) {
      this.$router.push({ path: "/login" });
    }
  },
});
</script>
