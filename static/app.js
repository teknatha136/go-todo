const API_URL = "/todos";

document.addEventListener("DOMContentLoaded", () => {
    fetchTodos();

    // Add event listener for Enter key on the input field
    const todoInput = document.getElementById("todo-input");
    todoInput.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
            addTodo();
        }
    });
});

async function fetchTodos() {
    const response = await fetch(API_URL);
    const todos = await response.json();
    const todoList = document.getElementById("todo-list");

    todoList.innerHTML = "";
    todos.forEach(todo => {
        const todoItem = document.createElement("div");
        todoItem.className = "todo-item";
        todoItem.innerHTML = `
            <input type="checkbox" ${todo.completed ? "checked" : ""} onclick="toggleComplete(${todo.ID})">
            <span id="title-${todo.ID}" class="${todo.completed ? 'completed' : ''}">${todo.title}</span>
            <input type="text" id="edit-${todo.ID}" class="edit-input" style="display:none;" placeholder="Edit to-do" value="${todo.title}" />
            <button onclick="toggleEdit(${todo.ID})" id="edit-btn-${todo.ID}">Edit</button>
            <button onclick="deleteTodo(${todo.ID})">Delete</button>
        `;
        todoList.appendChild(todoItem);
    });
}

async function addTodo() {
    const title = document.getElementById("todo-input").value;
    if (title === "") return;

    await fetch(API_URL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ title, completed: false })
    });

    document.getElementById("todo-input").value = "";
    fetchTodos();
}

async function toggleComplete(id) {
    const todo = await fetch(`${API_URL}/${id}`).then(res => res.json());
    await fetch(`${API_URL}/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ title: todo.title, completed: !todo.completed })
    });

    fetchTodos();
}

async function deleteTodo(id) {
    await fetch(`${API_URL}/${id}`, {
        method: "DELETE"
    });

    fetchTodos();
}

async function editTodoTitle(id, newTitle) {
    const todo = await fetch(`${API_URL}/${id}`).then(res => res.json());
    await fetch(`${API_URL}/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ title: newTitle, completed: todo.completed })
    });

    fetchTodos();
}

function toggleEdit(id) {
    const titleSpan = document.getElementById(`title-${id}`);
    const editInput = document.getElementById(`edit-${id}`);
    const editButton = document.getElementById(`edit-btn-${id}`);

    if (editInput.style.display === "none") {
        // Switch to edit mode
        titleSpan.style.display = "none";
        editInput.style.display = "inline";
        editButton.innerText = "Save";
    } else {
        // Save changes
        const newTitle = editInput.value;
        editTodoTitle(id, newTitle);
        titleSpan.innerText = newTitle;

        // Switch back to view mode
        titleSpan.style.display = "inline";
        editInput.style.display = "none";
        editButton.innerText = "Edit";
    }
}
