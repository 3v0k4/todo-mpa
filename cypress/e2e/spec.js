// ***********************************************
// All of these tests are written to implement
// the official TodoMVC tests written for Selenium.
//
// The Cypress tests cover the exact same functionality,
// and match the same test names as TodoMVC.
// Please read our getting started guide
// https://on.cypress.io/introduction-to-cypress
//
// You can find the original TodoMVC tests here:
// https://github.com/tastejs/todomvc/blob/master/tests/test.js
// ***********************************************

import createTodoCommands from "../support/e2e";

// to find flaky tests we are running the entire suite N times
const N = parseFloat(Cypress.env("times") || "1");
console.log("Running tests %d time(s)", N);
if (!Cypress._.isFinite(N)) {
  throw new Error(
    `Invalid number of tests ${N} from env "${Cypress.env("times")}"`
  );
}

Cypress._.times(N, () => {
  describe("TodoMPA", function () {
    // setup these constants to match what TodoMVC does
    let TODO_ITEM_ONE = "buy some cheese";
    let TODO_ITEM_TWO = "feed the cat";
    let TODO_ITEM_THREE = "book a doctors appointment";

    const selectors = {
      newTodo: ".new-todo",
      todoList: ".todo-list",
      todoItems: ".todo-list li",
      todoItemsVisible: ".todo-list li:visible",
      count: "span.todo-count",
      main: ".main",
      footer: ".footer",
      toggleAll: ".toggle-all",
      clearCompleted: ".clear-completed",
      filters: ".filters",
      filterItems: ".filters li a",
    };

    const visibleTodos = () => cy.get(selectors.todoItemsVisible);

    beforeEach(function () {
      cy.wrap(fetch("/reset", { method: "POST" }));

      cy.visit("/");

      cy.contains("h1", "todos").should("be.visible");

      cy.document().then(() => {
        createTodoCommands();
      });
    });

    context("When page is initially opened", function () {
      it("should focus on the todo input field", function () {
        cy.get(".new-todo").should("have.attr", "autofocus");
      });
    });

    context("No Todos", function () {
      it("should hide #main and #footer", function () {
        cy.get(selectors.todoItems).should("not.exist");
        cy.get(selectors.main).should("not.be.visible");
        cy.get(selectors.footer).should("not.be.visible");
      });
    });

    context("New Todo", function () {
      it("should allow me to add todo items", function () {
        cy.get(selectors.newTodo).type(`${TODO_ITEM_ONE}{enter}`);

        cy.get(selectors.todoItems)
          .eq(0)
          .find("label")
          .should("contain", TODO_ITEM_ONE);

        cy.get(selectors.newTodo).type(`${TODO_ITEM_TWO}{enter}`);

        cy.get(selectors.todoItems)
          .eq(1)
          .find("label")
          .should("contain", TODO_ITEM_TWO);
      });

      it("should clear text input field when an item is added", function () {
        cy.get(selectors.newTodo).type(`${TODO_ITEM_ONE}{enter}`);

        cy.get(selectors.newTodo).should("have.text", "");
      });

      it("should append new items to the bottom of the list", function () {
        cy.createDefaultTodos();
        cy.wait(50);

        cy.get(selectors.count).contains("3");
        cy.get(".todo-list li")
          .eq(0)
          .find("label")
          .should("contain", TODO_ITEM_ONE);
        cy.get(".todo-list li")
          .eq(1)
          .find("label")
          .should("contain", TODO_ITEM_TWO);
        cy.get(".todo-list li")
          .eq(2)
          .find("label")
          .should("contain", TODO_ITEM_THREE);
      });

      it("should trim text input", function () {
        cy.createTodo(`    ${TODO_ITEM_ONE}    `);

        cy.get(selectors.todoItems)
          .eq(0)
          .find("label")
          .should("have.text", TODO_ITEM_ONE);
      });

      it("should show #main and #footer when items added", function () {
        cy.createTodo(TODO_ITEM_ONE);

        cy.get(selectors.main).should("be.visible");
        cy.get(selectors.footer).should("be.visible");
      });
    });

    context("Mark all as completed", function () {
      beforeEach(function () {
        cy.createDefaultTodos();
      });

      it("should allow me to mark all items as completed", function () {
        cy.get(selectors.toggleAll).click();
        cy.wait(50);

        cy.get(".todo-list li").eq(0).should("have.class", "completed");
        cy.get(".todo-list li").eq(1).should("have.class", "completed");
        cy.get(".todo-list li").eq(2).should("have.class", "completed");
      });

      it("should allow me to clear the complete state of all items", function () {
        cy.get(selectors.toggleAll).click();
        cy.wait(50);
        cy.get(selectors.toggleAll).click();
        cy.wait(50);

        cy.get(".todo-list li").eq(0).should("not.have.class", "completed");
        cy.get(".todo-list li").eq(1).should("not.have.class", "completed");
        cy.get(".todo-list li").eq(2).should("not.have.class", "completed");
      });

      it("complete all checkbox should update state when items are completed / cleared", function () {
        cy.get(selectors.toggleAll).click();
        cy.wait(50);

        cy.get(selectors.toggleAll).should("be.checked");

        cy.get(selectors.todoItems).eq(0).find(".toggle").click();
        cy.wait(50);

        cy.get(selectors.toggleAll).should("not.be.checked");

        cy.get(".todo-list li").eq(0).find(".toggle").click();

        cy.get(selectors.toggleAll).should("be.checked");
      });
    });

    context("Item", function () {
      it("should allow me to mark items as complete", function () {
        cy.createTodo(TODO_ITEM_ONE);
        cy.createTodo(TODO_ITEM_TWO);

        cy.get(".todo-list li").eq(0).find(".toggle").check();

        cy.get(".todo-list li").eq(0).should("have.class", "completed");
        cy.get(".todo-list li").eq(1).should("not.have.class", "completed");

        cy.wait(100);

        cy.get(".todo-list li").eq(1).find(".toggle").check();

        cy.get(".todo-list li").eq(0).should("have.class", "completed");
        cy.get(".todo-list li").eq(1).should("have.class", "completed");
      });

      it("should allow me to un-mark items as complete", function () {
        cy.createTodo(TODO_ITEM_ONE);
        cy.createTodo(TODO_ITEM_TWO);

        cy.get(".todo-list li").eq(0).find(".toggle").click();

        cy.get(".todo-list li").eq(0).should("have.class", "completed");
        cy.get(".todo-list li").eq(1).should("not.have.class", "completed");

        cy.wait(100);

        cy.get(".todo-list li").eq(0).find(".toggle").click();

        cy.get(".todo-list li").eq(0).should("not.have.class", "completed");
        cy.get(".todo-list li").eq(1).should("not.have.class", "completed");
      });

      it("should allow me to edit an item", function () {
        cy.createDefaultTodos();

        visibleTodos().eq(1).find("label").dblclick();

        visibleTodos().eq(1).find(".edit").should("have.value", TODO_ITEM_TWO);

        cy.get(".todo-list li").eq(1).find(".edit").clear();
        cy.get(".todo-list li")
          .eq(1)
          .find(".edit")
          .type("buy some sausages{enter}");

        visibleTodos().eq(0).should("contain", TODO_ITEM_ONE);
        visibleTodos().eq(1).should("contain", "buy some sausages");
        visibleTodos().eq(2).should("contain", TODO_ITEM_THREE);
      });
    });

    context("Editing", function () {
      beforeEach(function () {
        cy.createDefaultTodos();
      });

      it("should hide other controls when editing", function () {
        cy.get(".todo-list li").eq(1).find("label").dblclick();

        cy.get(selectors.todoItems)
          .eq(1)
          .find(".toggle")
          .should("not.be.visible");
        cy.get(selectors.todoItems)
          .eq(1)
          .find("label")
          .should("not.be.visible");
      });

      it("should save edits on blur", function () {
        cy.get(".todo-list li").eq(1).find("label").dblclick();
        cy.get(selectors.todoItems).eq(1).find(".edit").clear();

        cy.get(selectors.todoItems)
          .eq(1)
          .find(".edit")
          .type("buy some sausages")
          .blur();

        visibleTodos().eq(0).should("contain", TODO_ITEM_ONE);
        visibleTodos().eq(1).should("contain", "buy some sausages");
        visibleTodos().eq(2).should("contain", TODO_ITEM_THREE);
      });

      it("should trim entered text", function () {
        cy.get(".todo-list li").eq(1).find("label").dblclick();
        cy.get(selectors.todoItems)
          .eq(1)
          .find(".edit")
          .type("{selectall}{backspace}    buy some sausages    {enter}");

        visibleTodos().eq(0).should("contain", TODO_ITEM_ONE);
        visibleTodos().eq(1).should("contain", "buy some sausages");
        visibleTodos().eq(2).should("contain", TODO_ITEM_THREE);
      });

      it("should remove the item if an empty text string was entered", function () {
        cy.get(".todo-list li").eq(1).find("label").dblclick();
        cy.get(selectors.todoItems).eq(1).find(".edit").clear();
        cy.get(selectors.todoItems).eq(1).find(".edit").type("{enter}");

        visibleTodos().should("have.length", 2);
      });

      it("should cancel edits on escape", function () {
        visibleTodos().eq(1).find("label").dblclick();
        cy.get(selectors.todoItems)
          .eq(1)
          .find(".edit")
          .type("{selectall}{backspace}foo{esc}");

        visibleTodos().eq(0).should("contain", TODO_ITEM_ONE);
        visibleTodos().eq(1).should("contain", TODO_ITEM_TWO);
        visibleTodos().eq(2).should("contain", TODO_ITEM_THREE);
      });
    });

    context("Counter", function () {
      it("should display the current number of todo items", function () {
        cy.createTodo(TODO_ITEM_ONE);
        cy.wait(50);

        cy.get(selectors.count).contains("1");

        cy.createTodo(TODO_ITEM_TWO);
        cy.wait(50);

        cy.get(selectors.count).contains("2");
      });
    });

    context("Clear completed button", function () {
      beforeEach(function () {
        cy.createDefaultTodos();
      });

      it("should display the correct text", function () {
        cy.get(".todo-list li").eq(0).find(".toggle").check();

        cy.get(selectors.clearCompleted).contains("Clear completed");
      });

      it("should remove completed items when clicked", function () {
        cy.get(".todo-list li").eq(1).find(".toggle").check();
        cy.get(selectors.clearCompleted).click();

        cy.get(".todo-list li").should("have.length", 2);
        cy.get(".todo-list li").eq(0).should("contain", TODO_ITEM_ONE);
        cy.get(".todo-list li").eq(1).should("contain", TODO_ITEM_THREE);
      });

      it("should be hidden when there are no items that are completed", function () {
        cy.get(".todo-list li").eq(1).find(".toggle").check();
        cy.get(selectors.clearCompleted).should("be.visible").click();

        cy.get(selectors.clearCompleted).should("not.exist");
      });
    });

    context("Routing", function () {
      beforeEach(function () {
        cy.createDefaultTodos();
      });

      it("should allow me to display active items", function () {
        cy.get(".todo-list li").eq(1).find(".toggle").check();
        cy.contains(selectors.filterItems, "Active").click();

        visibleTodos()
          .should("have.length", 2)
          .eq(0)
          .should("contain", TODO_ITEM_ONE);
        visibleTodos().eq(1).should("contain", TODO_ITEM_THREE);
      });

      it("should respect the back button", function () {
        cy.get(".todo-list li").eq(1).find(".toggle").check();

        cy.log("Showing all items");
        cy.contains(selectors.filterItems, "All").click();
        visibleTodos().should("have.length", 3);

        cy.log("Showing active items");
        cy.contains(selectors.filterItems, "Active").click();
        cy.wait(50);

        visibleTodos().should("have.length", 2);

        cy.log("Showing completed items");
        cy.contains(selectors.filterItems, "Completed").click();
        cy.wait(50);

        visibleTodos().should("have.length", 1);

        cy.log("Back to active items");
        cy.go("back");
        cy.wait(50);

        visibleTodos().should("have.length", 2);

        cy.log("Back to all items");
        cy.go("back");
        cy.wait(50);

        visibleTodos().should("have.length", 3);
      });

      it("should allow me to display completed items", function () {
        visibleTodos().eq(1).find(".toggle").check();
        cy.get(selectors.filters).contains("Completed").click();

        visibleTodos().should("have.length", 1);
      });

      it("should allow me to display all items", function () {
        visibleTodos().eq(1).find(".toggle").check();
        cy.get(selectors.filters).contains("Active").click();
        cy.get(selectors.filters).contains("Completed").click();
        cy.get(selectors.filters).contains("All").click();

        visibleTodos().should("have.length", 3);
      });

      it("should highlight the currently applied filter", function () {
        cy.contains(selectors.filterItems, "All").should(
          "have.class",
          "selected"
        );

        cy.contains(selectors.filterItems, "Active").click();

        cy.contains(selectors.filterItems, "Active").should(
          "have.class",
          "selected"
        );

        cy.contains(selectors.filterItems, "Completed").click();

        cy.contains(selectors.filterItems, "Completed").should(
          "have.class",
          "selected"
        );
      });
    });
  });
});
