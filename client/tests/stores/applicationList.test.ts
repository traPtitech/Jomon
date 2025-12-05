import { useApplicationListStore } from "@/stores/applicationList";
import axios from "axios";
import { createPinia, setActivePinia } from "pinia";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("axios");

describe("ApplicationList Store", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("fetches application list and updates state", async () => {
    const store = useApplicationListStore();
    const mockData = [{ id: 1, title: "Test App" }];
    vi.mocked(axios.get).mockResolvedValue({ data: mockData });

    await store.fetchApplicationList({});

    expect(axios.get).toHaveBeenCalledWith("/api/applications", { params: {} });
    expect(store.applicationList).toEqual(mockData);
  });

  it("cleans parameters before making API call", async () => {
    const store = useApplicationListStore();
    vi.mocked(axios.get).mockResolvedValue({ data: [] });

    const params = {
      valid: "value",
      empty: "",
      nullValue: null,
      undefinedValue: undefined
    };

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    await store.fetchApplicationList(params as any);

    expect(axios.get).toHaveBeenCalledWith("/api/applications", {
      params: { valid: "value" }
    });
  });
});
