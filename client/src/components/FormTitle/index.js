import React from 'react';
import CardTitle from 'react-toolbox/lib/card/CardTitle';

const FormTitle = ({ title, subtitle }) => (
  <CardTitle style={{ paddingBottom: "10px" }}>
    <div style={{ fontSize: "24px" }}>{title}</div>
    <div style={{ color: "red" }}>{subtitle}</div>
  </CardTitle>
);

export default FormTitle;
