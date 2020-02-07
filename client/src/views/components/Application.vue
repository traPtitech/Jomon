<template>
  <v-container>
    <v-row class="md5">
      <v-col cols="2" class="text-lg-right">
        <v-chip
          v-if="list.state !== ''"
          :width="10"
          :class="chip"
          label
          style="width:7em"
          >{{ list.state }}</v-chip
        >
      </v-col>
      <v-col cols="4">{{ list.title }}</v-col>
      <v-col cols="4" class="text-lg-left">{{ list.name }}</v-col>
      <v-col cols="2">{{ list.money }}</v-col>
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
