import { ApplicationDetail, User } from "@/types/application";
import { CommentLog, RefundLog, StatusLog } from "@/types/log";
import axios from "axios";
import { defineStore } from "pinia";
import { computed, reactive, ref } from "vue";

// Local interfaces for API response structure if they differ from shared types
// Or reuse shared types if they match.

// API response for comments seems to match CommentLog content
type Comment = CommentLog["content"];

// API response for state logs
// StatusLog in types/log.ts has content: { to_state, reason, created_at, update_user }
// Store StateLog has: { update_user, to_state, reason, created_at }
// They match.
type StateLog = StatusLog["content"];

// ApplicationDetailLog in store matches ApplicationDetail in types/application.ts
type ApplicationDetailLog = ApplicationDetail;

// RepaymentLog in store matches RefundLog content
type RepaymentLog = RefundLog["content"];

interface ApplicationCore {
  application_id: string;
  created_at: string;
  applicant: User;
  current_detail: ApplicationDetail & {
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
          is_admin: false
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
      images: [],
      comments: [],
      state_logs: [],
      application_detail_logs: [],
      repayment_logs: []
    });

    type Log =
      | {
          log_type: "comment";
          content: Comment;
          sort_date: Date;
        }
      | {
          log_type: "state";
          content: StateLog;
          sort_date: Date;
        }
      | {
          log_type: "application";
          content: { log: ApplicationDetailLog; pre_log: ApplicationDetailLog };
          sort_date: Date;
        }
      | {
          log_type: "repayment";
          content: RepaymentLog;
          sort_date: Date;
        };

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

    const updateApplicationState = async (
      to_state: string,
      reason: string = ""
    ) => {
      await axios.put("/api/applications/" + core.application_id + "/states", {
        to_state: to_state,
        reason: reason
      });
      await fetchApplicationDetail(core.application_id);
    };

    return {
      fix,
      core,
      logs,
      fetchApplicationDetail,
      changeFix,
      deleteFix,
      updateApplicationState
    };
  }
);
