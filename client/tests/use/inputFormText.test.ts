import { describe, expect, it } from "vitest";
import {
  remarksHint,
  remarksPlaceholder,
  titlePlaceholder
} from "../../src/use/inputFormText";

describe("inputFormText", () => {
  describe("titlePlaceholder", () => {
    it("returns correct placeholder for club", () => {
      expect(titlePlaceholder("club")).toBe("工大祭用ポスターの印刷代");
    });
    it("returns correct placeholder for contest", () => {
      expect(titlePlaceholder("contest")).toBe("ISUCON in 沖縄");
    });
    it("returns error message for invalid type", () => {
      expect(titlePlaceholder("invalid")).toBe("タイプが間違っています");
    });
  });

  describe("remarksPlaceholder", () => {
    it("returns correct placeholder for club", () => {
      expect(remarksPlaceholder("club")).toBe(
        "OO印刷所にA3サイズのポスターをXX部"
      );
    });
    it("returns correct placeholder for event", () => {
      expect(remarksPlaceholder("event")).toContain("東急大井町線急行");
    });
    it("returns error message for invalid type", () => {
      expect(remarksPlaceholder("invalid")).toBe("タイプが間違っています");
    });
  });

  describe("remarksHint", () => {
    it("returns correct hint for club", () => {
      expect(remarksHint("club")).toBe("具体的購入物、用途等を記入");
    });
    it("returns correct hint for contest", () => {
      expect(remarksHint("contest")).toBe("経路、参加費、宿場等を記入");
    });
    it("returns error message for invalid type", () => {
      expect(remarksHint("invalid")).toBe("タイプが間違っています");
    });
  });
});
