import { expect, test } from "@playwright/test";

test.describe("Login Flow", () => {
  test("should load home page when user is logged in", async ({ page }) => {
    // Mock the /api/users/me endpoint to return a valid user
    await page.route("*/**/api/users/me", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ trap_id: "test-user", is_admin: false })
      });
    });

    // Mock the /api/users endpoint (used in ApplicationListPage)
    await page.route("*/**/api/users", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([{ trap_id: "test-user" }])
      });
    });

    // Mock /api/applications (used in ApplicationListPage)
    await page.route("*/**/api/applications", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([])
      });
    });

    await page.goto("/");

    // Verify that the user is on the home page and logged in
    await expect(page).toHaveTitle(/Jomon/);
    // You might want to add more specific assertions here, e.g., checking for the user's icon or name
  });

  test("should redirect to auth endpoint when user is not logged in", async ({
    page
  }) => {
    // Mock /api/users/me to fail
    await page.route("*/**/api/users/me", async route => {
      await route.fulfill({
        status: 401
      });
    });

    // Mock /api/auth/genpkce
    await page.route("*/**/api/auth/genpkce", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          client_id: "test-client-id",
          code_challenge: "test-challenge",
          code_challenge_method: "S256"
        })
      });
    });

    // We expect the page to navigate to the traQ auth URL
    const navigationPromise = page.waitForURL(/q\.trap\.jp/);

    await page.goto("/");

    await navigationPromise;
  });
});
