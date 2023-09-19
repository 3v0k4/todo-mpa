const { defineConfig } = require("cypress");

module.exports = defineConfig({
  viewportWidth: 890,
  numTestsKeptInMemory: 1,
  projectId: "n4ynap",
  e2e: {
    baseUrl: "http://127.0.0.1:3000",
    specPattern: "cypress/e2e/**/*.{js,jsx,ts,tsx}",
  },
});
