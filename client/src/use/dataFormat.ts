export const numberFormat = (price: number | string): string => {
  return price.toString().replace(/(\d)(?=(\d\d\d)+(?!\d))/g, "$1,");
};
export const dayPrint = (time: string | number | Date): string => {
  const d = new Date(time);
  const year = d.getFullYear();
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const res = year + "年" + month + "月" + day + "日";
  return res;
};
