import axios from 'axios';
import {Todo} from "../models/Todo";

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:8080';
const API_URL = `${BACKEND_URL}/api/v1/todos`;

export const getTodos = async (page: number, pageSize: number):Promise<Todo[]> =>
    await axios.get(API_URL, {
        params: { page, page_size: pageSize },
    });
export const getTodo = async (id: number):Promise<Todo> => await axios.get(`${API_URL}/${id}`);
export const createTodo = async (title: string) => await axios.post(API_URL, { title });
export const updateTodo = async (id: number, title: string, active: boolean) =>
    await axios.put(API_URL, { id, title, active });
export const deleteTodo = async (id: number) => await axios.delete(`${API_URL}/${id}`);
