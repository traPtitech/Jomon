<template>
  <div>
    <!-- <v-card>
      <v-row align="center">
        <v-form ref="form" v-model="valid" :lazy-validation="lazy">
          <v-text-field
            v-model="name"
            :counter="10"
            :rules="nameRules"
            label="Comment"
            required
          ></v-text-field>

          <v-btn
            :disabled="!valid"
            color="success"
            class="mr-4"
            @click="validate"
          >
            Validate
          </v-btn>
        </v-form>
      </v-row>
    </v-card> -->
    <v-row justify="space-between">
      <v-col cols="1">
        <Icon :user="this.$store.state.me.trap_id" :size="25" />
      </v-col>
      <v-col cols="11">
        <v-card ref="form" v-model="valid" :lazy-validation="lazy">
          <v-textarea
            v-model="name"
            :counter="10"
            :rules="nameRules"
            solo
            name="input-7-4"
            label="Leave a comment"
            required
          ></v-textarea>
          <v-btn
            :disabled="!valid"
            color="success"
            class="mr-4"
            @click="validate"
          >
            Comment
          </v-btn>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>
<script>
import Icon from "./Icon";
export default {
  components: {
    Icon
  },
  data: () => ({
    valid: true,
    name: "",
    nameRules: [
      v => !!v || "Name is required",
      v => (v && v.length <= 10) || "Name must be less than 10 characters"
    ]
  }),

  methods: {
    validate() {
      if (this.$refs.form.validate()) {
        this.snackbar = true;
      }
    }
  }
};
</script>
