import StateChip from "@/views/shared/StateChip.vue";
import { mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";

describe("StateChip.vue", () => {
  const mountStateChip = (state: string) => {
    return mount(StateChip, {
      global: {
        plugins: []
      },
      props: {
        state
      }
    });
  };

  it("renders 'submitted' with info color and black text", () => {
    const wrapper = mountStateChip("submitted");
    expect(wrapper.text()).toBe("承認待ち");
    const chip = wrapper.find(".v-chip");
    // Vuetify 3 classes for colors can be bg-<color> or text-<color>
    expect(chip.classes()).toContain("bg-info");
    expect(chip.classes()).toContain("text-white"); // Correction: We changed submitted to info (blue) + white text in previous steps
  });

  it("renders 'fix_required' with warning color and black text", () => {
    const wrapper = mountStateChip("fix_required");
    expect(wrapper.text()).toBe("要修正");
    const chip = wrapper.find(".v-chip");
    expect(chip.classes()).toContain("bg-warning");
    expect(chip.classes()).toContain("text-black");
  });

  it("renders 'accepted' with success color", () => {
    const wrapper = mountStateChip("accepted");
    expect(wrapper.text()).toBe("承認済み");
    const chip = wrapper.find(".v-chip");
    expect(chip.classes()).toContain("bg-success");
    expect(chip.classes()).toContain("text-white");
  });

  it("renders 'fully_repaid' with done color", () => {
    const wrapper = mountStateChip("fully_repaid");
    expect(wrapper.text()).toBe("返済完了");
    const chip = wrapper.find(".v-chip");
    // 'done' is a custom theme color, so it might appear as text-done or bg-done if defined in theme,
    // but in unit test without full theme setup, it might just apply the class.
    // Vuetify applies 'bg-<color>' for variant="flat"
    expect(chip.classes()).toContain("bg-done");
  });

  it("renders 'rejected' with grey color", () => {
    const wrapper = mountStateChip("rejected");
    expect(wrapper.text()).toBe("却下");
    const chip = wrapper.find(".v-chip");
    expect(chip.classes()).toContain("bg-grey");
  });
});
