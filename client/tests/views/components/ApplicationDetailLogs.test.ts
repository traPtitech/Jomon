import { useApplicationDetailStore } from "@/stores/applicationDetail";
import ApplicationDetailLogs from "@/views/components/ApplicationDetailLogs.vue";
import { createTestingPinia } from "@pinia/testing";
import { mount } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";

// Mock logs child components to avoid deep rendering issues, but keep them identifiable
vi.mock("@/views/components/CommentLog.vue", () => ({
  default: { template: '<div class="comment-log-stub"></div>' }
}));
vi.mock("@/views/components/StatusLog.vue", () => ({
  default: { template: '<div class="status-log-stub"></div>' }
}));
vi.mock("@/views/components/ChangeLog.vue", () => ({
  default: { template: '<div class="change-log-stub"></div>' }
}));
vi.mock("@/views/components/RefundLog.vue", () => ({
  default: { template: '<div class="refund-log-stub"></div>' }
}));
vi.mock("@/views/components/TimelineNewComment.vue", () => ({
  default: { template: '<div class="new-comment-stub"></div>' }
}));

describe("ApplicationDetailLogs.vue", () => {
  it("renders timeline items correctly", () => {
    const wrapper = mount(ApplicationDetailLogs, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              applicationDetail: {
                core: {
                  comments: [
                    {
                      comment_id: "1",
                      user: { trap_id: "user1" },
                      comment: "test comment",
                      created_at: new Date().toISOString(),
                      updated_at: new Date().toISOString()
                    }
                  ],
                  state_logs: [],
                  application_detail_logs: [],
                  repayment_logs: []
                }
              }
            }
          })
        ],
        stubs: {
          "v-timeline": false, // Render v-timeline
          "v-timeline-item": false // Render v-timeline-item
        }
      }
    });

    const store = useApplicationDetailStore();
    // Force re-evaluation of computed if needed, though initial state should work.
    // The computed 'logs' depends on 'core'.

    expect(wrapper.find(".v-timeline").exists()).toBe(true);
    // Check for v-timeline-item. Since we didn't stub it, it should render (or try to).
    // However, Vuetify components in unit tests can be tricky.
    // Ideally we verify that the list of logs maps to items.

    // Check if the comment log stub is present
    expect(wrapper.find(".comment-log-stub").exists()).toBe(true);
  });

  it("renders multiple types of logs", () => {
    const wrapper = mount(ApplicationDetailLogs, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              applicationDetail: {
                core: {
                  comments: [
                    {
                      comment_id: "1",
                      user: { trap_id: "user1" },
                      comment: "test comment",
                      created_at: "2023-01-01T10:00:00Z",
                      updated_at: "2023-01-01T10:00:00Z"
                    }
                  ],
                  state_logs: [
                    {
                      update_user: { trap_id: "admin" },
                      to_state: "accepted",
                      reason: "ok",
                      created_at: "2023-01-01T11:00:00Z"
                    }
                  ],
                  application_detail_logs: [],
                  repayment_logs: []
                }
              }
            }
          })
        ],
        stubs: {
          "v-timeline": false,
          "v-timeline-item": false
        }
      }
    });

    expect(wrapper.findAll(".comment-log-stub").length).toBe(1);
    expect(wrapper.findAll(".status-log-stub").length).toBe(1);
  });
});
