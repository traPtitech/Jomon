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
    expect(chip.classes()).toContain("bg-state_submitted");
  });

  it("renders 'fix_required' with warning color and black text", () => {
    const wrapper = mount(StateChip, {
      props: {
        state: "fix_required"
      }
    });

    expect(wrapper.text()).toBe("要修正");
    const chip = wrapper.find(".v-chip");
    expect(chip.classes()).toContain("bg-state_fix_required");
    expect(chip.classes()).toContain("text-black");
  });

  it("renders 'accepted' with success color", () => {
    const wrapper = mount(StateChip, {
      props: {
        state: "accepted"
      }
    });

    expect(wrapper.text()).toBe("承認済み");
    const chip = wrapper.find(".v-chip");
    expect(chip.classes()).toContain("bg-state_accepted");
  });

  it("renders 'fully_repaid' with done color", () => {
    const wrapper = mount(StateChip, {
      props: {
        state: "fully_repaid"
      }
    });

    expect(wrapper.text()).toBe("返済完了");
    const chip = wrapper.find(".v-chip");
    expect(chip.classes()).toContain("bg-state_fully_repaid");
  });

  it("renders 'rejected' with grey color", () => {
    const wrapper = mount(StateChip, {
      props: {
        state: "rejected"
      }
    });

    expect(wrapper.text()).toBe("却下");
    const chip = wrapper.find(".v-chip");
    expect(chip.classes()).toContain("bg-state_rejected");
  });
});
