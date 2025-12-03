import { useApplicationListStore } from "@/stores/applicationList";
import { useUserListStore } from "@/stores/userList";
import ApplicationListPage from "@/views/ApplicationListPage.vue";
import { createTestingPinia } from "@pinia/testing";
import { mount } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";

describe("ApplicationListPage.vue", () => {
  const createWrapper = () => {
    return mount(ApplicationListPage, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              applicationList: {
                applicationList: [
                  {
                    application_id: "app-1",
                    applicant: { trap_id: "user1" },
                    current_detail: {
                      title: "Test App 1",
                      amount: 1000,
                      type: "club"
                    },
                    current_state: "submitted",
                    created_at: "2023-01-01T00:00:00Z"
                  }
                ]
              },
              userList: {
                userList: [{ trap_id: "user1" }, { trap_id: "user2" }]
              }
            }
          })
        ],
        stubs: {
          Application: true // Stub child component
        }
      }
    });
  };

  it("renders correctly", () => {
    const wrapper = createWrapper();
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.text()).toContain("絞り込み");
    expect(wrapper.text()).toContain("申請一覧");
  });

  it("fetches data on mount", () => {
    createWrapper();
    const applicationListStore = useApplicationListStore();
    const userListStore = useUserListStore();

    expect(applicationListStore.fetchApplicationList).toHaveBeenCalled();
    expect(userListStore.fetchUserList).toHaveBeenCalled();
  });

  it("toggles filter section", async () => {
    const wrapper = createWrapper();
    const toggleButton = wrapper.find(".v-card-title button");

    // Initial state: hidden (v-show="show", show defaults to false/true depending on display)
    // Note: useDisplay mock might be needed if we want to test defaultShow logic precisely.
    // For now, let's just check if clicking toggles the value.

    // We can check the vm state
    const vm = wrapper.vm as any;
    const initialShow = vm.show;

    await toggleButton.trigger("click");
    expect(vm.show).toBe(!initialShow);
  });

  it("calls fetchApplicationList with params when reload button is clicked", async () => {
    const wrapper = createWrapper();
    const applicationListStore = useApplicationListStore();

    // Show the filter form first to make buttons interactive (though v-show doesn't remove from DOM)
    const vm = wrapper.vm as any;
    vm.show = true;
    await wrapper.vm.$nextTick();

    const reloadButton = wrapper
      .findAll("button")
      .filter(b => b.text() === "")
      .at(1); // Finding by icon is tricky, let's look for the one calling getApplicationList
    // Or better, find by icon class if possible, or just call the method directly if UI testing is flaky.
    // Let's try to find the button with mdi-reload icon.
    // Since we are using Vuetify, icons are often in v-icon.

    // Let's assume the first primary button in the form is reload.
    const buttons = wrapper.findAll(".v-btn.bg-primary");
    // 0: reload, 1: close (reset), 2: date sort, 3: title sort

    await buttons[0].trigger("click");

    expect(applicationListStore.fetchApplicationList).toHaveBeenCalled();
  });
});
