export function createTodoCommands() {
  const newTodoSelector = ".new-todo";
  const todoListSelector = ".todo-list";
  const todoItemsSelector = ".todo-list li";

  Cypress.Commands.add("createDefaultTodos", function () {
    let TODO_ITEM_ONE = "buy some cheese";
    let TODO_ITEM_TWO = "feed the cat";
    let TODO_ITEM_THREE = "book a doctors appointment";

    let cmd = Cypress.log({
      name: "create default todos",
      message: [],
      consoleProps() {
        return {
          "Inserted Todos": [TODO_ITEM_ONE, TODO_ITEM_TWO, TODO_ITEM_THREE],
        };
      },
    });

    cy.get(newTodoSelector, { log: false }).type(`${TODO_ITEM_ONE}{enter}`, {
      log: false,
    });

    cy.get(todoItemsSelector, { log: false }).should("have.length", 1);

    cy.get(newTodoSelector, { log: false }).type(`${TODO_ITEM_TWO}{enter}`, {
      log: false,
    });

    cy.get(todoItemsSelector, { log: false }).should("have.length", 2);

    cy.get(newTodoSelector, { log: false }).type(`${TODO_ITEM_THREE}{enter}`, {
      log: false,
    });

    cy.get(todoItemsSelector, { log: false }).should("have.length", 3);

    const combinedSelector = todoItemsSelector + ":visible";
    cy.get(combinedSelector, { log: false }).then(function ($listItems) {
      cmd.set({ $el: $listItems }).snapshot().end();
    });
  });

  Cypress.Commands.add("createTodo", function (todo) {
    let cmd = Cypress.log({
      name: "create todo",
      message: todo,
      consoleProps() {
        return {
          "Inserted Todo": todo,
        };
      },
    });

    cy.get(newTodoSelector, { log: false }).type(`${todo}{enter}`, {
      log: false,
    });

    cy.get(todoListSelector, { log: false })
      .contains("li", todo.trim(), { log: false })
      .then(function ($li) {
        cmd.set({ $el: $li }).snapshot().end();
      });
  });
}
