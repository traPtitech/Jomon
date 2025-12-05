export interface User {
  trap_id: string;
  is_admin: boolean;
}

export interface ApplicationDetail {
  update_user: User;
  type: string;
  title: string;
  remarks: string;
  amount: number;
  paid_at: string;
  updated_at: string;
}

export interface ApplicationList {
  application_id: string;
  created_at: string;
  applicant: User;
  current_detail: ApplicationDetail;
  current_state: string;
}
