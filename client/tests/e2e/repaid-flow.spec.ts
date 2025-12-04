import { expect, test } from "@playwright/test";

test.describe("Repaid Application Flow", () => {
  test.beforeEach(async ({ page }) => {
    // Mock login as admin
    await page.route("*/**/api/users/me", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ trap_id: "admin-user", is_admin: true })
      });
    });

    // Mock user list
    await page.route("*/**/api/users", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([
          { trap_id: "admin-user" },
          { trap_id: "applicant-user" },
          { trap_id: "target-user" }
        ])
      });
    });
  });

  test("should mark user as repaid", async ({ page }) => {
    const appId = "repaid-app-id";

    // Mock application detail (accepted state)
    await page.route(`*/**/api/applications/${appId}`, async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          id: appId,
          application_id: appId,
          created_at: new Date().toISOString(),
          applicant: { trap_id: "applicant-user" },
          current_state: "accepted",
          latest_detail: {
            id: "detail-id",
            application_id: appId,
            update_user: { trap_id: "applicant-user" },
            type: "club",
            title: "Test Application",
            remarks: "Test Remarks",
            amount: 1000,
            paid_at: new Date().toISOString(),
            repaid_to: []
          },
          comments: [],
          images: [],
          histories: [],
          state_logs: [],
          application_detail_logs: [],
          // Repayment logs indicate who needs to be repaid
          repayment_logs: [
            {
              repaid_to_user: { trap_id: "target-user" },
              repaid_at: null // Not repaid yet
            }
          ]
        })
      });
    });

    // Mock repayment API
    let repaidCalled = false;
    await page.route(
      `*/**/api/applications/${appId}/states/repaid/target-user`,
      async route => {
        if (route.request().method() === "PUT") {
          repaidCalled = true;
          await route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify({})
          });
        } else {
          await route.continue();
        }
      }
    );

    // Wait for load - setup waiter before navigation or reload
    const responsePromise = page.waitForResponse(
      response =>
        response.url().includes(`/api/applications/${appId}`) &&
        response.status() === 200
    );
    
    await page.goto(`/applications/${appId}`);
    await responsePromise;

    // Open Dialog
    const openDialogBtn = page.getByRole("button", {
      name: "払い戻し済みのユーザーを選択"
    });
    await expect(openDialogBtn).toBeVisible();
    await openDialogBtn.click();

    // Dialog should be visible
    const dialog = page.locator(".v-dialog"); // Vuetify dialog class
    await expect(dialog).toBeVisible();

    // Select User
    // Click the v-select to open the menu. Use a more robust selector.
    // .v-field is the clickable area for v-select
    const select = dialog
      .locator(".v-field")
      .filter({ hasText: "払い戻し済みのユーザーを選択" });
    await select.click();

    // Select "target-user" from the list
    // The menu is attached to root (zIndex 3000), so it's outside the dialog.
    // We look for the menu content.
    const menu = page.locator(".v-overlay__content");
    const option = menu.getByText("target-user", { exact: true });
    await option.waitFor({ state: "visible" });
    await option.click();

    // Close the menu by clicking the dialog title or background
    // Since v-select multiple stays open.
    await dialog.getByText("払い戻し日").click();

    // Verify selection (chip should appear or value change)
    // Ideally, the chip with "target-user" should be visible in the input
    await expect(
      page.locator(".v-chip").getByText("target-user")
    ).toBeVisible();

    // Click OK
    const okBtn = page.getByRole("button", { name: "OK" });
    await expect(okBtn).toBeEnabled(); // Should be enabled now
    await okBtn.click();

    // Wait for API call
    await expect.poll(() => repaidCalled).toBe(true);
  });
});
