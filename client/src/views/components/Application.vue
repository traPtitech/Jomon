<template>
  <v-row>
    <v-col cols="2" class="text-lg-right">
      <div>
        <state-chip
          v-if="list.current_state !== ''"
          :state="list.current_state"
        ></state-chip>
      </div>
    </v-col>

    <v-col cols="4">{{ list.current_detail.title }}</v-col>

    <v-col cols="4" class="text-lg-left">{{ list.applicant.trap_id }}</v-col>

    <v-col cols="2">{{ list.current_detail.amount }}</v-col>
  </v-row>
</template>

<script>
import StateChip from "@/views/shared/StateChip";

export default {
  name: "ApplicationItem",
  components: {
    StateChip
  },
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
  }
};
</script>
