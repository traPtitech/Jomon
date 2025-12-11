export const getStateLabel = (state: string): string => {
  switch (state) {
    case "submitted":
      return "承認待ち";
    case "fix_required":
      return "要修正";
    case "accepted":
      return "承認済み";
    case "fully_repaid":
      return "返済完了";
    case "rejected":
      return "却下";
    default:
      return "ERROR";
  }
};

export const getStateColor = (state: string): string => {
  switch (state) {
    case "submitted":
      return "state_submitted";
    case "rejected":
      return "state_rejected";
    case "fix_required":
      return "state_fix_required";
    case "accepted":
      return "state_accepted";
    case "fully_repaid":
      return "state_fully_repaid";
    default:
      return "white";
  }
};

export const getStateTextColor = (state: string): string => {
  switch (state) {
    case "fix_required":
      return "text-black";
    case "submitted":
    case "rejected":
    case "accepted":
    case "fully_repaid":
    default:
      return "text-white";
  }
};
