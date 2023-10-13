import React, { useState } from 'react';
import { createTodo } from '../services/TodoService';
import { Button, TextField, Grid } from '@mui/material';
import { makeStyles } from '@mui/styles';
import {Todo} from "../models/Todo";

const useStyles = makeStyles((theme) => ({
    formContainer: {
        display: 'flex',
        alignItems: 'center',
    },
    inputField: {
        marginRight: theme.spacing(2),
    },
}));

interface TodoFormProps {
    todos: Todo[];
    setTodos: React.Dispatch<React.SetStateAction<Todo[]>>;
}

const TodoForm: React.FC<TodoFormProps> = ({todos, setTodos}) => {
    const [title, setTitle] = useState<string>('');
    const classes = useStyles();


    const handleAddTodo = () => {
        createTodo(title)
            .then((newTodoResponse) => {
                // Add the new todo to the list
                setTodos([...todos, newTodoResponse.data]);

                setTitle(''); // Clear the input field
            })
            .catch((error) => console.error('Error adding todo:', error));
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleAddTodo(); // Call handleAddTodo when the Enter key is pressed
        }
    };

    return (
        <Grid container className={classes.formContainer}>
            <Grid item className={classes.inputField}>
                <TextField
                    type="text"
                    label="New Todo"
                    variant="outlined"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    onKeyPress={handleKeyPress}
                />
            </Grid>
            <Grid item>
                <Button variant="contained" color="primary" onClick={handleAddTodo}>
                    Add Todo
                </Button>
            </Grid>
        </Grid>
    );
};

export default TodoForm;
