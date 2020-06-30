export const numberFormat = price => {
  return price.toString().replace(/(\d)(?=(\d\d\d)+(?!\d))/g, "$1,");
};
export const dayPrint = time => {
  let d = new Date(time);
  let year = d.getFullYear();
  let month = d.getMonth() + 1;
  let day = d.getDate();
  let res = year + "年" + month + "月" + day + "日";
  return res;
};
