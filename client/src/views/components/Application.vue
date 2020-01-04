<template>
  <v-container>
    <v-row class="md5">
      <v-col cols="4">
        {{ list.title }}
      </v-col>
      <v-col cols="3">
        {{ list.name }}
      </v-col>
      <v-col cols="2">
        {{ list.money }}
      </v-col>
      <v-col cols="3">
        {{ list.state }}
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "Applicaton",
  props: {
    list: {
      id: Number,
      title: String,
      name: String,
      money: Number,
      state: String
    }
  }
};
</script>
