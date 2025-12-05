export const remarksTitle = (type: string): string => {
  switch (type) {
    case "club":
      return "購入物の概要";
    case "contest":
      return "旅程";
    case "event":
      return "乗車区間";
    case "public":
      return "乗車区間";
    default:
      return "タイプが間違っています";
  }
};
export const applicationType = (type: string): string => {
  switch (type) {
    case "club":
      return "部費利用";
    case "contest":
      return "大会等旅費補助";
    case "event":
      return "イベント交通費補助";
    case "public":
      return "渉外交通費補助";
    default:
      return "error";
  }
};
