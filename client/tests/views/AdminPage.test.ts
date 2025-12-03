import { useUserListStore } from "@/stores/userList";
import AdminPage from "@/views/AdminPage.vue";
import { createTestingPinia } from "@pinia/testing";
import { mount } from "@vue/test-utils";
import axios from "axios";
import { describe, expect, it, vi } from "vitest";

// Mock axios
vi.mock("axios", () => ({
  default: {
    get: vi.fn(),
    put: vi.fn().mockResolvedValue({})
  }
}));

describe("AdminPage.vue", () => {
  const createWrapper = (isAdmin = true) => {
    return mount(AdminPage, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              me: {
                trapId: "admin-user",
                isAdmin: isAdmin
              },
              userList: {
                userList: [
                  { trap_id: "admin-user", is_admin: true },
                  { trap_id: "normal-user", is_admin: false }
                ]
              }
            }
          })
        ],
        stubs: {
          SimpleButton: {
            template: "<button @click=\"$emit('click')\">{{ label }}</button>",
            props: ["label"]
          }
        }
      }
    });
  };

  it("renders permission denied message when not admin", () => {
    const wrapper = createWrapper(false);
    expect(wrapper.text()).toContain("権限がありません");
  });

  it("renders admin controls when admin", () => {
    const wrapper = createWrapper(true);
    expect(wrapper.text()).not.toContain("権限がありません");
    expect(wrapper.text()).toContain("admin-user"); // Chip
    expect(wrapper.findAll("button").length).toBe(2); // Add and Remove buttons
  });

  it("calls addAdmin API when add button is clicked", async () => {
    const wrapper = createWrapper(true);
    const userListStore = useUserListStore();

    // Select user to add
    // Since we are using VAutocomplete, interacting with it directly in unit test is hard.
    // We can set the model value directly.
    const vm = wrapper.vm as any;
    vm.addAdminUsers = ["normal-user"];

    // Find the second button (Add Admin)
    const buttons = wrapper.findAll("button");
    await buttons[1].trigger("click");

    expect(axios.put).toHaveBeenCalledWith("api/users/admins", {
      trap_id: "normal-user",
      to_admin: true
    });
    expect(userListStore.fetchUserList).toHaveBeenCalled();
  });

  it("calls removeAdmin API when remove button is clicked", async () => {
    const wrapper = createWrapper(true);
    const userListStore = useUserListStore();

    // Select user to remove
    const vm = wrapper.vm as any;
    vm.removeAdminUsers = ["admin-user"];

    // Find the first button (Remove Admin)
    const buttons = wrapper.findAll("button");
    await buttons[0].trigger("click");

    expect(axios.put).toHaveBeenCalledWith("api/users/admins", {
      trap_id: "admin-user",
      to_admin: false
    });
    expect(userListStore.fetchUserList).toHaveBeenCalled();
  });
});
