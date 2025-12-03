import axios from "axios";
import { defineStore } from "pinia";
import { computed, reactive, ref } from "vue";

interface User {
  trap_id: string;
  is_admin: boolean;
}

interface Comment {
  comment_id: number;
  user: User;
  comment: string;
  created_at: string;
  updated_at: string;
}

interface StateLog {
  update_user: User;
  to_state: string;
  reason: string;
  created_at: string;
}

interface ApplicationDetailLog {
  update_user: User;
  type: string;
  title: string;
  remarks: string;
  amount: number;
  paid_at: string;
  updated_at: string;
}

interface RepaymentLog {
  repaid_by_user: User;
  repaid_to_user: User;
  created_at: string;
  repaid_at: string;
}

interface ApplicationCore {
  application_id: string;
  created_at: string;
  applicant: User;
  current_detail: {
    update_user: {
      trap_id: string;
      is_admin: string; // Note: In original code it was string, but likely boolean? Keeping as string to match
    };
    type: string;
    title: string;
    remarks: string;
    amount: number;
    paid_at: string;
    updated_at: string;
    repaid_to_id: string[];
  };
  current_state: string;
  images: string[];
  comments: Comment[];
  state_logs: StateLog[];
  application_detail_logs: ApplicationDetailLog[];
  repayment_logs: RepaymentLog[];
}

export const useApplicationDetailStore = defineStore(
  "applicationDetail",
  () => {
    const fix = ref(false);
    const core = reactive<ApplicationCore>({
      application_id: "",
      created_at: "",
      applicant: {
        trap_id: "",
        is_admin: false
      },
      current_detail: {
        update_user: {
          trap_id: "",
          is_admin: ""
        },
        type: "",
        title: "",
        remarks: "",
        amount: 0,
        paid_at: "",
        updated_at: "",
        repaid_to_id: []
      },
      current_state: "",
      images: [""],
      comments: [],
      state_logs: [],
      application_detail_logs: [],
      repayment_logs: []
    });

    interface Log {
      log_type: string;
      content:
        | Comment
        | StateLog
        | { log: ApplicationDetailLog; pre_log: ApplicationDetailLog }
        | RepaymentLog;
      sort_date: Date;
    }

    const logs = computed(() => {
      const logs: Log[] = [];
      core.comments.forEach(log => {
        logs.push({
          log_type: "comment",
          content: log,
          sort_date: new Date(log.created_at)
        });
      });
      core.state_logs.forEach(log => {
        logs.push({
          log_type: "state",
          content: log,
          sort_date: new Date(log.created_at)
        });
      });
      let pre_record: ApplicationDetailLog | "" = "";
      core.application_detail_logs.forEach(log => {
        if (pre_record !== "") {
          logs.push({
            log_type: "application",
            content: {
              log: log,
              pre_log: pre_record
            },
            sort_date: new Date(log.updated_at)
          });
        }
        pre_record = log;
      });
      core.repayment_logs.forEach(log => {
        logs.push({
          log_type: "repayment",
          content: log,
          sort_date: new Date(log.created_at)
        });
      });
      logs.sort((a, b) => {
        if (a.sort_date > b.sort_date) {
          return 1;
        } else {
          return -1;
        }
      });
      return logs;
    });

    const fetchApplicationDetail = async (applicationId: string) => {
      const response = await axios.get("/api/applications/" + applicationId);
      Object.assign(core, response.data);
    };

    const changeFix = () => {
      fix.value = !fix.value;
    };

    const deleteFix = () => {
      fix.value = false;
    };

    return {
      fix,
      core,
      logs,
      fetchApplicationDetail,
      changeFix,
      deleteFix
    };
  }
);
