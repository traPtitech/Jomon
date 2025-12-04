import RepaidButton from "@/views/components/RepaidButton.vue";
import SimpleButton from "@/views/shared/SimpleButton.vue";
import { createTestingPinia } from "@pinia/testing";
import { mount } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";

// Mock SimpleButton
vi.mock("@/views/shared/SimpleButton.vue", () => ({
  default: {
    name: "SimpleButton",
    template:
      '<button class="simple-button" @click="$emit(\'click\')"><slot/></button>',
    props: ["label", "disabled", "variant"]
  }
}));

// Mock v-dialog to render content directly
const VDialog = {
  template: '<div><slot name="activator" :props="{}"></slot><slot></slot></div>'
};

describe("RepaidButton.vue", () => {
  it("renders button correctly", () => {
    const wrapper = mount(RepaidButton, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          VDialog: VDialog
        }
      }
    });

    const btn = wrapper.findComponent(SimpleButton);
    expect(btn.exists()).toBe(true);
    expect(btn.props("label")).toBe("払い戻し済みのユーザーを選択");
    expect(btn.props("variant")).toBe("done");
  });

  it("disables OK button when no user is selected", async () => {
    const wrapper = mount(RepaidButton, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              applicationDetail: {
                core: {
                  repayment_logs: [
                    { repaid_to_user: { trap_id: "user1" }, repaid_at: null }
                  ]
                }
              }
            }
          })
        ],
        stubs: {
          VDialog: VDialog
        }
      }
    });

    const buttons = wrapper.findAllComponents(SimpleButton);
    const okBtn = buttons.find(b => b.props("label") === "OK");

    expect(okBtn).toBeDefined();
    expect(okBtn?.props("disabled")).toBe(true);

    // Select user
    const select = wrapper.findComponent({ name: "VSelect" });
    await select.setValue(["user1"]);

    expect(okBtn?.props("disabled")).toBe(false);
  });
});
