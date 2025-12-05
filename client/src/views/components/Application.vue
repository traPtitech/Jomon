<template>
  <v-row>
    <v-col cols="2" class="text-lg-right">
      <div>
        <state-chip
          v-if="list.current_state !== ''"
          :state="list.current_state"
        />
      </div>
    </v-col>

    <v-col cols="4">
      {{ list.current_detail.title }}
    </v-col>

    <v-col cols="4" class="text-lg-left">
      {{ list.applicant.trap_id }}
    </v-col>

    <v-col cols="2">
      {{ list.current_detail.amount }}
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ApplicationList } from "@/types/application";
import StateChip from "@/views/shared/StateChip.vue";

withDefaults(
  defineProps<{
    list: ApplicationList;
  }>(),
  {
    list: () => ({
      application_id: "",
      created_at: "",
      applicant: {
        trap_id: "",
        is_admin: false
      },
      current_detail: {
        update_user: {
          trap_id: "",
          is_admin: false
        },
        type: "",
        title: "",
        remarks: "",
        amount: 0,
        paid_at: "",
        updated_at: ""
      },
      current_state: ""
    })
  }
);
</script>
