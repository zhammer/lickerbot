/* global cy */
/// <reference types="cypress" />
import { Given, Then, When } from "cypress-cucumber-preprocessor/steps";

// uses BUFFALO grifts from bin (built using buffalo build)
// *much* faster than using buffalo directly.
const BUFFALO = "./bin/lickerbot";

beforeEach(() => {
  // clear and reseed tables before each run
  cy.exec(`${BUFFALO} task db:seed`, {
    env: { GO_ENV: "test" },
  });
});

When(`I visit {string}`, (path) => {
  cy.visit(path);
});

When(`I visit {string} expecting a non-200 status`, (path) => {
  cy.visit(path, { failOnStatusCode: false });
});

When(`I click the link {string}`, (text) => {
  cy.get("a").contains(text).click();
});

When(`I click the button {string}`, (text) => {
  // { force: true } because cypress thinks it can't see the button sometimes
  // due to hidden overlay.
  cy.get("button").contains(text).click({ force: true });
});

When(`I select the amount {string}`, (amount) => {
  cy.get("select").select(amount);
});

When("I refresh the page", () => {
  cy.reload();
});

Then(`I am on {string}`, (path) => {
  cy.location("pathname").should("eq", path);
});

Then(`I am on the hash {string}`, (hash) => {
  cy.location("hash").should("eq", hash);
});

Then(`I see the text {string}`, (text) => {
  cy.contains(text);
});

Then(`I see {int} embedded tweets`, (numberOfTweets) => {
  cy.get("twitter-widget").should("have.length", numberOfTweets);
});

Then(`the meta tag {string} {} {string}`, (tag, assertion, content) => {
  cy.get(`meta[${metaAttr(tag)}="${tag}"]`)
    .should("have.attr", "content")
    // contains -> contain, equals -> equal
    .should(assertion.slice(0, -1), content);
});

Then(`the meta tag {string} contains a resource that exists`, (tag) => {
  cy.get(`meta[${metaAttr(tag)}="${tag}"]`)
    .should("have.attr", "content")
    .then((content) => {
      cy.request(content);
    });
});

Then(`the title is {string}`, (text) => {
  cy.get("title").should("have.text", text);
});

// <meta ~name~="twitter:card" content="summary" />
// <meta ~property~="og:type" content="website" />
function metaAttr(name) {
  if (name.startsWith("twitter:")) {
    return "name";
  }
  return "property";
}
