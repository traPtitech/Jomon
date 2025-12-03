import { createTestingPinia } from "@pinia/testing";
import { mount } from "@vue/test-utils";
import { describe, expect, it, vi, type Mock } from "vitest";
import StateWithReasonButton from "../../../src/views/components/StateWithReasonButton.vue";

import axios from "axios";

// Mock axios
vi.mock("axios");

describe("StateWithReasonButton.vue", () => {
  it("handles input correctly without [object Object] bug and stack overflow", async () => {
    // Mock alert
    global.alert = vi.fn();

    // Mock axios response
    (axios.put as Mock).mockResolvedValue({});

    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        applicationDetail: {
          core: {
            application_id: "test-app-id",
            current_state: "accepted"
          }
        }
      }
    });

    const wrapper = mount(StateWithReasonButton, {
      global: {
        plugins: [pinia]
      },
      props: {
        toState: "submitted"
      }
    });

    // Open dialog
    const activator = wrapper.findComponent({ name: "SimpleButton" });
    await activator.trigger("click");

    // Wait for dialog to open and nextTick in watch
    await wrapper.vm.$nextTick();
    await new Promise(resolve => setTimeout(resolve, 100));

    // Check input value
    const input = document.querySelector("input");
    expect(input).not.toBeNull();
    if (input) {
      expect(input.value).toBe("");

      // Set reason via component emit
      const textField = wrapper.findComponent({ name: "VTextField" });
      await textField.vm.$emit("update:modelValue", "test reason");

      // Trigger validation
      await textField.trigger("blur");
      await wrapper.vm.$nextTick();
    }

    // Click submit
    const submitBtn = wrapper.findAllComponents({ name: "SimpleButton" })[2]; // 0: activator, 1: back, 2: submit

    // Check if enabled

    // Should not throw error
    await submitBtn.trigger("click");

    // Verify axios call
    expect(axios.put).toHaveBeenCalledWith(
      "../api/applications/test-app-id/states",
      expect.objectContaining({
        to_state: "submitted",
        reason: "test reason"
      })
    );
  });
});
