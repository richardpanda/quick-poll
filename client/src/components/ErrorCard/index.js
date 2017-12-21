import React from 'react';
import Card from 'react-toolbox/lib/card/Card';
import CardTitle from 'react-toolbox/lib/card/CardTitle';

const ErrorCard = ({ message }) => (
  <Card style={{ margin: "16px auto", maxWidth: "600px" }}>
    <CardTitle style={{ display: "flex", justifyContent: "center" }}>
      <div style={{ color: "red", fontSize: "16px", textAlign: "center" }}>
        {message}
      </div>
    </CardTitle>
  </Card>
);

export default ErrorCard;
