import { describe, expect, it } from "vitest";
import { dayPrint, numberFormat } from "../../src/use/dataFormat";

describe("dataFormat", () => {
  describe("numberFormat", () => {
    it("formats numbers with commas", () => {
      expect(numberFormat(1000)).toBe("1,000");
      expect(numberFormat(1000000)).toBe("1,000,000");
      expect(numberFormat(123456789)).toBe("123,456,789");
    });

    it("formats strings with commas", () => {
      expect(numberFormat("1000")).toBe("1,000");
      expect(numberFormat("1000000")).toBe("1,000,000");
    });

    it("does not add commas to small numbers", () => {
      expect(numberFormat(100)).toBe("100");
      expect(numberFormat(0)).toBe("0");
      expect(numberFormat("999")).toBe("999");
    });
  });

  describe("dayPrint", () => {
    it("formats Date object correctly", () => {
      const date = new Date("2023-01-01T00:00:00");
      expect(dayPrint(date)).toBe("2023年1月1日");
    });

    it("formats date string correctly", () => {
      expect(dayPrint("2023-12-31")).toBe("2023年12月31日");
    });

    it("formats timestamp correctly", () => {
      const timestamp = new Date("2023-05-05T00:00:00").getTime();
      expect(dayPrint(timestamp)).toBe("2023年5月5日");
    });
  });
});
