import { expect, test } from "@playwright/test";

test.describe("Admin Page Flow", () => {
  test.beforeEach(async ({ page }) => {
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
          { trap_id: "test-user", is_admin: false },
          { trap_id: "new-admin", is_admin: false }
        ])
      });
    });

    // Mock admin update API
    await page.route("*/**/api/users/admins", async route => {
      const method = route.request().method();
      if (method === "PUT") {
        const postData = JSON.parse(route.request().postData() || "{}");
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(postData)
        });
      }
    });
  });

  test("should display admin page for admin user", async ({ page }) => {
    const userListPromise = page.waitForResponse(
      resp => resp.url().includes("/api/users") && resp.status() === 200
    );
    await page.goto("/admin");
    await userListPromise;
    await expect(page.getByText("admin-user")).toBeVisible();
    await expect(page.getByLabel("管理権限を削除")).toBeVisible();
    await expect(page.getByLabel("管理権限を追加")).toBeVisible();
  });

  test("should add admin privileges", async ({ page }) => {
    await page.goto("/admin");

    // Select user to add as admin
    // Use .v-field selector for Vuetify v-autocomplete
    await page
      .locator(".v-input")
      .filter({ hasText: "管理権限を追加" })
      .locator(".v-field")
      .click();
    await page.getByRole("option", { name: "new-admin" }).click();

    // Close the dropdown by pressing Escape
    await page.keyboard.press("Escape");

    // Wait for API call
    const requestPromise = page.waitForRequest(
      request =>
        request.url().includes("/api/users/admins") &&
        request.method() === "PUT" &&
        request.postDataJSON().to_admin === true &&
        request.postDataJSON().trap_id === "new-admin"
    );

    // Click the second "設定" button (for adding)
    await page
      .locator(".v-input")
      .filter({ hasText: "管理権限を追加" })
      .locator("xpath=following-sibling::button | following::button")
      .first()
      .click({ force: true });

    await requestPromise;
  });

  test("should remove admin privileges", async ({ page }) => {
    await page.goto("/admin");

    // Select user to remove from admin
    await page
      .locator(".v-input")
      .filter({ hasText: "管理権限を削除" })
      .locator(".v-field")
      .click();
    await page.getByRole("option", { name: "admin-user" }).click();

    // Close dropdown by pressing Escape
    await page.keyboard.press("Escape");

    // Wait for API call
    const requestPromise = page.waitForRequest(
      request =>
        request.url().includes("/api/users/admins") &&
        request.method() === "PUT" &&
        request.postDataJSON().to_admin === false &&
        request.postDataJSON().trap_id === "admin-user"
    );

    // Click the first "設定" button (for removing)
    await page
      .locator(".v-input")
      .filter({ hasText: "管理権限を削除" })
      .locator("xpath=following-sibling::button | following::button")
      .first()
      .click({ force: true });

    await requestPromise;
  });
});

test.describe("Admin Page Access Control", () => {
  test("should show error message for non-admin user", async ({ page }) => {
    // Mock user API (Non-admin user)
    await page.route("*/**/api/users/me", async route => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ trap_id: "test-user", is_admin: false })
      });
    });

    await page.goto("/admin");
    await expect(page.getByText("権限がありません")).toBeVisible();
    await expect(page.getByLabel("管理権限を追加")).not.toBeVisible();
  });
});
