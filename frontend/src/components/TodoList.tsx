import React, { useEffect, useState } from 'react';
import { Todo } from '../models/Todo';
import { getTodos } from '../services/TodoService';
import TodoItem from './TodoItem';

const TodoList: React.FC = () => {
    const [page, setPage] = useState<number>(0);
    const [pageSize, setPageSize] = useState<number>(10); // Default page size
    const [todos, setTodos] = useState<Todo[]>([]);
    const [loading, setLoading] = useState(true);

    const loadTodos = () => {
        getTodos(page, pageSize)
            .then((response) => {
                setTodos(response.data.items || []);
                setPageSize(response.data.page_size);
                setLoading(false); // Mark the loading as complete
            })
            .catch((error) => {
                console.error('Error fetching todos:', error);
                setLoading(false); // Ensure loading is marked as complete in case of an error
            });
    };

    useEffect(() => {
        loadTodos();
    }, [page, pageSize]);

    if (loading) {
        return <div>Loading...</div>;
    }

    const handleTodoUpdate = (updatedTodo: Todo) => {
        const updatedTodos = todos.map((todo) =>
            todo.id === updatedTodo.id ? updatedTodo : todo
        );
        setTodos(updatedTodos);
    };

    const handleTodoDelete = (deletedTodoId: number) => {
        const updatedTodos = todos.filter((todo) => todo.id !== deletedTodoId);
        setTodos(updatedTodos);
    };

    return (
        <div>
            <ul>
                {todos.map((todo) => (
                    <TodoItem
                        key={todo.id}
                        todo={todo}
                        onTodoUpdate={handleTodoUpdate}
                        onTodoDelete={handleTodoDelete}
                    />
                ))}
            </ul>
            <div>
                <button
                    onClick={() => setPage(page - 1)}
                    disabled={page === 1}
                >
                    Previous Page
                </button>
                <button onClick={() => setPage(page + 1)}>Next Page</button>
            </div>
        </div>
    );
};

export default TodoList;
