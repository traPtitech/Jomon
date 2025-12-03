import { useMeStore } from "@/stores/me";
import axios from "axios";
import { createPinia, setActivePinia } from "pinia";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("axios");

describe("Me Store", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("fetches me and updates state", async () => {
    const store = useMeStore();
    const mockData = { trap_id: "user1", is_admin: true };
    vi.mocked(axios.get).mockResolvedValue({ data: mockData });

    await store.fetchMe();

    expect(axios.get).toHaveBeenCalledWith("/api/users/me");
    expect(store.trapId).toBe("user1");
    expect(store.isAdmin).toBe(true);
  });
});
