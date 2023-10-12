import React, { useState } from 'react';
import { createTodo } from '../services/TodoService';

const TodoForm: React.FC = () => {
    const [title, setTitle] = useState<string>('');

    const handleAddTodo = () => {
        createTodo(title)
            .then(() => {
                // Handle successful creation, e.g., update the list of todos
                setTitle(''); // Clear the input field
            })
            .catch((error) => console.error('Error adding todo:', error));
    };

    return (
        <div>
            <input
                type="text"
                placeholder="New Todo"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />
            <button onClick={handleAddTodo}>Add Todo</button>
        </div>
    );
};

export default TodoForm;
