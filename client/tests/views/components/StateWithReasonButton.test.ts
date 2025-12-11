import StateWithReasonButton from "@/views/components/StateWithReasonButton.vue";
import SimpleButton from "@/views/shared/SimpleButton.vue"; // Import mocked component
import { createTestingPinia } from "@pinia/testing";
import { mount } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";

// Mock SimpleButton
vi.mock("@/views/shared/SimpleButton.vue", () => ({
  default: {
    name: "SimpleButton",
    template: '<button class="simple-button"><slot/></button>',
    props: ["label", "variant", "disabled"]
  }
}));

describe("StateWithReasonButton.vue", () => {
  it("renders correctly with info variant for 'submitted'", () => {
    const wrapper = mount(StateWithReasonButton, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })]
      },
      props: {
        toState: "submitted"
      }
    });

    const btn = wrapper.findComponent(SimpleButton);
    expect(btn.exists()).toBe(true);
    expect(btn.props("variant")).toBe("state_submitted");
    expect(btn.props("label")).toContain("承認待ちに戻す");
  });

  it("renders correctly with warning variant for 'fix_required'", () => {
    const wrapper = mount(StateWithReasonButton, {
      props: {
        toState: "fix_required"
      },
      global: {
        plugins: [createTestingPinia()]
      }
    });

    const btn = wrapper.findComponent(SimpleButton);
    expect(btn.props("variant")).toBe("state_fix_required");
    expect(btn.props("label")).toContain("要修正");
  });

  it("renders correctly with error variant for 'rejected'", () => {
    const wrapper = mount(StateWithReasonButton, {
      props: {
        toState: "rejected"
      },
      global: {
        plugins: [createTestingPinia()]
      }
    });

    const btn = wrapper.findComponent(SimpleButton);
    expect(btn.props("variant")).toBe("state_rejected");
    expect(btn.props("label")).toContain("却下");
  });
});
