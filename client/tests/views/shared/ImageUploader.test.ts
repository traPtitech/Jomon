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

    const originalFileReader = global.FileReader;
    vi.stubGlobal("FileReader", MockFileReader);

    const fileInput = wrapper.findComponent(components.VFileInput);
    expect(fileInput.exists()).toBe(true);

    await fileInput.vm.$emit("update:modelValue", [file]);

    await new Promise(resolve => setTimeout(resolve, 0));

    expect(wrapper.emitted("update:modelValue")).toBeTruthy();
    expect(wrapper.emitted("update:modelValue")![0][0]).toEqual([file]);

    vi.stubGlobal("FileReader", originalFileReader);
  });

  it("appends new images instead of replacing them", async () => {
    const wrapper = mount(ImageUploader);

    const file1 = new File(["content1"], "test1.png", { type: "image/png" });
    const file2 = new File(["content2"], "test2.png", { type: "image/png" });

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

    const originalFileReader = global.FileReader;
    vi.stubGlobal("FileReader", MockFileReader);

    const fileInput = wrapper.findComponent(components.VFileInput);

    // Simulate first upload
    await fileInput.vm.$emit("update:modelValue", [file1]);

    let emitted = wrapper.emitted("update:modelValue");
    expect(emitted).toBeTruthy();
    expect(emitted!.slice(-1)[0][0]).toEqual([file1]);

    // Simulate second upload
    await fileInput.vm.$emit("update:modelValue", [file2]);

    emitted = wrapper.emitted("update:modelValue");
    const lastEmit = emitted!.slice(-1)[0][0] as File[];

    expect(lastEmit).toHaveLength(2);
    expect(lastEmit).toContain(file1);
    expect(lastEmit).toContain(file2);

    vi.stubGlobal("FileReader", originalFileReader);
  });
});
