import React, { useState } from 'react';
import { Todo } from '../models/Todo';
import { updateTodo, deleteTodo } from '../services/TodoService';

interface TodoProps {
    todo: Todo;
    onTodoUpdate: (updatedTodo: Todo) => void;
    onTodoDelete: (deletedTodoId: number) => void;
}

const TodoItem: React.FC<TodoProps> = ({ todo, onTodoUpdate, onTodoDelete }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [newTitle, setNewTitle] = useState(todo.title);

    const handleToggleActive = () => {
        const updatedTodo = { ...todo, active: !todo.active };
        updateTodo(updatedTodo.id, updatedTodo.title, updatedTodo.active)
            .then(() => onTodoUpdate(updatedTodo))
            .catch((error) => console.error('Error updating todo:', error));
    };

    const handleDeleteTodo = () => {
        deleteTodo(todo.id)
            .then(() => onTodoDelete(todo.id))
            .catch((error) => console.error('Error deleting todo:', error));
    };

    const handleEdit = () => {
        setIsEditing(true);
    };

    const handleSaveEdit = () => {
        const updatedTodo = { ...todo, title: newTitle };
        updateTodo(updatedTodo.id, updatedTodo.title, updatedTodo.active)
            .then(() => {
                onTodoUpdate(updatedTodo);
                setIsEditing(false); // Disable editing mode after saving
            })
            .catch((error) => console.error('Error updating todo:', error));
    };

    const handleCancelEdit = () => {
        setNewTitle(todo.title); // Reset the title input
        setIsEditing(false);
    };

    return (
        <li key={todo.id}>
            {isEditing ? (
                <div>
                    <input
                        type="text"
                        value={newTitle}
                        onChange={(e) => setNewTitle(e.target.value)}
                    />
                    <button onClick={handleSaveEdit}>Save</button>
                    <button onClick={handleCancelEdit}>Cancel</button>
                </div>
            ) : (
                <div>
                    {todo.title} (Active: {todo.active ? 'Yes' : 'No'})
                    <button onClick={handleToggleActive}>
                        {todo.active ? 'Deactivate' : 'Activate'}
                    </button>
                    <button onClick={handleEdit}>Edit</button>
                    <button onClick={handleDeleteTodo}>Delete</button>
                </div>
            )}
        </li>
    );
};

export default TodoItem;
