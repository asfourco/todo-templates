import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1/todos';

export const getTodos = async (page: number, pageSize: number) =>
    await axios.get(API_URL, {
        params: { page, page_size: pageSize },
    });
export const getTodo = async (id: number) => await axios.get(`${API_URL}/${id}`);
export const createTodo = async (title: string) => await axios.post(API_URL, { title });
export const updateTodo = async (id: number, title: string, active: boolean) =>
    await axios.put(API_URL, { id, title, active });
export const deleteTodo = async (id: number) => await axios.delete(`${API_URL}/${id}`);
