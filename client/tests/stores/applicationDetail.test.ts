import { useApplicationDetailStore } from "@/stores/applicationDetail";
import axios from "axios";
import { createPinia, setActivePinia } from "pinia";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("axios");

describe("ApplicationDetail Store", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("fetches application detail and updates state", async () => {
    const store = useApplicationDetailStore();
    const mockData = {
      application_id: "1",
      title: "Test App",
      applicant: { trap_id: "user1", is_admin: false }
    };
    vi.mocked(axios.get).mockResolvedValue({ data: mockData });

    await store.fetchApplicationDetail("1");

    expect(axios.get).toHaveBeenCalledWith("/api/applications/1");
    expect(store.core.application_id).toBe("1");
    expect(store.core.title).toBe("Test App");
  });

  it("toggles fix state", () => {
    const store = useApplicationDetailStore();
    expect(store.fix).toBe(false);

    store.changeFix();
    expect(store.fix).toBe(true);

    store.changeFix();
    expect(store.fix).toBe(false);

    store.fix = true;
    store.deleteFix();
    expect(store.fix).toBe(false);
  });

  it("computes logs correctly", () => {
    const store = useApplicationDetailStore();
    const mockCore = {
      comments: [
        {
          comment_id: 1,
          user: { trap_id: "user1", is_admin: false },
          comment: "test comment",
          created_at: "2023-01-01T10:00:00Z",
          updated_at: "2023-01-01T10:00:00Z"
        }
      ],
      state_logs: [
        {
          update_user: { trap_id: "user2", is_admin: true },
          to_state: "approved",
          reason: "ok",
          created_at: "2023-01-02T10:00:00Z"
        }
      ],
      application_detail_logs: [
        {
          update_user: { trap_id: "user1", is_admin: false },
          type: "club",
          title: "old title",
          remarks: "",
          amount: 100,
          paid_at: "2023-01-01",
          updated_at: "2023-01-01T09:00:00Z"
        },
        {
          update_user: { trap_id: "user1", is_admin: false },
          type: "club",
          title: "new title",
          remarks: "",
          amount: 100,
          paid_at: "2023-01-01",
          updated_at: "2023-01-03T10:00:00Z"
        }
      ],
      repayment_logs: []
    };

    // Manually populate state for testing computed
    Object.assign(store.core, mockCore);

    const logs = store.logs;
    expect(logs).toHaveLength(3); // 1 comment, 1 state log, 1 application log (diff)

    // Sort order check (by date)
    expect(logs[0].log_type).toBe("comment"); // Jan 1 10:00
    expect(logs[1].log_type).toBe("state"); // Jan 2 10:00
    expect(logs[2].log_type).toBe("application"); // Jan 3 10:00 (second log diff)
  });
});
