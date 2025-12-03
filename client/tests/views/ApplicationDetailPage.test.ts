import { useApplicationDetailStore } from "@/stores/applicationDetail";
import ApplicationDetailPage from "@/views/ApplicationDetailPage.vue";
import { createTestingPinia } from "@pinia/testing";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";

// Mock child components
vi.mock("@/views/components/ApplicationDetail.vue", () => ({
  default: { template: '<div data-testid="application-paper"></div>' }
}));
vi.mock("@/views/components/ApplicationDetailLogs.vue", () => ({
  default: { template: '<div data-testid="application-logs"></div>' }
}));
vi.mock("@/views/components/FixApplicationDetail.vue", () => ({
  default: { template: '<div data-testid="fix-application-paper"></div>' }
}));

// Mock useRoute
vi.mock("vue-router", () => ({
  useRoute: () => ({
    params: { id: "123" }
  })
}));

describe("ApplicationDetailPage.vue", () => {
  it("renders loading initially", () => {
    const wrapper = mount(ApplicationDetailPage, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })]
      }
    });
    expect(wrapper.text()).toContain("loading...");
  });

  it("fetches data on mount and renders content", async () => {
    const wrapper = mount(ApplicationDetailPage, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })]
      }
    });

    const store = useApplicationDetailStore();

    await flushPromises();

    expect(store.fetchApplicationDetail).toHaveBeenCalledWith("123");
    expect(store.deleteFix).toHaveBeenCalled();
    expect(wrapper.text()).not.toContain("loading...");
    expect(wrapper.find('[data-testid="application-logs"]').exists()).toBe(
      true
    );
  });

  it("renders ApplicationPaper when fix is false", async () => {
    const wrapper = mount(ApplicationDetailPage, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              applicationDetail: { fix: false }
            }
          })
        ]
      }
    });
    await flushPromises();

    expect(wrapper.find('[data-testid="application-paper"]').exists()).toBe(
      true
    );
    expect(wrapper.find('[data-testid="fix-application-paper"]').exists()).toBe(
      false
    );
  });

  it("renders FixApplicationPaper when fix is true", async () => {
    const wrapper = mount(ApplicationDetailPage, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              applicationDetail: { fix: true }
            }
          })
        ]
      }
    });
    await flushPromises();

    expect(wrapper.find('[data-testid="application-paper"]').exists()).toBe(
      false
    );
    expect(wrapper.find('[data-testid="fix-application-paper"]').exists()).toBe(
      true
    );
  });
});
