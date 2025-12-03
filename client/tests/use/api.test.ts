import axios from "axios";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { redirectAuthEndpoint, traQBaseURL } from "../../src/use/api";

vi.mock("axios");

describe("api", () => {
  describe("redirectAuthEndpoint", () => {
    const originalLocation = window.location;

    beforeEach(() => {
      // Mock window.location.assign
      Object.defineProperty(window, "location", {
        writable: true,
        value: {
          ...originalLocation,
          assign: vi.fn()
        }
      });
    });

    afterEach(() => {
      Object.defineProperty(window, "location", {
        writable: true,
        value: originalLocation
      });
      vi.clearAllMocks();
    });

    it("redirects to the correct URL with PKCE parameters", async () => {
      const mockData = {
        client_id: "test-client-id",
        code_challenge: "test-challenge",
        code_challenge_method: "S256"
      };

      (axios.get as any).mockResolvedValue({ data: mockData });

      await redirectAuthEndpoint();

      expect(axios.get).toHaveBeenCalledWith("/api/auth/genpkce");

      const expectedUrl = new URL(`${traQBaseURL}/oauth2/authorize`);
      expectedUrl.search = new URLSearchParams({
        response_type: "code",
        client_id: mockData.client_id,
        code_challenge: mockData.code_challenge,
        code_challenge_method: mockData.code_challenge_method
      }).toString();

      expect(window.location.assign).toHaveBeenCalledWith(
        expectedUrl.toString()
      );
    });
  });
});
