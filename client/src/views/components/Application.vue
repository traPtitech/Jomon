<template>
  <v-container>
    <v-row class="md5">
      <v-col cols="2" class="text-lg-right">
        <div>
          <StateChip
            v-if="list.current_state !== ''"
            :width="10"
            :state="list.current_state"
            style="width: 7em"
          ></StateChip>
        </div>
      </v-col>
      <v-col cols="4">{{ list.current_detail.title }}</v-col>
      <v-col cols="4" class="text-lg-left">{{ list.applicant.trap_id }}</v-col>
      <v-col cols="2">{{ list.current_detail.amount }}</v-col>
    </v-row>
  </v-container>
</template>

<script>
import StateChip from "../shared/StateChip";

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
  components: {
    StateChip
  }
};
</script>
