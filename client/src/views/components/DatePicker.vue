<template>
  <v-menu
    v-model="menu"
    :close-on-content-click="false"
    transition="scale-transition"
    offset-y
    max-width="290px"
    min-width="290px"
  >
    <template #activator="{ props }">
      <v-text-field
        v-model="computedDateFormatted"
        :rules="nullRules"
        label="支払日"
        filled
        readonly
        placeholder="2020年5月2日"
        height="10"
        v-bind="props"
      />
    </template>
    <v-date-picker
      v-model="date"
      no-title
      color="primary"
      @update:model-value="menu = false"
    />
  </v-menu>
</template>

<script lang="ts">
export default {
  props: {},
  data() {
    return {
      menu: false,
      date: null,
      nullRules: [(v: unknown) => !!v || ""]
    };
  },
  computed: {
    computedDateFormatted() {
      if (!this.date) return null;
      const d = new Date(this.date);
      return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`;
    }
  }
};
</script>

<style lang="scss" module></style>
