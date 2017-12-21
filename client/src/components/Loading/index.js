import React from 'react';
import ProgressBar from 'react-toolbox/lib/progress_bar/ProgressBar';

import './style.css';

const Loading = (props) => (
  <div {...props}>
    <ProgressBar className="loading" type="circular" mode="indeterminate" />
  </div>
);

export default Loading;
