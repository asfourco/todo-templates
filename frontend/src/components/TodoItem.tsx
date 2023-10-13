import React, { useState } from 'react';
import { Todo } from '../models/Todo';
import { updateTodo, deleteTodo } from '../services/TodoService';
import {
    TextField,
    Typography,
    Checkbox,
    ListItem,
    ListItemButton,
    ListItemText, ListItemIcon
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import SaveIcon from '@mui/icons-material/Save';
import CancelIcon from '@mui/icons-material/Cancel';

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
        <ListItem
            key={todo.id}
            alignItems="flex-start"
            disablePadding
        >
            {isEditing ? (
                <ListItemButton>
                    <ListItemText>
                    <TextField
                        label="Edit Title"
                        value={newTitle}
                        onChange={(e) => setNewTitle(e.target.value)}
                        fullWidth
                    />
                    </ListItemText>

                    <ListItemIcon onClick={handleSaveEdit}>
                                <SaveIcon />
                        </ListItemIcon>
                    <ListItemIcon onClick={handleCancelEdit}>
                                <CancelIcon />
                        </ListItemIcon>
                </ListItemButton>
            ) : (

                    <ListItemButton>
                        <ListItemText>
                            <Typography variant="h6">{todo.title}</Typography>
                        </ListItemText>
                        <ListItemIcon>
                            <Checkbox
                                checked={!todo.active}
                                onChange={handleToggleActive}
                                inputProps={{ 'aria-label': 'Toggle Active' }}
                            />
                        </ListItemIcon>
                        <ListItemIcon onClick={handleEdit}>
                            <EditIcon />
                        </ListItemIcon>
                        <ListItemIcon color="error" onClick={handleDeleteTodo}>
                            <DeleteIcon />
                        </ListItemIcon>
                    </ListItemButton>
            )}
        </ListItem>
    );
};

export default TodoItem;
