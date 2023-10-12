// src/App.tsx
import React from 'react';
import TodoList from './components/TodoList';
import TodoForm from './components/TodoForm';

function App() {
    return (
        <div className="App">
            <h1>Todo List</h1>
            <TodoForm />
            <TodoList />
        </div>
    );
}

export default App;
