import React from 'react';
import TodoList from './components/TodoList';
import {Paper} from "@mui/material";
import { ThemeProvider} from "@mui/material/styles";
import theme from "./theme";

function App() {
    return (
        <ThemeProvider theme={theme}>
            <div className="App">
                    <h1>Todo List</h1>
                    <TodoList />
            </div>
        </ThemeProvider>
    );
}

export default App;
