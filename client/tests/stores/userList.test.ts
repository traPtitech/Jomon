import { useUserListStore } from "@/stores/userList";
import axios from "axios";
import { createPinia, setActivePinia } from "pinia";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("axios");

describe("UserList Store", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("fetches user list and updates state", async () => {
    const store = useUserListStore();
    const mockData = [
      { trap_id: "user1", is_admin: true },
      { trap_id: "user2", is_admin: false }
    ];
    vi.mocked(axios.get).mockResolvedValue({ data: mockData });

    await store.fetchUserList();

    expect(axios.get).toHaveBeenCalledWith("/api/users");
    expect(store.userList).toEqual(mockData);
  });

  it("computes derived lists correctly", () => {
    const store = useUserListStore();
    store.userList = [
      { trap_id: "user1", is_admin: true },
      { trap_id: "user2", is_admin: false },
      { trap_id: "user3", is_admin: true }
    ];

    expect(store.trapIds).toEqual(["user1", "user2", "user3"]);
    expect(store.adminList).toEqual(["user1", "user3"]);
    expect(store.notAdminList).toEqual(["user2"]);
  });
});
