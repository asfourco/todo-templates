import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {Todo} from '../models/todo';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export default function TodoApp() {
    const [todos, setTodos] = useState<Todo[]>([]);
    const [newTodo, setNewTodo] = useState<string>('');
    const [updatedTodo, setUpdatedTodo] = useState<Todo | null>(null);

    useEffect(() => {
        // Fetch the list of todo items
        axios.get(`${API_BASE_URL}/todos`)
            .then((response) => {
                console.log(response.data);
                setTodos(response.data.items ?? []);
            })
            .catch((error) => {
                console.error('Error fetching todos:', error);
            });
    }, []);

    const handleAddTodo = () => {
        // Add a new todo item
        axios.post(`${API_BASE_URL}/todos`, { title: newTodo })
            .then((response) => {
                setTodos([...todos, response.data]);
                setNewTodo('');
            })
            .catch((error) => {
                console.error('Error adding todo:', error);
            });
    };

    const handleUpdateTodo = (id: number, title: string, active: boolean) => {
        // Update a todo item
        axios.patch(`${API_BASE_URL}/todos`, { id, title, active })
            .then((response) => {
                const updatedTodos = todos.map((todo) => {
                    if (todo.id === id) {
                        return response.data;
                    }
                    return todo;
                });
                setTodos(updatedTodos);
                setUpdatedTodo(null);
            })
            .catch((error) => {
                console.error('Error updating todo:', error);
            });
    };

    return (
        <div className="">
            <h1>Todo List</h1>
            <ul>
                {todos.map((todo) => (
                    <li key={todo.id}>
                        {todo.title}
                        <button onClick={() => setUpdatedTodo(todo)}>Edit</button>
                        <button onClick={() => handleUpdateTodo(todo.id, todo.title, !todo.active)}>
                            {todo.active ? 'Deactivate' : 'Activate'}
                        </button>
                    </li>
                ))}
            </ul>
            <input
                type="text"
                placeholder="New Todo"
                value={newTodo}
                onChange={(e) => setNewTodo(e.target.value)}
            />
            <button onClick={handleAddTodo}>Add Todo</button>
            {updatedTodo && (
                <div>
                    <input
                        type="text"
                        value={updatedTodo.title}
                        onChange={(e) => setUpdatedTodo({ ...updatedTodo, title: e.target.value })}
                    />
                    <button onClick={() => handleUpdateTodo(updatedTodo.id, updatedTodo.title, updatedTodo.active)}>
                        Save
                    </button>
                </div>
            )}
        </div>
    );
}
