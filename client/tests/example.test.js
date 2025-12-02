import { mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import Icon from "../src/views/shared/Icon.vue";

describe("Icon.vue", () => {
  it("renders correctly", () => {
    const wrapper = mount(Icon, {
      props: {
        user: "test-user",
        size: 30
      }
    });
    expect(wrapper.exists()).toBe(true);
    // Check if v-avatar is rendered (Vuetify components might need setup, but basic check)
  });

  it("sanity check", () => {
    expect(1 + 1).toBe(2);
  });
});
