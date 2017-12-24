import React from 'react';
import ProgressBar from 'react-toolbox/lib/progress_bar/ProgressBar';

import './style.css';

const Loading = ({ center }) => (
  <div className={center ? "loading__center" : ""}>
    <ProgressBar className="loading" type="circular" mode="indeterminate" />
  </div>
);

export default Loading;
