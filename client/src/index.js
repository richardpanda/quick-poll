import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import ThemeProvider from 'react-toolbox/lib/ThemeProvider';

import './toolbox/theme.css';

import App from './App';
import theme from './toolbox/theme.js';
import registerServiceWorker from './registerServiceWorker';

ReactDOM.render(
  <ThemeProvider theme={theme}>
    <BrowserRouter> 
      <App />
    </BrowserRouter>
  </ThemeProvider>,
  document.getElementById('root')
);
registerServiceWorker();
