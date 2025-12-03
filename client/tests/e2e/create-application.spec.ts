import { expect, test } from "@playwright/test";

test.describe("Create Application Flow", () => {
  test.beforeEach(async ({ page }) => {
    // Mock login
    await page.route("*/**/api/users/me", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ trap_id: "test-user", is_admin: false })
      });
    });

    // Mock user list
    await page.route("*/**/api/users", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([
          { trap_id: "test-user" },
          { trap_id: "other-user" }
        ])
      });
    });
  });

  test("should create a new application successfully", async ({ page }) => {
    // Mock POST /api/applications
    await page.route("*/**/api/applications", async route => {
      const request = route.request();
      if (request.method() === "POST") {
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({
            application_id: "new-app-id",
            applicant: { trapid: "test-user" },
            created_at: new Date().toISOString(),
            current_detail: {
              title: "Test Application",
              type: "club",
              amount: 1000,
              remarks: "Test Remarks",
              created_at: new Date().toISOString(),
              paid_at: new Date().toISOString()
            }
          })
        });
      } else {
        await route.continue();
      }
    });

    // Mock redirection target page (Application Detail)
    await page.route("*/**/api/applications/new-app-id", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          id: "new-app-id",
          created_at: new Date().toISOString(),
          applicant: { trap_id: "test-user" },
          latest_detail: {
            id: "detail-id",
            application_id: "new-app-id",
            update_user: { trap_id: "test-user" },
            type: "club",
            title: "Test Application",
            remarks: "Test Remarks",
            amount: 1000,
            paid_at: new Date().toISOString(),
            repaid_to: []
          },
          comments: [],
          images: [],
          histories: []
        })
      });
    });

    await page.goto("/applications/new/club");

    // Fill form
    const summaryInput = page.getByLabel("概要", { exact: true });
    await summaryInput.waitFor({ state: "visible" });
    await summaryInput.fill("Test Application");

    // Date picker interaction
    await page.getByLabel("支払日").click();
    // Click a day (e.g., 15th) to select it. Using first() to avoid ambiguity if multiple calendars (unlikely) or other 15s.
    // We assume the current month has a 15th.
    await page.getByText("15", { exact: true }).first().click();

    await page.getByLabel("支払金額").fill("1000");

    // Select Repayment Target (Autocomplete)
    // Click to open menu
    await page.getByLabel("返金対象者").click();
    // Select item
    await page.getByText("test-user", { exact: true }).click();
    // Click away to close menu (optional, but good practice)
    await page.getByLabel("概要", { exact: true }).click();

    await page.getByLabel("購入物の概要").fill("Test Remarks");

    // Submit
    await page.getByRole("button", { name: "作成する" }).click();

    // Verify success message
    await expect(page.getByText("作成できました")).toBeVisible();

    const okButton = page.getByRole("link", { name: "OK" });

    // Verify redirection button and click it
    await okButton.click();

    // Verify redirection
    await expect(page).toHaveURL(/\/applications\/new-app-id/);
  });
});
