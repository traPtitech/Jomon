import ImageUploader from "@/views/shared/ImageUploader.vue";
import { mount } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";
import * as components from "vuetify/components";
describe("ImageUploader.vue", () => {
  it("emits update:modelValue when images are selected", async () => {
    const wrapper = mount(ImageUploader);

    const file = new File(["content"], "test.png", { type: "image/png" });

    // Mock FileReader
    class MockFileReader {
      result = "data:image/png;base64,test";
      readAsDataURL = vi.fn();
      addEventListener = vi.fn((event, callback) => {
        if (event === "load") {
          callback();
        }
      });
    }

    // Stub global FileReader
    const originalFileReader = global.FileReader;
    vi.stubGlobal("FileReader", MockFileReader);

    // Trigger image change
    // We can't easily trigger v-file-input change event directly in unit test without full mount
    // But we can call the handler if we could access it.
    // Since it's <script setup>, we can't access `imageChange` directly.
    // We will simulate the event on the component if possible, or verify structure.

    // Ideally we should find the v-file-input and trigger update:modelValue
    const fileInput = wrapper.findComponent(components.VFileInput);
    expect(fileInput.exists()).toBe(true);

    // Trigger the event that VFileInput would emit
    await fileInput.vm.$emit("update:modelValue", [file]);

    // Wait for any async operations
    await new Promise(resolve => setTimeout(resolve, 0));

    expect(wrapper.emitted("update:modelValue")).toBeTruthy();
    expect(wrapper.emitted("update:modelValue")![0][0]).toEqual([file]);

    // Restore
    vi.stubGlobal("FileReader", originalFileReader);
  });
});
