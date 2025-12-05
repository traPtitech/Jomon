import { ApplicationDetail, User } from "./application";

export interface StatusLog {
  log_id: string;
  update_user: User;
  updated_at: string;
  content: {
    to_state: string;
    reason: string;
    created_at: string;
    update_user: User;
  };
}

export interface CommentLog {
  log_id: string;
  update_user: User;
  updated_at: string;
  content: {
    comment: string;
    comment_id: string;
    user: User;
    created_at: string;
    updated_at: string;
  };
}

export interface RefundLog {
  log_id: string;
  update_user: User;
  updated_at: string;
  content: {
    repaid_at: string;
    repaid_by_user: User;
    repaid_to_user: User;
    created_at: string;
  };
}

export interface ChangeLog {
  log_id: string;
  update_user: User;
  updated_at: string;
  content: {
    log: ApplicationDetail;
    pre_log: ApplicationDetail;
  };
}

export type Log = StatusLog | CommentLog | RefundLog | ChangeLog;
