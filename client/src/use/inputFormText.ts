export const titlePlaceholder = (type: string): string => {
  switch (type) {
    case "club":
      return "工大祭用ポスターの印刷代";
    case "contest":
      return "ISUCON in 沖縄";
    case "event":
      return "コミックマーケット98";
    case "public":
      return "OO株式会社様との打ち合わせ";
    default:
      return "タイプが間違っています";
  }
};
export const remarksPlaceholder = (type: string): string => {
  switch (type) {
    case "club":
      return "OO印刷所にA3サイズのポスターをXX部";
    case "contest":
      return "大岡山から沖縄まで片道XX円の往復\n沖縄でOO民宿に二泊三日でxxx円";
    case "event":
    case "public":
      return "東急大井町線急行\n大岡山 から 大井町:160円\n\n東京臨海高速鉄道臨海線\n大井町 から 東京テレポート:280円\n\n各区間往復一人分";
    default:
      return "タイプが間違っています";
  }
};
export const remarksHint = (type: string): string => {
  switch (type) {
    case "club":
      return "具体的購入物、用途等を記入";
    case "contest":
      return "経路、参加費、宿場等を記入";
    case "event":
    case "public":
      return "駅名、駅間の料金、往復or片道、複数人の時は誰がどの区間か、記入";
    default:
      return "タイプが間違っています";
  }
};
