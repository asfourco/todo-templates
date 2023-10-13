import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

describe('Todo List App', () => {
  it('renders app title on page load', () => {
    render(<App />);
    const appTitle = screen.getByText(/Todo List/i);
    expect(appTitle).toBeInTheDocument();
  });

});
