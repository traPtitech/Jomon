<template>
  <v-container>
    <v-row class="md5">
      <v-col cols="2" class="text-lg-right">
        <v-chip
          v-if="list.current_state !== ''"
          :width="10"
          :class="chip"
          label
          style="width:7em"
          >{{ list.current_state }}</v-chip
        >
      </v-col>
      <v-col cols="4">{{ list.current_detail.title }}</v-col>
      <v-col cols="4" class="text-lg-left">{{ list.applicant.trap_id }}</v-col>
      <v-col cols="2">{{ list.current_detail.amount }}</v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "Applicaton",
  props: {
    list: {
      application_id: String,
      created_at: String,
      applicant: {
        trap_id: String,
        is_admin: Boolean
      },
      current_detail: {
        update_user: {
          trap_id: String,
          is_admin: Boolean
        },
        type: String,
        title: String,
        remarks: String,
        amount: Number,
        paid_at: String,
        updated_at: String
      },
      current_state: String
    }
  },
  computed: {
    chip: function() {
      switch (this.list.current_state) {
        case "submitted":
          return "black white--text";
        case "rejected":
          return "yellow red--text";
        case "fix_required":
          return "orange white--text";
        case "accepted":
          return "red white--text";
        case "fully_repaid":
          return "gray black--text";
        default:
          return "white red--text";
      }
    }
  }
};
</script>
