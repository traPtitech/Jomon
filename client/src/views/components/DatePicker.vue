<template>
  <v-menu
    v-model="menu"
    :close-on-content-click="false"
    transition="scale-transition"
    offset-y
    max-width="290px"
    min-width="290px"
  >
    <template v-slot:activator="{ on }">
      <v-text-field
        v-model="computedDateFormatted"
        :rules="nullRules"
        label="支払日"
        filled
        readonly
        placeholder="2020年5月2日"
        v-on="on"
        height="10"
      ></v-text-field>
    </template>
    <v-date-picker v-model="date" no-title color="primary" @input="menu = false"></v-date-picker>
  </v-menu>
</template>

<script>
export default {
  props: {}
};
</script>

<style lang="scss" module></style>
