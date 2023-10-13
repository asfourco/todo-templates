import React, { useEffect, useState } from 'react';
import { Todo } from '../models/Todo';
import { getTodos } from '../services/TodoService';
import TodoItem from './TodoItem';
import TodoForm from "./TodoForm";
import { List, ListItem, Button, Paper } from '@mui/material';

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
        <Paper elevation={1} style={{ padding: '16px', width: '80%', margin: '0 auto' }}>
            <h1> My Lists of things ...</h1>
            <TodoForm todos={todos} setTodos={setTodos} />
            <List>
                {todos.map((todo) => (
                    <ListItem key={todo.id} alignItems="flex-start">
                        <TodoItem
                            todo={todo}
                            onTodoUpdate={handleTodoUpdate}
                            onTodoDelete={handleTodoDelete}
                        />
                    </ListItem>
                ))}
            </List>
            <div>
                <Button
                    variant="contained"
                    onClick={() => setPage(page - 1)}
                    disabled={page === 0}
                >
                    Previous Page
                </Button>
                <Button
                    variant="contained"
                    onClick={() => setPage(page + 1)}
                    disabled={todos.length < pageSize}
                >
                    Next Page
                </Button>
            </div>
        </Paper>
    );
};

export default TodoList;
