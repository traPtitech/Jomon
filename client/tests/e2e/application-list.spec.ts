import { expect, test } from "@playwright/test";

test.describe("Application List Flow", () => {
  test.beforeEach(async ({ page, browserName }) => {
    // Disable animations for stability (Firefox only)
    if (browserName === "firefox") {
      await page.addStyleTag({
        content: `
          *, *::before, *::after {
            transition: none !important;
            animation: none !important;
          }
        `
      });
    }
    // Mock user API (Admin user)
    await page.route("*/**/api/users/me", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ trap_id: "admin-user", is_admin: true })
      });
    });

    // Mock users list
    await page.route("*/**/api/users", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([
          { trap_id: "admin-user", is_admin: true },
          { trap_id: "test-user", is_admin: false }
        ])
      });
    });

    // Mock applications list (Default)
    await page.route("*/**/api/applications*", async route => {
      const url = new URL(route.request().url());
      const type = url.searchParams.get("type");
      const state = url.searchParams.get("current_state");

      let applications = [
        {
          application_id: "app-1",
          applicant: { trap_id: "test-user" },
          current_detail: { title: "Club App", amount: 1000, type: "club" },
          current_state: "submitted",
          created_at: "2023-01-01T00:00:00Z"
        },
        {
          application_id: "app-2",
          applicant: { trap_id: "test-user" },
          current_detail: {
            title: "Contest App",
            amount: 2000,
            type: "contest"
          },
          current_state: "accepted",
          created_at: "2023-01-02T00:00:00Z"
        }
      ];

      if (type) {
        applications = applications.filter(
          app => app.current_detail.type === type
        );
      }
      if (state) {
        applications = applications.filter(app => app.current_state === state);
      }

      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(applications)
      });
    });
  });

  test("should display application list", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByText("Club App")).toBeVisible();
    await expect(page.getByText("Contest App")).toBeVisible();
  });

  test("should filter by application type", async ({ page, browserName }) => {
    await page.goto("/");

    // Ensure filter menu is open
    // Ensure filter menu is open
    const filterCard = page.locator(".v-card").filter({ hasText: "絞り込み" });
    const filterButton = filterCard.locator(".v-card-title button");
    const filterForm = filterCard.locator(".v-card-text");

    // Check if the button shows 'mdi-chevron-down' (meaning it's closed)
    if (await filterButton.getByText("mdi-chevron-down").isVisible()) {
      await filterButton.click({ force: browserName === "firefox" });
    }
    // Wait for animation
    await page.waitForTimeout(500);
    await expect(filterForm).toBeVisible();

    // Select type 'club'
    // Click the v-field containing the label
    const typeField = filterForm
      .locator(".v-field")
      .filter({ hasText: "申請タイプ" });
    await typeField.click({ force: browserName === "firefox" });
    if (browserName === "firefox") await page.waitForTimeout(200);
    await page
      .getByRole("option", { name: "部費利用申請" })
      .click({ force: browserName === "firefox" });
    // Wait for selection to be applied
    await expect(filterForm.getByText("部費利用申請")).toBeVisible();

    // Click reload button
    await filterForm.locator("button.bg-primary").first().click();

    await expect(page.getByText("Club App")).toBeVisible();
    await expect(page.getByText("Contest App")).not.toBeVisible();
  });

  test("should filter by status", async ({ page, browserName }) => {
    await page.goto("/");

    // Ensure filter menu is open
    // Ensure filter menu is open
    const filterCard = page.locator(".v-card").filter({ hasText: "絞り込み" });
    const filterButton = filterCard.locator(".v-card-title button");
    const filterForm = filterCard.locator(".v-card-text");

    // Check if the button shows 'mdi-chevron-down' (meaning it's closed)
    if (await filterButton.getByText("mdi-chevron-down").isVisible()) {
      await filterButton.click({ force: browserName === "firefox" });
    }
    // Wait for animation
    await page.waitForTimeout(500);
    await expect(filterForm).toBeVisible();

    // Select state 'accepted'
    const statusField = filterForm
      .locator(".v-field")
      .filter({ hasText: "現在の状態" });
    await statusField.click({ force: browserName === "firefox" });
    if (browserName === "firefox") await page.waitForTimeout(200);
    await page
      .getByRole("option", { name: "払い戻し待ち" })
      .click({ force: browserName === "firefox" });
    // Wait for selection to be applied
    await expect(filterForm.getByText("払い戻し待ち")).toBeVisible();

    // Click reload button
    await filterForm.locator("button.bg-primary").first().click();

    await expect(page.getByText("Contest App")).toBeVisible();
    await expect(page.getByText("Club App")).not.toBeVisible();
  });
});
